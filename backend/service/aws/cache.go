package aws

import (
	"context"
	"errors"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/autoscaling"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/kinesis"
	"github.com/golang/protobuf/ptypes"
	"go.uber.org/zap"

	topologyv1 "github.com/lyft/clutch/backend/api/topology/v1"
)

const topologyObjectChanBufferSize = 5000
const topologyLockId = 1

func (c *client) CacheEnabled() bool {
	return true
}

func (c *client) StartTopologyCaching(ctx context.Context) (<-chan *topologyv1.UpdateCacheRequest, error) {
	if !c.topologyLock.TryAcquire(topologyLockId) {
		return nil, errors.New("TopologyCaching is already in progress")
	}

	go func() {
		<-ctx.Done()
		close(c.topologyObjectChan)
		c.topologyLock.Release(topologyLockId)
	}()

	c.processRegionTopologyObjects(ctx)
	return c.topologyObjectChan, nil
}

// We process each resource and region individually, creating tickers for each.
// If a single resource and region takes longer to process it will not block other resources.
// Additionally this give us the flexibility to tune the frequency based on resource.
func (c *client) processRegionTopologyObjects(ctx context.Context) {
	for name, client := range c.clients {
		c.log.Info("processing topology objects for region", zap.String("region", name))
		go c.startTickerForCacheResource(ctx, time.Duration(time.Minute*5), client, c.processAllEC2Instances)
		go c.startTickerForCacheResource(ctx, time.Duration(time.Minute*10), client, c.processAllAutoScalingGroups)
		go c.startTickerForCacheResource(ctx, time.Duration(time.Minute*10), client, c.processAllKinesisStreams)
	}
}

func (c *client) startTickerForCacheResource(ctx context.Context, duration time.Duration, client *regionalClient, resourceFunc func(context.Context, *regionalClient)) {
	ticker := time.NewTicker(duration)
	for {
		resourceFunc(ctx, client)

		select {
		case <-ticker.C:
			continue
		case <-ctx.Done():
			ticker.Stop()
			return
		}
	}
}

func (c *client) processAllAutoScalingGroups(ctx context.Context, client *regionalClient) {
	c.log.Info("starting to process auto scaling groups for region", zap.String("region", client.region))
	// 100 is the maximum amount of records per page allowed for this API
	input := autoscaling.DescribeAutoScalingGroupsInput{
		MaxRecords: aws.Int64(100),
	}

	err := client.autoscaling.DescribeAutoScalingGroupsPagesWithContext(ctx, &input,
		func(page *autoscaling.DescribeAutoScalingGroupsOutput, lastPage bool) bool {
			for _, asg := range page.AutoScalingGroups {
				protoAsg, err := ptypes.MarshalAny(newProtoForAutoscalingGroup(asg))
				if err != nil {
					c.log.Error("unable to marshal asg proto", zap.Error(err))
					continue
				}

				c.topologyObjectChan <- &topologyv1.UpdateCacheRequest{
					Resource: &topologyv1.Resource{
						Id: *asg.AutoScalingGroupName,
						Pb: protoAsg,
					},
					Action: topologyv1.UpdateCacheRequest_CREATE_OR_UPDATE,
				}
			}

			return !lastPage
		})

	if err != nil {
		c.log.Error("unable to process auto scaling groups for region", zap.String("region", client.region), zap.Error(err))
	}
}

func (c *client) processAllEC2Instances(ctx context.Context, client *regionalClient) {
	c.log.Info("starting to process ec2 instances for region", zap.String("region", client.region))
	// 1000 is the maximum amount of records per page allowed for this API
	input := ec2.DescribeInstancesInput{
		MaxResults: aws.Int64(1000),
	}

	err := client.ec2.DescribeInstancesPagesWithContext(ctx, &input,
		func(page *ec2.DescribeInstancesOutput, lastPage bool) bool {
			for _, reservation := range page.Reservations {
				for _, instance := range reservation.Instances {
					protoInstance, err := ptypes.MarshalAny(newProtoForInstance(instance))
					if err != nil {
						c.log.Error("unable to marshal instance proto", zap.Error(err))
						continue
					}
					c.topologyObjectChan <- &topologyv1.UpdateCacheRequest{
						Resource: &topologyv1.Resource{
							Id: *instance.InstanceId,
							Pb: protoInstance,
						},
						Action: topologyv1.UpdateCacheRequest_CREATE_OR_UPDATE,
					}
				}
			}

			return !lastPage
		})

	if err != nil {
		c.log.Error("unable to process ec2 instances for region", zap.String("region", client.region), zap.Error(err))
	}
}

func (c *client) processAllKinesisStreams(ctx context.Context, client *regionalClient) {
	c.log.Info("starting to process auto scaling groups for region", zap.String("region", client.region))
	// 100 is arbatrary, currently this API does not have a per page limit,
	// looking at other aws API limits this value felt safe.
	input := kinesis.ListStreamsInput{
		Limit: aws.Int64(100),
	}

	err := client.kinesis.ListStreamsPagesWithContext(ctx, &input,
		func(page *kinesis.ListStreamsOutput, lastPage bool) bool {
			for _, stream := range page.StreamNames {
				v1Stream, err := c.DescribeKinesisStream(ctx, client.region, *stream)
				if err != nil {
					c.log.Error("unable to describe kinesis stream", zap.Error(err))
					continue
				}

				protoStream, _ := ptypes.MarshalAny(v1Stream)
				if err != nil {
					c.log.Error("unable to marshal kinesis stream", zap.Error(err))
					continue
				}

				c.topologyObjectChan <- &topologyv1.UpdateCacheRequest{
					Resource: &topologyv1.Resource{
						Id: v1Stream.StreamName,
						Pb: protoStream,
					},
					Action: topologyv1.UpdateCacheRequest_CREATE_OR_UPDATE,
				}
			}

			return !lastPage
		})

	if err != nil {
		c.log.Error("unable to process kinesis streams for region", zap.String("region", client.region), zap.Error(err))
	}
}

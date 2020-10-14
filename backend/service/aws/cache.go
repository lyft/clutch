package aws

import (
	"context"
	"errors"
	"sync"
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
		return nil, errors.New("TopologyCahing is already in progress")
	}

	go func() {
		<-ctx.Done()
		c.topologyLock.Release(topologyLockId)
	}()

	go c.processRegionTopologyObjects(ctx)
	return c.topologyObjectChan, nil
}

func (c *client) processRegionTopologyObjects(ctx context.Context) {
	ticker := time.NewTicker(time.Minute * 10)
	for ; true; <-ticker.C {
		wg := &sync.WaitGroup{}
		for name, client := range c.clients {
			c.log.Info("processing topology objects for region", zap.String("region", name))
			go c.processAllEC2Instances(ctx, client, wg)
			go c.processAllAutoScalingGroups(ctx, client, wg)
			go c.processAllKinesisStreams(ctx, client, wg)
		}
		wg.Wait()
		c.log.Info("finished processing all aws topology objects, waiting for next tick.")
	}
}

func (c *client) processAllAutoScalingGroups(ctx context.Context, client *regionalClient, wg *sync.WaitGroup) {
	wg.Add(1)

	// 100 is the maximum amount of records allowed for this API
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
		c.log.Error("unable to process auto scaling groups", zap.Error(err))
	}
	wg.Done()
}

func (c *client) processAllEC2Instances(ctx context.Context, client *regionalClient, wg *sync.WaitGroup) {
	wg.Add(1)

	// 1000 is the maximum amount of records allowed for this API
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
		c.log.Error("unable to process ec2 instances", zap.Error(err))
	}
	wg.Done()
}

func (c *client) processAllKinesisStreams(ctx context.Context, client *regionalClient, wg *sync.WaitGroup) {
	wg.Add(1)

	// 100 is arbatrary, currently this API does not have a limit,
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
		c.log.Error("unable to process kinesis streams", zap.Error(err))
	}
	wg.Done()
}

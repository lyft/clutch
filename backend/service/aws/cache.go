package aws

import (
	"context"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/autoscaling"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/golang/protobuf/ptypes"
	topologyv1 "github.com/lyft/clutch/backend/api/topology/v1"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
)

func (c *client) CacheEnabled() bool {
	return true
}

func (c *client) StartTopologyCaching(ctx context.Context) <-chan *topologyv1.UpdateCacheRequest {
	for name, client := range c.clients {
		c.log.Info("starting topology caching for region", zap.String("region", name))
		go c.processRegionTopologyObjects(ctx, client)
	}

	return c.topologyObjectChan
}

func (c *client) processRegionTopologyObjects(ctx context.Context, client *regionalClient) {
	ticker := time.NewTicker(time.Minute * 30)

	for ; true; <-ticker.C {
		eg, _ := errgroup.WithContext(context.Background())

		eg.Go(func() error { return c.processAllEC2Instances(ctx, client) })
		eg.Go(func() error { return c.processAllAutoScalingGroups(ctx, client) })
		// eg.Go(func() error { return c.getAllKinesisStreams(ctx, client) })

		if err := eg.Wait(); err != nil {
			c.log.Error("errors", zap.Error(err))
		}
	}
}

func (c *client) processAllAutoScalingGroups(ctx context.Context, client *regionalClient) error {
	input := autoscaling.DescribeAutoScalingGroupsInput{
		MaxRecords: aws.Int64(100),
	}

	return client.autoscaling.DescribeAutoScalingGroupsPagesWithContext(ctx, &input,
		func(page *autoscaling.DescribeAutoScalingGroupsOutput, lastPage bool) bool {

			for _, asg := range page.AutoScalingGroups {
				protoAsg, _ := ptypes.MarshalAny(newProtoForAutoscalingGroup(asg))
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
}

func (c *client) processAllEC2Instances(ctx context.Context, client *regionalClient) error {
	input := ec2.DescribeInstancesInput{
		MaxResults: aws.Int64(1000),
	}

	return client.ec2.DescribeInstancesPagesWithContext(ctx, &input,
		func(page *ec2.DescribeInstancesOutput, lastPage bool) bool {

			for _, reservation := range page.Reservations {
				for _, instance := range reservation.Instances {
					protoInstance, _ := ptypes.MarshalAny(newProtoForInstance(instance))
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
}

// func (c *client) processAllKinesisStreams(ctx context.Context, client *regionalClient) error {
// 	input := kinesis.DescribeStreamInput{
// 		Max
// 	}
// 	return nil
// }

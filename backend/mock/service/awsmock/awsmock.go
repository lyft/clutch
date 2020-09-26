package awsmock

import (
	"context"
	"fmt"
	"io"
	"math/rand"

	"github.com/golang/protobuf/ptypes/any"
	"github.com/uber-go/tally"
	"go.uber.org/zap"

	ec2v1 "github.com/lyft/clutch/backend/api/aws/ec2/v1"
	kinesisv1 "github.com/lyft/clutch/backend/api/aws/kinesis/v1"
	"github.com/lyft/clutch/backend/service"
	"github.com/lyft/clutch/backend/service/aws"
)

type svc struct{}

func (s *svc) DescribeKinesisStream(ctx context.Context, region string, streamName string) (*kinesisv1.Stream, error) {
	ret := &kinesisv1.Stream{
		StreamName:        streamName,
		Region:            region,
		CurrentShardCount: uint32(32),
	}
	return ret, nil
}

func (s *svc) UpdateKinesisShardCount(ctx context.Context, region string, streamName string, targetShardCount uint32) error {
	return nil
}

func (s *svc) ResizeAutoscalingGroup(ctx context.Context, region string, name string, size *ec2v1.AutoscalingGroupSize) error {
	return nil
}

func (s *svc) DescribeAutoscalingGroups(ctx context.Context, region string, names []string) ([]*ec2v1.AutoscalingGroup, error) {
	ret := make([]*ec2v1.AutoscalingGroup, len(names))
	for idx, name := range names {
		ret[idx] = &ec2v1.AutoscalingGroup{
			Name:   name,
			Region: region,
			Zones:  []string{region + "a", region + "b"},
			Size: &ec2v1.AutoscalingGroupSize{
				Min:     25,
				Max:     200,
				Desired: 100,
			},
			TerminationPolicies: []ec2v1.AutoscalingGroup_TerminationPolicy{
				ec2v1.AutoscalingGroup_OLDEST_INSTANCE,
			},
			Instances: []*ec2v1.AutoscalingGroup_Instance{
				{
					Id:                      "i-12345892010",
					Zone:                    region + "b",
					LaunchConfigurationName: "myLaunchConfig",
					Healthy:                 true,
					LifecycleState:          ec2v1.AutoscalingGroup_Instance_IN_SERVICE,
				},
				{
					Id:                      "i-22345892010",
					Zone:                    region + "a",
					LaunchConfigurationName: "myLaunchConfig",
					Healthy:                 false,
					LifecycleState:          ec2v1.AutoscalingGroup_Instance_PENDING,
				},
			},
		}
	}

	return ret, nil
}

func (s *svc) DescribeInstances(ctx context.Context, region string, ids []string) ([]*ec2v1.Instance, error) {
	var ret []*ec2v1.Instance
	for _, id := range ids {
		ret = append(ret, &ec2v1.Instance{
			InstanceId:       id,
			Region:           region,
			State:            ec2v1.Instance_State(rand.Intn(len(ec2v1.Instance_State_value))),
			InstanceType:     "c3.2xlarge",
			PublicIpAddress:  "8.1.1.1",
			PrivateIpAddress: "10.1.1.1",
			AvailabilityZone: fmt.Sprintf("%sz", region),
			Tags:             map[string]string{"Name": fmt.Sprintf("my-instance-%s", id)},
		})
	}
	return ret, nil
}

func (s *svc) TerminateInstances(ctx context.Context, region string, ids []string) error {
	return nil
}

func (s *svc) S3StreamingGet(ctx context.Context, region string, bucket string, key string) (io.ReadCloser, error) {
	return nil, nil
}

func (s *svc) Regions() []string {
	return []string{"us-mock-1"}
}

func (s *svc) GetCurrentRegion() string {
	return "us-east-1"
}

func New() aws.Client {
	return &svc{}
}

func NewAsService(*any.Any, *zap.Logger, tally.Scope) (service.Service, error) {
	return New(), nil
}

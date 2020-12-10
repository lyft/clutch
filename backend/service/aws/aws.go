package aws

// <!-- START clutchdoc -->
// description: Multi-region client for Amazon Web Services.
// <!-- END clutchdoc -->

import (
	"context"
	"fmt"
	"io"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/autoscaling"
	"github.com/aws/aws-sdk-go/service/autoscaling/autoscalingiface"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/ec2/ec2iface"
	"github.com/aws/aws-sdk-go/service/kinesis"
	"github.com/aws/aws-sdk-go/service/kinesis/kinesisiface"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3iface"
	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/any"
	"github.com/iancoleman/strcase"
	"github.com/uber-go/tally"
	"go.uber.org/zap"
	"golang.org/x/sync/semaphore"

	ec2v1 "github.com/lyft/clutch/backend/api/aws/ec2/v1"
	kinesisv1 "github.com/lyft/clutch/backend/api/aws/kinesis/v1"
	awsv1 "github.com/lyft/clutch/backend/api/config/service/aws/v1"
	topologyv1 "github.com/lyft/clutch/backend/api/topology/v1"
	"github.com/lyft/clutch/backend/service"
)

const (
	Name = "clutch.service.aws"
)

func New(cfg *any.Any, logger *zap.Logger, scope tally.Scope) (service.Service, error) {
	ac := &awsv1.Config{}
	err := ptypes.UnmarshalAny(cfg, ac)
	if err != nil {
		return nil, err
	}

	c := &client{
		clients:            make(map[string]*regionalClient, len(ac.Regions)),
		topologyObjectChan: make(chan *topologyv1.UpdateCacheRequest, topologyObjectChanBufferSize),
		topologyLock:       semaphore.NewWeighted(1),
		log:                logger,
		scope:              scope,
	}

	for _, region := range ac.Regions {
		regionCfg := aws.NewConfig().WithRegion(region).WithMaxRetries(50)
		awsSession, err := session.NewSession(regionCfg)
		if err != nil {
			return nil, err
		}

		c.clients[region] = &regionalClient{
			region:      region,
			s3:          s3.New(awsSession),
			kinesis:     kinesis.New(awsSession),
			ec2:         ec2.New(awsSession),
			autoscaling: autoscaling.New(awsSession),
		}
	}

	return c, nil
}

type Client interface {
	DescribeInstances(ctx context.Context, region string, ids []string) ([]*ec2v1.Instance, error)
	TerminateInstances(ctx context.Context, region string, ids []string) error
	RebootInstances(ctx context.Context, region string, ids []string) error

	DescribeAutoscalingGroups(ctx context.Context, region string, names []string) ([]*ec2v1.AutoscalingGroup, error)
	ResizeAutoscalingGroup(ctx context.Context, region string, name string, size *ec2v1.AutoscalingGroupSize) error

	DescribeKinesisStream(ctx context.Context, region string, streamName string) (*kinesisv1.Stream, error)
	UpdateKinesisShardCount(ctx context.Context, region string, streamName string, targetShardCount uint32) error

	S3StreamingGet(ctx context.Context, region string, bucket string, key string) (io.ReadCloser, error)

	Regions() []string
}

type client struct {
	clients            map[string]*regionalClient
	topologyObjectChan chan *topologyv1.UpdateCacheRequest
	topologyLock       *semaphore.Weighted
	log                *zap.Logger
	scope              tally.Scope
}

func (c *client) ResizeAutoscalingGroup(ctx context.Context, region string, name string, size *ec2v1.AutoscalingGroupSize) error {
	rc, ok := c.clients[region]
	if !ok {
		return fmt.Errorf("no client found for region '%s'", region)
	}

	input := &autoscaling.UpdateAutoScalingGroupInput{
		AutoScalingGroupName: aws.String(name),
		DesiredCapacity:      aws.Int64(int64(size.Desired)),
		MaxSize:              aws.Int64(int64(size.Max)),
		MinSize:              aws.Int64(int64(size.Min)),
	}

	_, err := rc.autoscaling.UpdateAutoScalingGroupWithContext(ctx, input)
	return err
}

func (c *client) DescribeAutoscalingGroups(ctx context.Context, region string, names []string) ([]*ec2v1.AutoscalingGroup, error) {
	cl, ok := c.clients[region]
	if !ok {
		return nil, fmt.Errorf("no client found for region '%s'", region)
	}

	input := &autoscaling.DescribeAutoScalingGroupsInput{
		AutoScalingGroupNames: aws.StringSlice(names),
	}
	result, err := cl.autoscaling.DescribeAutoScalingGroupsWithContext(ctx, input)
	if err != nil {
		return nil, err
	}

	ret := make([]*ec2v1.AutoscalingGroup, len(result.AutoScalingGroups))
	for idx, group := range result.AutoScalingGroups {
		ret[idx] = newProtoForAutoscalingGroup(group)
	}
	return ret, nil
}

// Shave off the trailing zone identifier to get the region
func zoneToRegion(zone string) string {
	if zone == "" {
		return "UNKNOWN"
	}
	return zone[:len(zone)-1]
}

func protoForTerminationPolicy(policy string) ec2v1.AutoscalingGroup_TerminationPolicy {
	policy = strcase.ToScreamingSnake(policy)
	val, ok := ec2v1.AutoscalingGroup_TerminationPolicy_value[policy]
	if !ok {
		return ec2v1.AutoscalingGroup_UNKNOWN
	}
	return ec2v1.AutoscalingGroup_TerminationPolicy(val)
}

func protoForAutoscalingGroupInstanceLifecycleState(state string) ec2v1.AutoscalingGroup_Instance_LifecycleState {
	state = strcase.ToScreamingSnake(strings.ReplaceAll(state, ":", ""))
	val, ok := ec2v1.AutoscalingGroup_Instance_LifecycleState_value[state]
	if !ok {
		return ec2v1.AutoscalingGroup_Instance_UNKNOWN
	}
	return ec2v1.AutoscalingGroup_Instance_LifecycleState(val)
}

func newProtoForAutoscalingGroupInstance(instance *autoscaling.Instance) *ec2v1.AutoscalingGroup_Instance {
	return &ec2v1.AutoscalingGroup_Instance{
		Id:                      aws.StringValue(instance.InstanceId),
		Zone:                    aws.StringValue(instance.AvailabilityZone),
		LaunchConfigurationName: aws.StringValue(instance.LaunchConfigurationName),
		Healthy:                 aws.StringValue(instance.HealthStatus) == "HEALTHY",
		LifecycleState:          protoForAutoscalingGroupInstanceLifecycleState(aws.StringValue(instance.LifecycleState)),
	}
}

func newProtoForAutoscalingGroup(group *autoscaling.Group) *ec2v1.AutoscalingGroup {
	pb := &ec2v1.AutoscalingGroup{
		Name:  aws.StringValue(group.AutoScalingGroupName),
		Zones: aws.StringValueSlice(group.AvailabilityZones),
		Size: &ec2v1.AutoscalingGroupSize{
			Min:     uint32(aws.Int64Value(group.MinSize)),
			Max:     uint32(aws.Int64Value(group.MaxSize)),
			Desired: uint32(aws.Int64Value(group.DesiredCapacity)),
		},
	}

	if len(pb.Zones) > 0 {
		pb.Region = zoneToRegion(pb.Zones[0])
	}

	pb.TerminationPolicies = make([]ec2v1.AutoscalingGroup_TerminationPolicy, len(group.TerminationPolicies))
	for idx, p := range aws.StringValueSlice(group.TerminationPolicies) {
		pb.TerminationPolicies[idx] = protoForTerminationPolicy(p)
	}

	pb.Instances = make([]*ec2v1.AutoscalingGroup_Instance, len(group.Instances))
	for idx, i := range group.Instances {
		pb.Instances[idx] = newProtoForAutoscalingGroupInstance(i)
	}

	return pb
}

func (c *client) Regions() []string {
	regions := make([]string, 0, len(c.clients))
	for region := range c.clients {
		regions = append(regions, region)
	}
	return regions
}

type regionalClient struct {
	region string

	s3          s3iface.S3API
	kinesis     kinesisiface.KinesisAPI
	ec2         ec2iface.EC2API
	autoscaling autoscalingiface.AutoScalingAPI
}

func (c *client) DescribeInstances(ctx context.Context, region string, ids []string) ([]*ec2v1.Instance, error) {
	cl, ok := c.clients[region]
	if !ok {
		return nil, fmt.Errorf("no client found for region '%s'", region)
	}

	input := &ec2.DescribeInstancesInput{InstanceIds: aws.StringSlice(ids)}
	// 20 transactions per second account wide
	result, err := cl.ec2.DescribeInstancesWithContext(ctx, input)

	if err != nil {
		return nil, err
	}

	var ret []*ec2v1.Instance
	for _, r := range result.Reservations {
		for _, i := range r.Instances {
			ret = append(ret, newProtoForInstance(i))
		}
	}

	return ret, nil
}

func (c *client) TerminateInstances(ctx context.Context, region string, ids []string) error {
	cl, ok := c.clients[region]
	if !ok {
		return fmt.Errorf("no client found for region '%s'", region)
	}

	input := &ec2.TerminateInstancesInput{InstanceIds: aws.StringSlice(ids)}
	_, err := cl.ec2.TerminateInstancesWithContext(ctx, input)

	return err
}

func (c *client) RebootInstances(ctx context.Context, region string, ids []string) error {
	cl, ok := c.clients[region]
	if !ok {
		return fmt.Errorf("no client found for region '%s'", region)
	}

	input := &ec2.RebootInstancesInput{InstanceIds: aws.StringSlice(ids)}
	_, err := cl.ec2.RebootInstancesWithContext(ctx, input)

	return err
}

func protoForInstanceState(state string) ec2v1.Instance_State {
	// Transform kebab case 'shutting-down' to upper snake case 'SHUTTING_DOWN'.
	state = strings.ReplaceAll(strings.ToUpper(state), "-", "_")
	// Look up value in generated enum map.
	val, ok := ec2v1.Instance_State_value[state]
	if !ok {
		return ec2v1.Instance_UNKNOWN
	}

	return ec2v1.Instance_State(val)
}

func newProtoForInstance(i *ec2.Instance) *ec2v1.Instance {
	ret := &ec2v1.Instance{
		InstanceId:       aws.StringValue(i.InstanceId),
		State:            protoForInstanceState(aws.StringValue(i.State.Name)),
		InstanceType:     aws.StringValue(i.InstanceType),
		PublicIpAddress:  aws.StringValue(i.PublicIpAddress),
		PrivateIpAddress: aws.StringValue(i.PrivateIpAddress),
		AvailabilityZone: aws.StringValue(i.Placement.AvailabilityZone),
	}

	ret.Region = zoneToRegion(ret.AvailabilityZone)

	// Transform tag list to map.
	ret.Tags = make(map[string]string, len(i.Tags))
	for _, tag := range i.Tags {
		ret.Tags[aws.StringValue(tag.Key)] = aws.StringValue(tag.Value)
	}

	return ret
}

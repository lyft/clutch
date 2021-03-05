package aws

// <!-- START clutchdoc -->
// description: Multi-region client for Amazon Web Services.
// <!-- END clutchdoc -->

import (
	"context"
	"io"
	"net/http"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/aws/retry"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/autoscaling"
	astypes "github.com/aws/aws-sdk-go-v2/service/autoscaling/types"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	ec2types "github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/aws/aws-sdk-go-v2/service/kinesis"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/iancoleman/strcase"
	"github.com/uber-go/tally"
	"go.uber.org/zap"
	"golang.org/x/sync/semaphore"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/anypb"

	ec2v1 "github.com/lyft/clutch/backend/api/aws/ec2/v1"
	kinesisv1 "github.com/lyft/clutch/backend/api/aws/kinesis/v1"
	awsv1 "github.com/lyft/clutch/backend/api/config/service/aws/v1"
	topologyv1 "github.com/lyft/clutch/backend/api/topology/v1"
	"github.com/lyft/clutch/backend/service"
)

const (
	Name = "clutch.service.aws"
)

func New(cfg *anypb.Any, logger *zap.Logger, scope tally.Scope) (service.Service, error) {
	ac := &awsv1.Config{}
	err := cfg.UnmarshalTo(ac)
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

	clientRetries := 0
	if ac.ClientConfig != nil && ac.ClientConfig.Retries >= 0 {
		clientRetries = int(ac.ClientConfig.Retries)
	}

	awsHTTPClient := &http.Client{}

	for _, region := range ac.Regions {
		regionCfg, err := config.LoadDefaultConfig(context.TODO(),
			config.WithHTTPClient(awsHTTPClient),
			config.WithRegion(region),
			config.WithRetryer(func() aws.Retryer {
				customRetryer := retry.NewStandard(func(so *retry.StandardOptions) {
					so.MaxAttempts = clientRetries
				})
				return customRetryer
			}),
		)
		if err != nil {
			return nil, err
		}

		c.clients[region] = &regionalClient{
			region:      region,
			s3:          s3.NewFromConfig(regionCfg),
			kinesis:     kinesis.NewFromConfig(regionCfg),
			ec2:         ec2.NewFromConfig(regionCfg),
			autoscaling: autoscaling.NewFromConfig(regionCfg),
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
	UpdateKinesisShardCount(ctx context.Context, region string, streamName string, targetShardCount int32) error

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

type regionalClient struct {
	region string

	s3          s3Client
	kinesis     kinesisClient
	ec2         ec2Client
	autoscaling autoscalingClient
}

// Implement the interface provided by errorintercept, so errors are caught at middleware and converted to gRPC status.
func (c *client) InterceptError(e error) error {
	return ConvertError(e)
}

func (c *client) getRegionalClient(region string) (*regionalClient, error) {
	rc, ok := c.clients[region]
	if !ok {
		return nil, status.Errorf(codes.NotFound, "no client found for region '%s'", region)
	}
	return rc, nil
}

func (c *client) ResizeAutoscalingGroup(ctx context.Context, region string, name string, size *ec2v1.AutoscalingGroupSize) error {
	rc, err := c.getRegionalClient(region)
	if err != nil {
		return err
	}

	input := &autoscaling.UpdateAutoScalingGroupInput{
		AutoScalingGroupName: aws.String(name),
		DesiredCapacity:      aws.Int32(int32(size.Desired)),
		MaxSize:              aws.Int32(int32(size.Max)),
		MinSize:              aws.Int32(int32(size.Min)),
	}

	_, err = rc.autoscaling.UpdateAutoScalingGroup(ctx, input)
	return err
}

func (c *client) DescribeAutoscalingGroups(ctx context.Context, region string, names []string) ([]*ec2v1.AutoscalingGroup, error) {
	cl, err := c.getRegionalClient(region)
	if err != nil {
		return nil, err
	}

	input := &autoscaling.DescribeAutoScalingGroupsInput{
		AutoScalingGroupNames: names,
	}
	result, err := cl.autoscaling.DescribeAutoScalingGroups(ctx, input)
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

func newProtoForAutoscalingGroupInstance(instance astypes.Instance) *ec2v1.AutoscalingGroup_Instance {
	return &ec2v1.AutoscalingGroup_Instance{
		Id:                      aws.ToString(instance.InstanceId),
		Zone:                    aws.ToString(instance.AvailabilityZone),
		LaunchConfigurationName: aws.ToString(instance.LaunchConfigurationName),
		Healthy:                 aws.ToString(instance.HealthStatus) == "HEALTHY",
		LifecycleState:          protoForAutoscalingGroupInstanceLifecycleState(string(instance.LifecycleState)),
	}
}

func newProtoForAutoscalingGroup(group astypes.AutoScalingGroup) *ec2v1.AutoscalingGroup {
	pb := &ec2v1.AutoscalingGroup{
		Name:  aws.ToString(group.AutoScalingGroupName),
		Zones: group.AvailabilityZones,
		Size: &ec2v1.AutoscalingGroupSize{
			Min:     uint32(aws.ToInt32(group.MinSize)),
			Max:     uint32(aws.ToInt32(group.MaxSize)),
			Desired: uint32(aws.ToInt32(group.DesiredCapacity)),
		},
	}

	if len(pb.Zones) > 0 {
		pb.Region = zoneToRegion(pb.Zones[0])
	}

	pb.TerminationPolicies = make([]ec2v1.AutoscalingGroup_TerminationPolicy, len(group.TerminationPolicies))
	for idx, p := range group.TerminationPolicies {
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

func (c *client) DescribeInstances(ctx context.Context, region string, ids []string) ([]*ec2v1.Instance, error) {
	cl, err := c.getRegionalClient(region)
	if err != nil {
		return nil, err
	}

	input := &ec2.DescribeInstancesInput{InstanceIds: ids}
	result, err := cl.ec2.DescribeInstances(ctx, input)

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
	cl, err := c.getRegionalClient(region)
	if err != nil {
		return err
	}

	input := &ec2.TerminateInstancesInput{InstanceIds: ids}
	_, err = cl.ec2.TerminateInstances(ctx, input)

	return err
}

func (c *client) RebootInstances(ctx context.Context, region string, ids []string) error {
	cl, err := c.getRegionalClient(region)
	if err != nil {
		return err
	}

	input := &ec2.RebootInstancesInput{InstanceIds: ids}
	_, err = cl.ec2.RebootInstances(ctx, input)

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

func newProtoForInstance(i ec2types.Instance) *ec2v1.Instance {
	ret := &ec2v1.Instance{
		InstanceId:       aws.ToString(i.InstanceId),
		State:            protoForInstanceState(string(i.State.Name)),
		InstanceType:     string(i.InstanceType),
		PublicIpAddress:  aws.ToString(i.PublicIpAddress),
		PrivateIpAddress: aws.ToString(i.PrivateIpAddress),
		AvailabilityZone: aws.ToString(i.Placement.AvailabilityZone),
	}

	ret.Region = zoneToRegion(ret.AvailabilityZone)

	// Transform tag list to map.
	ret.Tags = make(map[string]string, len(i.Tags))
	for _, tag := range i.Tags {
		ret.Tags[aws.ToString(tag.Key)] = aws.ToString(tag.Value)
	}

	return ret
}

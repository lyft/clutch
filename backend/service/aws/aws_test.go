package aws

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/autoscaling"
	astypes "github.com/aws/aws-sdk-go-v2/service/autoscaling/types"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	ec2types "github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/aws/smithy-go"
	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/any"
	"github.com/stretchr/testify/assert"
	"github.com/uber-go/tally"
	"go.uber.org/zap/zaptest"
	"google.golang.org/grpc/status"

	ec2v1 "github.com/lyft/clutch/backend/api/aws/ec2/v1"
	awsv1 "github.com/lyft/clutch/backend/api/config/service/aws/v1"
)

func TestNew(t *testing.T) {
	regions := []string{"us-east-1", "us-west25"}

	cfg, _ := ptypes.MarshalAny(&awsv1.Config{
		Regions: regions,
		ClientConfig: &awsv1.ClientConfig{
			Retries: 10,
		},
	})
	log := zaptest.NewLogger(t)
	scope := tally.NewTestScope("", nil)
	s, err := New(cfg, log, scope)
	assert.NoError(t, err)

	// Test conformance to public interface.
	_, ok := s.(Client)
	assert.True(t, ok)

	// Test private interface.
	c, ok := s.(*client)
	assert.True(t, ok)

	assert.NotNil(t, c.log)
	assert.NotNil(t, c.scope)

	assert.Len(t, c.clients, len(regions))
	addedRegions := make([]string, 0, len(regions))
	for key, rc := range c.clients {
		addedRegions = append(addedRegions, key)
		assert.Equal(t, key, rc.region)
	}
	assert.ElementsMatch(t, addedRegions, regions)
}

func TestNewWithWrongConfigType(t *testing.T) {
	_, err := New(&any.Any{TypeUrl: "foo"}, nil, nil)
	assert.EqualError(t, err, "mismatched message type: got \"foo\" want \"clutch.config.service.aws.v1.Config\"")
}

func TestRegions(t *testing.T) {
	c := client{}
	c.clients = map[string]*regionalClient{"us-east-1": nil, "us-west-2": nil, "us-north-5": nil}
	assert.ElementsMatch(t, c.Regions(), []string{"us-east-1", "us-west-2", "us-north-5"})
}

func TestMissingRegionOnEachServiceCall(t *testing.T) {
	c := &client{
		clients: map[string]*regionalClient{"us-east-1": nil},
	}

	_, err := c.DescribeInstances(context.Background(), "us-north-5", nil)
	assert.EqualError(t, err, "no client found for region 'us-north-5'")

	err = c.TerminateInstances(context.Background(), "us-north-5", nil)
	assert.EqualError(t, err, "no client found for region 'us-north-5'")

	err = c.RebootInstances(context.Background(), "us-north-5", nil)
	assert.EqualError(t, err, "no client found for region 'us-north-5'")
}

var testInstance = ec2types.Instance{
	InstanceId: aws.String("i-123456789abcdef0"),
	Tags: []ec2types.Tag{
		{Key: aws.String("Name"), Value: aws.String("locations-staging-iad")},
		{Key: aws.String("Canary"), Value: aws.String("false")},
	},
	LaunchTime:       aws.Time(time.Unix(1449952498, 0)),
	State:            &ec2types.InstanceState{Name: "running"},
	PrivateIpAddress: aws.String("192.168.0.1"),
	PublicIpAddress:  aws.String("35.40.0.1"),
	InstanceType:     "c5.xlarge",
	Placement:        &ec2types.Placement{AvailabilityZone: aws.String("us-east-1f")},
}

var testInstanceProto = ec2v1.Instance{
	InstanceId:       "i-123456789abcdef0",
	Region:           "us-east-1",
	State:            ec2v1.Instance_RUNNING,
	InstanceType:     "c5.xlarge",
	PublicIpAddress:  "35.40.0.1",
	PrivateIpAddress: "192.168.0.1",
	AvailabilityZone: "us-east-1f",
	Tags:             map[string]string{"Name": "locations-staging-iad", "Canary": "false"},
}

var testAutoscalingGroupInstance = astypes.Instance{
	InstanceId:              aws.String("asginstance"),
	AvailabilityZone:        aws.String("us-east-1"),
	LaunchConfigurationName: aws.String("launch-config-name"),
	HealthStatus:            aws.String("HEALTHY"),
	LifecycleState:          astypes.LifecycleStatePending,
}

var testAutoscalingGroupInstnaceProto = ec2v1.AutoscalingGroup_Instance{
	Id:                      "asginstance",
	Zone:                    "us-east-1",
	LaunchConfigurationName: "launch-config-name",
	Healthy:                 true,
	LifecycleState:          ec2v1.AutoscalingGroup_Instance_PENDING,
}

var testAutoscalingGroup = astypes.AutoScalingGroup{
	AutoScalingGroupName: aws.String("asgname"),
	AvailabilityZones:    []string{"us-east-1a", "us-east-1b"},
	MinSize:              aws.Int32(1),
	MaxSize:              aws.Int32(10),
	DesiredCapacity:      aws.Int32(5),
	TerminationPolicies:  []string{"oldest-instance"},
	Tags: []astypes.TagDescription{
		{
			Key:   aws.String("key"),
			Value: aws.String("value"),
		},
	},
	Instances: []astypes.Instance{
		testAutoscalingGroupInstance,
	},
}

var testAutoscalingGroupProto = ec2v1.AutoscalingGroup{
	Name:   "asgname",
	Zones:  []string{"us-east-1a", "us-east-1b"},
	Region: "us-east-1",
	Size: &ec2v1.AutoscalingGroupSize{
		Min:     uint32(1),
		Max:     uint32(10),
		Desired: uint32(5),
	},
	TerminationPolicies: []ec2v1.AutoscalingGroup_TerminationPolicy{
		ec2v1.AutoscalingGroup_OLDEST_INSTANCE,
	},
	Instances: []*ec2v1.AutoscalingGroup_Instance{
		&testAutoscalingGroupInstnaceProto,
	},
}

func TestNewProtoForInstance(t *testing.T) {
	pb := newProtoForInstance(testInstance)
	assert.Equal(t, &testInstanceProto, pb)
}

func TestNewProtoForAutoscalingGroupInstance(t *testing.T) {
	pb := newProtoForAutoscalingGroupInstance(testAutoscalingGroupInstance)
	assert.Equal(t, &testAutoscalingGroupInstnaceProto, pb)
}

func TestNewProtoForAutoscalingGroup(t *testing.T) {
	pb := newProtoForAutoscalingGroup(testAutoscalingGroup)
	assert.Equal(t, &testAutoscalingGroupProto, pb)
}

func TestProtoForInstanceState(t *testing.T) {
	assert.Equal(t, ec2v1.Instance_UNKNOWN, protoForInstanceState("foo"))
	assert.Equal(t, ec2v1.Instance_UNKNOWN, protoForInstanceState(""))
	assert.Equal(t, ec2v1.Instance_RUNNING, protoForInstanceState("running"))
	assert.Equal(t, ec2v1.Instance_SHUTTING_DOWN, protoForInstanceState("shutting-down"))
}

func TestProtoForTerminationPolicy(t *testing.T) {
	assert.Equal(t, ec2v1.AutoscalingGroup_UNKNOWN, protoForTerminationPolicy("foo"))
	assert.Equal(t, ec2v1.AutoscalingGroup_UNKNOWN, protoForTerminationPolicy(""))
	assert.Equal(t, ec2v1.AutoscalingGroup_DEFAULT, protoForTerminationPolicy("default"))
	assert.Equal(t, ec2v1.AutoscalingGroup_OLDEST_INSTANCE, protoForTerminationPolicy("oldest-instance"))
	assert.Equal(t, ec2v1.AutoscalingGroup_OLDEST_LAUNCH_CONFIGURATION, protoForTerminationPolicy("oldest-launch-configuration"))
}

func TestProtoForAutoscalingGroupInstanceLifecycleState(t *testing.T) {
	assert.Equal(t, ec2v1.AutoscalingGroup_Instance_UNKNOWN, protoForAutoscalingGroupInstanceLifecycleState("foo"))
	assert.Equal(t, ec2v1.AutoscalingGroup_Instance_UNKNOWN, protoForAutoscalingGroupInstanceLifecycleState(""))
	assert.Equal(t, ec2v1.AutoscalingGroup_Instance_TERMINATING_WAIT, protoForAutoscalingGroupInstanceLifecycleState("terminating-wait"))
	assert.Equal(t, ec2v1.AutoscalingGroup_Instance_ENTERING_STANDBY, protoForAutoscalingGroupInstanceLifecycleState("entering-standby"))
}

func TestDescribeInstances(t *testing.T) {
	m := &mockEC2{}
	c := &client{
		clients: map[string]*regionalClient{"us-east-1": {region: "us-east-1", ec2: m}},
	}

	results, err := c.DescribeInstances(context.Background(), "us-east-1", nil)
	assert.NoError(t, err)
	assert.Len(t, results, 0)

	m.instances = []ec2types.Instance{testInstance}

	results, err = c.DescribeInstances(context.Background(), "us-east-1", []string{"i-12345"})
	assert.NoError(t, err)
	assert.Len(t, results, 1)
	assert.Equal(t, &testInstanceProto, results[0])

	m.instancesErr = errors.New("whoops")
	_, err = c.DescribeInstances(context.Background(), "us-east-1", nil)
	assert.EqualError(t, err, "whoops")
}

func TestTerminateInstances(t *testing.T) {
	m := &mockEC2{}
	c := &client{
		clients: map[string]*regionalClient{"us-east-1": {region: "us-east-1", ec2: m}},
	}

	err := c.TerminateInstances(context.Background(), "us-east-1", nil)
	assert.NoError(t, err)

	m.terminateErr = errors.New("yikes")
	err = c.TerminateInstances(context.Background(), "us-east-1", nil)
	assert.EqualError(t, err, "yikes")
}

func TestRebootInstances(t *testing.T) {
	m := &mockEC2{}
	c := &client{
		clients: map[string]*regionalClient{"us-east-1": {region: "us-east-1", ec2: m}},
	}

	err := c.RebootInstances(context.Background(), "us-east-1", nil)
	assert.NoError(t, err)

	m.rebootErr = errors.New("ohno")
	err = c.RebootInstances(context.Background(), "us-east-1", nil)
	assert.EqualError(t, err, "ohno")
}

func TestZoneToRegion(t *testing.T) {
	tests := []struct {
		input  string
		expect string
	}{
		{
			input:  "us-east-1a",
			expect: "us-east-1",
		},
		{
			input:  "",
			expect: "UNKNOWN",
		},
	}

	for _, test := range tests {
		output := zoneToRegion(test.input)
		assert.Equal(t, test.expect, output)
	}
}

func TestResizeAutoscalingGroupErrorHandling(t *testing.T) {
	autoscalingClient := &mockAutoscaling{
		updateASGErr: fmt.Errorf("error"),
	}
	c := &client{
		clients: map[string]*regionalClient{"us-east-1": {region: "us-east-1", autoscaling: autoscalingClient}},
	}

	err1 := c.ResizeAutoscalingGroup(context.Background(), "us-east-1", "asgname", &ec2v1.AutoscalingGroupSize{
		Min:     1,
		Max:     10,
		Desired: 5,
	})
	assert.Error(t, err1)

	// Test unknown region
	err2 := c.ResizeAutoscalingGroup(context.Background(), "choice-region-1", "clutch", &ec2v1.AutoscalingGroupSize{})
	assert.Error(t, err2)
}

func TestDescribeAutoScalingGroups(t *testing.T) {
	autoscalingClient := &mockAutoscaling{
		describeASGOutput: &autoscaling.DescribeAutoScalingGroupsOutput{
			AutoScalingGroups: []astypes.AutoScalingGroup{
				{
					AutoScalingGroupName: aws.String("asg-one"),
					AvailabilityZones:    []string{"us-east-1a", "us-east-1b"},
					DesiredCapacity:      aws.Int32(5),
					MaxSize:              aws.Int32(10),
					MinSize:              aws.Int32(1),
				},
				{
					AutoScalingGroupName: aws.String("asg-two"),
					AvailabilityZones:    []string{"us-east-1c", "us-east-1d"},
					DesiredCapacity:      aws.Int32(1),
					MaxSize:              aws.Int32(2),
					MinSize:              aws.Int32(1),
				},
			},
		},
	}
	c := &client{
		clients: map[string]*regionalClient{"us-east-1": {region: "us-east-1", autoscaling: autoscalingClient}},
	}

	asgs, err := c.DescribeAutoscalingGroups(context.Background(), "us-east-1", []string{"asg-one", "asg-two"})
	assert.NoError(t, err)
	assert.Len(t, asgs, 2)

	for i, asg := range asgs {
		assert.Equal(t, aws.ToString(autoscalingClient.describeASGOutput.AutoScalingGroups[i].AutoScalingGroupName), asg.Name)
		assert.Equal(t, autoscalingClient.describeASGOutput.AutoScalingGroups[i].AvailabilityZones, asg.Zones)
		assert.Equal(t, uint32(aws.ToInt32(autoscalingClient.describeASGOutput.AutoScalingGroups[i].DesiredCapacity)), asg.Size.Desired)
		assert.Equal(t, uint32(aws.ToInt32(autoscalingClient.describeASGOutput.AutoScalingGroups[i].MaxSize)), asg.Size.Max)
		assert.Equal(t, uint32(aws.ToInt32(autoscalingClient.describeASGOutput.AutoScalingGroups[i].MinSize)), asg.Size.Min)
	}
}

func TestDescribeAutoscalingGroupsErrorHandling(t *testing.T) {
	autoscalingClient := &mockAutoscaling{
		describeASGErr: fmt.Errorf("error"),
	}

	c := &client{
		clients: map[string]*regionalClient{"us-east-1": {region: "us-east-1", autoscaling: autoscalingClient}},
	}

	asg1, err1 := c.DescribeAutoscalingGroups(context.Background(), "us-east-1", []string{"asgname"})
	assert.Nil(t, asg1)
	assert.Error(t, err1)

	// Test unknown region
	asg2, err2 := c.DescribeAutoscalingGroups(context.Background(), "unknown-region", []string{"asgname"})
	assert.Nil(t, asg2)
	assert.Error(t, err2)
}

func TestErrorIntercept(t *testing.T) {
	c := &client{}
	{
		origErr := newResponseError(400, &smithy.GenericAPIError{Code: "whoopsie", Message: "bad"})
		err := c.InterceptError(origErr)
		_, ok := status.FromError(err)
		assert.True(t, ok)
	}
	{
		origErr := errors.New("foo")
		err := c.InterceptError(origErr)
		assert.Equal(t, origErr, err)
	}
}

type mockEC2 struct {
	ec2Client

	instancesErr error
	instances    []ec2types.Instance

	terminateResult []ec2types.InstanceStateChange
	terminateErr    error
	rebootErr       error
}

func (m *mockEC2) DescribeInstances(ctx context.Context, params *ec2.DescribeInstancesInput, optFns ...func(*ec2.Options)) (*ec2.DescribeInstancesOutput, error) {
	if m.instancesErr != nil {
		return nil, m.instancesErr
	}

	if len(params.InstanceIds) != len(m.instances) {
		panic(fmt.Sprintf("DescribeInstances mock count mismatch input %d, expected %d", len(params.InstanceIds), len(m.instances)))
	}

	ret := &ec2.DescribeInstancesOutput{
		Reservations: []ec2types.Reservation{
			{Instances: m.instances},
		},
	}
	return ret, nil
}

func (m *mockEC2) TerminateInstances(ctx context.Context, params *ec2.TerminateInstancesInput, optFns ...func(*ec2.Options)) (*ec2.TerminateInstancesOutput, error) {
	if m.terminateErr != nil {
		return nil, m.terminateErr
	}
	ret := &ec2.TerminateInstancesOutput{
		TerminatingInstances: m.terminateResult,
	}
	return ret, nil
}

func (m *mockEC2) RebootInstances(ctx context.Context, params *ec2.RebootInstancesInput, optFns ...func(*ec2.Options)) (*ec2.RebootInstancesOutput, error) {
	if m.rebootErr != nil {
		return nil, m.rebootErr
	}
	return &ec2.RebootInstancesOutput{}, nil
}

type mockAutoscaling struct {
	autoscalingClient

	describeASGErr    error
	describeASGOutput *autoscaling.DescribeAutoScalingGroupsOutput

	updateASGErr error
}

func (ma mockAutoscaling) DescribeAutoScalingGroups(ctx context.Context, params *autoscaling.DescribeAutoScalingGroupsInput, optFns ...func(*autoscaling.Options)) (*autoscaling.DescribeAutoScalingGroupsOutput, error) {
	if ma.describeASGErr != nil {
		return nil, ma.describeASGErr
	}

	output := &autoscaling.DescribeAutoScalingGroupsOutput{}
	if ma.describeASGOutput != nil {
		output = ma.describeASGOutput
	}

	return output, nil
}

func (ma mockAutoscaling) UpdateAutoScalingGroup(ctx context.Context, params *autoscaling.UpdateAutoScalingGroupInput, optFns ...func(*autoscaling.Options)) (*autoscaling.UpdateAutoScalingGroupOutput, error) {
	if ma.updateASGErr != nil {
		return nil, ma.updateASGErr
	}

	return &autoscaling.UpdateAutoScalingGroupOutput{}, nil
}

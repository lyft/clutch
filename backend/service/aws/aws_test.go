package aws

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	ec2types "github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/any"
	"github.com/stretchr/testify/assert"
	"github.com/uber-go/tally"
	"go.uber.org/zap/zaptest"

	ec2v1 "github.com/lyft/clutch/backend/api/aws/ec2/v1"
	awsv1 "github.com/lyft/clutch/backend/api/config/service/aws/v1"
)

func TestNew(t *testing.T) {
	regions := []string{"us-east-1", "us-west25"}

	cfg, _ := ptypes.MarshalAny(&awsv1.Config{
		Regions: regions,
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

func TestNewProtoForInstance(t *testing.T) {
	pb := newProtoForInstance(testInstance)
	assert.Equal(t, testInstanceProto, pb)
}

func TestProtoForInstanceState(t *testing.T) {
	assert.Equal(t, ec2v1.Instance_UNKNOWN, protoForInstanceState("foo"))
	assert.Equal(t, ec2v1.Instance_UNKNOWN, protoForInstanceState(""))
	assert.Equal(t, ec2v1.Instance_RUNNING, protoForInstanceState("running"))
	assert.Equal(t, ec2v1.Instance_SHUTTING_DOWN, protoForInstanceState("shutting-down"))
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
	assert.Equal(t, testInstanceProto, results[0])

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

type mockEC2 struct {
	*ec2.Client // satisfies interface

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

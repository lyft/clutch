package meta

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/proto"

	apiv1 "github.com/lyft/clutch/backend/api/api/v1"
	auditv1 "github.com/lyft/clutch/backend/api/audit/v1"
	ec2v1 "github.com/lyft/clutch/backend/api/aws/ec2/v1"
	healthcheckv1 "github.com/lyft/clutch/backend/api/healthcheck/v1"
	k8sapiv1 "github.com/lyft/clutch/backend/api/k8s/v1"

	"github.com/lyft/clutch/backend/module"
	"github.com/lyft/clutch/backend/module/healthcheck"
)

type mockRegistrar struct {
	s *grpc.Server
}

func (m *mockRegistrar) GRPCServer() *grpc.Server { return m.s }

func (m *mockRegistrar) RegisterJSONGateway(handlerFunc module.GatewayRegisterAPIHandlerFunc) error {
	return nil
}

func TestGetAction(t *testing.T) {
	hc, err := healthcheck.New(nil, nil, nil)
	assert.NoError(t, err)

	r := &mockRegistrar{s: grpc.NewServer()}
	err = hc.Register(r)
	assert.NoError(t, err)
	err = GenerateGRPCMetadata(r.s)
	assert.NoError(t, err)

	_ = &healthcheckv1.HealthcheckRequest{} // Ensure it's imported.
	action := GetAction("/clutch.healthcheck.v1.HealthcheckAPI/Healthcheck")
	assert.Equal(t, apiv1.ActionType_READ, action)

	action = GetAction("/nonexistent/doesnotexist")
	assert.Equal(t, apiv1.ActionType_UNSPECIFIED, action)
}

func TestResourceNames(t *testing.T) {
	t.Parallel()

	instanceID := "i-1234567890abcdef0"
	region := "iad"
	expectedName := "iad/i-1234567890abcdef0"
	typeUrl := "clutch.aws.ec2.v1.Instance"

	tests := []struct {
		id      string
		message proto.Message
		names   []*auditv1.Resource
	}{
		{
			id:      "nil names",
			message: &ec2v1.TerminateInstanceResponse{},
		},
		{
			id: "single pattern",
			message: &ec2v1.TerminateInstanceRequest{
				InstanceId: instanceID,
				Region:     region,
			},
			names: []*auditv1.Resource{{
				TypeUrl: typeUrl,
				Id:      expectedName,
			}},
		},
		{
			id: "resolve through field",
			message: &ec2v1.GetInstanceResponse{
				Instance: &ec2v1.Instance{
					InstanceId: instanceID,
					Region:     region,
				},
			},
			names: []*auditv1.Resource{{
				TypeUrl: typeUrl,
				Id:      expectedName,
			}},
		},
		{
			id: "resolve through repeated field",
			message: &auditv1.RequestEvent{
				Resources: []*auditv1.Resource{
					{TypeUrl: typeUrl},
				},
			},
			names: []*auditv1.Resource{{
				TypeUrl: "clutch.audit.v1.Resource",
				Id:      typeUrl,
			}},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.id, func(t *testing.T) {
			t.Parallel()
			resolvedNames := ResourceNames(tt.message)
			assert.Len(t, resolvedNames, len(tt.names))
			for i := range resolvedNames {
				assert.Equal(t, tt.names[i], resolvedNames[i])
			}
		})
	}
}

// proto.message with nil value
var m = (*ec2v1.Instance)(nil)

// test cases used in both TestIsValidInterface and TestAPIBody
var tests = []struct {
	input     interface{}
	isValid   bool // checked in TestIsValidInterface
	expectNil bool // checked in TestAPIBody
}{
	// case: type is proto.message and value is not nil
	{
		input:     &ec2v1.Instance{InstanceId: "i-123456789abcdef0"},
		isValid:   true,
		expectNil: false,
	},
	// case: type is proto.message and value is not nil
	{
		input:     &k8sapiv1.ResizeHPAResponse{},
		isValid:   true,
		expectNil: false,
	},
	// case: type is proto.message and value is nil
	{
		input:     m,
		isValid:   false,
		expectNil: true,
	},
	// case: type is struct
	{
		input:     ec2v1.Instance{InstanceId: "i-123456789abcdef0"},
		isValid:   false,
		expectNil: true,
	},
	// case: type/value is nil
	{
		input:     nil,
		isValid:   false,
		expectNil: true,
	},
	// case: type is string
	{
		input:     "foo",
		isValid:   false,
		expectNil: true,
	},
}

func TestIsValidInterface(t *testing.T) {
	for _, test := range tests {
		output := isValidInterface(test.input)
		assert.Equal(t, test.isValid, output)
	}
}

func TestAPIBody(t *testing.T) {
	for _, test := range tests {
		result, err := APIBody(test.input)
		if test.expectNil {
			assert.Nil(t, result)
			assert.NoError(t, err)
		} else {
			assert.NotNil(t, result)
			assert.NoError(t, err)
		}
	}
}

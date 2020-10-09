package meta

import (
	"testing"

	"github.com/golang/protobuf/descriptor"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"

	apiv1 "github.com/lyft/clutch/backend/api/api/v1"
	auditv1 "github.com/lyft/clutch/backend/api/audit/v1"
	ec2v1 "github.com/lyft/clutch/backend/api/aws/ec2/v1"
	healthcheckv1 "github.com/lyft/clutch/backend/api/healthcheck/v1"
	k8sv1 "github.com/lyft/clutch/backend/api/k8s/v1"
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
}

func TestResourceNames(t *testing.T) {
	t.Parallel()

	instanceID := "i-1234567890abcdef0"
	region := "iad"
	expectedName := "iad/i-1234567890abcdef0"
	typeUrl := "clutch.aws.ec2.v1.Instance"

	tests := []struct {
		id      string
		message descriptor.Message
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

func TestAPIMetadata(t *testing.T) {
	// nil value
	result := APIMetadata(nil)
	assert.Nil(t, result)

	// non Ptr value
	result = APIMetadata([]string{"foo"})
	assert.Nil(t, result)

	// Ptr value
	result = APIMetadata(&k8sv1.DeletePodRequest{})
	assert.NotNil(t, result)
}

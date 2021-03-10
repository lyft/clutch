package meta

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"

	apiv1 "github.com/lyft/clutch/backend/api/api/v1"
	auditv1 "github.com/lyft/clutch/backend/api/audit/v1"
	authnv1 "github.com/lyft/clutch/backend/api/authn/v1"
	ec2v1 "github.com/lyft/clutch/backend/api/aws/ec2/v1"
	healthcheckv1 "github.com/lyft/clutch/backend/api/healthcheck/v1"
	k8sapiv1 "github.com/lyft/clutch/backend/api/k8s/v1"
	modulemock "github.com/lyft/clutch/backend/mock/module"
	"github.com/lyft/clutch/backend/module/assets"
	"github.com/lyft/clutch/backend/module/healthcheck"
)

func TestGetAction(t *testing.T) {
	hc, err := healthcheck.New(nil, nil, nil)
	assert.NoError(t, err)

	r := &modulemock.MockRegistrar{Server: grpc.NewServer()}
	err = hc.Register(r)
	assert.NoError(t, err)

	// Register a non-Clutch endpoint with no annotations.
	grpc_health_v1.RegisterHealthServer(r.GRPCServer(), &grpc_health_v1.UnimplementedHealthServer{})

	err = GenerateGRPCMetadata(r.GRPCServer())
	assert.NoError(t, err)

	action := GetAction("/clutch.healthcheck.v1.HealthcheckAPI/Healthcheck")
	assert.Equal(t, apiv1.ActionType_READ, action)

	action = GetAction("/grpc.health.v1.Health/Check")
	assert.Equal(t, apiv1.ActionType_UNSPECIFIED, action)

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

func TestIsRedacted(t *testing.T) {
	assert.True(t, IsRedacted(&authnv1.CallbackResponse{}))
	assert.False(t, IsRedacted(&healthcheckv1.HealthcheckRequest{}))
}

func TestAPIBodyRedaction(t *testing.T) {
	b, err := APIBody(&authnv1.CallbackResponse{AccessToken: "secret"})
	assert.NoError(t, err)

	m, err := b.UnmarshalNew()
	assert.NoError(t, err)
	assert.IsType(t, (*apiv1.Redacted)(nil), m)

	r := m.(*apiv1.Redacted)
	assert.Equal(t, r.RedactedTypeUrl, "type.googleapis.com/clutch.authn.v1.CallbackResponse")
}

func TestAPIBody(t *testing.T) {
	// proto.message with nil value
	m := (*ec2v1.Instance)(nil)

	var tests = []struct {
		input     interface{}
		expectNil bool
	}{
		// case: type is proto.message
		{input: &ec2v1.Instance{InstanceId: "i-123456789abcdef0"}, expectNil: false},
		// case: type is proto.message
		{input: &k8sapiv1.ResizeHPAResponse{}, expectNil: false},
		// case: type is proto.message
		// anypb.New gracefully handles (*myProtoStruct)(nil)
		{input: m, expectNil: false},
		// case: type is struct
		{input: ec2v1.Instance{InstanceId: "i-123456789abcdef0"}, expectNil: true},
		// case: untyped nil
		{input: nil, expectNil: true},
		// case: type is string
		{input: "foo", expectNil: true},
	}

	for _, test := range tests {
		result, err := APIBody(test.input)
		if test.expectNil {
			assert.Nil(t, result)
		} else {
			assert.NotNil(t, result)
		}

		assert.NoError(t, err)
	}
}

func TestAuditDisabled(t *testing.T) {
	r := &modulemock.MockRegistrar{Server: grpc.NewServer()}

	hc, err := healthcheck.New(nil, nil, nil)
	assert.NoError(t, err)
	assert.NoError(t, hc.Register(r))

	a, err := assets.New(nil, nil, nil)
	assert.NoError(t, err)
	assert.NoError(t, a.Register(r))

	assert.NoError(t, GenerateGRPCMetadata(r.GRPCServer()))

	result := IsAuditDisabled("/clutch.healthcheck.v1.HealthcheckAPI/Healthcheck")
	assert.True(t, result)

	result = IsAuditDisabled("/clutch.assets.v1.AssetsAPI/Fetch")
	assert.False(t, result)

	result = IsAuditDisabled("/nonexistent/doesnotexist")
	assert.False(t, result)
}

func TestExtractProtoPatternsValues(t *testing.T) {
	t.Parallel()
	tests := []struct {
		id          string
		pb          proto.Message
		expect      string
		shouldError bool
	}{
		{
			id: "deployment",
			pb: &k8sapiv1.Deployment{
				Cluster:   "foo",
				Namespace: "bar",
				Name:      "cat",
			},
			expect:      "foo/bar/cat",
			shouldError: false,
		},
		{
			id: "ec2 instance",
			pb: &ec2v1.Instance{
				Region:     "us-east-1",
				InstanceId: "i-000000000",
			},
			expect:      "us-east-1/i-000000000",
			shouldError: false,
		},
		{
			id:          "no pattern found",
			pb:          &anypb.Any{},
			expect:      "",
			shouldError: true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.id, func(t *testing.T) {
			t.Parallel()

			result, err := HydratedPatternForProto(tt.pb)
			if tt.shouldError {
				assert.Empty(t, result)
				assert.Error(t, err)
			} else {
				assert.Equal(t, tt.expect, result)
				assert.NoError(t, err)
			}
		})
	}
}

func TestPatternValueMapping(t *testing.T) {
	t.Parallel()
	tests := []struct {
		id     string
		pb     proto.Message
		search string
		result map[string]string
		ok     bool
	}{
		{
			id:     "aws asg",
			pb:     &ec2v1.AutoscalingGroup{},
			search: "us-east-1/my-asg-name",
			result: map[string]string{
				"region": "us-east-1",
				"name":   "my-asg-name",
			},
			ok: true,
		},
		{
			id:     "aws instance",
			pb:     &ec2v1.Instance{},
			search: "us-east-1/i-0000000",
			result: map[string]string{
				"region":      "us-east-1",
				"instance_id": "i-0000000",
			},
			ok: true,
		},
		{
			id:     "test for partial match",
			pb:     &ec2v1.Instance{},
			search: "us-east-1/i-0000000/meow",
			result: map[string]string{
				"region":      "us-east-1/i-0000000",
				"instance_id": "meow",
			},
			ok: true,
		},
		{
			id:     "k8s deployment",
			pb:     &k8sapiv1.Deployment{},
			search: "mycluster/mynamespace/deploymentname",
			result: map[string]string{
				"cluster":   "mycluster",
				"namespace": "mynamespace",
				"name":      "deploymentname",
			},
			ok: true,
		},
		{
			id:     "k8s deployment failed pattern match",
			pb:     &k8sapiv1.Deployment{},
			search: "nothecorrectpattern",
			result: map[string]string{},
			ok:     false,
		},
		{
			id:     "failed match results should have nil values",
			pb:     &k8sapiv1.Deployment{},
			search: "cluster/namespace",
			result: map[string]string{},
			ok:     false,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.id, func(t *testing.T) {
			t.Parallel()

			result, ok, err := ExtractPatternValuesFromString(tt.pb, tt.search)
			assert.NoError(t, err)
			assert.Equal(t, tt.ok, ok)
			assert.Equal(t, tt.result, result)
		})
	}
}

func TestExtractProtoPatternFieldNames(t *testing.T) {
	t.Parallel()
	tests := []struct {
		id      string
		pattern *apiv1.Pattern
		expect  []string
	}{
		{
			id:      "3 fields",
			pattern: &apiv1.Pattern{Pattern: "{name}/{of}/{fields}"},
			expect:  []string{"name", "of", "fields"},
		},
		{
			id:      "2 fields",
			pattern: &apiv1.Pattern{Pattern: "{name}/{of}"},
			expect:  []string{"name", "of"},
		},
		{
			id:      "1 fields",
			pattern: &apiv1.Pattern{Pattern: "{name}"},
			expect:  []string{"name"},
		},
		{
			id:      "different delimiters",
			pattern: &apiv1.Pattern{Pattern: "{cat}/{meow}-{nom}_{food}--{tasty}"},
			expect:  []string{"cat", "meow", "nom", "food", "tasty"},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.id, func(t *testing.T) {
			t.Parallel()

			actual := extractProtoPatternFieldNames(tt.pattern)
			assert.Equal(t, tt.expect, actual)
		})
	}
}

package audit

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/anypb"

	auditv1 "github.com/lyft/clutch/backend/api/audit/v1"
	healthcheckv1 "github.com/lyft/clutch/backend/api/healthcheck/v1"
	k8sapiv1 "github.com/lyft/clutch/backend/api/k8s/v1"
	"github.com/lyft/clutch/backend/gateway/meta"
	modulemock "github.com/lyft/clutch/backend/mock/module"
	"github.com/lyft/clutch/backend/module/healthcheck"
	"github.com/lyft/clutch/backend/service/audit"
)

func TestEventFromResponse(t *testing.T) {
	m := &mid{}
	err := errors.New("error")
	a := (*anypb.Any)(nil)

	// case: err, passed to eventFromResponse, does not equal nil
	event, err := m.eventFromResponse(nil, err)
	assert.NoError(t, err)
	assert.NotEmpty(t, event)
	assert.Equal(t, "error", event.Status.Message)
	assert.Equal(t, 0, len(event.Resources))
	assert.Equal(t, a, event.ResponseMetadata.Body)

	resp := &k8sapiv1.Pod{Cluster: "kind-clutch", Namespace: "envoy-staging", Name: "envoy-main-579848cc64-cxnqm"}
	// case: err, passed to eventFromResponse, equals nil
	event, err = m.eventFromResponse(resp, nil)
	assert.NoError(t, err)
	assert.NotEmpty(t, event)
	assert.Equal(t, int32(0), event.Status.Code)
	assert.Equal(t, "", event.Status.Message)
	assert.Equal(t, 1, len(event.Resources))
	assert.Equal(t, "type.googleapis.com/clutch.k8s.v1.Pod", event.ResponseMetadata.Body.TypeUrl)
}

type mockAuditor struct {
	audit.Auditor

	writeCount  uint64
	updateCount uint64
}

func (m *mockAuditor) WriteRequestEvent(ctx context.Context, req *auditv1.RequestEvent) (int64, error) {
	m.writeCount += 1
	return 0, nil
}

func (m *mockAuditor) UpdateRequestEvent(ctx context.Context, id int64, update *auditv1.RequestEvent) error {
	m.updateCount += 1
	return nil
}

func fakeHandler(ctx context.Context, req interface{}) (interface{}, error) {
	return &healthcheckv1.HealthcheckResponse{}, nil
}

func TestInterceptor(t *testing.T) {
	a := &mockAuditor{}
	m := &mid{
		audit: a,
	}

	interceptor := m.UnaryInterceptor()
	resp, err := interceptor(
		context.Background(),
		&healthcheckv1.HealthcheckRequest{},
		&grpc.UnaryServerInfo{FullMethod: "/foo/bar"},
		fakeHandler)

	assert.NotNil(t, resp)
	assert.NoError(t, err)

	assert.EqualValues(t, 1, a.writeCount)
	assert.EqualValues(t, 1, a.updateCount)
}

func TestInterceptorShortCircuitDisabled(t *testing.T) {
	a := &mockAuditor{}
	m := &mid{
		audit: a,
	}

	server := grpc.NewServer()
	r := &modulemock.MockRegistrar{Server: server}
	hc, _ := healthcheck.New(nil, nil, nil)
	assert.NoError(t, hc.Register(r))
	assert.NoError(t, meta.GenerateGRPCMetadata(server))

	interceptor := m.UnaryInterceptor()
	resp, err := interceptor(
		context.Background(),
		&healthcheckv1.HealthcheckRequest{},
		&grpc.UnaryServerInfo{FullMethod: "/clutch.healthcheck.v1.HealthcheckAPI/Healthcheck"},
		fakeHandler)

	assert.NotNil(t, resp)
	assert.NoError(t, err)

	assert.EqualValues(t, 0, a.writeCount)
	assert.EqualValues(t, 0, a.updateCount)
}

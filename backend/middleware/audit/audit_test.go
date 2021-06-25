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
	testpb "github.com/lyft/clutch/backend/internal/test/pb"
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

func TestInterceptorWithRedaction(t *testing.T) {
	a := &mockAuditor{}
	m := &mid{
		audit: a,
	}

	logTestProto := &testpb.LogOptionsTester{
		StrLogFalse:      "test",
		StrLogTrue:       "test",
		StrWithoutOption: "test",
		NestedNoLog: &testpb.NestedLogOptionTester{
			StrWithoutOption: "test",
		},
		Nested: &testpb.NestedLogOptionTester{
			StrLogFalse:      "test",
			StrWithoutOption: "test",
		},
		MessageMap: map[string]*testpb.NestedLogOptionTester{
			"test": {
				StrLogFalse:      "test",
				StrWithoutOption: "test",
			},
			"nil": nil,
		},
		RepeatedMessage: []*testpb.NestedLogOptionTester{
			{
				StrLogFalse:      "test",
				StrWithoutOption: "test",
			},
			nil,
		},
	}

	interceptor := m.UnaryInterceptor()
	resp, err := interceptor(
		context.Background(),
		logTestProto,
		&grpc.UnaryServerInfo{FullMethod: "/foo/bar"},
		func(ctx context.Context, req interface{}) (interface{}, error) {
			// We must assert that the fields that were redacted are still present on the request object
			assert.Equal(t, "test", req.(*testpb.LogOptionsTester).StrLogFalse)
			assert.Equal(t, "test", req.(*testpb.LogOptionsTester).Nested.StrLogFalse)
			assert.Equal(t, "test", req.(*testpb.LogOptionsTester).MessageMap["test"].StrLogFalse)
			assert.Equal(t, "test", req.(*testpb.LogOptionsTester).RepeatedMessage[0].StrLogFalse)

			// Passthrough the orginal request to the response to assert again
			return req, nil
		})

	assert.NotNil(t, resp)
	assert.NoError(t, err)

	// We must assert that the fields that were redacted are still present on the response object
	assert.Equal(t, "test", resp.(*testpb.LogOptionsTester).StrLogFalse)
	assert.Equal(t, "test", resp.(*testpb.LogOptionsTester).Nested.StrLogFalse)
	assert.Equal(t, "test", resp.(*testpb.LogOptionsTester).MessageMap["test"].StrLogFalse)
	assert.Equal(t, "test", resp.(*testpb.LogOptionsTester).RepeatedMessage[0].StrLogFalse)

	assert.EqualValues(t, 1, a.writeCount)
	assert.EqualValues(t, 1, a.updateCount)
}

func BenchmarkUnaryInterceptor(b *testing.B) {
	a := &mockAuditor{}
	m := &mid{
		audit: a,
	}
	server := grpc.NewServer()
	r := &modulemock.MockRegistrar{Server: server}
	hc, _ := healthcheck.New(nil, nil, nil)
	assert.NoError(b, hc.Register(r))
	assert.NoError(b, meta.GenerateGRPCMetadata(server))

	interceptor := m.UnaryInterceptor()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = interceptor(
			context.Background(),
			&healthcheckv1.HealthcheckRequest{},
			&grpc.UnaryServerInfo{FullMethod: "/foo/bar"},
			fakeHandler)
	}
}

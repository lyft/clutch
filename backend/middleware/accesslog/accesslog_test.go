package log

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap/zaptest"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"

	middlewarev1 "github.com/lyft/clutch/backend/api/config/middleware/v1"
	healthcheckv1 "github.com/lyft/clutch/backend/api/healthcheck/v1"
)

func TestNew(t *testing.T) {
	tests := []struct {
		config *middlewarev1.AccessLog
	}{
		{config: nil},
		{config: &middlewarev1.AccessLog{
			StatusCodeFilter: &middlewarev1.AccessLog_StatusCodeFilter{
				Operator: middlewarev1.AccessLog_StatusCodeFilter_GE,
				Value:    uint32(codes.OK),
			},
		},
		},
	}

	for idx, tt := range tests {
		tt := tt
		t.Run(fmt.Sprintf("%d", idx), func(t *testing.T) {
			t.Parallel()

			m, err := New(tt.config, nil, nil)
			assert.NoError(t, err)
			assert.NotNil(t, m)
		})
	}
}

func fakeHandler(ctx context.Context, req interface{}) (interface{}, error) {
	return &healthcheckv1.HealthcheckResponse{}, nil
}

func TestInterceptor(t *testing.T) {
	m := &mid{
		logger:     zaptest.NewLogger(t),
		operator:   middlewarev1.AccessLog_StatusCodeFilter_GE,
		statusCode: codes.OK,
	}

	interceptor := m.UnaryInterceptor()
	resp, err := interceptor(
		context.Background(),
		&healthcheckv1.HealthcheckRequest{},
		&grpc.UnaryServerInfo{FullMethod: "/foo/bar"},
		fakeHandler)

	assert.NotNil(t, resp)
	assert.NoError(t, err)
}

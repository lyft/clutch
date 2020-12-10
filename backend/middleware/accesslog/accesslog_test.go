package log

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zaptest"
	"go.uber.org/zap/zaptest/observer"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	accesslogv1 "github.com/lyft/clutch/backend/api/config/middleware/accesslog/v1"
	healthcheckv1 "github.com/lyft/clutch/backend/api/healthcheck/v1"
)

func TestNew(t *testing.T) {
	tests := []struct {
		config *accesslogv1.Config
	}{
		{config: nil},
		{config: &accesslogv1.Config{
			StatusCodeFilters: []*accesslogv1.Config_StatusCodeFilter{
				&accesslogv1.Config_StatusCodeFilter{
					FilterType: &accesslogv1.Config_StatusCodeFilter_Equals{
						Equals: 1,
					},
				},
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
		logger:      zaptest.NewLogger(t),
		statusCodes: []codes.Code{codes.OK},
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

func TestStatusCodeFilter(t *testing.T) {
	filterTests := []struct {
		// status code filter value
		equals []int
		// length of recorded log
		expectedLogLength int
	}{
		{
			equals:            []int{0},
			expectedLogLength: 1,
		},
		{
			equals:            []int{5},
			expectedLogLength: 0,
		},
		{
			equals:            []int{0, 12, 13, 14},
			expectedLogLength: 1,
		},
	}
	for _, test := range filterTests {
		core, recorded := observer.New(zapcore.DebugLevel)
		log := zap.New(core)

		var statusCodeFilters []*accesslogv1.Config_StatusCodeFilter
		for _, value := range test.equals {
			statusCodeFilters = append(statusCodeFilters,
				&accesslogv1.Config_StatusCodeFilter{
					FilterType: &accesslogv1.Config_StatusCodeFilter_Equals{
						Equals: uint32(value),
					},
				})
		}
		m, err := New(&accesslogv1.Config{
			StatusCodeFilters: statusCodeFilters,
		}, log, nil)
		assert.NoError(t, err)

		fakeHandler := func(ctx context.Context, req interface{}) (interface{}, error) {
			return &healthcheckv1.HealthcheckResponse{}, nil
		}

		midFn := m.UnaryInterceptor()
		resp, err := midFn(context.Background(), nil, &grpc.UnaryServerInfo{FullMethod: "/foo/bar"}, fakeHandler)
		assert.NotNil(t, resp)

		s := status.Convert(err)
		assert.Equal(t, codes.OK, s.Code())

		// ensure the recorded log message matches status code filter
		assert.Equal(t, test.expectedLogLength, recorded.Len())
	}
}

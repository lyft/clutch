package timeouts

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zaptest/observer"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/durationpb"

	gatewayv1 "github.com/lyft/clutch/backend/api/config/gateway/v1"
	healthcheckv1 "github.com/lyft/clutch/backend/api/healthcheck/v1"
)

func TestNew(t *testing.T) {
	tests := []struct {
		config *gatewayv1.Timeouts
	}{
		{config: nil},
		{config: &gatewayv1.Timeouts{Default: durationpb.New(time.Second)}},
		{config: &gatewayv1.Timeouts{
			Default: durationpb.New(time.Second),
			Overrides: []*gatewayv1.Timeouts_Entry{
				{Service: "foo", Method: "bar", Timeout: durationpb.New(time.Minute)},
			},
		}},
	}

	for idx, tt := range tests {
		tt := tt
		t.Run(fmt.Sprintf("%d", idx), func(t *testing.T) {
			t.Parallel()
			assert.NoError(t, tt.config.Validate())

			m, err := New(tt.config, nil, nil)
			assert.NoError(t, err)
			assert.NotNil(t, m)
		})
	}
}

func TestJoin(t *testing.T) {
	assert.Equal(t, "/foo/bar", join("foo", "bar"))
}

func TestOverride(t *testing.T) {
	cfg := &gatewayv1.Timeouts{
		Default: durationpb.New(time.Second),
		Overrides: []*gatewayv1.Timeouts_Entry{
			{Service: "foo", Method: "bar", Timeout: durationpb.New(time.Minute)},
		},
	}

	m, err := New(cfg, nil, nil)
	assert.NoError(t, err)
	tm := m.(*mid)

	expected := []struct {
		service, method string
		expected        time.Duration
	}{
		{"foo", "bar", time.Minute},
		{"foo", "blue", time.Second},
		{"bar", "bing", time.Second},
	}
	for _, tt := range expected {
		assert.Equal(t, tt.expected, tm.getDuration(tt.service, tt.method))
	}
}

// Return within timeout.
func TestTimeoutOkay(t *testing.T) {
	m, err := New(&gatewayv1.Timeouts{Default: durationpb.New(time.Second)}, nil, nil)
	assert.NoError(t, err)

	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return &healthcheckv1.HealthcheckResponse{}, nil
	}

	midFn := m.UnaryInterceptor()
	resp, err := midFn(context.Background(), nil, &grpc.UnaryServerInfo{FullMethod: "/zip/zoom"}, handler)
	assert.NoError(t, err)
	assert.NotNil(t, resp)
}

// Return after timeout but before boost period ends.
func TestTimeoutWithinBoost(t *testing.T) {
	m, err := New(&gatewayv1.Timeouts{Default: durationpb.New(time.Nanosecond)}, nil, nil)
	assert.NoError(t, err)

	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		time.Sleep(time.Millisecond) // Sleep for longer than timeout but less than boost.
		select {
		case <-ctx.Done():
			return nil, errors.New("i was done")
		default:
		}
		return &healthcheckv1.HealthcheckResponse{}, nil
	}

	ctx, cancelFn := context.WithCancel(context.Background())
	cancelFn() // Immediately cancel so handler returns its error.

	midFn := m.UnaryInterceptor()
	resp, err := midFn(ctx, nil, &grpc.UnaryServerInfo{FullMethod: "/zip/zoom"}, handler)
	assert.Nil(t, resp)
	assert.EqualError(t, err, "i was done")
}

// Return after timeout and boost.
func TestTimeoutBlockForMoreThanBoost(t *testing.T) {
	core, recorded := observer.New(zapcore.DebugLevel)
	log := zap.New(core)

	m, err := New(&gatewayv1.Timeouts{Default: durationpb.New(time.Nanosecond)}, log, nil)
	assert.NoError(t, err)

	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		time.Sleep(boost * 4)
		return &healthcheckv1.HealthcheckResponse{}, nil
	}

	midFn := m.UnaryInterceptor()
	resp, err := midFn(context.Background(), nil, &grpc.UnaryServerInfo{FullMethod: "/zip/zoom"}, handler)
	assert.Nil(t, resp)
	s := status.Convert(err)
	assert.Equal(t, codes.DeadlineExceeded, s.Code())

	// Wait for the log or timeout and fail.
	timeout := time.After(boost * 5)
	for recorded.Len() == 0 {
		select {
		case <-timeout:
			t.Fatal("timed out waiting for log")
		default:
			time.Sleep(time.Millisecond)
		}
	}
	assert.Equal(t, "handler completed after timeout", recorded.All()[0].Message)
}

func TestInfiniteTimeout(t *testing.T) {
	cfg := &gatewayv1.Timeouts{Default: durationpb.New(0)}

	m, err := New(cfg, nil, nil)
	assert.NoError(t, err)

	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		time.Sleep(boost + time.Millisecond)
		return &healthcheckv1.HealthcheckResponse{}, nil
	}

	midFn := m.UnaryInterceptor()
	resp, err := midFn(context.Background(), nil, &grpc.UnaryServerInfo{FullMethod: "/zip/zoom"}, handler)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
}

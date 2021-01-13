package timeouts

// <!-- START clutchdoc -->
// description: Adds a deadline to all request contexts.
// <!-- END clutchdoc -->

import (
	"context"
	"fmt"
	"time"

	"github.com/uber-go/tally"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/durationpb"

	gatewayv1 "github.com/lyft/clutch/backend/api/config/gateway/v1"
	"github.com/lyft/clutch/backend/middleware"
)

const (
	DefaultTimeout = time.Second * 15

	// Boost is added to the timeout to give a handler that respects the deadline the opportunity to return an error.
	boost = time.Millisecond * 50
)

func New(config *gatewayv1.Timeouts, logger *zap.Logger, scope tally.Scope) (middleware.Middleware, error) {
	if config == nil {
		config = &gatewayv1.Timeouts{Default: durationpb.New(DefaultTimeout)}
	}

	defaultTimeout := config.Default.AsDuration()

	overrides := make(map[string]time.Duration, len(config.Overrides))
	for _, entry := range config.Overrides {
		overrides[join(entry.Service, entry.Method)] = entry.Timeout.AsDuration()
	}

	return &mid{
		logger:         logger,
		defaultTimeout: defaultTimeout,
		overrides:      overrides,
	}, nil
}

type mid struct {
	logger         *zap.Logger
	defaultTimeout time.Duration
	overrides      map[string]time.Duration
}

func (m *mid) getDuration(service, method string) time.Duration {
	if override, ok := m.overrides[join(service, method)]; ok {
		return override
	}
	return m.defaultTimeout
}

type unaryHandlerReturn struct {
	resp interface{}
	err  error
}

func (m *mid) UnaryInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		service, method, ok := middleware.SplitFullMethod(info.FullMethod)
		if !ok {
			m.logger.Warn("could not parse gRPC method", zap.String("fullMethod", info.FullMethod))
		}

		// Create a return channel for the goroutine.
		resultChan := make(chan unaryHandlerReturn)
		defer close(resultChan)

		// Compute timeout and set-up a context with timeout.
		timeout := m.getDuration(service, method)
		if timeout != 0 {
			var cancel context.CancelFunc
			ctx, cancel = context.WithTimeout(ctx, timeout)
			defer cancel()
		}

		// Create a channel to track when the timeout error has already been returned and the return channel is closed.
		done := make(chan struct{})
		if timeout != 0 {
			// If timeout is not infinite, return after timeout plus boost
			timer := time.AfterFunc(timeout+boost, func() { close(done) })
			defer timer.Stop()
		}

		// Spawn the handler in a goroutine so we can return early on timeout if it doesn't complete.
		go func() {
			resp, err := handler(ctx, req)
			select {
			case <-done:
				m.logger.Error(
					"handler completed after timeout",
					zap.String("service", service),
					zap.String("method", method),
					zap.Error(err))
			default:
				resultChan <- unaryHandlerReturn{resp: resp, err: err}
			}
		}()

		select {
		case ret := <-resultChan:
			return ret.resp, ret.err
		case <-done:
			return nil, status.New(codes.DeadlineExceeded, "timeout exceeded").Err()
		}
	}
}

func join(service, method string) string {
	const pattern = "/%s/%s"
	return fmt.Sprintf(pattern, service, method)
}

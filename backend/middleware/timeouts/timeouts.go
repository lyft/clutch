package timeouts

// <!-- START clutchdoc -->
// description: Adds a deadline to all request contexts.
// <!-- END clutchdoc -->

import (
	"context"
	"fmt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"time"

	"github.com/golang/protobuf/ptypes"
	"github.com/uber-go/tally"
	"go.uber.org/zap"
	"google.golang.org/grpc"

	gatewayv1 "github.com/lyft/clutch/backend/api/config/gateway/v1"
	"github.com/lyft/clutch/backend/middleware"
)

func New(config *gatewayv1.Timeouts, logger *zap.Logger, scope tally.Scope) (middleware.Middleware, error) {
	if config == nil {
		config = &gatewayv1.Timeouts{Default: ptypes.DurationProto(time.Second * 15)}
	}

	defaultTimeout, err := ptypes.Duration(config.Default)
	if err != nil {
		return nil, fmt.Errorf("error in deadline config: %w", err)
	}

	overrides := make(map[string]time.Duration, len(config.Overrides))
	for i, entry := range config.Overrides {
		override, err := ptypes.Duration(entry.Timeout)
		if err != nil {
			return nil, fmt.Errorf("error in %d entry in deadline config: %w", i, err)
		}
		overrides[join(entry.Service, entry.Method)] = override
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

type ret struct {
	resp interface{}
	err error
}

func (m *mid) UnaryInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		service, method, ok := middleware.SplitFullMethod(info.FullMethod)
		if !ok {
			m.logger.Warn("could not parse gRPC method", zap.String("fullMethod", info.FullMethod))
		}

		timeout := m.getDuration(service, method)
		ctx, cancel := context.WithTimeout(ctx, timeout)
		defer cancel()

		resultChan := make(chan ret)
		defer close(resultChan)

		go func() {
			resp, err := handler(ctx, req)
			select {
			case <- ctx.Done():
			default:
				resultChan <- ret{resp: resp, err: err}
			}
		}()

		select {
		case ret := <- resultChan:
			return ret.resp, ret.err
		case <- time.After(timeout):
			return nil, status.New(codes.DeadlineExceeded, "deadline exceeded").Err()
		}
	}
}

func join(service, method string) string {
	const pattern = "/%s/%s"
	return fmt.Sprintf(pattern, service, method)
}

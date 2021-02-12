package unaryerror

// <!-- START clutchdoc -->
// description: Allows services to register error interceptors in order to rewrite errors with additional information.
// <!-- END clutchdoc -->

import (
	"context"

	"google.golang.org/grpc"

	"github.com/lyft/clutch/backend/middleware"
)

// Implement this interface in a service to have the ConvertError function registered into middleware.
type Interceptor interface {
	InterceptError(error) error
}

type errorInterceptorFunc func(error) error

func New() (middleware.Middleware, error) {
	return &Middleware{}, nil
}

type Middleware struct {
	interceptors []errorInterceptorFunc
}

func (m *Middleware) AddInterceptor(fn errorInterceptorFunc) {
	m.interceptors = append(m.interceptors, fn)
}

func (m *Middleware) UnaryInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		resp, err := handler(ctx, req)
		if err != nil {
			for _, fn := range m.interceptors {
				err = fn(err)
			}
		}
		return resp, err
	}
}

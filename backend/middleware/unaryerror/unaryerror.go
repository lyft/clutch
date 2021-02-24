package unaryerror

// <!-- START clutchdoc -->
// description: Allows services to register error interceptors in order to rewrite errors with additional information.
// <!-- END clutchdoc -->

import (
	"context"

	"google.golang.org/grpc"

	"github.com/lyft/clutch/backend/middleware"
)

const Name = "clutch.middleware.unaryerror"

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

// AddInterceptor is used during gateway start to register any service that implements the Interceptor interface.
func (m *Middleware) AddInterceptor(fn errorInterceptorFunc) {
	m.interceptors = append(m.interceptors, fn)
}

func (m *Middleware) UnaryInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		// Invoke handler.
		resp, err := handler(ctx, req)

		// Attempt to transform error if there was one.
		if err != nil {
			// Iterate in reverse order over each interceptor so the 'significant' foundational service's interceptors get applied last.
			for i := len(m.interceptors) - 1; i >= 0; i-- {
				// Apply interceptor and overwrite error.
				err = m.interceptors[i](err)
			}
		}
		return resp, err
	}
}

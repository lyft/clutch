package unaryerror

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"

	healthcheckv1 "github.com/lyft/clutch/backend/api/healthcheck/v1"
)

func TestNew(t *testing.T) {
	m, err := New()
	assert.NoError(t, err)
	assert.NotNil(t, m)
}

func TestEmpty(t *testing.T) {
	m := &Middleware{}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return &healthcheckv1.HealthcheckResponse{}, nil
	}

	midFn := m.UnaryInterceptor()

	resp, err := midFn(context.Background(), nil, &grpc.UnaryServerInfo{FullMethod: "/zip/zoom"}, handler)
	assert.NoError(t, err)
	assert.NotNil(t, resp)
}

func TestSkipsIfErrNil(t *testing.T) {
	m := &Middleware{}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return &healthcheckv1.HealthcheckResponse{}, nil
	}

	m.AddInterceptor(func(error) error {
		panic("should not happen")
	})

	midFn := m.UnaryInterceptor()

	resp, err := midFn(context.Background(), nil, &grpc.UnaryServerInfo{FullMethod: "/zip/zoom"}, handler)
	assert.NoError(t, err)
	assert.NotNil(t, resp)
}

func TestCallsIfErr(t *testing.T) {
	m := &Middleware{}

	originalErr := errors.New("whoopsie daisy")
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return nil, originalErr
	}

	callCount := 0

	newErr := errors.New("this is a rewritten error")
	m.AddInterceptor(func(e error) error {
		callCount += 1
		assert.Equal(t, originalErr, e)
		return newErr
	})

	midFn := m.UnaryInterceptor()
	resp, err := midFn(context.Background(), nil, &grpc.UnaryServerInfo{FullMethod: "/zip/zoom"}, handler)
	assert.Nil(t, resp)
	assert.Error(t, err)
	assert.Equal(t, 1, callCount)
	assert.Equal(t, err, newErr)
}

func TestMultipleInterceptors(t *testing.T) {
	m := &Middleware{}

	e := errors.New("whoopsie daisy")
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return nil, e
	}

	var calls []string

	m.AddInterceptor(func(e error) error {
		calls = append(calls, "a")
		return errors.New("a")
	})

	m.AddInterceptor(func(e error) error {
		calls = append(calls, "b")
		return errors.New("b")
	})

	midFn := m.UnaryInterceptor()
	resp, err := midFn(context.Background(), nil, &grpc.UnaryServerInfo{FullMethod: "/zip/zoom"}, handler)
	assert.Nil(t, resp)
	assert.Error(t, err)
	assert.Equal(t, []string{"b", "a"}, calls)
	assert.Equal(t, err, errors.New("a"))
}

package authz

import (
	"context"
	"sync/atomic"
	"testing"

	"github.com/dgrijalva/jwt-go"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"

	authzv1 "github.com/lyft/clutch/backend/api/authz/v1"
	healthcheckv1 "github.com/lyft/clutch/backend/api/healthcheck/v1"
	"github.com/lyft/clutch/backend/middleware"
	"github.com/lyft/clutch/backend/mock/service/authzmock"
	"github.com/lyft/clutch/backend/service"
	"github.com/lyft/clutch/backend/service/authn"
	"github.com/lyft/clutch/backend/service/authz"
)

func newWithMock(c authz.Client) (middleware.Middleware, error) {
	if c == nil {
		c = authzmock.New()
	}
	service.Registry["clutch.service.authz"] = c
	return New(nil, nil, nil)
}

func TestNew(t *testing.T) {
	m, err := newWithMock(nil)
	assert.NoError(t, err)
	assert.NotNil(t, m)
}

func TestNoClaims(t *testing.T) {
	m, _ := newWithMock(nil)
	interceptor := m.UnaryInterceptor()

	ctx := context.Background()

	req := &healthcheckv1.HealthcheckRequest{}
	info := &grpc.UnaryServerInfo{FullMethod: "/clutch.HealthcheckAPI/Healthcheck"}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return nil, nil
	}

	resp, err := interceptor(ctx, req, info, handler)
	assert.Nil(t, resp)
	assert.Error(t, err)
}

func TestAllowlist(t *testing.T) {
	s := &svcMock{}
	m, _ := newWithMock(s)
	interceptor := m.UnaryInterceptor()

	ctx := context.Background()

	req := &healthcheckv1.HealthcheckRequest{}
	info := &grpc.UnaryServerInfo{FullMethod: "/clutch.authn.v1.AuthnAPI/Callback"}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return &healthcheckv1.HealthcheckResponse{}, nil
	}

	resp, err := interceptor(ctx, req, info, handler)
	assert.NotNil(t, resp)
	assert.NoError(t, err)
	assert.EqualValues(t, 0, s.called)
}

type svcMock struct {
	authz.Client

	called      int32
	lastSubject *authzv1.Subject
}

func (m *svcMock) Check(ctx context.Context, req *authzv1.CheckRequest) (*authzv1.CheckResponse, error) {
	atomic.AddInt32(&m.called, 1)
	m.lastSubject = req.Subject
	return &authzv1.CheckResponse{
		Decision: authzv1.Decision_ALLOW,
	}, nil
}

func TestNoResources(t *testing.T) {
	s := &svcMock{}
	var m, _ = newWithMock(s)
	interceptor := m.UnaryInterceptor()

	claims := &authn.Claims{
		StandardClaims: &jwt.StandardClaims{Subject: "foo@example.com"},
		Groups:         []string{"group-a", "group-b"},
	}
	ctx := authn.ContextWithClaims(context.Background(), claims)

	req := &healthcheckv1.HealthcheckRequest{}
	info := &grpc.UnaryServerInfo{FullMethod: "/clutch.foo/Bar"}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return &healthcheckv1.HealthcheckResponse{}, nil
	}

	resp, err := interceptor(ctx, req, info, handler)
	assert.NotNil(t, resp)
	assert.NoError(t, err)

	assert.Equal(t, claims.Subject, s.lastSubject.User)
	assert.EqualValues(t, claims.Groups, s.lastSubject.Groups)
}

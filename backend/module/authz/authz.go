package authz

// <!-- START clutchdoc -->
// description: Exposes an endpoint for RBAC policy checks on arbitrary resources and users.
// <!-- END clutchdoc -->

import (
	"context"
	"errors"

	"github.com/golang/protobuf/ptypes/any"
	"github.com/uber-go/tally"
	"go.uber.org/zap"

	authzv1 "github.com/lyft/clutch/backend/api/authz/v1"
	"github.com/lyft/clutch/backend/module"
	"github.com/lyft/clutch/backend/service"
	"github.com/lyft/clutch/backend/service/authz"
)

const Name = "clutch.module.authz"

func New(*any.Any, *zap.Logger, tally.Scope) (module.Module, error) {
	svc, ok := service.Registry["clutch.service.authz"]
	if !ok {
		return nil, errors.New("unable to get authz service")
	}

	p, ok := svc.(authz.Client)
	if !ok {
		return nil, errors.New("authz service was not the correct type")
	}

	return &mod{
		authzv1: &api{zclient: p},
	}, nil
}

type mod struct {
	authzv1 authzv1.AuthzAPIServer
}

func (m *mod) Register(r module.Registrar) error {
	authzv1.RegisterAuthzAPIServer(r.GRPCServer(), m.authzv1)
	return r.RegisterJSONGateway(authzv1.RegisterAuthzAPIHandler)
}

type api struct {
	zclient authz.Client
}

func (a *api) Check(ctx context.Context, request *authzv1.CheckRequest) (*authzv1.CheckResponse, error) {
	return a.zclient.Check(ctx, request)
}

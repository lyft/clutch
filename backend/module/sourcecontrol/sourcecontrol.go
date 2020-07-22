package sourcecontrol

// <!-- START clutchdoc -->
// description: Creates repositories in a source control system.
// <!-- END clutchdoc -->

import (
	"context"
	"errors"

	"github.com/golang/protobuf/ptypes/any"
	"github.com/uber-go/tally"
	"go.uber.org/zap"

	sourcecontrolv1 "github.com/lyft/clutch/backend/api/sourcecontrol/v1"
	"github.com/lyft/clutch/backend/module"
	"github.com/lyft/clutch/backend/service"
	"github.com/lyft/clutch/backend/service/github"
)

const Name = "clutch.module.sourcecontrol"

func New(*any.Any, *zap.Logger, tally.Scope) (module.Module, error) {
	svc, ok := service.Registry["clutch.service.github"]
	if !ok {
		return nil, errors.New("could not find service")
	}

	c, ok := svc.(github.Client)
	if !ok {
		return nil, errors.New("service was not the correct type")
	}

	return &mod{c}, nil
}

type mod struct {
	github github.Client
}

func (m *mod) Register(r module.Registrar) error {
	sourcecontrolv1.RegisterSourceControlAPIServer(r.GRPCServer(), m)
	return r.RegisterJSONGateway(sourcecontrolv1.RegisterSourceControlAPIHandler)
}

func (m *mod) CreateRepository(ctx context.Context, req *sourcecontrolv1.CreateRepositoryRequest) (*sourcecontrolv1.CreateRepositoryResponse, error) {
	return m.github.CreateRepository(ctx, req)
}

package project

import (
	"context"
	"errors"

	"github.com/golang/protobuf/ptypes/any"
	projectv1 "github.com/lyft/clutch/backend/api/project/v1"
	"github.com/lyft/clutch/backend/module"
	"github.com/lyft/clutch/backend/service"
	projectservice "github.com/lyft/clutch/backend/service/project"
	"github.com/uber-go/tally"
	"go.uber.org/zap"
)

const (
	Name = "clutch.module.project"
)

func New(cfg *any.Any, log *zap.Logger, scope tally.Scope) (module.Module, error) {
	client, ok := service.Registry["clutch.service.project"]
	if !ok {
		return nil, errors.New("could not find project service")
	}

	svc, ok := client.(projectservice.Service)
	if !ok {
		return nil, errors.New("service was no the correct type")
	}

	return &mod{
		project: newProjectAPI(svc),
		logger:  log,
		scope:   scope,
	}, nil
}

type mod struct {
	project projectv1.ProjectAPIServer
	logger  *zap.Logger
	scope   tally.Scope
}

func (m *mod) Register(r module.Registrar) error {
	projectv1.RegisterProjectAPIServer(r.GRPCServer(), m.project)
	return r.RegisterJSONGateway(projectv1.RegisterProjectAPIHandler)
}

type projectAPI struct {
	project projectservice.Service
}

func newProjectAPI(svc projectservice.Service) projectv1.ProjectAPIServer {
	return &projectAPI{
		project: svc,
	}
}

func (p *projectAPI) GetProjects(ctx context.Context, req *projectv1.GetProjectsRequest) (*projectv1.GetProjectsResponse, error) {
	return p.project.GetProjects(ctx, req)
}

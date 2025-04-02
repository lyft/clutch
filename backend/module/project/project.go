package project

import (
	"context"
	"fmt"

	"github.com/uber-go/tally/v4"
	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/anypb"

	projectv1cfg "github.com/lyft/clutch/backend/api/config/module/project/v1"
	projectv1 "github.com/lyft/clutch/backend/api/project/v1"
	"github.com/lyft/clutch/backend/module"
	"github.com/lyft/clutch/backend/service"
	projectservice "github.com/lyft/clutch/backend/service/project"
)

const (
	Name                  = "clutch.module.project"
	defaultProjectService = "clutch.service.project"
)

func New(cfg *anypb.Any, log *zap.Logger, scope tally.Scope) (module.Module, error) {
	config := &projectv1cfg.Config{}
	err := cfg.UnmarshalTo(config)
	if err != nil {
		return nil, err
	}

	projectSvc := defaultProjectService
	if len(config.ProjectServiceOverride) > 0 {
		projectSvc = config.ProjectServiceOverride
	}

	client, ok := service.Registry[projectSvc]
	if !ok {
		return nil, fmt.Errorf("could not find project service [%s]", projectSvc)
	}

	svc, ok := client.(projectservice.Service)
	if !ok {
		return nil, fmt.Errorf("service [%s] was not the correct type ", projectSvc)
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

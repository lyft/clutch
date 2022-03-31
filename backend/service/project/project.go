package project

import (
	"context"

	"github.com/gogo/status"
	"github.com/golang/protobuf/ptypes/any"
	"github.com/uber-go/tally/v4"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"

	projectv1 "github.com/lyft/clutch/backend/api/project/v1"
	"github.com/lyft/clutch/backend/service"
)

const Name = "clutch.service.project"

type Service interface {
	GetProjects(ctx context.Context, req *projectv1.GetProjectsRequest) (*projectv1.GetProjectsResponse, error)
}

type client struct {
	log   *zap.Logger
	scope tally.Scope
}

func New(cfg *any.Any, logger *zap.Logger, scope tally.Scope) (service.Service, error) {
	c := &client{
		log:   logger,
		scope: scope,
	}

	return c, nil
}

func (c *client) GetProjects(ctx context.Context, req *projectv1.GetProjectsRequest) (*projectv1.GetProjectsResponse, error) {
	return nil, status.Error(codes.Unimplemented, "not implemented")
}

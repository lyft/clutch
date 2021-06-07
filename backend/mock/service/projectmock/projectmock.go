package projectmock

import (
	"context"

	"github.com/golang/protobuf/ptypes/any"
	projectv1core "github.com/lyft/clutch/backend/api/core/project/v1"
	projectv1 "github.com/lyft/clutch/backend/api/project/v1"
	"github.com/lyft/clutch/backend/service"
	projectservice "github.com/lyft/clutch/backend/service/project"
	"github.com/uber-go/tally"
	"go.uber.org/zap"
)

type svc struct{}

func New() projectservice.Service {
	return &svc{}
}

func NewAsService(*any.Any, *zap.Logger, tally.Scope) (service.Service, error) {
	return New(), nil
}

func (s *svc) GetProjects(ctx context.Context, req *projectv1.GetProjectsRequest) (*projectv1.GetProjectsResponse, error) {
	return &projectv1.GetProjectsResponse{
		Projects: []*projectv1core.Project{
			{
				Name: "foo",
			},
			{
				Name: "bar",
			},
		},
	}, nil
}

package topologymock

import (
	"context"
	"errors"

	"github.com/golang/protobuf/ptypes/any"
	"github.com/uber-go/tally"
	"go.uber.org/zap"

	topologyv1 "github.com/lyft/clutch/backend/api/topology/v1"
	"github.com/lyft/clutch/backend/service"
	topologyservice "github.com/lyft/clutch/backend/service/topology"
)

type svc struct{}

func New() topologyservice.Service {
	return &svc{}
}

func NewAsService(*any.Any, *zap.Logger, tally.Scope) (service.Service, error) {
	return New(), nil
}

func (s *svc) GetTopology(ctx context.Context) error {
	return errors.New("Not implemented")
}

func (s *svc) Search(context.Context, *topologyv1.SearchRequest) ([]*topologyv1.Resource, uint64, error) {
	return []*topologyv1.Resource{
		{
			Id: "pod-123",
			Pb: &any.Any{},
			Metadata: map[string]string{
				"label": "value",
			},
		},
	}, 1, nil
}

package topologymock

import (
	"context"
	"errors"

	"github.com/golang/protobuf/ptypes/any"
	"github.com/uber-go/tally"
	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/structpb"

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

func (s *svc) Search(context.Context, *topologyv1.SearchRequest) ([]*topologyv1.Resource, string, error) {
	return []*topologyv1.Resource{
		{
			Id: "pod-123",
			Pb: &any.Any{},
			Metadata: map[string]*structpb.Value{
				"label": &structpb.Value{},
			},
		},
	}, "1", nil
}

func (s *svc) Autocomplete(ctx context.Context, typeURL, search string, limit uint64) ([]*topologyv1.Resource, error) {
	return []*topologyv1.Resource{
		{
			Id: "autocomplete-result",
			Pb: &any.Any{},
			Metadata: map[string]*structpb.Value{
				"label": &structpb.Value{},
			},
		},
	}, nil
}

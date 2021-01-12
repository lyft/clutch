package topology

import (
	"context"
	"errors"

	"github.com/golang/protobuf/ptypes/any"
	"github.com/uber-go/tally"
	"go.uber.org/zap"

	topologyv1 "github.com/lyft/clutch/backend/api/topology/v1"
	"github.com/lyft/clutch/backend/module"
	"github.com/lyft/clutch/backend/service"
	topologyservice "github.com/lyft/clutch/backend/service/topology"
)

const (
	Name = "clutch.module.topology"
)

func New(*any.Any, *zap.Logger, tally.Scope) (module.Module, error) {
	client, ok := service.Registry["clutch.service.topology"]
	if !ok {
		return nil, errors.New("could not find topology service")
	}

	svc, ok := client.(topologyservice.Service)
	if !ok {
		return nil, errors.New("service was no the correct type")
	}

	mod := &mod{
		topology: newTopologyAPI(svc),
	}

	return mod, nil
}

type mod struct {
	topology topologyv1.TopologyAPIServer
}

func (m *mod) Register(r module.Registrar) error {
	topologyv1.RegisterTopologyAPIServer(r.GRPCServer(), m.topology)
	return r.RegisterJSONGateway(topologyv1.RegisterTopologyAPIHandler)
}

type topologyAPI struct {
	topology topologyservice.Service
}

func newTopologyAPI(svc topologyservice.Service) topologyv1.TopologyAPIServer {
	return &topologyAPI{
		topology: svc,
	}
}

func (t *topologyAPI) GetTopology(ctx context.Context, req *topologyv1.GetTopologyRequest) (*topologyv1.GetTopologyResponse, error) {
	return nil, errors.New("not implemented")
}

func (t *topologyAPI) SearchTopology(ctx context.Context, req *topologyv1.SearchTopologyRequest) (*topologyv1.SearchTopologyResponse, error) {
	resources, err := t.topology.SearchTopology(ctx, req)
	if err != nil {
		return nil, err
	}

	return &topologyv1.SearchTopologyResponse{
		Resources: resources,
	}, nil
}

package topology

import (
	"errors"

	"github.com/golang/protobuf/ptypes/any"
	topologyv1 "github.com/lyft/clutch/backend/api/topology/v1"
	"github.com/lyft/clutch/backend/module"
	"github.com/lyft/clutch/backend/service"
	"github.com/lyft/clutch/backend/service/topology"
	"github.com/uber-go/tally"
	"go.uber.org/zap"
)

const (
	Name = "clutch.module.topology"
)

func New(*any.Any, *zap.Logger, tally.Scope) (module.Module, error) {
	client, ok := service.Registry["clutch.service.topology"]
	if !ok {
		return nil, errors.New("could not find service")
	}

	svc, ok := client.(topology.Service)
	if !ok {
		return nil, errors.New("service was not the correct type")
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
	topology topology.Service
}

func newTopologyAPI(svc topology.Service) topologyv1.TopologyAPIServer {
	return &topologyAPI{
		topology: svc,
	}
}

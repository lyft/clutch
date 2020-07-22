package k8s

// <!-- START clutchdoc -->
// description: Endpoints for interacting with Kubernetes resources.
// <!-- END clutchdoc -->

import (
	"errors"

	"github.com/golang/protobuf/ptypes/any"
	"github.com/uber-go/tally"
	"go.uber.org/zap"

	k8sv1 "github.com/lyft/clutch/backend/api/k8s/v1"
	"github.com/lyft/clutch/backend/module"
	"github.com/lyft/clutch/backend/service"
	k8sservice "github.com/lyft/clutch/backend/service/k8s"
)

const (
	Name = "clutch.module.k8s"
)

func New(*any.Any, *zap.Logger, tally.Scope) (module.Module, error) {
	client, ok := service.Registry["clutch.service.k8s"]
	if !ok {
		return nil, errors.New("could not find service")
	}

	svc, ok := client.(k8sservice.Service)
	if !ok {
		return nil, errors.New("service was not the correct type")
	}

	mod := &mod{
		k8s: newK8sAPI(svc),
	}

	return mod, nil
}

type mod struct {
	k8s k8sv1.K8SAPIServer
}

func (m *mod) Register(r module.Registrar) error {
	k8sv1.RegisterK8SAPIServer(r.GRPCServer(), m.k8s)
	return r.RegisterJSONGateway(k8sv1.RegisterK8SAPIHandler)
}

type k8sAPI struct {
	k8s k8sservice.Service
}

func newK8sAPI(svc k8sservice.Service) k8sv1.K8SAPIServer {
	return &k8sAPI{
		k8s: svc,
	}
}

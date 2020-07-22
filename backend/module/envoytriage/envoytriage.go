package envoytriage

// <!-- START clutchdoc -->
// description: Remote debugging APIs for Envoy Proxy.
// <!-- END clutchdoc -->

import (
	"context"
	"errors"

	"github.com/golang/protobuf/ptypes/any"
	"github.com/uber-go/tally"
	"go.uber.org/zap"

	envoytriagev1 "github.com/lyft/clutch/backend/api/envoytriage/v1"
	"github.com/lyft/clutch/backend/module"
	"github.com/lyft/clutch/backend/service"
	"github.com/lyft/clutch/backend/service/envoyadmin"
)

const (
	Name = "clutch.module.envoytriage"
)

func New(*any.Any, *zap.Logger, tally.Scope) (module.Module, error) {
	client, ok := service.Registry["clutch.service.envoyadmin"]
	if !ok {
		return nil, errors.New("could not find service")
	}

	c, ok := client.(envoyadmin.Client)
	if !ok {
		return nil, errors.New("service was not the correct type")
	}

	mod := &mod{
		api: newAPI(c),
	}
	return mod, nil
}

type mod struct {
	api envoytriagev1.EnvoyTriageAPIServer
}

func (m *mod) Register(r module.Registrar) error {
	envoytriagev1.RegisterEnvoyTriageAPIServer(r.GRPCServer(), m.api)
	return r.RegisterJSONGateway(envoytriagev1.RegisterEnvoyTriageAPIHandler)
}

func newAPI(client envoyadmin.Client) envoytriagev1.EnvoyTriageAPIServer {
	return &api{client: client}
}

type api struct {
	client envoyadmin.Client
}

func (a *api) Read(ctx context.Context, request *envoytriagev1.ReadRequest) (*envoytriagev1.ReadResponse, error) {
	resp := &envoytriagev1.ReadResponse{
		Results: make([]*envoytriagev1.Result, len(request.Operations)),
	}

	// TODO(maybe): fan-out, fan-in (but limit concurrency)
	for idx, op := range request.Operations {
		res, err := a.client.Get(ctx, op)
		if err != nil {
			return nil, err
		}
		resp.Results[idx] = res
	}

	return resp, nil
}

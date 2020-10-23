package healthcheck

// <!-- START clutchdoc -->
// description: Simple healthcheck endpoint.
// <!-- END clutchdoc -->

import (
	"context"
	"time"

	"github.com/golang/protobuf/ptypes/any"
	"github.com/uber-go/tally"
	"go.uber.org/zap"

	healthcheckv1 "github.com/lyft/clutch/backend/api/healthcheck/v1"
	"github.com/lyft/clutch/backend/module"
)

const (
	Name = "clutch.module.healthcheck"
)

func New(*any.Any, *zap.Logger, tally.Scope) (module.Module, error) {
	mod := &mod{
		api: newAPI(),
	}
	return mod, nil
}

type mod struct {
	api healthcheckv1.HealthcheckAPIServer
}

func (m *mod) Register(r module.Registrar) error {
	healthcheckv1.RegisterHealthcheckAPIServer(r.GRPCServer(), m.api)
	return r.RegisterJSONGateway(healthcheckv1.RegisterHealthcheckAPIHandler)
}

func newAPI() healthcheckv1.HealthcheckAPIServer {
	return &healthcheckAPI{}
}

type healthcheckAPI struct{}

func (a *healthcheckAPI) Healthcheck(ctx context.Context, _ *healthcheckv1.HealthcheckRequest) (*healthcheckv1.HealthcheckResponse, error) {
	time.Sleep(30 * time.Second)
	return &healthcheckv1.HealthcheckResponse{}, nil
}

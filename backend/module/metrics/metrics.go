package metrics

import (
	"context"
	"fmt"

	"github.com/golang/protobuf/ptypes/any"
	"github.com/uber-go/tally/v4"
	"go.uber.org/zap"

	metricsv1 "github.com/lyft/clutch/backend/api/metrics/v1"
	"github.com/lyft/clutch/backend/module"
	"github.com/lyft/clutch/backend/service"
	"github.com/lyft/clutch/backend/service/metrics"
)

const (
	Name = "clutch.module.metrics"
)

type mod struct {
	metricsservice metrics.Service
	logger         *zap.Logger
	scope          tally.Scope
}

func New(_ *any.Any, log *zap.Logger, scope tally.Scope) (module.Module, error) {
	svc, ok := service.Registry[metrics.Name]
	if !ok {
		return nil, fmt.Errorf("could not find metrics service [%s]", metrics.Name)
	}

	metricsservice, ok := svc.(metrics.Service)
	if !ok {
		return nil, fmt.Errorf("metrics service [%s] is not a metrics service", metrics.Name)
	}

	return &mod{
		metricsservice: metricsservice,
		logger:         log,
		scope:          scope,
	}, nil
}

func (m *mod) Register(r module.Registrar) error {
	metricsv1.RegisterMetricsAPIServer(r.GRPCServer(), m)
	return r.RegisterJSONGateway(metricsv1.RegisterMetricsAPIHandler)
}

func (m *mod) GetMetrics(ctx context.Context, req *metricsv1.GetMetricsRequest) (*metricsv1.GetMetricsResponse, error) {
	return m.metricsservice.GetMetrics(ctx, req)
}

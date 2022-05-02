package metrics

import (
	"context"

	"github.com/golang/protobuf/ptypes/any"
	"github.com/uber-go/tally/v4"
	"go.uber.org/zap"

	"net/http"

	metricsv1cfg "github.com/lyft/clutch/backend/api/config/service/metrics/v1"
	metricsv1 "github.com/lyft/clutch/backend/api/metrics/v1"
	"github.com/lyft/clutch/backend/service"
)

const Name = "clutch.service.metrics"

type Service interface {
	GetMetrics(context.Context, *metricsv1.GetMetricsRequest) (*metricsv1.GetMetricsResponse, error)
}

type client struct {
	prometheusAPI         *http.Client
	prometheusAPIEndpoint string
	log                   *zap.Logger
	scope                 tally.Scope
}

func New(cfg *any.Any, logger *zap.Logger, scope tally.Scope) (service.Service, error) {
	return NewWithHTTPClient(cfg, logger, scope, &http.Client{})
}

func NewWithHTTPClient(cfg *any.Any, logger *zap.Logger, scope tally.Scope, prometheusAPI *http.Client) (service.Service, error) {
	config := &metricsv1cfg.Config{}
	if err := cfg.UnmarshalTo(config); err != nil {
		return nil, err
	}

	return &client{
		prometheusAPI:         prometheusAPI,
		prometheusAPIEndpoint: config.PrometheusApiEndpoint,
		log:                   logger,
		scope:                 scope,
	}, nil
}

func (c *client) GetMetrics(context.Context, *metricsv1.GetMetricsRequest) (*metricsv1.GetMetricsResponse, error) {
	return nil, nil
}

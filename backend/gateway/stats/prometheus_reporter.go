package stats

import (
	gatewayv1 "github.com/lyft/clutch/backend/api/config/gateway/v1"
	tally_prom "github.com/uber-go/tally/prometheus"
)

func NewPrometheusReporter(cfg *gatewayv1.Stats_PrometheusReporter) (tally_prom.Reporter, error) {
	promCfg := tally_prom.Configuration{
		HandlerPath: cfg.HandlerPath,
	}
	return promCfg.NewReporter(tally_prom.ConfigurationOptions{})
}

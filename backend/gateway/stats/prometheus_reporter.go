package stats

import (
	tallyprom "github.com/uber-go/tally/prometheus"

	gatewayv1 "github.com/lyft/clutch/backend/api/config/gateway/v1"
)

func NewPrometheusReporter(cfg *gatewayv1.Stats_PrometheusReporter) (tallyprom.Reporter, error) {
	promCfg := tallyprom.Configuration{
		HandlerPath: cfg.HandlerPath,
	}
	return promCfg.NewReporter(tallyprom.ConfigurationOptions{})
}

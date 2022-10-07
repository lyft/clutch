package stats

import (
	gatewayv1 "github.com/lyft/clutch/backend/api/config/gateway/v1"
	tallyprom "github.com/uber-go/tally/v4/prometheus"
)

func NewPrometheusReporter(cfg *gatewayv1.Stats_PrometheusReporter) (tallyprom.Reporter, error) {
	promCfg := tallyprom.Configuration{
		HandlerPath: cfg.HandlerPath,
	}
	return promCfg.NewReporter(tallyprom.ConfigurationOptions{})
}

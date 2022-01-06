package temporal

import (
	"github.com/uber-go/tally"
	"go.temporal.io/sdk/client"
)

type metricsHandler struct {
	tally.Scope
}

func (m *metricsHandler) Counter(name string) client.MetricsCounter {
	return m.Scope.Counter(name)
}

func (m *metricsHandler) Gauge(name string) client.MetricsGauge {
	return m.Scope.Gauge(name)
}

func (m *metricsHandler) Timer(name string) client.MetricsTimer {
	return m.Scope.Timer(name)
}

func (m *metricsHandler) WithTags(tags map[string]string) client.MetricsHandler {
	return &metricsHandler{Scope: m.Scope.Tagged(tags)}
}

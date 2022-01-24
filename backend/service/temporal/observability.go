package temporal

import (
	"github.com/uber-go/tally"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/log"
	"go.uber.org/zap"
)

type metricsHandler struct {
	scope tally.Scope
}

func newMetricsHandler(scope tally.Scope) client.MetricsHandler {
	return &metricsHandler{scope: scope}
}

func (m *metricsHandler) Counter(name string) client.MetricsCounter {
	return m.scope.Counter(name)
}

func (m *metricsHandler) Gauge(name string) client.MetricsGauge {
	return m.scope.Gauge(name)
}

func (m *metricsHandler) Timer(name string) client.MetricsTimer {
	return m.scope.Timer(name)
}

func (m *metricsHandler) WithTags(tags map[string]string) client.MetricsHandler {
	return &metricsHandler{scope: m.scope.Tagged(tags)}
}

func newTemporalLogger(logger *zap.Logger) log.Logger {
	return &temporalLogger{
		sl: logger.Sugar(),
	}
}

type temporalLogger struct {
	sl *zap.SugaredLogger
}

func (t *temporalLogger) Debug(msg string, keyvals ...interface{}) {
	t.sl.Debugw(msg, keyvals...)
}

func (t *temporalLogger) Info(msg string, keyvals ...interface{}) {
	t.sl.Infow(msg, keyvals...)
}

func (t *temporalLogger) Warn(msg string, keyvals ...interface{}) {
	t.sl.Warnw(msg, keyvals...)
}

func (t *temporalLogger) Error(msg string, keyvals ...interface{}) {
	t.sl.Errorw(msg, keyvals...)
}

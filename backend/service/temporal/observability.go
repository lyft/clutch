package temporal

import (
	"github.com/uber-go/tally"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/log"
	"go.uber.org/zap"
)

type metricsHandler struct {
	tally.Scope
}

func newMetricsHandler(scope tally.Scope) client.MetricsHandler {
	return &metricsHandler{Scope: scope}
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

func newTemporalLogger(logger *zap.Logger) log.Logger {
	return &temporalLogger{
		SugaredLogger: logger.Sugar(),
	}
}

type temporalLogger struct {
	*zap.SugaredLogger
}

func (t *temporalLogger) Debug(msg string, keyvals ...interface{}) {
	t.SugaredLogger.Debugw(msg, keyvals...)
}

func (t *temporalLogger) Info(msg string, keyvals ...interface{}) {
	t.SugaredLogger.Infow(msg, keyvals...)
}

func (t *temporalLogger) Warn(msg string, keyvals ...interface{}) {
	t.SugaredLogger.Warnw(msg, keyvals...)
}

func (t *temporalLogger) Error(msg string, keyvals ...interface{}) {
	t.SugaredLogger.Errorw(msg, keyvals...)
}

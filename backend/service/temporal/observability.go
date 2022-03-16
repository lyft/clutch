package temporal

import (
	"go.temporal.io/sdk/log"
	"go.uber.org/zap"
)

func newTemporalLogger(logger *zap.Logger) log.Logger {
	return &temporalLogger{
		sl: logger.Sugar(),
	}
}

type temporalLogger struct {
	sl *zap.SugaredLogger
}

func (t *temporalLogger) Debug(msg string, keyvals ...any) {
	t.sl.Debugw(msg, keyvals...)
}

func (t *temporalLogger) Info(msg string, keyvals ...any) {
	t.sl.Infow(msg, keyvals...)
}

func (t *temporalLogger) Warn(msg string, keyvals ...any) {
	t.sl.Warnw(msg, keyvals...)
}

func (t *temporalLogger) Error(msg string, keyvals ...any) {
	t.sl.Errorw(msg, keyvals...)
}

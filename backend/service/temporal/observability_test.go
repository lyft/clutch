package temporal

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"go.uber.org/zap/zaptest/observer"
)

func TestLogWrapper(t *testing.T) {
	core, logs := observer.New(zap.DebugLevel)
	tl := newTemporalLogger(zap.New(core))
	args := []interface{}{zap.String("foo", "bar"), zap.String("blip", "pop")}
	tl.Debug("debugz", args...)
	tl.Info("infoz", args...)
	tl.Warn("warnz", args...)
	tl.Error("errorz", args...)

	all := logs.All()
	assert.Equal(t, "debugz", all[0].Message)
	assert.Equal(t, zap.DebugLevel, all[0].Level)
	assert.Equal(t, "infoz", all[1].Message)
	assert.Equal(t, zap.InfoLevel, all[1].Level)
	assert.Equal(t, "warnz", all[2].Message)
	assert.Equal(t, zap.WarnLevel, all[2].Level)
	assert.Equal(t, "errorz", all[3].Message)
	assert.Equal(t, zap.ErrorLevel, all[3].Level)

	for _, l := range all {
		assert.Equal(t, "bar", l.ContextMap()["foo"])
		assert.Equal(t, "pop", l.ContextMap()["blip"])
	}
}

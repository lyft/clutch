package temporal

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/uber-go/tally"
	"go.uber.org/zap"
	"go.uber.org/zap/zaptest/observer"
)

func TestMetricsWrapper(t *testing.T) {
	s := tally.NewTestScope("foo", map[string]string{"test_tag": "tagged"})
	mh := newMetricsHandler(s)

	mh.Counter("num_occurrences").Inc(3)
	mh.Timer("elapsed").Record(time.Second * 5)
	mh.Gauge("num_active").Update(6)

	taggedScope := mh.WithTags(map[string]string{"test_tag": "retagged"})
	taggedScope.Counter("num_occurrences").Inc(4)

	snap := s.Snapshot()
	assert.EqualValues(t, snap.Counters()["foo.num_occurrences+test_tag=tagged"].Value(), 3)
	assert.EqualValues(t, snap.Counters()["foo.num_occurrences+test_tag=retagged"].Value(), 4)
	assert.Len(t, snap.Timers()["foo.elapsed+test_tag=tagged"].Values(), 1)
	assert.Equal(t, snap.Timers()["foo.elapsed+test_tag=tagged"].Values()[0], 5*time.Second)
	assert.EqualValues(t, snap.Gauges()["foo.num_active+test_tag=tagged"].Value(), 6)
}

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

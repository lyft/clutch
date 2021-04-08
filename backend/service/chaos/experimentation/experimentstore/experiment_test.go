package experimentstore

import (
	"fmt"
	"testing"
	"time"

	"github.com/golang/protobuf/ptypes/any"
	"github.com/stretchr/testify/assert"
)

func TestExperimentRunConfigToExperiment(t *testing.T) {
	now := time.Now()

	tests := []struct {
		startTime time.Time
		endTime   *time.Time
	}{
		{
			startTime: now,
			endTime:   nil,
		},
		{
			startTime: now,
			endTime:   &now,
		},
	}

	for idx, tt := range tests {
		tt := tt
		t.Run(fmt.Sprintf("%d", idx), func(t *testing.T) {
			t.Parallel()

			run := &ExperimentRun{Id: "1", StartTime: tt.startTime, EndTime: tt.endTime}
			config := &ExperimentConfig{id: "2", Config: &any.Any{}}
			p := Experiment{Run: run, Config: config}
			_, err := p.toProto()
			assert.NoError(t, err)
		})
	}
}

package experimentstore

import (
	"fmt"
	"testing"
	"time"

	"github.com/golang/protobuf/ptypes/any"
	"github.com/stretchr/testify/assert"

	experimentationv1 "github.com/lyft/clutch/backend/api/chaos/experimentation/v1"
)

func TestExperimentRunConfigToExperiment(t *testing.T) {
	now := time.Now()
	future := now.Add(time.Minute)
	past := now.Add(-time.Minute)

	tests := []struct {
		startTime        time.Time
		endTime          *time.Time
		cancellationTime *time.Time
		now              time.Time
		expectedStatus   experimentationv1.Experiment_Status
	}{
		{
			startTime:      past,
			endTime:        nil,
			now:            now,
			expectedStatus: experimentationv1.Experiment_STATUS_RUNNING,
		},
		{
			startTime:      past,
			endTime:        &future,
			now:            now,
			expectedStatus: experimentationv1.Experiment_STATUS_RUNNING,
		},
		{
			startTime:      now,
			endTime:        &future,
			now:            past,
			expectedStatus: experimentationv1.Experiment_STATUS_SCHEDULED,
		},
		{
			startTime:      past,
			endTime:        &now,
			now:            future,
			expectedStatus: experimentationv1.Experiment_STATUS_COMPLETED,
		},
		{
			startTime:        now,
			endTime:          &future,
			cancellationTime: &now,
			now:              past,
			expectedStatus:   experimentationv1.Experiment_STATUS_CANCELED,
		},
		{
			startTime:        past,
			endTime:          &future,
			cancellationTime: &now,
			now:              now,
			expectedStatus:   experimentationv1.Experiment_STATUS_CANCELED,
		},
	}

	for idx, tt := range tests {
		tt := tt
		t.Run(fmt.Sprintf("%d", idx), func(t *testing.T) {
			t.Parallel()

			run := &ExperimentRun{Id: "1", StartTime: tt.startTime, EndTime: tt.endTime, CancellationTime: tt.cancellationTime}
			config := &ExperimentConfig{Id: "2", Config: &any.Any{}}
			e := Experiment{Run: run, Config: config}
			proto, err := e.Proto(tt.now)
			assert.NoError(t, err)
			assert.Equal(t, tt.expectedStatus, proto.Status)
		})
	}
}

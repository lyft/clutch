package experimentstore

import (
	"fmt"
	"testing"
	"time"

	"github.com/golang/protobuf/ptypes/any"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap/zaptest"

	experimentation "github.com/lyft/clutch/backend/api/chaos/experimentation/v1"
)

func TestExperimentRunDetailsStatus(t *testing.T) {
	now := time.Now()
	past := now.Add(-1 * time.Hour)
	future := now.Add(1 * time.Hour)
	farFuture := future.Add(1 * time.Hour)

	tests := []struct {
		startTime        time.Time
		endTime          *time.Time
		cancellationTime *time.Time
		creationTime     time.Time
		now              time.Time
		expectedStatus   experimentation.Experiment_Status
	}{
		{
			startTime:        future,
			endTime:          nil,
			cancellationTime: nil,
			creationTime:     now,
			now:              now,
			expectedStatus:   experimentation.Experiment_STATUS_SCHEDULED,
		},
		{
			startTime:        now,
			endTime:          nil,
			cancellationTime: &past,
			creationTime:     past,
			now:              now,
			expectedStatus:   experimentation.Experiment_STATUS_CANCELED,
		},
		{
			startTime:        now,
			endTime:          nil,
			cancellationTime: nil,
			creationTime:     now,
			now:              future,
			expectedStatus:   experimentation.Experiment_STATUS_RUNNING,
		},
		{
			startTime:        now,
			endTime:          &farFuture,
			cancellationTime: &future,
			creationTime:     now,
			now:              future,
			expectedStatus:   experimentation.Experiment_STATUS_STOPPED,
		},
		{
			startTime:        now,
			endTime:          &future,
			cancellationTime: nil,
			creationTime:     now,
			now:              farFuture,
			expectedStatus:   experimentation.Experiment_STATUS_COMPLETED,
		},
	}

	for idx, tt := range tests {
		tt := tt
		t.Run(fmt.Sprintf("%d", idx), func(t *testing.T) {
			t.Parallel()

			run := &ExperimentRun{
				Id:               "1",
				StartTime:        tt.startTime,
				EndTime:          tt.endTime,
				CancellationTime: tt.cancellationTime,
				CreationTime:     tt.creationTime,
			}
			config := &ExperimentConfig{Id: "2", Config: &any.Any{}}
			tr := NewTransformer(zaptest.NewLogger(t).Sugar())
			rd, err := NewRunDetails(run, config, &tr, tt.now)
			assert.NoError(t, err)
			assert.Equal(t, tt.expectedStatus, rd.Status)
		})
	}
}

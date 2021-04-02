package experimentstore

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestRunningExperimentRunProperties(t *testing.T) {
	now := time.Now()
	past := now.Add(-1 * time.Hour)
	future := now.Add(1 * time.Hour)

	tests := []struct {
		startTime           time.Time
		endTime             *time.Time
		cancellationTime    *time.Time
		creationTime        time.Time
		now                 time.Time
		terminationReason   string
		expectedPropertyIds []string
	}{
		{
			startTime:         now,
			endTime:           &future,
			cancellationTime:  nil,
			creationTime:      now,
			now:               now,
			terminationReason: "",
			expectedPropertyIds: []string{
				"run_identifier",
				"status",
				"run_creation_time",
				"start_time",
				"end_time",
			},
		},
		{
			startTime:         now,
			endTime:           nil,
			cancellationTime:  nil,
			creationTime:      now,
			now:               now,
			terminationReason: "",
			expectedPropertyIds: []string{
				"run_identifier",
				"status",
				"run_creation_time",
				"start_time",
				"end_time",
			},
		},
		{
			startTime:         past,
			endTime:           &future,
			cancellationTime:  &now,
			creationTime:      past,
			now:               past,
			terminationReason: "",
			expectedPropertyIds: []string{
				"run_identifier",
				"status",
				"run_creation_time",
				"start_time",
				"end_time",
				"stopped_at",
				"termination_reason",
			},
		},
		{
			startTime:         now,
			endTime:           &future,
			cancellationTime:  &past,
			creationTime:      past,
			now:               now,
			terminationReason: "",
			expectedPropertyIds: []string{
				"run_identifier",
				"status",
				"run_creation_time",
				"start_time",
				"end_time",
				"canceled_at",
				"termination_reason",
			},
		},
	}

	for idx, tt := range tests {
		tt := tt
		t.Run(fmt.Sprintf("%d", idx), func(t *testing.T) {
			t.Parallel()
			run := &ExperimentRun{
				Id:                "1",
				StartTime:         tt.startTime,
				EndTime:           tt.endTime,
				CancellationTime:  tt.cancellationTime,
				CreationTime:      tt.creationTime,
				TerminationReason: tt.terminationReason,
			}

			properties, err := run.CreateProperties(tt.now)
			assert.NoError(t, err)
			var ids []string
			for _, p := range properties {
				ids = append(ids, p.Id)
			}

			assert.Equal(t, tt.expectedPropertyIds, ids)
		})
	}
}

package experimentstore

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	experimentationv1 "github.com/lyft/clutch/backend/api/chaos/experimentation/v1"
)

func TestExperimentRunStatus(t *testing.T) {
	now := time.Now()
	past := now.Add(-1 * time.Hour)
	future := now.Add(1 * time.Hour)

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
			startTime:      now,
			endTime:        nil,
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
			startTime:        past,
			endTime:          nil,
			cancellationTime: &now,
			now:              now,
			expectedStatus:   experimentationv1.Experiment_STATUS_COMPLETED,
		},
		{
			startTime:        now,
			endTime:          &future,
			cancellationTime: &now,
			now:              past,
			expectedStatus:   experimentationv1.Experiment_STATUS_CANCELED,
		},
		{
			startTime:        now,
			endTime:          nil,
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
			s := run.Status(tt.now)
			assert.Equal(t, tt.expectedStatus, s)
		})
	}
}

func TestExperimentRunProperties(t *testing.T) {
	now := time.Now()
	past := now.Add(-1 * time.Hour)
	future := now.Add(1 * time.Hour)

	tests := []struct {
		startTime              time.Time
		endTime                *time.Time
		cancellationTime       *time.Time
		creationTime           time.Time
		now                    time.Time
		terminationReason      string
		expectedPropertyIds    []string
		expectedPropertyValues map[string]string
	}{
		{
			startTime:         now,
			endTime:           &future,
			cancellationTime:  nil,
			creationTime:      now,
			now:               past,
			terminationReason: "foo",
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
			now:               future,
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
			expectedPropertyValues: map[string]string{
				"termination_reason": "N/A",
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
			expectedPropertyValues: map[string]string{
				"termination_reason": "N/A",
			},
		},
		{
			startTime:         now,
			endTime:           &future,
			cancellationTime:  &past,
			creationTime:      past,
			now:               now,
			terminationReason: "foo",
			expectedPropertyIds: []string{
				"run_identifier",
				"status",
				"run_creation_time",
				"start_time",
				"end_time",
				"canceled_at",
				"termination_reason",
			},
			expectedPropertyValues: map[string]string{
				"termination_reason": "foo",
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
				if val, ok := tt.expectedPropertyValues[p.Id]; ok {
					assert.Equal(t, val, p.GetStringValue())
				}
			}

			assert.Equal(t, tt.expectedPropertyIds, ids)
		})
	}
}

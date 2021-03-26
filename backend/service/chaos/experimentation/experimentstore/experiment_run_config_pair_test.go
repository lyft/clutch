package experimentstore

import (
	"database/sql"
	"fmt"
	"testing"
	"time"

	"github.com/golang/protobuf/ptypes/any"
	"github.com/stretchr/testify/assert"
)

func TestExperimentRunConfigToExperiment(t *testing.T) {
	tests := []struct {
		startTime time.Time
		endTime   sql.NullTime
	}{
		{
			startTime: time.Now(),
			endTime:   sql.NullTime{Valid: false},
		},
		{
			startTime: time.Now(),
			endTime:   sql.NullTime{Time: time.Now(), Valid: true},
		},
	}

	a := assert.New(t)
	for idx, tt := range tests {
		tt := tt
		t.Run(fmt.Sprintf("%d", idx), func(t *testing.T) {
			t.Parallel()

			run := &ExperimentRun{Id: "1", StartTime: tt.startTime, EndTime: tt.endTime}
			config := &ExperimentConfig{id: "2", Config: &any.Any{}}
			p := experimentRunConfigPair{run: run, config: config}
			_, err := p.toExperiment()
			a.NoError(err)
		})
	}
}
package experimentstore

import (
	"database/sql"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestRunningExperimentRunProperties(t *testing.T) {
	run := &ExperimentRun{Id: 1, StartTime: startTime, EndTime: sql.NullTime{}, CancellationTime: sql.NullTime{}, creationTime: creationTime}
	properties, err := run.CreateProperties(time.Now())

	assert := assert.New(t)
	assert.NoError(err)

	var identifiers []string
	for _, p := range properties {
		identifiers = append(identifiers, p.Id)
	}
	expectedIdentifiers := []string{
		"run_identifier",
		"status",
		"run_creation_time",
		"start_time",
		"end_time",
	}
	assert.Equal(expectedIdentifiers, identifiers)
	assert.Equal(int64(1), properties[0].GetIntValue())
}

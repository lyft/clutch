package experimentstore

import (
	"database/sql"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	experimentation "github.com/lyft/clutch/backend/api/chaos/experimentation/v1"
)

var startTime = time.Date(2020, 11, 17, 20, 34, 58, 651387237, time.UTC)
var endTime = time.Date(2020, 11, 20, 20, 34, 58, 651387237, time.UTC)

func TestScheduledExperimentStatus(t *testing.T) {
	start := sql.NullTime{
		Time:  startTime,
		Valid: true,
	}
	end := sql.NullTime{Valid: false}
	cancellation := sql.NullTime{Valid: false}

	now := startTime.AddDate(0, 0, -1)

	status, err := TimesToStatus(start, end, cancellation, now)

	assert := assert.New(t)
	assert.NoError(err)
	assert.Equal(experimentation.ExperimentRunDetails_SCHEDULED, status)
}

func TestCanceledExperimentStatus(t *testing.T) {
	start := sql.NullTime{
		Time:  startTime,
		Valid: true,
	}
	end := sql.NullTime{Valid: false}
	cancellation := sql.NullTime{
		Time:  startTime.AddDate(0, 0, -1),
		Valid: true,
	}

	status, err := TimesToStatus(start, end, cancellation, time.Now())

	assert := assert.New(t)
	assert.NoError(err)
	assert.Equal(experimentation.ExperimentRunDetails_CANCELED, status)
}

func TestRunningExperimentStatus(t *testing.T) {
	now := startTime.AddDate(0, 0, 1)
	start := sql.NullTime{
		Time:  startTime,
		Valid: true,
	}
	end := sql.NullTime{Valid: false}
	cancellation := sql.NullTime{Valid: false}

	status, err := TimesToStatus(start, end, cancellation, now)

	assert := assert.New(t)
	assert.NoError(err)
	assert.Equal(experimentation.ExperimentRunDetails_RUNNING, status)
}

func TestStoppedExperimentStatus(t *testing.T) {
	start := sql.NullTime{
		Time:  startTime,
		Valid: true,
	}
	end := sql.NullTime{
		Time:  endTime,
		Valid: true,
	}
	cancellation := sql.NullTime{
		Time:  startTime.AddDate(0, 0, 1),
		Valid: true,
	}

	status, err := TimesToStatus(start, end, cancellation, time.Now())

	assert := assert.New(t)
	assert.NoError(err)
	assert.Equal(experimentation.ExperimentRunDetails_STOPPED, status)
}

func TestCompletedExperimentStatus(t *testing.T) {
	start := sql.NullTime{
		Time:  startTime,
		Valid: true,
	}
	end := sql.NullTime{
		Time:  endTime,
		Valid: true,
	}
	cancellation := sql.NullTime{
		Valid: false,
	}
	now := endTime.AddDate(0, 0, 1)

	status, err := TimesToStatus(start, end, cancellation, now)

	assert := assert.New(t)
	assert.NoError(err)
	assert.Equal(experimentation.ExperimentRunDetails_COMPLETED, status)
}

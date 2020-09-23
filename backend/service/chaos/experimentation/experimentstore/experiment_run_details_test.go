package experimentstore

import (
	"database/sql"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	experimentation "github.com/lyft/clutch/backend/api/chaos/experimentation/v1"
)

var startTime = time.Date(2020, 11, 17, 20, 34, 58, 651387237, time.UTC)
var creationTime = startTime
var endTime = time.Date(2020, 11, 20, 20, 34, 58, 651387237, time.UTC)
var details = `{"@type":"type.googleapis.com/clutch.chaos.serverexperimentation.v1.TestConfig","clusterPair":{"downstreamCluster":"upstreamCluster","upstreamCluster":"downstreamCluster"},"abort":{"percent":100,"httpStatus":401}}`

func TestScheduledExperiment(t *testing.T) {
	start := sql.NullTime{
		Time:  startTime,
		Valid: true,
	}
	end := sql.NullTime{Valid: false}
	cancellation := sql.NullTime{Valid: false}

	now := startTime.AddDate(0, 0, -1)

	runDetails, err := NewRunDetails(uint64(1), creationTime, start, end, cancellation, now, details)

	assert := assert.New(t)
	assert.NoError(err)
	assert.Equal(experimentation.ExperimentRunDetails_SCHEDULED, runDetails.Status)
}

func TestCanceledExperiment(t *testing.T) {
	start := sql.NullTime{
		Time:  startTime,
		Valid: true,
	}
	end := sql.NullTime{Valid: false}
	cancellation := sql.NullTime{
		Time:  startTime.AddDate(0, 0, -1),
		Valid: true,
	}

	runDetails, err := NewRunDetails(uint64(1), creationTime, start, end, cancellation, time.Now(), details)

	assert := assert.New(t)
	assert.NoError(err)
	assert.Equal(experimentation.ExperimentRunDetails_CANCELED, runDetails.Status)
}

func TestRunningExperiment(t *testing.T) {
	now := startTime.AddDate(0, 0, 1)
	start := sql.NullTime{
		Time:  startTime,
		Valid: true,
	}
	end := sql.NullTime{Valid: false}
	cancellation := sql.NullTime{Valid: false}

	runDetails, err := NewRunDetails(uint64(1), creationTime, start, end, cancellation, now, details)

	assert := assert.New(t)
	assert.NoError(err)
	assert.Equal(experimentation.ExperimentRunDetails_RUNNING, runDetails.Status)
}

func TestStoppedExperiment(t *testing.T) {
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

	runDetails, err := NewRunDetails(uint64(1), creationTime, start, end, cancellation, time.Now(), details)

	assert := assert.New(t)
	assert.NoError(err)
	assert.Equal(experimentation.ExperimentRunDetails_STOPPED, runDetails.Status)
}

func TestCompletedExperiment(t *testing.T) {
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

	runDetails, err := NewRunDetails(uint64(1), creationTime, start, end, cancellation, now, details)

	assert := assert.New(t)
	assert.NoError(err)
	assert.Equal(experimentation.ExperimentRunDetails_COMPLETED, runDetails.Status)
}

package experimentstore

import (
	"database/sql"
	"testing"
	"time"

	"github.com/golang/protobuf/ptypes/any"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap/zaptest"

	experimentation "github.com/lyft/clutch/backend/api/chaos/experimentation/v1"
)

var startTime = time.Date(2020, 11, 17, 20, 34, 58, 651387237, time.UTC)
var creationTime = startTime
var endTime = time.Date(2020, 11, 20, 20, 34, 58, 651387237, time.UTC)

func TestScheduledExperiment(t *testing.T) {
	logger := zaptest.NewLogger(t).Sugar()

	end := sql.NullTime{Valid: false}
	cancellation := sql.NullTime{Valid: false}
	now := startTime.AddDate(0, 0, -1)

	run := &ExperimentRun{id: 1, startTime: startTime, endTime: end, cancellationTime: cancellation, creationTime: creationTime}
	config := &ExperimentConfig{id: 2, Config: &any.Any{}}
	transformer := NewTransformer(logger)

	runDetails, err := NewRunDetails(run, config, &transformer, now)

	assert := assert.New(t)
	assert.NoError(err)
	assert.Equal(experimentation.Experiment_SCHEDULED, runDetails.Status)
}

func TestCanceledExperiment(t *testing.T) {
	logger := zaptest.NewLogger(t).Sugar()

	end := sql.NullTime{Valid: false}
	cancellation := sql.NullTime{
		Time:  startTime.AddDate(0, 0, -1),
		Valid: true,
	}

	run := &ExperimentRun{id: 1, startTime: startTime, endTime: end, cancellationTime: cancellation, creationTime: creationTime}
	config := &ExperimentConfig{id: 2, Config: &any.Any{}}
	transformer := NewTransformer(logger)

	runDetails, err := NewRunDetails(run, config, &transformer, time.Now())

	assert := assert.New(t)
	assert.NoError(err)
	assert.Equal(experimentation.Experiment_CANCELED, runDetails.Status)
}

func TestRunningExperiment(t *testing.T) {
	logger := zaptest.NewLogger(t).Sugar()

	end := sql.NullTime{Valid: false}
	cancellation := sql.NullTime{Valid: false}
	now := startTime.AddDate(0, 0, 1)

	run := &ExperimentRun{id: 1, startTime: startTime, endTime: end, cancellationTime: cancellation, creationTime: creationTime}
	config := &ExperimentConfig{id: 2, Config: &any.Any{}}
	transformer := NewTransformer(logger)

	runDetails, err := NewRunDetails(run, config, &transformer, now)

	assert := assert.New(t)
	assert.NoError(err)
	assert.Equal(experimentation.Experiment_RUNNING, runDetails.Status)
}

func TestStoppedExperiment(t *testing.T) {
	logger := zaptest.NewLogger(t).Sugar()

	end := sql.NullTime{
		Time:  endTime,
		Valid: true,
	}
	cancellation := sql.NullTime{
		Time:  startTime.AddDate(0, 0, 1),
		Valid: true,
	}

	run := &ExperimentRun{id: 1, startTime: startTime, endTime: end, cancellationTime: cancellation, creationTime: creationTime}
	config := &ExperimentConfig{id: 2, Config: &any.Any{}}
	transformer := NewTransformer(logger)

	runDetails, err := NewRunDetails(run, config, &transformer, time.Now())

	assert := assert.New(t)
	assert.NoError(err)
	assert.Equal(experimentation.Experiment_STOPPED, runDetails.Status)
}

func TestCompletedExperiment(t *testing.T) {
	logger := zaptest.NewLogger(t).Sugar()

	end := sql.NullTime{
		Time:  endTime,
		Valid: true,
	}
	cancellation := sql.NullTime{
		Valid: false,
	}
	now := endTime.AddDate(0, 0, 1)

	run := &ExperimentRun{id: 1, startTime: startTime, endTime: end, cancellationTime: cancellation, creationTime: creationTime}
	config := &ExperimentConfig{id: 2, Config: &any.Any{}}
	transformer := NewTransformer(logger)

	runDetails, err := NewRunDetails(run, config, &transformer, now)

	assert := assert.New(t)
	assert.NoError(err)
	assert.Equal(experimentation.Experiment_COMPLETED, runDetails.Status)
}

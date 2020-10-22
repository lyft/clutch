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

func TestListViewRunningExperimentRunProperties(t *testing.T) {
	startTime := time.Now()
	expectedRun := &ExperimentRun{id: 1, startTime: startTime, endTime: sql.NullTime{}, cancellationTime: sql.NullTime{}, creationTime: creationTime}
	expectedConfig := &ExperimentConfig{id: 2, Config: &any.Any{TypeUrl: "foo"}}

	expectedProperty := &experimentation.Property{
		Id:    "foo",
		Label: "bar",
		Value: &experimentation.Property_StringValue{StringValue: "dar"},
	}

	logger := zaptest.NewLogger(t).Sugar()
	transformer := NewTransformer(logger)
	transform := func(run *ExperimentRun, config *ExperimentConfig) ([]*experimentation.Property, error) {
		assert.Equal(t, expectedRun, run)
		assert.Equal(t, expectedConfig, config)
		return []*experimentation.Property{expectedProperty}, nil
	}

	transformation := Transformation{ConfigTypeUrl: "foo", RunTransform: transform}
	assert.NoError(t, transformer.Register(transformation))
	listView, err := NewRunListView(expectedRun, expectedConfig, &transformer, time.Now())

	assert.NoError(t, err)

	assert.Equal(t, uint64(1), listView.Id)
	assert.Equal(t, int64(1), listView.GetProperties().GetItems()["run_identifier"].GetIntValue())
	assert.Equal(t, int64(2), listView.GetProperties().GetItems()["config_identifier"].GetIntValue())
	assert.Equal(t, expectedProperty, listView.GetProperties().GetItems()["foo"])
}

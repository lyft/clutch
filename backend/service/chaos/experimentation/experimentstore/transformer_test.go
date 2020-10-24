package experimentstore

import (
	"testing"

	"github.com/golang/protobuf/ptypes/any"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap/zaptest"

	experimentation "github.com/lyft/clutch/backend/api/chaos/experimentation/v1"
)

func TestNoRegisteredTransformation(t *testing.T) {
	logger := zaptest.NewLogger(t).Sugar()
	transformer := NewTransformer(logger)
	config := &ExperimentConfig{id: 1, Config: &any.Any{}}
	_, err := transformer.CreateProperties(&ExperimentRun{Id: 123}, config)

	assert.NoError(t, err)
}

func TestNoMatchingRegisteredRunConfigTransform(t *testing.T) {
	logger := zaptest.NewLogger(t).Sugar()
	transformer := NewTransformer(logger)

	underlyingConfig := &any.Any{TypeUrl: "foo"}
	config := &ExperimentConfig{id: 1, Config: underlyingConfig}

	transform := func(run *ExperimentRun, config *ExperimentConfig) ([]*experimentation.Property, error) {
		assert.FailNow(t, "not matching transform should not be called")
		return []*experimentation.Property{}, nil
	}
	transformation := Transformation{ConfigTypeUrl: "bar", RunTransform: transform}
	assert.NoError(t, transformer.Register(transformation))
	properties, err := transformer.CreateProperties(&ExperimentRun{Id: 123}, config)

	assert.NoError(t, err)
	assert.Equal(t, 0, len(properties))
}

func TestMatchingRegisteredNullRunConfigTransform(t *testing.T) {
	logger := zaptest.NewLogger(t).Sugar()

	run := &ExperimentRun{Id: 123}
	config := &ExperimentConfig{id: 1, Config: &any.Any{TypeUrl: "test"}}

	transformation := Transformation{ConfigTypeUrl: "test"}
	transformer := NewTransformer(logger)
	assert.NoError(t, transformer.Register(transformation))

	properties, err := transformer.CreateProperties(run, config)
	assert.NoError(t, err)
	assert.Equal(t, 0, len(properties))
}

func TestMatchingRegisteredRunConfigTransform(t *testing.T) {
	logger := zaptest.NewLogger(t).Sugar()

	expectedRun := &ExperimentRun{Id: 123}
	expectedConfig := &ExperimentConfig{id: 1, Config: &any.Any{TypeUrl: "test"}}
	expectedProperty := &experimentation.Property{
		Id:    "foo",
		Label: "bar",
		Value: &experimentation.Property_StringValue{StringValue: "dar"},
	}

	transform := func(run *ExperimentRun, config *ExperimentConfig) ([]*experimentation.Property, error) {
		assert.Equal(t, expectedRun, run)
		assert.Equal(t, expectedConfig, config)
		return []*experimentation.Property{expectedProperty}, nil
	}
	transformation := Transformation{ConfigTypeUrl: "test", RunTransform: transform}
	transformer := NewTransformer(logger)
	assert.NoError(t, transformer.Register(transformation))

	properties, err := transformer.CreateProperties(expectedRun, expectedConfig)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(properties))
	assert.Equal(t, expectedProperty, properties[0])
}

func TestMatchingMultipleRegisteredRunConfigTransforms(t *testing.T) {
	logger := zaptest.NewLogger(t).Sugar()

	expectedRun := &ExperimentRun{Id: 123}
	expectedConfig := &ExperimentConfig{id: 1, Config: &any.Any{TypeUrl: "foo"}}
	expectedProperty1 := &experimentation.Property{
		Id:    "foo1",
		Label: "bar1",
		Value: &experimentation.Property_StringValue{StringValue: "dar1"},
	}

	expectedProperty2 := &experimentation.Property{
		Id:    "foo2",
		Label: "bar2",
		Value: &experimentation.Property_StringValue{StringValue: "dar2"},
	}

	transform1 := func(run *ExperimentRun, config *ExperimentConfig) ([]*experimentation.Property, error) {
		assert.Equal(t, expectedRun, run)
		assert.Equal(t, expectedConfig, config)
		return []*experimentation.Property{expectedProperty1}, nil
	}
	transform2 := func(run *ExperimentRun, config *ExperimentConfig) ([]*experimentation.Property, error) {
		assert.Equal(t, expectedRun, run)
		assert.Equal(t, expectedConfig, config)
		return []*experimentation.Property{expectedProperty2}, nil
	}

	transformation1 := Transformation{ConfigTypeUrl: "foo", RunTransform: transform1}
	transformation2 := Transformation{ConfigTypeUrl: "foo", RunTransform: transform2}
	transformer := NewTransformer(logger)
	assert.NoError(t, transformer.Register(transformation1))
	assert.NoError(t, transformer.Register(transformation2))

	properties, err := transformer.CreateProperties(expectedRun, expectedConfig)
	assert.NoError(t, err)
	assert.Equal(t, 2, len(properties))
	assert.Equal(t, []*experimentation.Property{expectedProperty1, expectedProperty2}, properties)
}

package experimentstore

import (
	"testing"

	"github.com/golang/protobuf/ptypes/any"
	"github.com/stretchr/testify/assert"

	experimentation "github.com/lyft/clutch/backend/api/chaos/experimentation/v1"
)

func TestNoRegisteredTransform(t *testing.T) {
	transformer := NewTransformer()
	config := &ExperimentConfig{id: 1, Config: &any.Any{}}
	_, err := transformer.CreateProperties(config)

	assert := assert.New(t)
	assert.NoError(err)
}

func TestRegisteredTransform(t *testing.T) {
	assert := assert.New(t)

	underlyingConfig := &any.Any{TypeUrl: "test"}
	config := &ExperimentConfig{id: 1, Config: underlyingConfig}

	expectedProperty := &experimentation.Property{
		Id:    "foo",
		Label: "bar",
		Value: &experimentation.Property_StringValue{StringValue: "dar"},
	}

	transform := func(config *ExperimentConfig) ([]*experimentation.Property, error) {
		assert.Equal(underlyingConfig, config.Config)
		return []*experimentation.Property{expectedProperty}, nil
	}
	transformation := Transformation{ConfigTypeUrl: "test", ConfigTransform: transform}
	transformer := NewTransformer()
	assert.NoError(transformer.Register(transformation))

	properties, err := transformer.CreateProperties(config)
	assert.NoError(err)
	assert.Equal(1, len(properties))
	assert.Equal(expectedProperty, properties[0])
}

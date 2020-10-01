package experimentstore

import (
	"testing"

	"github.com/golang/protobuf/ptypes/any"
	"github.com/stretchr/testify/assert"

	experimentation "github.com/lyft/clutch/backend/api/chaos/experimentation/v1"
)

func TestNoRegisteredTransform(t *testing.T) {
	transform := Transformer{nameToConfigTransformMap: map[string]func(*ExperimentConfig) ([]*experimentation.Property, error){}}
	config := &ExperimentConfig{id: 1, Config: &any.Any{}}
	_, err := transform.CreateProperties(config)

	assert := assert.New(t)
	assert.NoError(err)
}

func TestRegisteredTransform(t *testing.T) {
	assert := assert.New(t)

	underlyingConfig := &any.Any{TypeUrl: "test"}
	config := &ExperimentConfig{id: 1, Config: underlyingConfig}
	transformer := Transformer{nameToConfigTransformMap: map[string]func(*ExperimentConfig) ([]*experimentation.Property, error){}}

	expectedProperty := &experimentation.Property{
		Id:    "foo",
		Label: "bar",
		Value: &experimentation.Property_StringValue{StringValue: "dar"},
	}
	transformer.nameToConfigTransformMap["test"] = func(config *ExperimentConfig) ([]*experimentation.Property, error) {
		assert.Equal(underlyingConfig, config.Config)
		return []*experimentation.Property{expectedProperty}, nil
	}

	properties, err := transformer.CreateProperties(config)
	assert.NoError(err)
	assert.Equal(1, len(properties))
	assert.Equal(expectedProperty, properties[0])
}

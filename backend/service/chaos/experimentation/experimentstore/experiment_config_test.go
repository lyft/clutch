package experimentstore

import (
	"testing"

	"github.com/golang/protobuf/ptypes/any"
	"github.com/stretchr/testify/assert"

	experimentation "github.com/lyft/clutch/backend/api/chaos/experimentation/v1"
)

func TestConfigPropertiesWithNoRegisteredTransform(t *testing.T) {
	transformer := NewTransformer()
	config := &ExperimentConfig{id: 1, Config: &any.Any{}}

	properties, err := config.CreateProperties(&transformer)

	assert := assert.New(t)
	assert.NoError(err)
	assert.Equal(1, len(properties))
	assert.Equal("config_identifier", properties[0].Id)
	assert.Equal(int64(1), properties[0].GetIntValue())
}

func TestConfigPropertiesWithRegisteredTransform(t *testing.T) {
	assert := assert.New(t)

	transform := func(config *ExperimentConfig) ([]*experimentation.Property, error) {
		return []*experimentation.Property{
			{
				Id:    "foo",
				Label: "bar",
				Value: &experimentation.Property_StringValue{StringValue: "dar"},
			},
		}, nil
	}
	transformation := Transformation{ConfigTypeUrl: "type", ConfigTransform: transform}
	transformer := NewTransformer()
	assert.NoError(transformer.Register(transformation))

	config := &ExperimentConfig{id: 1, Config: &any.Any{TypeUrl: "type"}}
	properties, err := config.CreateProperties(&transformer)

	assert.NoError(err)
	assert.Equal(2, len(properties))
	assert.Equal("config_identifier", properties[0].Id)
	assert.Equal(int64(1), properties[0].GetIntValue())
	assert.Equal("foo", properties[1].Id)
	assert.Equal("dar", properties[1].GetStringValue())
}

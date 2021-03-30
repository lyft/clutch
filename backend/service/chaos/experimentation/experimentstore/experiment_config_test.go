package experimentstore

import (
	"testing"

	"github.com/golang/protobuf/ptypes/any"
	"github.com/stretchr/testify/assert"
)

func TestConfigProperties(t *testing.T) {
	config := &ExperimentConfig{id: "1", Config: &any.Any{}}

	properties, err := config.CreateProperties()

	assert := assert.New(t)
	assert.NoError(err)
	assert.Equal(1, len(properties))
	assert.Equal("config_identifier", properties[0].Id)
	assert.Equal("1", properties[0].GetStringValue())
}

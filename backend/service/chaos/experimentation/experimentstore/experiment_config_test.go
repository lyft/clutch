package experimentstore

import (
	"testing"

	"github.com/golang/protobuf/ptypes/any"
	"github.com/stretchr/testify/assert"
)

func TestConfigProperties(t *testing.T) {
	config := &ExperimentConfig{Id: "1", Config: &any.Any{}}

	properties, err := config.CreateProperties()

	assert.NoError(t, err)
	assert.Len(t, properties, 1)
	assert.Equal(t, "config_identifier", properties[0].Id)
	assert.Equal(t, "1", properties[0].GetStringValue())
}

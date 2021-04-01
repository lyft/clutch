package experimentstore

import (
	"github.com/golang/protobuf/ptypes/any"

	experimentation "github.com/lyft/clutch/backend/api/chaos/experimentation/v1"
)

type ExperimentConfig struct {
	id     string
	Config *any.Any
}

func (ec *ExperimentConfig) CreateProperties() ([]*experimentation.Property, error) {
	return []*experimentation.Property{
		{
			Id:    "config_identifier",
			Label: "Config Identifier",
			Value: &experimentation.Property_StringValue{StringValue: ec.id},
		},
	}, nil
}

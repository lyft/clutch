package experimentstore

import (
	"github.com/golang/protobuf/ptypes/any"

	experimentationv1 "github.com/lyft/clutch/backend/api/chaos/experimentation/v1"
)

type ExperimentConfig struct {
	Id     string
	Config *any.Any
}

func (ec *ExperimentConfig) CreateProperties() ([]*experimentationv1.Property, error) {
	return []*experimentationv1.Property{
		{
			Id:    "config_identifier",
			Label: "Config Identifier",
			Value: &experimentationv1.Property_StringValue{StringValue: ec.Id},
		},
	}, nil
}

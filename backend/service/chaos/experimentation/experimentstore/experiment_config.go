package experimentstore

import (
	"google.golang.org/protobuf/types/known/anypb"

	experimentation "github.com/lyft/clutch/backend/api/chaos/experimentation/v1"
)

type ExperimentConfig struct {
	id     uint64
	Config *anypb.Any
}

func (ec *ExperimentConfig) CreateProperties() ([]*experimentation.Property, error) {
	return []*experimentation.Property{
		{
			Id:    "config_identifier",
			Label: "Config Identifier",
			Value: &experimentation.Property_IntValue{IntValue: int64(ec.id)},
		},
	}, nil
}

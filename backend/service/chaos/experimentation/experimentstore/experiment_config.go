package experimentstore

import (
	"github.com/golang/protobuf/ptypes/any"

	experimentation "github.com/lyft/clutch/backend/api/chaos/experimentation/v1"
)

type ExperimentConfig struct {
	id     uint64
	Config *any.Any
}

func (ec *ExperimentConfig) CreateProperties(transformer *Transformer) ([]*experimentation.Property, error) {
	properties, err := transformer.CreateProperties(ec)
	if err != nil {
		return nil, err
	}

	idProperty := []*experimentation.Property{
		{
			Id:    "config_identifier",
			Label: "Config Identifier",
			Value: &experimentation.Property_IntValue{IntValue: int64(ec.id)},
		},
	}

	return append(idProperty, properties...), nil
}

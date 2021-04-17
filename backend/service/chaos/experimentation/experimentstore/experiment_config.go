package experimentstore

import (
	"github.com/golang/protobuf/ptypes/any"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"

	experimentationv1 "github.com/lyft/clutch/backend/api/chaos/experimentation/v1"
)

type ExperimentConfig struct {
	Id string
	// TODO(Augustyniak): Remove Config property once all of its existing usages are removed. Use Message property instead.
	Config  *any.Any
	Message proto.Message
}

func NewExperimentConfig(id string, stringifedData string) (*ExperimentConfig, error) {
	data := &any.Any{}
	if err := protojson.Unmarshal([]byte(stringifedData), data); err != nil {
		return nil, err
	}

	message, err := anypb.UnmarshalNew(data, proto.UnmarshalOptions{})
	if err != nil {
		return nil, err
	}

	return &ExperimentConfig{
		Id:      id,
		Config:  data,
		Message: message,
	}, nil
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

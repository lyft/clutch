package serverexperimentation

import (
	"errors"
	"fmt"

	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/any"
	"github.com/uber-go/tally"
	"go.uber.org/zap"

	experimentation "github.com/lyft/clutch/backend/api/chaos/experimentation/v1"
	serverexperimentation "github.com/lyft/clutch/backend/api/chaos/serverexperimentation/v1"
	"github.com/lyft/clutch/backend/module"
	"github.com/lyft/clutch/backend/service"
	"github.com/lyft/clutch/backend/service/chaos/experimentation/experimentstore"
)

const (
	Name = "clutch.module.chaos.serverexperimentation"
)

type Service struct {
	experimentStore experimentstore.ExperimentStore
}

// New instantiates a Service object.
func New(_ *any.Any, logger *zap.Logger, scope tally.Scope) (module.Module, error) {
	store, ok := service.Registry[experimentstore.Name]
	if !ok {
		return nil, errors.New("could not find experiment store service")
	}

	experimentStore, ok := store.(experimentstore.ExperimentStore)
	if !ok {
		return nil, errors.New("service was not the correct type")
	}

	return &Service{
		experimentStore: experimentStore,
	}, nil
}

func (s *Service) Register(r module.Registrar) error {
	s.experimentStore.RegisterTransformation("type.googleapis.com/clutch.chaos.serverexperimentation.v1.TestConfig", transform)
	return nil
}

func transform(config *experimentstore.ExperimentConfig) ([]*experimentation.Property, error) {
	var testConfig = serverexperimentation.TestConfig{}
	err := ptypes.UnmarshalAny(config.Config, &testConfig)
	if err != nil {
		return []*experimentation.Property{}, nil
	}

	upstream := testConfig.GetClusterPair().GetUpstreamCluster()
	downstream := testConfig.GetClusterPair().GetDownstreamCluster()

	var description string
	switch testConfig.GetFault().(type) {
	case *serverexperimentation.TestConfig_Abort:
		description = "Abort"
	case *serverexperimentation.TestConfig_Latency:
		description = "Latency"
	}

	return []*experimentation.Property{
		{
			Id:    "type",
			Label: "Type",
			Value: &experimentation.Property_StringValue{StringValue: "Server"},
		},
		{
			Id:    "target",
			Label: "Target",
			Value: &experimentation.Property_StringValue{StringValue: fmt.Sprintf("\"%s\" ➡️ \"%s\"", downstream, upstream)},
		},
		{
			Id:    "fault_types",
			Label: "Fault Types",
			Value: &experimentation.Property_StringValue{StringValue: description},
		},
	}, nil
}

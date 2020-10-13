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
	storer experimentstore.Storer
	logger *zap.SugaredLogger
}

// New instantiates a Service object.
func New(_ *any.Any, logger *zap.Logger, scope tally.Scope) (module.Module, error) {
	store, ok := service.Registry[experimentstore.Name]
	if !ok {
		return nil, errors.New("could not find experiment store service")
	}

	storer, ok := store.(experimentstore.Storer)
	if !ok {
		return nil, errors.New("service was not the correct type")
	}

	return &Service{
		storer: storer,
		logger: logger.Sugar(),
	}, nil
}

func (s *Service) Register(r module.Registrar) error {
	transformation := experimentstore.Transformation{ConfigTypeUrl: "type.googleapis.com/clutch.chaos.serverexperimentation.v1.TestConfig", ConfigTransform: s.transform}
	return s.storer.RegisterTransformation(transformation)
}

func (s *Service) transform(config *experimentstore.ExperimentConfig) ([]*experimentation.Property, error) {
	var experimentConfig = serverexperimentation.TestConfig{}
	if err := ptypes.UnmarshalAny(config.Config, &experimentConfig); err != nil {
		return []*experimentation.Property{}, err
	}

	upstream := experimentConfig.GetClusterPair().GetUpstreamCluster()
	downstream := experimentConfig.GetClusterPair().GetDownstreamCluster()

	faultsDescription, err := experimentConfigToString(&experimentConfig)
	if err != nil {
		return nil, err
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
			Value: &experimentation.Property_StringValue{StringValue: fmt.Sprintf("%s ➡️ %s", downstream, upstream)},
		},
		{
			Id:    "fault_types",
			Label: "Fault Types",
			Value: &experimentation.Property_StringValue{StringValue: faultsDescription},
		},
	}, nil
}

func experimentConfigToString(experiment *serverexperimentation.TestConfig) (string, error) {
	if experiment == nil {
		return "", errors.New("experiment is nil")
	}

	var description string
	switch experiment.GetFault().(type) {
	case *serverexperimentation.TestConfig_Abort:
		description = "Abort"
	case *serverexperimentation.TestConfig_Latency:
		description = "Latency"
	default:
		return nil, fmt.Errorf("unexpected fault type %v", fault)
	}

	return description
}

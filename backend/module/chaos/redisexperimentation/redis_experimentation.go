package redisexperimentation

import (
	"errors"
	"fmt"
	"strings"

	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/any"
	"github.com/uber-go/tally"
	"go.uber.org/zap"

	experimentation "github.com/lyft/clutch/backend/api/chaos/experimentation/v1"
	redisexperimentation "github.com/lyft/clutch/backend/api/chaos/redisexperimentation/v1"
	"github.com/lyft/clutch/backend/module"
	"github.com/lyft/clutch/backend/service"
	"github.com/lyft/clutch/backend/service/chaos/experimentation/experimentstore"
)

const (
	Name = "clutch.module.chaos.redisexperimentation"
)

type Service struct {
	storer experimentstore.Storer
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
	}, nil
}

func (s *Service) Register(r module.Registrar) error {
	transformation := experimentstore.Transformation{ConfigTypeUrl: "type.googleapis.com/clutch.chaos.redisexperimentation.v1.FaultConfig", RunTransform: s.transform}
	return s.storer.RegisterTransformation(transformation)
}

func (s *Service) transform(_ *experimentstore.ExperimentRun, config *experimentstore.ExperimentConfig) ([]*experimentation.Property, error) {
	var experimentConfig = redisexperimentation.FaultConfig{}
	if err := ptypes.UnmarshalAny(config.Config, &experimentConfig); err != nil {
		return []*experimentation.Property{}, err
	}

	faultsDescription, err := experimentConfigToFaultString(&experimentConfig)
	if err != nil {
		return nil, err
	}

	var downstream, upstream string
	downstream = experimentConfig.GetFaultTargeting().GetDownstreamCluster().GetName()
	upstream = experimentConfig.GetFaultTargeting().GetUpstreamCluster().GetName() +
		strings.Join(experimentConfig.GetFaultTargeting().GetRedisCommands(), ",")

	return []*experimentation.Property{
		{
			Id:    "type",
			Label: "Type",
			Value: &experimentation.Property_StringValue{StringValue: "Redis"},
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

func experimentConfigToFaultString(experiment *redisexperimentation.FaultConfig) (string, error) {
	if experiment == nil {
		return "", errors.New("experiment is nil")
	}

	switch experiment.GetFault().(type) {
	case *redisexperimentation.FaultConfig_ErrorFault:
		return "Error", nil
	case *redisexperimentation.FaultConfig_LatencyFault:
		return "Delay", nil
	default:
		return "", fmt.Errorf("unexpected fault type %v", experiment.GetFault())
	}
}

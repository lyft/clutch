package serverexperimentation

import (
	"errors"
	"fmt"
	"github.com/lyft/clutch/backend/module/chaos/serverexperimentation/xds"

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
	transformation := experimentstore.Transformation{
		ConfigTypeUrl: "type.googleapis.com/clutch.chaos.serverexperimentation.v1.HTTPFaultConfig",
		RunTransform: s.transform,
	}
	runtimeGeneration := experimentstore.RuntimeGeneration{
		ConfigTypeUrl: "type.googleapis.com/clutch.chaos.serverexperimentation.v1.HTTPFaultConfig",
		RuntimeKeysGeneration: xds.RuntimeKeysGeneration,
		GetEnforcingCluster: xds.GetEnforcingCluster,
	}
	err := s.storer.RegisterTransformation(transformation)
	if err != nil {
		return err
	}
	return s.storer.RegisterRuntimeGeneration(runtimeGeneration)
}

func (s *Service) transform(_ *experimentstore.ExperimentRun, config *experimentstore.ExperimentConfig) ([]*experimentation.Property, error) {
	var experimentConfig = serverexperimentation.HTTPFaultConfig{}
	if err := config.Config.UnmarshalTo(&experimentConfig); err != nil {
		return []*experimentation.Property{}, err
	}

	faultsDescription, err := experimentConfigToString(&experimentConfig)
	if err != nil {
		return nil, err
	}

	var downstream, upstream string
	switch experimentConfig.GetFaultTargeting().GetEnforcer().(type) {
	case *serverexperimentation.FaultTargeting_DownstreamEnforcing:
		downstreamEnforcing := experimentConfig.GetFaultTargeting().GetDownstreamEnforcing()

		switch downstreamEnforcing.GetDownstreamType().(type) {
		case *serverexperimentation.DownstreamEnforcing_DownstreamCluster:
			downstream = downstreamEnforcing.GetDownstreamCluster().GetName()
		default:
			return nil, fmt.Errorf("unknown downstream type of downstream enforcing %v", downstreamEnforcing.GetDownstreamType())
		}

		switch downstreamEnforcing.GetUpstreamType().(type) {
		case *serverexperimentation.DownstreamEnforcing_UpstreamCluster:
			upstream = downstreamEnforcing.GetUpstreamCluster().GetName()
		default:
			return nil, fmt.Errorf("unknown upstream type of downstream enforcing %v", downstreamEnforcing.GetUpstreamType())
		}

	case *serverexperimentation.FaultTargeting_UpstreamEnforcing:
		upstreamEnforcing := experimentConfig.GetFaultTargeting().GetUpstreamEnforcing()

		switch upstreamEnforcing.GetDownstreamType().(type) {
		case *serverexperimentation.UpstreamEnforcing_DownstreamCluster:
			downstream = upstreamEnforcing.GetDownstreamCluster().GetName()
		default:
			return nil, fmt.Errorf("unknown downstream type of upstream enforcing %v", upstreamEnforcing.GetDownstreamType())
		}

		switch upstreamEnforcing.GetUpstreamType().(type) {
		case *serverexperimentation.UpstreamEnforcing_UpstreamCluster:
			upstream = upstreamEnforcing.GetUpstreamCluster().GetName()
		case *serverexperimentation.UpstreamEnforcing_UpstreamPartialSingleCluster:
			upstream = upstreamEnforcing.GetUpstreamPartialSingleCluster().GetName()
		default:
			return nil, fmt.Errorf("unknown upstream type of upstream enforcing %v", upstreamEnforcing.GetUpstreamType())
		}

	default:
		return nil, fmt.Errorf("unknown enforcer %v", experimentConfig.GetFaultTargeting())
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

func experimentConfigToString(experiment *serverexperimentation.HTTPFaultConfig) (string, error) {
	if experiment == nil {
		return "", errors.New("experiment is nil")
	}

	switch experiment.GetFault().(type) {
	case *serverexperimentation.HTTPFaultConfig_AbortFault:
		return "Abort", nil
	case *serverexperimentation.HTTPFaultConfig_LatencyFault:
		return "Latency", nil
	default:
		return "", fmt.Errorf("unexpected fault type %v", experiment.GetFault())
	}
}

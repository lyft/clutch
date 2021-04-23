package serverexperimentation

import (
	"errors"
	"fmt"
	"github.com/golang/protobuf/ptypes/any"
	"github.com/uber-go/tally"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"

	experimentationv1 "github.com/lyft/clutch/backend/api/chaos/experimentation/v1"
	serverexperimentationv1 "github.com/lyft/clutch/backend/api/chaos/serverexperimentation/v1"
	configv1 "github.com/lyft/clutch/backend/api/config/module/chaos/serverexperimentation/v1"
	"github.com/lyft/clutch/backend/module"
	"github.com/lyft/clutch/backend/module/chaos/experimentation/xds"
	serverexperimentationxds "github.com/lyft/clutch/backend/module/chaos/serverexperimentation/xds"
	"github.com/lyft/clutch/backend/service"
	"github.com/lyft/clutch/backend/service/chaos/experimentation/experimentstore"
)

const (
	Name = "clutch.module.chaos.serverexperimentation"
)

type Service struct {
	storer experimentstore.Storer
}

func TypeUrl(message proto.Message) string {
	return "type.googleapis.com/" + string(message.ProtoReflect().Descriptor().FullName())
}

// New instantiates a Service object.
func New(untypedConfig *any.Any, logger *zap.Logger, scope tally.Scope) (module.Module, error) {
	config := &configv1.Config{}
	if err := untypedConfig.UnmarshalTo(config); err != nil {
		return nil, err
	}

	store, ok := service.Registry[experimentstore.Name]
	if !ok {
		return nil, errors.New("could not find experiment store service")
	}

	storer, ok := store.(experimentstore.Storer)
	if !ok {
		return nil, errors.New("service was not the correct type")
	}

	xds.RTDSGeneratorsByTypeUrl[TypeUrl(&serverexperimentationv1.HTTPFaultConfig{})] = &serverexperimentationxds.RTDSServerFaultsResourceGenerator{
		IngressFaultRuntimePrefix: config.IngressFaultRuntimePrefix,
		EgressFaultRuntimePrefix:  config.EgressFaultRuntimePrefix,
	}

	xds.ECDSGeneratorsByTypeUrl[TypeUrl(&serverexperimentationv1.HTTPFaultConfig{})] = &serverexperimentationxds.ECDSServerFaultsResourceGenerator{}

	return &Service{
		storer: storer,
	}, nil
}

func (s *Service) Register(r module.Registrar) error {
	transformation := experimentstore.Transformation{ConfigTypeUrl: "type.googleapis.com/clutch.chaos.serverexperimentation.v1.HTTPFaultConfig", RunTransform: s.transform}
	return s.storer.RegisterTransformation(transformation)
}

func (s *Service) transform(_ *experimentstore.ExperimentRun, config *experimentstore.ExperimentConfig) ([]*experimentationv1.Property, error) {
	var experimentConfig = serverexperimentationv1.HTTPFaultConfig{}
	if err := config.Config.UnmarshalTo(&experimentConfig); err != nil {
		return []*experimentationv1.Property{}, err
	}

	faultsDescription, err := experimentConfigToString(&experimentConfig)
	if err != nil {
		return nil, err
	}

	var downstream, upstream string
	switch experimentConfig.GetFaultTargeting().GetEnforcer().(type) {
	case *serverexperimentationv1.FaultTargeting_DownstreamEnforcing:
		downstreamEnforcing := experimentConfig.GetFaultTargeting().GetDownstreamEnforcing()

		switch downstreamEnforcing.GetDownstreamType().(type) {
		case *serverexperimentationv1.DownstreamEnforcing_DownstreamCluster:
			downstream = downstreamEnforcing.GetDownstreamCluster().GetName()
		default:
			return nil, fmt.Errorf("unknown downstream type of downstream enforcing %v", downstreamEnforcing.GetDownstreamType())
		}

		switch downstreamEnforcing.GetUpstreamType().(type) {
		case *serverexperimentationv1.DownstreamEnforcing_UpstreamCluster:
			upstream = downstreamEnforcing.GetUpstreamCluster().GetName()
		default:
			return nil, fmt.Errorf("unknown upstream type of downstream enforcing %v", downstreamEnforcing.GetUpstreamType())
		}

	case *serverexperimentationv1.FaultTargeting_UpstreamEnforcing:
		upstreamEnforcing := experimentConfig.GetFaultTargeting().GetUpstreamEnforcing()

		switch upstreamEnforcing.GetDownstreamType().(type) {
		case *serverexperimentationv1.UpstreamEnforcing_DownstreamCluster:
			downstream = upstreamEnforcing.GetDownstreamCluster().GetName()
		default:
			return nil, fmt.Errorf("unknown downstream type of upstream enforcing %v", upstreamEnforcing.GetDownstreamType())
		}

		switch upstreamEnforcing.GetUpstreamType().(type) {
		case *serverexperimentationv1.UpstreamEnforcing_UpstreamCluster:
			upstream = upstreamEnforcing.GetUpstreamCluster().GetName()
		case *serverexperimentationv1.UpstreamEnforcing_UpstreamPartialSingleCluster:
			upstream = upstreamEnforcing.GetUpstreamPartialSingleCluster().GetName()
		default:
			return nil, fmt.Errorf("unknown upstream type of upstream enforcing %v", upstreamEnforcing.GetUpstreamType())
		}

	default:
		return nil, fmt.Errorf("unknown enforcer %v", experimentConfig.GetFaultTargeting())
	}

	return []*experimentationv1.Property{
		{
			Id:    "type",
			Label: "Type",
			Value: &experimentationv1.Property_StringValue{StringValue: "Server"},
		},
		{
			Id:    "target",
			Label: "Target",
			Value: &experimentationv1.Property_StringValue{StringValue: fmt.Sprintf("%s ➡️ %s", downstream, upstream)},
		},
		{
			Id:    "fault_types",
			Label: "Fault Types",
			Value: &experimentationv1.Property_StringValue{StringValue: faultsDescription},
		},
	}, nil
}

func experimentConfigToString(experiment *serverexperimentationv1.HTTPFaultConfig) (string, error) {
	if experiment == nil {
		return "", errors.New("experiment is nil")
	}

	switch experiment.GetFault().(type) {
	case *serverexperimentationv1.HTTPFaultConfig_AbortFault:
		return "Abort", nil
	case *serverexperimentationv1.HTTPFaultConfig_LatencyFault:
		return "Latency", nil
	default:
		return "", fmt.Errorf("unexpected fault type %v", experiment.GetFault())
	}
}

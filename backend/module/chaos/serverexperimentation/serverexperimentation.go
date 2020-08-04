package serverexperimentation

import (
	"errors"
	"fmt"
	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/any"
	experimentation "github.com/lyft/clutch/backend/api/chaos/experimentation/v1"
	serverexperimentation "github.com/lyft/clutch/backend/api/chaos/serverexperimentation/v1"
	"github.com/lyft/clutch/backend/module"
	"github.com/lyft/clutch/backend/service"
	"github.com/lyft/clutch/backend/service/chaos/experimentation/experimentstore"
	"github.com/uber-go/tally"
	"go.uber.org/zap"
)

const (
	Name = "clutch.module.chaos.serverexperimentation"
)

type Service struct {}

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

	experimentStore.Register(
		"type.googleapis.com/clutch.chaos.serverexperimentation.v1.TestSpecification",
		convert)

	return &Service{}, nil
}

func (s *Service) Register(r module.Registrar) error {
	return nil
}

func convert(experiment *experimentation.Experiment) (*any.Any, error) {
	testSpecification := &serverexperimentation.TestSpecification{}
	ptypes.UnmarshalAny(experiment.GetTestConfig(), testSpecification)

	var clusterPair *serverexperimentation.ClusterPairTarget
	var description string

	switch testSpecification.GetConfig().(type) {
	case *serverexperimentation.TestSpecification_Abort:
		abort := testSpecification.GetAbort()
		clusterPair = abort.GetClusterPair()
		description = fmt.Sprintf("%d abort %.2f%%", abort.GetHttpStatus(), abort.GetPercent())
	case *serverexperimentation.TestSpecification_Latency:
		latency := testSpecification.GetLatency()
		clusterPair = latency.GetClusterPair()
		description = fmt.Sprintf("%dms delay %.2f%%", latency.GetDurationMs(), latency.GetPercent())
	default:
		break
	}

	viewModel := &serverexperimentation.ExperimentListViewModel{}
	viewModel.Identifier = string(experiment.GetId())
	viewModel.Clusters = fmt.Sprintf("%s cluster to %s cluster",
		clusterPair.GetUpstreamCluster(),
		clusterPair.GetDownstreamCluster())
	viewModel.Type = "Server-side Fault Injection Test"
	viewModel.Description = description
	return ptypes.MarshalAny(viewModel)
}
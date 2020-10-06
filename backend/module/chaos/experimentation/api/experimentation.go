package api

// <!-- START clutchdoc -->
// description: Experimentation Framework Service implementation. Supports a CRUD API for managing experiments.
// <!-- END clutchdoc -->

import (
	"errors"
	"time"

	"github.com/golang/protobuf/ptypes/any"
	"github.com/uber-go/tally"
	"go.uber.org/zap"
	"golang.org/x/net/context"

	experimentation "github.com/lyft/clutch/backend/api/chaos/experimentation/v1"
	"github.com/lyft/clutch/backend/module"
	"github.com/lyft/clutch/backend/service"
	"github.com/lyft/clutch/backend/service/chaos/experimentation/experimentstore"
)

const (
	Name = "clutch.module.chaos.experimentation.api"
)

// Service contains all dependencies for the API service.
type Service struct {
	experimentStore             experimentstore.ExperimentStore
	logger                      *zap.SugaredLogger
	createExperimentStat        tally.Counter
	cancelExperimentRunStat     tally.Counter
	getExperimentsStat          tally.Counter
	getListViewStat             tally.Counter
	getExperimentRunDetailsStat tally.Counter
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

	apiScope := scope.SubScope("experimentation")
	return &Service{
		experimentStore:             experimentStore,
		logger:                      logger.Sugar(),
		createExperimentStat:        apiScope.Counter("create_experiment"),
		cancelExperimentRunStat:     apiScope.Counter("cancel_experiment_run"),
		getExperimentsStat:          apiScope.Counter("get_experiments"),
		getListViewStat:             apiScope.Counter("get_list_view"),
		getExperimentRunDetailsStat: apiScope.Counter("get_experiment_run_config_pair_details"),
	}, nil
}

func (s *Service) Register(r module.Registrar) error {
	experimentation.RegisterExperimentsAPIServer(r.GRPCServer(), s)
	return r.RegisterJSONGateway(experimentation.RegisterExperimentsAPIHandler)
}

// CreateExperiments adds experiments to the experiment store.
func (s *Service) CreateExperiment(ctx context.Context, req *experimentation.CreateExperimentRequest) (*experimentation.CreateExperimentResponse, error) {
	s.createExperimentStat.Inc(1)

	// If start time is not provided, default to starting now
	now := time.Now()
	startTime := &now
	if req.StartTime != nil {
		s := req.StartTime.AsTime()
		startTime = &s
	}

	// If the end time is not provided, default to no end time
	var endTime *time.Time = nil
	if req.EndTime != nil {
		s := req.EndTime.AsTime()
		endTime = &s
	}

	experiment, err := s.experimentStore.CreateExperiment(ctx, req.Config, startTime, endTime)
	if err != nil {
		return nil, err
	}

	return &experimentation.CreateExperimentResponse{Experiment: experiment}, nil
}

// CancelExperimentRun cancels experiment that is currently running or is scheduled to be run in the future.
func (s *Service) CancelExperimentRun(ctx context.Context, req *experimentation.CancelExperimentRunRequest) (*experimentation.CancelExperimentRunResponse, error) {
	s.cancelExperimentRunStat.Inc(1)
	err := s.experimentStore.CancelExperimentRun(ctx, req.Id)
	if err != nil {
		return nil, err
	}

	return &experimentation.CancelExperimentRunResponse{}, nil
}

// GetExperiments returns all experiments from the experiment store.
func (s *Service) GetExperiments(ctx context.Context, request *experimentation.GetExperimentsRequest) (*experimentation.GetExperimentsResponse, error) {
	s.getExperimentsStat.Inc(1)
	experiments, err := s.experimentStore.GetExperiments(ctx, request.GetConfigType(), request.GetStatus())
	if err != nil {
		s.logger.Errorw("GetExperiments: Unable to retrieve experiments", "error", err)
		return &experimentation.GetExperimentsResponse{}, err
	}

	return &experimentation.GetExperimentsResponse{Experiments: experiments}, nil
}

func (s *Service) GetListView(ctx context.Context, _ *experimentation.GetListViewRequest) (*experimentation.GetListViewResponse, error) {
	s.getListViewStat.Inc(1)
	items, err := s.experimentStore.GetListView(ctx)
	if err != nil {
		return &experimentation.GetListViewResponse{}, err
	}

	return &experimentation.GetListViewResponse{Items: items}, nil
}

func (s *Service) GetExperimentRunDetails(ctx context.Context, request *experimentation.GetExperimentRunDetailsRequest) (*experimentation.GetExperimentRunDetailsResponse, error) {
	s.getExperimentRunDetailsStat.Inc(1)
	runDetails, err := s.experimentStore.GetExperimentRunDetails(ctx, request.Id)
	if err != nil {
		return &experimentation.GetExperimentRunDetailsResponse{}, err
	}

	return &experimentation.GetExperimentRunDetailsResponse{RunDetails: runDetails}, nil
}

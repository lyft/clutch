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
	storer                      experimentstore.Storer
	logger                      *zap.SugaredLogger
	createExperimentStat        tally.Counter
	createOrGetExperimentStats  tally.Scope
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

	experimentStore, ok := store.(experimentstore.Storer)
	if !ok {
		return nil, errors.New("could not find valid experiment store")
	}

	return &Service{
		storer:                      experimentStore,
		logger:                      logger.Sugar(),
		createExperimentStat:        scope.Counter("create_experiment"),
		createOrGetExperimentStats:  scope.SubScope("create_or_get_experiment"),
		cancelExperimentRunStat:     scope.Counter("cancel_experiment_run"),
		getExperimentsStat:          scope.Counter("get_experiments"),
		getListViewStat:             scope.Counter("get_list_view"),
		getExperimentRunDetailsStat: scope.Counter("get_experiment_run_config_pair_details"),
	}, nil
}

func (s *Service) Register(r module.Registrar) error {
	experimentation.RegisterExperimentsAPIServer(r.GRPCServer(), s)
	return r.RegisterJSONGateway(experimentation.RegisterExperimentsAPIHandler)
}

// CreateExperiments adds experiments to the experiment store.
func (s *Service) CreateExperiment(ctx context.Context, req *experimentation.CreateExperimentRequest) (*experimentation.CreateExperimentResponse, error) {
	s.createExperimentStat.Inc(1)

	es, err := experimentstore.NewExperimentSpecification(req.Data, time.Now())
	if err != nil {
		return nil, err
	}

	experiment, err := s.storer.CreateExperiment(ctx, es)
	if err != nil {
		return nil, err
	}

	return &experimentation.CreateExperimentResponse{Experiment: experiment}, nil
}

func (s *Service) CreateOrGetExperiment(ctx context.Context, req *experimentation.CreateOrGetExperimentRequest) (*experimentation.CreateOrGetExperimentResponse, error) {
	s.createOrGetExperimentStats.Counter("request").Inc(1)

	es, err := experimentstore.NewExperimentSpecification(req.Data, time.Now())
	if err != nil {
		return nil, err
	}
	createOrGetExperiment, err := s.storer.CreateOrGetExperiment(ctx, es)
	if err != nil {
		return nil, err
	}

	r := s.createOrGetExperimentStats.SubScope("result")
	switch createOrGetExperiment.Origin {
	case experimentation.CreateOrGetExperimentResponse_ORIGIN_UNSPECIFIED:
		r.Counter("unspecified").Inc(1)
	case experimentation.CreateOrGetExperimentResponse_ORIGIN_EXISTING:
		r.Counter("get").Inc(1)
	case experimentation.CreateOrGetExperimentResponse_ORIGIN_NEW:
		r.Counter("create").Inc(1)
	}

	return &experimentation.CreateOrGetExperimentResponse{
		Experiment: createOrGetExperiment.Experiment,
		Origin:     createOrGetExperiment.Origin,
	}, nil
}

// CancelExperimentRun cancels experiment that is currently running or is scheduled to be run in the future.
func (s *Service) CancelExperimentRun(ctx context.Context, req *experimentation.CancelExperimentRunRequest) (*experimentation.CancelExperimentRunResponse, error) {
	s.cancelExperimentRunStat.Inc(1)
	err := s.storer.CancelExperimentRun(ctx, req.Id, req.Reason)
	if err != nil {
		return nil, err
	}

	return &experimentation.CancelExperimentRunResponse{}, nil
}

// GetExperiments returns all experiments from the experiment store.
func (s *Service) GetExperiments(ctx context.Context, request *experimentation.GetExperimentsRequest) (*experimentation.GetExperimentsResponse, error) {
	s.getExperimentsStat.Inc(1)
	experiments, err := s.storer.GetExperiments(ctx, request.GetConfigType(), request.GetStatus())
	if err != nil {
		s.logger.Errorw("GetExperiments: Unable to retrieve experiments", "error", err)
		return &experimentation.GetExperimentsResponse{}, err
	}

	return &experimentation.GetExperimentsResponse{Experiments: experiments}, nil
}

func (s *Service) GetListView(ctx context.Context, _ *experimentation.GetListViewRequest) (*experimentation.GetListViewResponse, error) {
	s.getListViewStat.Inc(1)
	items, err := s.storer.GetListView(ctx)
	if err != nil {
		return &experimentation.GetListViewResponse{}, err
	}

	return &experimentation.GetListViewResponse{Items: items}, nil
}

func (s *Service) GetExperimentRunDetails(ctx context.Context, request *experimentation.GetExperimentRunDetailsRequest) (*experimentation.GetExperimentRunDetailsResponse, error) {
	s.getExperimentRunDetailsStat.Inc(1)
	runDetails, err := s.storer.GetExperimentRunDetails(ctx, request.Id)
	if err != nil {
		return &experimentation.GetExperimentRunDetailsResponse{}, err
	}

	return &experimentation.GetExperimentRunDetailsResponse{RunDetails: runDetails}, nil
}

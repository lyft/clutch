package experimentstoremock

import (
	"context"
	"strconv"
	"sync"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"

	experimentation "github.com/lyft/clutch/backend/api/chaos/experimentation/v1"
	"github.com/lyft/clutch/backend/service/chaos/experimentation/experimentstore"
)

type SimpleStorer struct {
	experiments []*experimentstore.Experiment
	idGenerator int

	sync.Mutex
}

func (s *SimpleStorer) CreateExperiment(ctx context.Context, es *experimentstore.ExperimentSpecification) (*experimentstore.Experiment, error) {
	s.Lock()
	defer s.Unlock()

	message, err := anypb.UnmarshalNew(es.Config, proto.UnmarshalOptions{})
	if err != nil {
		return nil, err
	}

	e := &experimentstore.Experiment{
		Run: &experimentstore.ExperimentRun{
			Id:        strconv.Itoa(s.idGenerator),
			StartTime: es.StartTime,
			EndTime:   es.EndTime,
		},
		Config: &experimentstore.ExperimentConfig{
			Id:      strconv.Itoa(s.idGenerator),
			Config:  es.Config,
			Message: message,
		},
	}
	s.experiments = append(s.experiments, e)
	s.idGenerator++
	return s.experiments[len(s.experiments)-1], nil
}

func (s *SimpleStorer) CreateOrGetExperiment(ctx context.Context, es *experimentstore.ExperimentSpecification) (*experimentstore.CreateOrGetExperimentResult, error) {
	experiment, err := s.CreateExperiment(ctx, es)
	if err != nil {
		return nil, err
	}

	return &experimentstore.CreateOrGetExperimentResult{
		Experiment: experiment,
		Origin:     experimentation.CreateOrGetExperimentResponse_ORIGIN_NEW,
	}, nil
}

func (s *SimpleStorer) CancelExperimentRun(ctx context.Context, runId string, reason string) error {
	s.Lock()
	defer s.Unlock()

	var newExperiments []*experimentstore.Experiment
	for _, e := range s.experiments {
		if e.Run.Id == runId {
			continue
		}

		newExperiments = append(newExperiments, e)
	}

	s.experiments = newExperiments
	return nil
}

func (s *SimpleStorer) GetExperiments(ctx context.Context, configType string, status experimentation.GetExperimentsRequest_Status) ([]*experimentstore.Experiment, error) {
	s.Lock()
	defer s.Unlock()

	var experiments []*experimentstore.Experiment
	experiments = append(experiments, s.experiments...)

	return experiments, nil
}
func (s *SimpleStorer) GetExperimentRunDetails(ctx context.Context, runId string) (*experimentation.ExperimentRunDetails, error) {
	return nil, nil
}
func (s *SimpleStorer) GetListView(ctx context.Context) ([]*experimentation.ListViewItem, error) {
	return nil, nil
}
func (s *SimpleStorer) RegisterTransformation(transformation experimentstore.Transformation) error {
	return nil
}

func (s *SimpleStorer) Close() {}

var _ experimentstore.Storer = &SimpleStorer{}

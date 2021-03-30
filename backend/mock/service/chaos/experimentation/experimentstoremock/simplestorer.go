package experimentstoremock

import (
	"context"
	"strconv"
	"sync"
	"time"

	"github.com/golang/protobuf/ptypes"
	"google.golang.org/protobuf/types/known/anypb"

	experimentation "github.com/lyft/clutch/backend/api/chaos/experimentation/v1"
	"github.com/lyft/clutch/backend/service/chaos/experimentation/experimentstore"
)

type SimpleExperiment struct {
	runId     string
	config    *anypb.Any
	startTime time.Time
	stopTime  *time.Time
}

func (s *SimpleExperiment) toProto() *experimentation.Experiment {
	startTime, _ := ptypes.TimestampProto(s.startTime)
	endTime, _ := ptypes.TimestampProto(*s.stopTime)
	return &experimentation.Experiment{
		RunId:     s.runId,
		Config:    s.config,
		StartTime: startTime,
		EndTime:   endTime,
	}
}

type SimpleStorer struct {
	experiments []SimpleExperiment
	idGenerator int

	sync.Mutex
}

func (s *SimpleStorer) CreateExperiment(ctx context.Context, es *experimentstore.ExperimentSpecification) (*experimentation.Experiment, error) {
	s.Lock()
	defer s.Unlock()

	s.experiments = append(s.experiments, SimpleExperiment{runId: strconv.Itoa(s.idGenerator), config: es.Config, startTime: es.StartTime, stopTime: es.EndTime})
	s.idGenerator++
	return s.experiments[len(s.experiments)-1].toProto(), nil
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

func (s *SimpleStorer) CancelExperimentRun(ctx context.Context, runId string) error {
	s.Lock()
	defer s.Unlock()

	newExperiments := []SimpleExperiment{}
	for _, e := range s.experiments {
		if e.runId == runId {
			continue
		}

		newExperiments = append(newExperiments, e)
	}

	s.experiments = newExperiments
	return nil
}

func (s *SimpleStorer) GetExperiments(ctx context.Context, configType string, status experimentation.GetExperimentsRequest_Status) ([]*experimentation.Experiment, error) {
	s.Lock()
	defer s.Unlock()

	experiments := []*experimentation.Experiment{}

	for _, e := range s.experiments {
		experiments = append(experiments, e.toProto())
	}

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

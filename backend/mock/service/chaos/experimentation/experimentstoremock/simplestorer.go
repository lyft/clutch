package experimentstoremock

import (
	"context"
	"github.com/golang/protobuf/ptypes"
	"github.com/lyft/clutch/backend/api/chaos/experimentation/v1"
	"github.com/lyft/clutch/backend/service/chaos/experimentation/experimentstore"
	"google.golang.org/protobuf/types/known/anypb"
	"sync"
	"time"
)

type SimpleExperiment struct {
	id        uint64
	config    *anypb.Any
	startTime *time.Time
	stopTime  *time.Time
}

func (s SimpleExperiment) toProto() *experimentationv1.Experiment {
	startTime, _ := ptypes.TimestampProto(*s.startTime)
	endTime, _ := ptypes.TimestampProto(*s.stopTime)
	return &experimentationv1.Experiment{
		Id:        s.id,
		Config:    s.config,
		StartTime: startTime,
		EndTime:   endTime,
	}
}

type SimpleStorer struct {
	experiments []SimpleExperiment
	idGenerator uint64

	sync.Mutex
}

func (s *SimpleStorer) CreateExperiment(ctx context.Context, config *anypb.Any, startTime *time.Time, stopTime *time.Time) (*experimentationv1.Experiment, error) {
	s.Lock()
	defer s.Unlock()

	s.experiments = append(s.experiments, SimpleExperiment{id: s.idGenerator, config: config, startTime: startTime, stopTime: stopTime})
	s.idGenerator++
	return s.experiments[len(s.experiments)-1].toProto(), nil
}

func (s *SimpleStorer) CancelExperimentRun(ctx context.Context, id uint64) error {
	s.Lock()
	defer s.Unlock()

	newExperiments := []SimpleExperiment{}
	for _, e := range s.experiments {
		if e.id == id {
			continue
		}

		newExperiments = append(newExperiments, e)
	}

	s.experiments = newExperiments
	return nil
}

func (s *SimpleStorer) GetExperiments(ctx context.Context, configType string, status experimentationv1.GetExperimentsRequest_Status) ([]*experimentationv1.Experiment, error) {
	s.Lock()
	defer s.Unlock()

	experiments := []*experimentationv1.Experiment{}

	for _, e := range s.experiments {
		experiments = append(experiments, e.toProto())
	}

	return experiments, nil
}
func (s *SimpleStorer) GetExperimentRunDetails(ctx context.Context, id uint64) (*experimentationv1.ExperimentRunDetails, error) {
	return nil, nil
}
func (s *SimpleStorer) GetListView(ctx context.Context) ([]*experimentationv1.ListViewItem, error) {
	return nil, nil
}
func (s *SimpleStorer) RegisterTransformation(transformation experimentstore.Transformation) error {
	return nil
}

func (s *SimpleStorer) Close() {}

var _ experimentstore.Storer = &SimpleStorer{}

package rtds

import (
	"context"
	"time"

	"github.com/golang/protobuf/ptypes"
	"google.golang.org/protobuf/types/known/anypb"

	experimentationv1 "github.com/lyft/clutch/backend/api/chaos/experimentation/v1"
	"github.com/lyft/clutch/backend/service/chaos/experimentation/experimentstore"
)

type simpleExperiment struct {
	id        uint64
	config    *anypb.Any
	startTime *time.Time
	stopTime  *time.Time
}

func (s simpleExperiment) toProto() *experimentationv1.Experiment {
	startTime, _ := ptypes.TimestampProto(*s.startTime)
	endTime, _ := ptypes.TimestampProto(*s.stopTime)
	return &experimentationv1.Experiment{
		Id:        s.id,
		Config:    s.config,
		StartTime: startTime,
		EndTime:   endTime,
	}
}

type simpleStorer struct {
	experiments []simpleExperiment
	idGenerator uint64
}

func (s *simpleStorer) CreateExperiment(ctx context.Context, config *anypb.Any, startTime *time.Time, stopTime *time.Time) (*experimentationv1.Experiment, error) {
	s.experiments = append(s.experiments, simpleExperiment{id: s.idGenerator, config: config, startTime: startTime, stopTime: stopTime})
	s.idGenerator++
	return s.experiments[len(s.experiments)-1].toProto(), nil
}

func (s *simpleStorer) CancelExperimentRun(ctx context.Context, id uint64) error {
	newExperiments := []simpleExperiment{}
	for _, e := range s.experiments {
		if e.id == id {
			continue
		}

		newExperiments = append(newExperiments, e)
	}

	s.experiments = newExperiments
	return nil
}

func (s *simpleStorer) GetExperiments(ctx context.Context, configType string, status experimentationv1.GetExperimentsRequest_Status) ([]*experimentationv1.Experiment, error) {
	experiments := []*experimentationv1.Experiment{}

	for _, e := range s.experiments {
		experiments = append(experiments, e.toProto())
	}

	return experiments, nil
}
func (s *simpleStorer) GetExperimentRunDetails(ctx context.Context, id uint64) (*experimentationv1.ExperimentRunDetails, error) {
	return nil, nil
}
func (s *simpleStorer) GetListView(ctx context.Context) ([]*experimentationv1.ListViewItem, error) {
	return nil, nil
}
func (s *simpleStorer) RegisterTransformation(transformation experimentstore.Transformation) error {
	return nil
}

func (s *simpleStorer) Close() {}

var _ experimentstore.Storer = &simpleStorer{}

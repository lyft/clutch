package terminator

import (
	"fmt"
	"time"

	"github.com/golang/protobuf/ptypes/any"
	"google.golang.org/protobuf/proto"

	experimentationv1 "github.com/lyft/clutch/backend/api/chaos/experimentation/v1"
	terminatorv1 "github.com/lyft/clutch/backend/api/config/service/chaos/experimentation/terminator/v1"
)

type maxTimeTerminationCriterion struct {
	maxTime time.Duration
}

func (m maxTimeTerminationCriterion) ShouldTerminate(experiment *experimentationv1.Experiment, experimentConfig proto.Message) (string, error) {
	startTime := experiment.StartTime.AsTime()

	if startTime.Add(m.maxTime).Before(time.Now()) {
		return fmt.Sprintf("timed out (max duration %s)", m.maxTime), nil
	}

	return "", nil
}

type maxTimeTerminationFactory struct{}

func (maxTimeTerminationFactory) Create(cfg *any.Any) (TerminationCriterion, error) {
	typedConfig := &terminatorv1.MaxTimeTerminationCriterion{}
	err := cfg.UnmarshalTo(typedConfig)
	if err != nil {
		return nil, err
	}

	maxTime := typedConfig.MaxDuration.AsDuration()

	return maxTimeTerminationCriterion{maxTime: maxTime}, nil
}

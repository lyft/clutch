package terminator

import (
	"fmt"
	"time"

	"github.com/golang/protobuf/ptypes/any"

	terminatorv1 "github.com/lyft/clutch/backend/api/config/service/chaos/experimentation/terminator/v1"
	"github.com/lyft/clutch/backend/service/chaos/experimentation/experimentstore"
)

type maxTimeTerminationCriterion struct {
	maxTime time.Duration
}

func (m maxTimeTerminationCriterion) ShouldTerminate(experiment *experimentstore.Experiment) (string, error) {
	startTime := experiment.Run.StartTime

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

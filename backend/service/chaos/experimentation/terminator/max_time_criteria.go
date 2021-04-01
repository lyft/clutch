package terminator

import (
	"fmt"
	"time"

	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/any"
	"google.golang.org/protobuf/proto"

	experimentationv1 "github.com/lyft/clutch/backend/api/chaos/experimentation/v1"
	terminatorv1 "github.com/lyft/clutch/backend/api/config/service/chaos/experimentation/terminator/v1"
)

type maxTimeTerminationCriteria struct {
	maxTime time.Duration
}

func (m maxTimeTerminationCriteria) ShouldTerminate(experiment *experimentationv1.Experiment, experimentConfig proto.Message) (string, error) {
	startTime, err := ptypes.Timestamp(experiment.StartTime)
	if err != nil {
		return "", err
	}

	if startTime.Add(m.maxTime).Before(time.Now()) {
		return fmt.Sprintf("timed out (max duration %s)", m.maxTime), nil
	}

	return "", nil
}

type maxTimeTerminationFactory struct{}

func (maxTimeTerminationFactory) Create(cfg *any.Any) (TerminationCriteria, error) {
	typedConfig := &terminatorv1.MaxTimeTerminationCriteria{}
	err := ptypes.UnmarshalAny(cfg, typedConfig)
	if err != nil {
		return nil, err
	}

	maxTime, err := ptypes.Duration(typedConfig.MaxDuration)
	if err != nil {
		return nil, err
	}

	return maxTimeTerminationCriteria{maxTime: maxTime}, nil
}

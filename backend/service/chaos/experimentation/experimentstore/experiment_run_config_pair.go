package experimentstore

import (
	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/timestamp"

	experimentation "github.com/lyft/clutch/backend/api/chaos/experimentation/v1"
)

type experimentRunConfigPair struct {
	run    *ExperimentRun
	config *ExperimentConfig
}

func NewExperimentFromRunConfigPair(rc *experimentRunConfigPair) (*experimentation.Experiment, error) {
	startTimestampProto, err := ptypes.TimestampProto(rc.run.StartTime)
	if err != nil {
		return nil, err
	}

	var endTimestampProto *timestamp.Timestamp
	if rc.run.EndTime.Valid {
		endTimestampProto, err = ptypes.TimestampProto(rc.run.EndTime.Time)
		if err != nil {
			return nil, err
		}
	}

	return &experimentation.Experiment{
		RunId:     rc.run.Id,
		StartTime: startTimestampProto,
		EndTime:   endTimestampProto,
		Config:    rc.config.Config,
	}, nil
}
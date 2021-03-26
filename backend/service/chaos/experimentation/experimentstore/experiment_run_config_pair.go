package experimentstore

import (
	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/timestamp"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	experimentation "github.com/lyft/clutch/backend/api/chaos/experimentation/v1"
)

type experimentRunConfigPair struct {
	run    *ExperimentRun
	config *ExperimentConfig
}

func (rc *experimentRunConfigPair) toExperiment() (*experimentation.Experiment, error) {
	startTimestampProto, err := ptypes.TimestampProto(rc.run.StartTime)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "%v", err)
	}

	var endTimestampProto *timestamp.Timestamp
	if rc.run.EndTime.Valid {
		endTimestampProto, err = ptypes.TimestampProto(rc.run.EndTime.Time)
		if err != nil {
			return nil, status.Errorf(codes.Internal, "%v", err)
		}
	}

	return &experimentation.Experiment{
		RunId:     rc.run.Id,
		StartTime: startTimestampProto,
		EndTime:   endTimestampProto,
		Config:    rc.config.Config,
	}, nil
}
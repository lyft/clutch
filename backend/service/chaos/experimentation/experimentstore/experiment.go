package experimentstore

import (
	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/timestamp"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	experimentationv1 "github.com/lyft/clutch/backend/api/chaos/experimentation/v1"
)

type Experiment struct {
	Run    *ExperimentRun
	Config *ExperimentConfig
}

func (rc *Experiment) toProto() (*experimentationv1.Experiment, error) {
	startTimestampProto, err := ptypes.TimestampProto(rc.Run.StartTime)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "%v", err)
	}

	var endTimestampProto *timestamp.Timestamp
	if rc.Run.EndTime != nil {
		endTimestampProto, err = ptypes.TimestampProto(*rc.Run.EndTime)
		if err != nil {
			return nil, status.Errorf(codes.Internal, "%v", err)
		}
	}

	return &experimentationv1.Experiment{
		RunId:     rc.Run.Id,
		StartTime: startTimestampProto,
		EndTime:   endTimestampProto,
		Config:    rc.Config.Config,
	}, nil
}

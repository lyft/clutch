package experimentstore

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	experimentationv1 "github.com/lyft/clutch/backend/api/chaos/experimentation/v1"
)

type Experiment struct {
	Run    *ExperimentRun
	Config *ExperimentConfig
}

func (rc Experiment) toProto() (*experimentationv1.Experiment, error) {
	startTimestampProto := timestamppb.New(rc.Run.StartTime)
	if err := startTimestampProto.CheckValid(); err != nil {
		return nil, status.Errorf(codes.Internal, "%v", err)
	}

	var endTimestampProto *timestamppb.Timestamp
	if rc.Run.EndTime != nil {
		endTimestampProto = timestamppb.New(*rc.Run.EndTime)
		if err := endTimestampProto.CheckValid(); err != nil {
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

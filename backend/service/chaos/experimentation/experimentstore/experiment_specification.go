package experimentstore

import (
	"time"

	"github.com/golang/protobuf/ptypes/any"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	experimentationv1 "github.com/lyft/clutch/backend/api/chaos/experimentation/v1"
	"github.com/lyft/clutch/backend/id"
)

type ExperimentSpecification struct {
	RunId     string
	ConfigId  string
	StartTime time.Time
	EndTime   *time.Time
	Config    *any.Any
}

// New returns a new experimentationSpecification instance.
func NewExperimentSpecification(ced *experimentationv1.CreateExperimentData, now time.Time) (*ExperimentSpecification, error) {
	var runId = ced.RunId
	// If the run identifier is not provided, generate one
	if runId == "" {
		runId = id.NewID().String()
	}

	// If the start time is not provided, default to current time
	startTime := now
	if ced.StartTime != nil {
		startTime = ced.StartTime.AsTime()
	}

	if now.After(startTime) {
		return nil, status.Errorf(codes.InvalidArgument, "experiment start time (%v) cannot be before current time (%v)", startTime, now)
	}

	// If the end time is not provided, default to no end time
	var endTime *time.Time = nil
	if ced.EndTime != nil {
		s := ced.EndTime.AsTime()
		endTime = &s
	}

	if endTime != nil && !endTime.After(startTime) {
		return nil, status.Errorf(codes.InvalidArgument, "experiment end time (%v) must be after experiment start time (%v)", *endTime, startTime)
	}

	return &ExperimentSpecification{
		RunId:     runId,
		ConfigId:  id.NewID().String(),
		StartTime: startTime,
		EndTime:   endTime,
		Config:    ced.Config,
	}, nil
}

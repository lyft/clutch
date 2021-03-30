package experimentstore

import (
	"regexp"
	"time"

	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/any"
	"github.com/golang/protobuf/ptypes/timestamp"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	experimentationv1 "github.com/lyft/clutch/backend/api/chaos/experimentation/v1"
	"github.com/lyft/clutch/backend/id"
)

const runIdRegExp = `^[A-Za-z0-9-._~]+$`

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
	} else if !regexp.MustCompile(runIdRegExp).MatchString(runId) {
		return nil, status.Errorf(codes.InvalidArgument, "provided experiment runId (%v) contained forbidden characters and was not matched by \"%v\" regular expresion", runId, runIdRegExp)
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

func (es *ExperimentSpecification) toExperiment() (*experimentationv1.Experiment, error) {
	startTimestamp, err := ptypes.TimestampProto(es.StartTime)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "%v", err)
	}

	var endTimestamp *timestamp.Timestamp
	if es.EndTime != nil {
		endTimestamp, err = ptypes.TimestampProto(*es.EndTime)
		if err != nil {
			return nil, status.Errorf(codes.InvalidArgument, "%v", err)
		}
	}

	return &experimentationv1.Experiment{
		RunId:     es.RunId,
		StartTime: startTimestamp,
		EndTime:   endTimestamp,
		Config:    es.Config,
	}, nil
}

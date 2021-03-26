package experimentstore

import (
	"errors"
	"fmt"
	"regexp"
	"time"

	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/any"
	"github.com/golang/protobuf/ptypes/timestamp"

	experimentation "github.com/lyft/clutch/backend/api/chaos/experimentation/v1"
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
func NewExperimentSpecification(ced *experimentation.CreateExperimentData, now time.Time) (*ExperimentSpecification, error) {
	var runId = ced.RunId
	// If the run identifier is not provided, generate one
	if runId == "" {
		runId = id.NewID().String()
	} else if !regexp.MustCompile(runIdRegExp).MatchString(runId) {
		return nil, fmt.Errorf("provided experiment runId (%v) contained unallowed characters and was not matched by \"%v\" regular expresion", runId, runIdRegExp)
	}

	if ced.Config == nil {
		return nil, errors.New("experiment config cannot be equal to nil")
	}

	// If the start time is not provided, default to current time
	startTime := now
	if ced.StartTime != nil {
		startTime = ced.StartTime.AsTime()
	}

	if now.After(startTime) {
		return nil, fmt.Errorf("current time (%v) must be equal to or after experiment start time (%v)", now, startTime)
	}

	// If the end time is not provided, default to no end time
	var endTime *time.Time = nil
	if ced.EndTime != nil {
		s := ced.EndTime.AsTime()
		endTime = &s
	}

	if endTime != nil && !endTime.After(startTime) {
		return nil, fmt.Errorf("experiment end time (%v) must be after experiment start time (%v)", *endTime, startTime)
	}

	return &ExperimentSpecification{
		RunId:     runId,
		ConfigId:  id.NewID().String(),
		StartTime: startTime,
		EndTime:   endTime,
		Config:    ced.Config,
	}, nil
}

func (es *ExperimentSpecification) toExperiment() (*experimentation.Experiment, error) {
	startTimestamp, err := ptypes.TimestampProto(es.StartTime)
	if err != nil {
		return nil, err
	}

	var endTimestamp *timestamp.Timestamp
	if es.EndTime != nil {
		endTimestamp, err = ptypes.TimestampProto(*es.EndTime)
		if err != nil {
			return nil, err
		}
	}

	return &experimentation.Experiment{
		RunId:     es.RunId,
		StartTime: startTimestamp,
		EndTime:   endTimestamp,
		Config:    es.Config,
	}, nil
}

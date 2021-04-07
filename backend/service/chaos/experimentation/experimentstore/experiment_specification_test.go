package experimentstore

import (
	"fmt"
	"testing"
	"time"

	"github.com/golang/protobuf/ptypes/any"
	"github.com/golang/protobuf/ptypes/timestamp"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	experimentationv1 "github.com/lyft/clutch/backend/api/chaos/experimentation/v1"
)

func TestExperimentSpecificationInitialization(t *testing.T) {
	a := assert.New(t)

	now := time.Date(2011, 0, 0, 0, 0, 0, 0, time.UTC)
	past := now.Add(-1 * time.Hour)
	pastTimestamp := timestamppb.New(past)
	future := now.Add(1 * time.Hour)
	futureTimestamp := timestamppb.New(future)
	farFutureTimestamp := timestamppb.New(future.Add(1 * time.Hour))

	tests := []struct {
		runId             string
		startTime         *timestamp.Timestamp
		endTime           *timestamp.Timestamp
		now               time.Time
		Config            *any.Any
		expectedError     error
		expectedStartTime time.Time
	}{
		{
			runId:             "",
			startTime:         futureTimestamp,
			endTime:           nil,
			now:               now,
			Config:            &any.Any{},
			expectedError:     nil,
			expectedStartTime: future,
		},
		{
			runId:             "aA0-._~",
			startTime:         futureTimestamp,
			endTime:           nil,
			now:               now,
			Config:            &any.Any{},
			expectedError:     nil,
			expectedStartTime: future,
		},
		{
			runId:             "1231231231^",
			startTime:         futureTimestamp,
			endTime:           farFutureTimestamp,
			now:               now,
			Config:            &any.Any{},
			expectedError:     status.Error(codes.InvalidArgument, "provided experiment runId (1231231231^) contained forbidden characters and was not matched by \"^[A-Za-z0-9-._~]+$\" regular expresion"),
			expectedStartTime: future,
		},
		{
			runId:             "1",
			startTime:         futureTimestamp,
			endTime:           farFutureTimestamp,
			now:               now,
			Config:            &any.Any{},
			expectedError:     nil,
			expectedStartTime: future,
		},
		{
			runId:             "1",
			startTime:         futureTimestamp,
			endTime:           farFutureTimestamp,
			now:               now,
			Config:            nil,
			expectedError:     status.Error(codes.InvalidArgument, "experiment config cannot be equal to nil"),
			expectedStartTime: future,
		},
		{
			runId:         "1",
			startTime:     pastTimestamp,
			endTime:       futureTimestamp,
			now:           now,
			Config:        &any.Any{},
			expectedError: status.Error(codes.InvalidArgument, "experiment start time (2010-11-29 23:00:00 +0000 UTC) cannot be before current time (2010-11-30 00:00:00 +0000 UTC)"),
		},
		{
			runId:         "1",
			startTime:     futureTimestamp,
			endTime:       pastTimestamp,
			now:           now,
			Config:        &any.Any{},
			expectedError: status.Error(codes.InvalidArgument, "experiment end time (2010-11-29 23:00:00 +0000 UTC) must be after experiment start time (2010-11-30 01:00:00 +0000 UTC)"),
		},
		{
			runId:         "1",
			startTime:     futureTimestamp,
			endTime:       futureTimestamp,
			now:           now,
			Config:        &any.Any{},
			expectedError: status.Error(codes.InvalidArgument, "experiment end time (2010-11-30 01:00:00 +0000 UTC) must be after experiment start time (2010-11-30 01:00:00 +0000 UTC)"),
		},
		{
			runId:         "",
			startTime:     farFutureTimestamp,
			endTime:       futureTimestamp,
			now:           now,
			Config:        &any.Any{},
			expectedError: status.Error(codes.InvalidArgument, "experiment end time (2010-11-30 01:00:00 +0000 UTC) must be after experiment start time (2010-11-30 02:00:00 +0000 UTC)"),
		},
	}

	for idx, tt := range tests {
		tt := tt
		t.Run(fmt.Sprintf("%d", idx), func(t *testing.T) {
			t.Parallel()
			ced := &experimentationv1.CreateExperimentData{RunId: tt.runId, StartTime: tt.startTime, EndTime: tt.endTime, Config: tt.Config}
			es, err := NewExperimentSpecification(ced, tt.now)
			if err != nil {
				a.Equal(tt.expectedError, err)
			} else {
				a.True(len(es.RunId) > 0)
				a.True(len(es.ConfigId) > 0)
				a.NotEqual(es.RunId, es.ConfigId)
				a.Equal(es.StartTime, tt.expectedStartTime)
			}
		})
	}
}

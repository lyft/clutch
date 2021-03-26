package experimentstore

import (
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/any"
	"github.com/golang/protobuf/ptypes/timestamp"
	"github.com/stretchr/testify/assert"

	experimentation "github.com/lyft/clutch/backend/api/chaos/experimentation/v1"
)

func TestExperimentSpecificationInitialization(t *testing.T) {
	a := assert.New(t)

	now := time.Date(2011, 0, 0, 0, 0, 0, 0, time.UTC)
	past := time.Date(2010, 0, 0, 0, 0, 0, 0, time.UTC)
	pastTimestamp, err := ptypes.TimestampProto(past)
	a.NoError(err)
	future1 := time.Date(2030, 0, 0, 0, 0, 0, 0, time.UTC)
	futureTimestamp1, err := ptypes.TimestampProto(future1)
	a.NoError(err)
	future2 := time.Date(2031, 0, 0, 0, 0, 0, 0, time.UTC)
	futureTimestamp2, err := ptypes.TimestampProto(future2)
	a.NoError(err)

	config := &any.Any{}

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
			startTime:         futureTimestamp1,
			endTime:           nil,
			now:               now,
			Config:            config,
			expectedError:     nil,
			expectedStartTime: future1,
		},
		{
			runId:             "aA0-._~",
			startTime:         futureTimestamp1,
			endTime:           nil,
			now:               now,
			Config:            config,
			expectedError:     nil,
			expectedStartTime: future1,
		},
		{
			runId:             "1231231231^",
			startTime:         futureTimestamp1,
			endTime:           futureTimestamp2,
			now:               now,
			Config:            config,
			expectedError:     errors.New("provided experiment runId (1231231231^) contained unallowed characters and was not matched by \"^[A-Za-z0-9-._~]+$\" regular expresion"),
			expectedStartTime: future1,
		},
		{
			runId:             "1",
			startTime:         futureTimestamp1,
			endTime:           futureTimestamp2,
			now:               now,
			Config:            config,
			expectedError:     nil,
			expectedStartTime: future1,
		},
		{
			runId:             "1",
			startTime:         futureTimestamp1,
			endTime:           futureTimestamp2,
			now:               now,
			Config:            nil,
			expectedError:     errors.New("experiment config cannot be equal to nil"),
			expectedStartTime: future1,
		},
		{
			runId:         "1",
			startTime:     pastTimestamp,
			endTime:       futureTimestamp1,
			now:           now,
			Config:        config,
			expectedError: errors.New("current time (2010-11-30 00:00:00 +0000 UTC) must be equal to or after experiment start time (2009-11-30 00:00:00 +0000 UTC)"),
		},
		{
			runId:         "1",
			startTime:     futureTimestamp1,
			endTime:       pastTimestamp,
			now:           now,
			Config:        config,
			expectedError: errors.New("experiment end time (2009-11-30 00:00:00 +0000 UTC) must be after experiment start time (2029-11-30 00:00:00 +0000 UTC)"),
		},
		{
			runId:         "1",
			startTime:     futureTimestamp1,
			endTime:       futureTimestamp1,
			now:           now,
			Config:        config,
			expectedError: errors.New("experiment end time (2029-11-30 00:00:00 +0000 UTC) must be after experiment start time (2029-11-30 00:00:00 +0000 UTC)"),
		},
		{
			runId:         "",
			startTime:     futureTimestamp2,
			endTime:       futureTimestamp1,
			now:           now,
			Config:        config,
			expectedError: errors.New("experiment end time (2029-11-30 00:00:00 +0000 UTC) must be after experiment start time (2030-11-30 00:00:00 +0000 UTC)"),
		},
	}

	for idx, tt := range tests {
		tt := tt
		t.Run(fmt.Sprintf("%d", idx), func(t *testing.T) {
			t.Parallel()
			ced := &experimentation.CreateExperimentData{RunId: tt.runId, StartTime: tt.startTime, EndTime: tt.endTime, Config: tt.Config}
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
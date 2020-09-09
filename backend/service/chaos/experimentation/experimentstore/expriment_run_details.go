package experimentstore

import (
	"database/sql"
	"strconv"
	"strings"
	"time"

	"github.com/golang/protobuf/jsonpb"
	"github.com/golang/protobuf/ptypes/any"

	experimentation "github.com/lyft/clutch/backend/api/chaos/experimentation/v1"
)

func NewRunDetails(fetchedID uint64, startTime sql.NullTime, endTime sql.NullTime, scheduledEndTime sql.NullTime, creationTime time.Time, details string) (*experimentation.ExperimentRunDetails, error) {
	runConfigPair := &experimentation.ExperimentRunDetails{}
	runConfigPair.RunId = fetchedID
	runConfigPair.Properties = &experimentation.Properties{}
	runConfigPair.Status = timesToStatus(startTime, endTime, scheduledEndTime)
	runConfigPair.GetProperties().Items = []*experimentation.Property{
		{Label: "Run Identifier", Value: strconv.FormatUint(fetchedID, 10)},
		{Label: "Status", Value: statusToString(runConfigPair.Status)},
		{Label: "Start Time", Value: timeToString(startTime)},
	}

	if runConfigPair.Status == experimentation.Status_COMPLETED {
		endTimeField := &experimentation.Property{Label: "Scheduled End Time", Value: timeToString(endTime)}
		runConfigPair.GetProperties().Items = append(runConfigPair.GetProperties().Items, endTimeField)
	}

	scheduledEndTimeField := &experimentation.Property{Label: "Scheduled End Time", Value: timeToString(endTime)}
	runConfigPair.GetProperties().Items = append(runConfigPair.GetProperties().Items, scheduledEndTimeField)

	anyConfig := &any.Any{}
	err := jsonpb.Unmarshal(strings.NewReader(details), anyConfig)
	if err != nil {
		return nil, err
	}

	runConfigPair.Config = anyConfig
	return runConfigPair, err
}

func timeToString(t sql.NullTime) string {
	layout := "01/02/06 15:04:05"
	if t.Valid {
		return t.Time.Format(layout)
	} else {
		return "Undefined"
	}
}

func timesToStatus(startTime sql.NullTime, endTime sql.NullTime, scheduledEndTime sql.NullTime) experimentation.Status {
	now := time.Now()
	switch {
	case startTime.Valid && !endTime.Valid:
		if now.After(startTime.Time) {
			return experimentation.Status_RUNNING
		} else {
			return experimentation.Status_SCHEDULED
		}
	case !startTime.Valid && endTime.Valid:
		if now.Before(endTime.Time) {
			return experimentation.Status_RUNNING
		} else {
			return experimentation.Status_COMPLETED
		}
	default:
		if now.Before(startTime.Time) {
			return experimentation.Status_SCHEDULED
		} else if now.After(startTime.Time) && now.Before(endTime.Time) {
			return experimentation.Status_RUNNING
		}
		return experimentation.Status_COMPLETED
	}
}

func statusToString(status experimentation.Status) string {
	switch status {
	case experimentation.Status_UNSPECIFIED:
		return "Unspecified"
	case experimentation.Status_SCHEDULED:
		return "Scheduled"
	case experimentation.Status_RUNNING:
		return "Running"
	case experimentation.Status_COMPLETED:
		return "Completed"
	}

	return status.String()
}

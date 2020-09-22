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

func NewRunDetails(fetchedID uint64, creationTime time.Time, startTime sql.NullTime, endTime sql.NullTime, cancellationTime sql.NullTime, now time.Time, details string) (*experimentation.ExperimentRunDetails, error) {
	runConfigPair := &experimentation.ExperimentRunDetails{}
	runConfigPair.RunId = fetchedID
	runConfigPair.Properties = &experimentation.Properties{}
	runConfigPair.Status = timesToStatus(startTime, endTime, cancellationTime, now)
	runConfigPair.GetProperties().Items = []*experimentation.Property{
		{Label: "Run Identifier", Value: strconv.FormatUint(fetchedID, 10)},
		{Label: "Status", Value: statusToString(runConfigPair.Status)},
		{Label: "Start Time", Value: timeToString(startTime)},
	}

	if runConfigPair.Status == experimentation.Status_COMPLETED {
		var time sql.NullTime
		if endTime.Valid {
			time = endTime
		} else if cancellationTime.Valid {
			time = cancellationTime
		}

		endTimeField := &experimentation.Property{Label: "End Time", Value: timeToString(time)}
		runConfigPair.GetProperties().Items = append(runConfigPair.GetProperties().Items, endTimeField)
	}

	if runConfigPair.Status == experimentation.Status_STOPPED {
		stoppedTimeField := &experimentation.Property{Label: "Stopped At", Value: timeToString(cancellationTime)}
		runConfigPair.GetProperties().Items = append(runConfigPair.GetProperties().Items, stoppedTimeField)
	} else if runConfigPair.Status == experimentation.Status_CANCELED {
		cancellationTimeField := &experimentation.Property{Label: "Canceled At", Value: timeToString(cancellationTime)}
		runConfigPair.GetProperties().Items = append(runConfigPair.GetProperties().Items, cancellationTimeField)
	}

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

func timesToStatus(startTime sql.NullTime, endTime sql.NullTime, cancellationTime sql.NullTime, now time.Time) experimentation.Status {
	if !startTime.Valid {
		return experimentation.Status_UNSPECIFIED
	}

	if cancellationTime.Valid {
		if cancellationTime.Time.After(startTime.Time) {
			if endTime.Valid {
				return experimentation.Status_STOPPED
			} else {
				return experimentation.Status_COMPLETED
			}
		} else {
			return experimentation.Status_CANCELED
		}
	} else {
		if now.Before(startTime.Time) {
			return experimentation.Status_SCHEDULED
		} else if now.After(startTime.Time) && (!endTime.Valid || now.Before(endTime.Time)) {
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
	case experimentation.Status_CANCELED:
		return "Canceled"
	case experimentation.Status_STOPPED:
		return "Stopped"
	default:
		return status.String()
	}
}

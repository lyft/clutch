package experimentstore

import (
	"database/sql"
	"errors"
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
	status, err := TimesToStatus(startTime, endTime, cancellationTime, now)
	if err != nil {
		return nil, err
	}

	runConfigPair.Status = status
	runConfigPair.GetProperties().Items = []*experimentation.Property{
		{Label: "Run Identifier", Value: strconv.FormatUint(fetchedID, 10)},
		{Label: "Status", Value: statusToString(runConfigPair.Status)},
		{Label: "Start Time", Value: timeToString(startTime)},
	}

	if runConfigPair.Status == experimentation.ExperimentRunDetails_COMPLETED {
		var time sql.NullTime
		if endTime.Valid {
			time = endTime
		} else if cancellationTime.Valid {
			time = cancellationTime
		}

		endTimeField := &experimentation.Property{Label: "End Time", Value: timeToString(time)}
		runConfigPair.GetProperties().Items = append(runConfigPair.GetProperties().Items, endTimeField)
	}

	if runConfigPair.Status == experimentation.ExperimentRunDetails_STOPPED {
		stoppedTimeField := &experimentation.Property{Label: "Stopped At", Value: timeToString(cancellationTime)}
		runConfigPair.GetProperties().Items = append(runConfigPair.GetProperties().Items, stoppedTimeField)
	} else if runConfigPair.Status == experimentation.ExperimentRunDetails_CANCELED {
		cancellationTimeField := &experimentation.Property{Label: "Canceled At", Value: timeToString(cancellationTime)}
		runConfigPair.GetProperties().Items = append(runConfigPair.GetProperties().Items, cancellationTimeField)
	}

	anyConfig := &any.Any{}
	err = jsonpb.Unmarshal(strings.NewReader(details), anyConfig)
	if err != nil {
		return nil, err
	}

	runConfigPair.Config = anyConfig
	return runConfigPair, err
}

func TimesToStatus(startTime sql.NullTime, endTime sql.NullTime, cancellationTime sql.NullTime, now time.Time) (experimentation.ExperimentRunDetails_Status, error) {
	if !startTime.Valid {
		return experimentation.ExperimentRunDetails_UNSPECIFIED, errors.New("experiment doesn't have startTime")
	}

	if cancellationTime.Valid {
		if cancellationTime.Time.After(startTime.Time) {
			if endTime.Valid {
				return experimentation.ExperimentRunDetails_STOPPED, nil
			} else {
				return experimentation.ExperimentRunDetails_COMPLETED, nil
			}
		} else {
			return experimentation.ExperimentRunDetails_CANCELED, nil
		}
	} else {
		if now.Before(startTime.Time) {
			return experimentation.ExperimentRunDetails_SCHEDULED, nil
		} else if now.After(startTime.Time) && (!endTime.Valid || now.Before(endTime.Time)) {
			return experimentation.ExperimentRunDetails_RUNNING, nil
		}
		return experimentation.ExperimentRunDetails_COMPLETED, nil
	}
}

func timeToString(t sql.NullTime) string {
	layout := "01/02/06 15:04:05"
	if t.Valid {
		return t.Time.Format(layout)
	} else {
		return "Undefined"
	}
}

func statusToString(status experimentation.ExperimentRunDetails_Status) string {
	switch status {
	case experimentation.ExperimentRunDetails_UNSPECIFIED:
		return "Unspecified"
	case experimentation.ExperimentRunDetails_SCHEDULED:
		return "Scheduled"
	case experimentation.ExperimentRunDetails_RUNNING:
		return "Running"
	case experimentation.ExperimentRunDetails_COMPLETED:
		return "Completed"
	case experimentation.ExperimentRunDetails_CANCELED:
		return "Canceled"
	case experimentation.ExperimentRunDetails_STOPPED:
		return "Stopped"
	default:
		return status.String()
	}
}

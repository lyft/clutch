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
	status, err := timesToStatus(startTime, endTime, cancellationTime, now)
	if err != nil {
		return nil, err
	}

	runConfigPair.Status = status
	runConfigPair.GetProperties().Items = []*experimentation.Property{
		{Label: "Run Identifier", Value: strconv.FormatUint(fetchedID, 10)},
		{Label: "Status", Value: statusToString(runConfigPair.Status)},
		{Label: "Start Time", Value: timeToString(startTime)},
	}

	if runConfigPair.Status == experimentation.Experiment_COMPLETED {
		var time sql.NullTime
		if endTime.Valid {
			time = endTime
		} else if cancellationTime.Valid {
			time = cancellationTime
		}

		endTimeField := &experimentation.Property{Label: "End Time", Value: timeToString(time)}
		runConfigPair.GetProperties().Items = append(runConfigPair.GetProperties().Items, endTimeField)
	}

	if runConfigPair.Status == experimentation.Experiment_STOPPED {
		stoppedTimeField := &experimentation.Property{Label: "Stopped At", Value: timeToString(cancellationTime)}
		runConfigPair.GetProperties().Items = append(runConfigPair.GetProperties().Items, stoppedTimeField)
	} else if runConfigPair.Status == experimentation.Experiment_CANCELED {
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

func timesToStatus(startTime sql.NullTime, endTime sql.NullTime, cancellationTime sql.NullTime, now time.Time) (experimentation.Experiment_Status, error) {
	if !startTime.Valid {
		return experimentation.Experiment_UNSPECIFIED, errors.New("experiment doesn't have startTime")
	}

	if cancellationTime.Valid {
		if cancellationTime.Time.After(startTime.Time) {
			if endTime.Valid {
				return experimentation.Experiment_STOPPED, nil
			} else {
				return experimentation.Experiment_COMPLETED, nil
			}
		} else {
			return experimentation.Experiment_CANCELED, nil
		}
	} else {
		if now.Before(startTime.Time) {
			return experimentation.Experiment_SCHEDULED, nil
		} else if now.After(startTime.Time) && (!endTime.Valid || now.Before(endTime.Time)) {
			return experimentation.Experiment_RUNNING, nil
		}
		return experimentation.Experiment_COMPLETED, nil
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

func statusToString(status experimentation.Experiment_Status) string {
	switch status {
	case experimentation.Experiment_UNSPECIFIED:
		return "Unspecified"
	case experimentation.Experiment_SCHEDULED:
		return "Scheduled"
	case experimentation.Experiment_RUNNING:
		return "Running"
	case experimentation.Experiment_COMPLETED:
		return "Completed"
	case experimentation.Experiment_CANCELED:
		return "Canceled"
	case experimentation.Experiment_STOPPED:
		return "Stopped"
	default:
		return status.String()
	}
}

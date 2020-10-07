package experimentstore

import (
	"database/sql"
	"errors"
	"strings"
	"time"

	"github.com/golang/protobuf/jsonpb"
	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/any"
	"github.com/golang/protobuf/ptypes/timestamp"

	experimentation "github.com/lyft/clutch/backend/api/chaos/experimentation/v1"
)

func NewRunDetails(fetchedID uint64, creationTime time.Time, startTime sql.NullTime, endTime sql.NullTime, cancellationTime sql.NullTime, now time.Time, details string) (*experimentation.ExperimentRunDetails, error) {
	status, err := timesToStatus(startTime, endTime, cancellationTime, now)
	if err != nil {
		return nil, err
	}

	properties, err := NewProperties(fetchedID, creationTime, startTime, endTime, cancellationTime, now, details)
	if err != nil {
		return nil, err
	}

	anyConfig := &any.Any{}
	err = jsonpb.Unmarshal(strings.NewReader(details), anyConfig)
	if err != nil {
		return nil, err
	}

	return &experimentation.ExperimentRunDetails{
		RunId:      fetchedID,
		Status:     status,
		Properties: &experimentation.PropertiesList{Items: properties},
		Config:     anyConfig}, nil
}

func NewProperties(fetchedID uint64, creationTime time.Time, startTime sql.NullTime, endTime sql.NullTime, cancellationTime sql.NullTime, now time.Time, details string) ([]*experimentation.Property, error) {
	status, err := timesToStatus(startTime, endTime, cancellationTime, now)
	if err != nil {
		return nil, err
	}

	startTimeTimestamp, err := timeToTimestamp(startTime)
	if err != nil {
		return nil, err
	}

	properties := []*experimentation.Property{
		{
			Id:    "run_identifier",
			Label: "Run Identifier",
			Value: &experimentation.Property_IntValue{IntValue: int64(fetchedID)},
		},
		{
			Id:    "status",
			Label: "Status",
			Value: &experimentation.Property_StringValue{StringValue: statusToString(status)},
		},
		{
			Id:    "start_time",
			Label: "Start Time",
			Value: &experimentation.Property_DateValue{DateValue: startTimeTimestamp},
		},
	}

	var time sql.NullTime
	if endTime.Valid {
		time = endTime
	} else if cancellationTime.Valid {
		time = cancellationTime
	}

	endTimeTimestamp, err := timeToTimestamp(time)
	if err != nil {
		return nil, err
	}

	properties = append(properties, &experimentation.Property{
		Id:    "end_time",
		Label: "End Time",
		Value: &experimentation.Property_DateValue{DateValue: endTimeTimestamp},
	})

	cancelationTimeTimestamp, err := timeToTimestamp(cancellationTime)
	if err != nil {
		return nil, err
	}

	if status == experimentation.Experiment_STOPPED {
		properties = append(properties, &experimentation.Property{
			Id:    "stopped_at",
			Label: "Stopped At",
			Value: &experimentation.Property_DateValue{DateValue: cancelationTimeTimestamp},
		})
	} else if status == experimentation.Experiment_CANCELED {
		properties = append(properties, &experimentation.Property{
			Id:    "canceled_at",
			Label: "Canceled At",
			Value: &experimentation.Property_DateValue{DateValue: cancelationTimeTimestamp},
		})
	}

	return properties, nil
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

func timeToTimestamp(t sql.NullTime) (*timestamp.Timestamp, error) {
	if t.Valid {
		return ptypes.TimestampProto(t.Time)
	} else {
		return nil, nil
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

package experimentstore

import (
	"database/sql"
	"time"

	"github.com/golang/protobuf/ptypes"

	experimentation "github.com/lyft/clutch/backend/api/chaos/experimentation/v1"
)

type ExperimentRun struct {
	id               uint64
	startTime        time.Time
	endTime          sql.NullTime
	cancellationTime sql.NullTime
	creationTime     time.Time
}

func (er *ExperimentRun) CreateProperties(now time.Time) ([]*experimentation.Property, error) {
	status := timesToStatus(er.startTime, er.endTime, er.cancellationTime, now)
	startTimeTimestamp, err := ptypes.TimestampProto(er.startTime)
	if err != nil {
		return nil, err
	}

	creationTimeTimestamp, err := ptypes.TimestampProto(er.creationTime)
	if err != nil {
		return nil, err
	}

	properties := []*experimentation.Property{
		{
			Id:    "run_identifier",
			Label: "Run Identifier",
			Value: &experimentation.Property_IntValue{IntValue: int64(er.id)},
		},
		{
			Id:    "status",
			Label: "Status",
			Value: &experimentation.Property_StringValue{StringValue: statusToString(status)},
		},
		{
			Id:    "run_creation_time",
			Label: "Created At",
			Value: &experimentation.Property_DateValue{DateValue: creationTimeTimestamp},
		},
		{
			Id:    "start_time",
			Label: "Start Time",
			Value: &experimentation.Property_DateValue{DateValue: startTimeTimestamp},
		},
	}

	var time sql.NullTime
	if er.endTime.Valid {
		time = er.endTime
	} else if er.cancellationTime.Valid {
		time = er.cancellationTime
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

	cancelationTimeTimestamp, err := timeToTimestamp(er.cancellationTime)
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

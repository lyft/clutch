package experimentstore

import (
	"database/sql"
	"time"

	"github.com/golang/protobuf/ptypes"

	experimentation "github.com/lyft/clutch/backend/api/chaos/experimentation/v1"
)

type ExperimentRun struct {
	Id               uint64
	StartTime        time.Time
	EndTime          sql.NullTime
	CancellationTime sql.NullTime
	creationTime     time.Time
}

func (er *ExperimentRun) CreateProperties(now time.Time) ([]*experimentation.Property, error) {
	status := timesToStatus(er.StartTime, er.EndTime, er.CancellationTime, now)
	startTimeTimestamp, err := ptypes.TimestampProto(er.StartTime)
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
			Value: &experimentation.Property_IntValue{IntValue: int64(er.Id)},
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
	if er.EndTime.Valid {
		time = er.EndTime
	} else if er.CancellationTime.Valid {
		time = er.CancellationTime
	}

	endTimeTimestamp, err := TimeToPropertyDateValue(time)
	if err != nil {
		return nil, err
	}

	properties = append(properties, &experimentation.Property{
		Id:    "end_time",
		Label: "End Time",
		Value: endTimeTimestamp,
	})

	cancelationTimeTimestamp, err := TimeToPropertyDateValue(er.CancellationTime)
	if err != nil {
		return nil, err
	}

	if status == experimentation.Experiment_STATUS_STOPPED {
		properties = append(properties, &experimentation.Property{
			Id:    "stopped_at",
			Label: "Stopped At",
			Value: cancelationTimeTimestamp,
		})
	} else if status == experimentation.Experiment_STATUS_CANCELED {
		properties = append(properties, &experimentation.Property{
			Id:    "canceled_at",
			Label: "Canceled At",
			Value: cancelationTimeTimestamp,
		})
	}

	return properties, nil
}

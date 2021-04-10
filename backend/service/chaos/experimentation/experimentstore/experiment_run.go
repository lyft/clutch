package experimentstore

import (
	"time"

	"google.golang.org/protobuf/types/known/timestamppb"

	experimentation "github.com/lyft/clutch/backend/api/chaos/experimentation/v1"
)

type ExperimentRun struct {
	Id                string
	StartTime         time.Time
	EndTime           *time.Time
	CancellationTime  *time.Time
	CreationTime      time.Time
	TerminationReason string
}

func (er *ExperimentRun) Status(now time.Time) experimentation.Experiment_Status {
	if er.CancellationTime.Valid {
		if er.CancellationTime.Time.After(er.StartTime) {
			if er.EndTime.Valid {
				return experimentation.Experiment_STATUS_STOPPED
			} else {
				return experimentation.Experiment_STATUS_COMPLETED
			}
		} else {
			return experimentation.Experiment_STATUS_CANCELED
		}
	} else {
		if now.Before(er.StartTime) {
			return experimentation.Experiment_STATUS_SCHEDULED
		} else if now.After(er.StartTime) && (!er.EndTime.Valid || now.Before(er.EndTime.Time)) {
			return experimentation.Experiment_STATUS_RUNNING
		}
		return experimentation.Experiment_STATUS_COMPLETED
	}
}

func (er *ExperimentRun) CreateProperties(now time.Time) ([]*experimentation.Property, error) {
	status := er.Status(now)
	startTimeTimestamp := timestamppb.New(er.StartTime)
	if err := startTimeTimestamp.CheckValid(); err != nil {
		return nil, err
	}

	creationTimeTimestamp := timestamppb.New(er.CreationTime)
	if err := creationTimeTimestamp.CheckValid(); err != nil {
		return nil, err
	}

	properties := []*experimentation.Property{
		{
			Id:    "run_identifier",
			Label: "Run Identifier",
			Value: &experimentation.Property_StringValue{StringValue: er.Id},
		},
		{
			Id:    "status",
			Label: "Status",
			Value: &experimentation.Property_StringValue{StringValue: StatusToString(status)},
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

	var time *time.Time
	if er.EndTime != nil {
		time = er.EndTime
	} else if er.CancellationTime != nil {
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

	cancellationTimeTimestamp, err := TimeToPropertyDateValue(er.CancellationTime)
	if err != nil {
		return nil, err
	}

	if status == experimentation.Experiment_STATUS_STOPPED || status == experimentation.Experiment_STATUS_CANCELED {
		if status == experimentation.Experiment_STATUS_STOPPED {
			properties = append(properties, &experimentation.Property{
				Id:    "stopped_at",
				Label: "Stopped At",
				Value: cancellationTimeTimestamp,
			})
		} else if status == experimentation.Experiment_STATUS_CANCELED {
			properties = append(properties, &experimentation.Property{
				Id:    "canceled_at",
				Label: "Canceled At",
				Value: cancellationTimeTimestamp,
			})
		}

		terminationReason := "Unknown"
		if er.TerminationReason != "" {
			terminationReason = er.TerminationReason
		}

		properties = append(properties, &experimentation.Property{
			Id:    "termination_reason",
			Label: "Termination Reason",
			Value: &experimentation.Property_StringValue{StringValue: terminationReason},
		})
	}

	return properties, nil
}

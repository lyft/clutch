package experimentstore

import (
	"time"

	"google.golang.org/protobuf/types/known/timestamppb"

	experimentationv1 "github.com/lyft/clutch/backend/api/chaos/experimentation/v1"
)

type ExperimentRun struct {
	Id                string
	StartTime         time.Time
	EndTime           *time.Time
	CancellationTime  *time.Time
	CreationTime      time.Time
	TerminationReason string
}

func (er *ExperimentRun) Status(now time.Time) experimentationv1.Experiment_Status {
	if er.CancellationTime != nil {
		if er.CancellationTime.After(er.StartTime) {
			if er.EndTime != nil {
				return experimentationv1.Experiment_STATUS_STOPPED
			} else {
				return experimentationv1.Experiment_STATUS_COMPLETED
			}
		} else {
			return experimentationv1.Experiment_STATUS_CANCELED
		}
	} else {
		if now.Before(er.StartTime) {
			return experimentationv1.Experiment_STATUS_SCHEDULED
		} else if now.After(er.StartTime) && (er.EndTime == nil || now.Before(*er.EndTime)) {
			return experimentationv1.Experiment_STATUS_RUNNING
		}
		return experimentationv1.Experiment_STATUS_COMPLETED
	}
}

func (er *ExperimentRun) CreateProperties(now time.Time) ([]*experimentationv1.Property, error) {
	status := er.Status(now)
	startTimeTimestamp := timestamppb.New(er.StartTime)
	if err := startTimeTimestamp.CheckValid(); err != nil {
		return nil, err
	}

	creationTimeTimestamp := timestamppb.New(er.CreationTime)
	if err := creationTimeTimestamp.CheckValid(); err != nil {
		return nil, err
	}

	properties := []*experimentationv1.Property{
		{
			Id:    "run_identifier",
			Label: "Run Identifier",
			Value: &experimentationv1.Property_StringValue{StringValue: er.Id},
		},
		{
			Id:    "status",
			Label: "Status",
			Value: &experimentationv1.Property_StringValue{StringValue: StatusToString(status)},
		},
		{
			Id:    "run_creation_time",
			Label: "Created At",
			Value: &experimentationv1.Property_DateValue{DateValue: creationTimeTimestamp},
		},
		{
			Id:    "start_time",
			Label: "Start Time",
			Value: &experimentationv1.Property_DateValue{DateValue: startTimeTimestamp},
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

	properties = append(properties, &experimentationv1.Property{
		Id:    "end_time",
		Label: "End Time",
		Value: endTimeTimestamp,
	})

	cancellationTimeTimestamp, err := TimeToPropertyDateValue(er.CancellationTime)
	if err != nil {
		return nil, err
	}

	if status == experimentationv1.Experiment_STATUS_STOPPED || status == experimentationv1.Experiment_STATUS_CANCELED {
		if status == experimentationv1.Experiment_STATUS_STOPPED {
			properties = append(properties, &experimentationv1.Property{
				Id:    "stopped_at",
				Label: "Stopped At",
				Value: cancellationTimeTimestamp,
			})
		} else if status == experimentationv1.Experiment_STATUS_CANCELED {
			properties = append(properties, &experimentationv1.Property{
				Id:    "canceled_at",
				Label: "Canceled At",
				Value: cancellationTimeTimestamp,
			})
		}

		terminationReason := "Unknown"
		if er.TerminationReason != "" {
			terminationReason = er.TerminationReason
		}

		properties = append(properties, &experimentationv1.Property{
			Id:    "termination_reason",
			Label: "Termination Reason",
			Value: &experimentationv1.Property_StringValue{StringValue: terminationReason},
		})
	}

	return properties, nil
}

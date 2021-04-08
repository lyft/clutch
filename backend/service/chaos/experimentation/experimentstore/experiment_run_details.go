package experimentstore

import (
	"time"

	experimentation "github.com/lyft/clutch/backend/api/chaos/experimentation/v1"
)

func NewRunDetails(run *ExperimentRun, config *ExperimentConfig, transformer *Transformer, now time.Time) (*experimentation.ExperimentRunDetails, error) {
	status := timesToStatus(run.StartTime, run.EndTime, run.CancellationTime, now)

	runProperties, err := run.CreateProperties(time.Now())
	if err != nil {
		return nil, err
	}

	configProperties, err := config.CreateProperties()
	if err != nil {
		return nil, err
	}

	transformerProperties, err := transformer.CreateProperties(run, config)
	if err != nil {
		return nil, err
	}

	properties := append(runProperties, configProperties...)
	properties = append(properties, transformerProperties...)

	if err != nil {
		return nil, err
	}

	return &experimentation.ExperimentRunDetails{
		RunId:      run.Id,
		Status:     status,
		Properties: &experimentation.PropertiesList{Items: properties},
		Config:     config.Config,
	}, nil
}

func timesToStatus(startTime time.Time, endTime *time.Time, cancellationTime *time.Time, now time.Time) experimentation.Experiment_Status {
	if cancellationTime != nil {
		if cancellationTime.After(startTime) {
			if endTime != nil {
				return experimentation.Experiment_STATUS_STOPPED
			} else {
				return experimentation.Experiment_STATUS_COMPLETED
			}
		} else {
			return experimentation.Experiment_STATUS_CANCELED
		}
	} else {
		if now.Before(startTime) {
			return experimentation.Experiment_STATUS_SCHEDULED
		} else if now.After(startTime) && (endTime == nil || now.Before(*endTime)) {
			return experimentation.Experiment_STATUS_RUNNING
		}
		return experimentation.Experiment_STATUS_COMPLETED
	}
}

func statusToString(status experimentation.Experiment_Status) string {
	switch status {
	case experimentation.Experiment_STATUS_UNSPECIFIED:
		return "Unspecified"
	case experimentation.Experiment_STATUS_SCHEDULED:
		return "Scheduled"
	case experimentation.Experiment_STATUS_RUNNING:
		return "Running"
	case experimentation.Experiment_STATUS_COMPLETED:
		return "Completed"
	case experimentation.Experiment_STATUS_CANCELED:
		return "Canceled"
	case experimentation.Experiment_STATUS_STOPPED:
		return "Stopped"
	default:
		return status.String()
	}
}

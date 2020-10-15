package experimentstore

import (
	"database/sql"
	"time"

	experimentation "github.com/lyft/clutch/backend/api/chaos/experimentation/v1"
)

func NewRunDetails(run *ExperimentRun, config *ExperimentConfig, transformer *Transformer, now time.Time) (*experimentation.ExperimentRunDetails, error) {
	status := timesToStatus(run.startTime, run.endTime, run.cancellationTime, now)

	runProperties, err := run.CreateProperties(time.Now())
	if err != nil {
		return nil, err
	}

	configProperties, err := config.CreateProperties(transformer)
	if err != nil {
		return nil, err
	}

	properties := append(runProperties, configProperties...)

	if err != nil {
		return nil, err
	}

	return &experimentation.ExperimentRunDetails{
		RunId:      run.id,
		Status:     status,
		Properties: &experimentation.PropertiesList{Items: properties},
		Config:     config.Config,
	}, nil
}

func timesToStatus(startTime time.Time, endTime sql.NullTime, cancellationTime sql.NullTime, now time.Time) experimentation.Experiment_Status {
	if cancellationTime.Valid {
		if cancellationTime.Time.After(startTime) {
			if endTime.Valid {
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
		} else if now.After(startTime) && (!endTime.Valid || now.Before(endTime.Time)) {
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

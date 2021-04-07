package experimentstore

import experimentation "github.com/lyft/clutch/backend/api/chaos/experimentation/v1"

func StatusToString(status experimentation.Experiment_Status) string {
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
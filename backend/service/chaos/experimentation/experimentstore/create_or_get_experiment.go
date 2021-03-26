package experimentstore

import experimentation "github.com/lyft/clutch/backend/api/chaos/experimentation/v1"

type CreateOrGetExperiment struct {
	Experiment *experimentation.Experiment
	Origin     experimentation.CreateOrGetExperimentResponse_Origin
}

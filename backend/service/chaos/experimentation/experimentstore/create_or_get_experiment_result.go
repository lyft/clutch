package experimentstore

import experimentationv1 "github.com/lyft/clutch/backend/api/chaos/experimentation/v1"

type CreateOrGetExperimentResult struct {
	Experiment *Experiment
	Origin     experimentationv1.CreateOrGetExperimentResponse_Origin
}

package experimentstore

import experimentationv1 "github.com/lyft/clutch/backend/api/chaos/experimentation/v1"

type CreateOrGetExperiment struct {
	Experiment *experimentationv1.Experiment
	Origin     experimentationv1.CreateOrGetExperimentResponse_Origin
}

package experimentstoremock

import (
	"context"

	"github.com/golang/protobuf/jsonpb"
	"github.com/golang/protobuf/ptypes/any"
	"github.com/uber-go/tally"
	"go.uber.org/zap"

	experimentation "github.com/lyft/clutch/backend/api/chaos/experimentation/v1"
	"github.com/lyft/clutch/backend/service"
)

type mockExperimentStore struct {
}

func NewMock(_ *any.Any, _ *zap.Logger, _ tally.Scope) (service.Service, error) {
	return &mockExperimentStore{}, nil
}
func (fs *mockExperimentStore) CreateExperiments(context.Context, []*experimentation.Experiment) error {
	// does nothing
	return nil
}

// DeleteExperiments deletes the specified experiments from the store.
func (fs *mockExperimentStore) DeleteExperiments(context.Context, []uint64) error {
	// does nothing
	return nil
}

// GetExperiments gets all experiments
func (fs *mockExperimentStore) GetExperiments(context.Context) ([]*experimentation.Experiment, error) {
	var experiments []*experimentation.Experiment
	var experiment experimentation.Experiment

	details :=
		`{"abort":{
            "clusterPair":{"downstreamCluster":"dCluster","upstreamCluster":"uCluster"},
            "percent":100,
            "httpStatus":401}}`

	spec := &experimentation.TestSpecification{}
	err := jsonpb.UnmarshalString(details, spec)
	if err != nil {
		return nil, err
	}

	experiment.Id = 1
	experiment.TestSpecification = spec

	experiments = append(experiments, &experiment)

	return experiments, nil
}

func (fs *mockExperimentStore) Close() {
}

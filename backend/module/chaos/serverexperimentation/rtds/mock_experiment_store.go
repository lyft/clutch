package rtds

import (
	"context"
	"github.com/lyft/clutch/backend/service/chaos/experimentation/experimentstore"
	"time"

	"github.com/golang/protobuf/ptypes/any"

	experimentation "github.com/lyft/clutch/backend/api/chaos/experimentation/v1"
)

type mockExperimentStore struct {
	getExperimentArguments getExperimentArguments
}

type getExperimentArguments struct {
	configType string
}

func (fs *mockExperimentStore) CreateExperiment(ctx context.Context, config *any.Any, startTime *time.Time, endTime *time.Time) (*experimentation.Experiment, error) {
	return nil, nil
}

func (fs *mockExperimentStore) CancelExperimentRun(ctx context.Context, id uint64) error {
	return nil
}

func (fs *mockExperimentStore) GetExperiments(ctx context.Context, configTypes string, status experimentation.GetExperimentsRequest_Status) ([]*experimentation.Experiment, error) {
	fs.getExperimentArguments = getExperimentArguments{configType: configTypes}
	return nil, nil
}

func (fs *mockExperimentStore) GetListView(ctx context.Context) ([]*experimentation.ListViewItem, error) {
	return nil, nil
}

func (fs *mockExperimentStore) GetListView(ctx context.Context) ([]*experimentation.ListViewItem, error) {
	return nil, nil
}

func (fs *mockExperimentStore) GetExperimentRunDetails(ctx context.Context, id uint64) (*experimentation.ExperimentRunDetails, error) {
	return nil, nil
}

func (fs *mockExperimentStore) RegisterTransformation(typeUrl string, transformation func(config *experimentstore.ExperimentConfig)([]*experimentation.Property, error)) {
}


func (fs *mockExperimentStore) Close() {
}

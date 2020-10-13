package rtds

import (
	"context"
	"fmt"
	"time"

	"github.com/golang/protobuf/ptypes/any"

	experimentation "github.com/lyft/clutch/backend/api/chaos/experimentation/v1"
	"github.com/lyft/clutch/backend/service/chaos/experimentation/experimentstore"
)

type mockStorer struct {
	getExperimentArguments getExperimentArguments
}

type getExperimentArguments struct {
	configType string
}

func (fs *mockStorer) CreateExperiment(ctx context.Context, config *any.Any, startTime *time.Time, endTime *time.Time) (*experimentation.Experiment, error) {
	return nil, nil
}

func (fs *mockStorer) CancelExperimentRun(ctx context.Context, id uint64) error {
	return nil
}

func (fs *mockStorer) GetExperiments(ctx context.Context, configTypes string, status experimentation.GetExperimentsRequest_Status) ([]*experimentation.Experiment, error) {
	fs.getExperimentArguments = getExperimentArguments{configType: configTypes}
	return nil, nil
}

func (fs *mockStorer) GetListView(ctx context.Context) ([]*experimentation.ListViewItem, error) {
	return nil, nil
}

func (fs *mockStorer) GetExperimentRunDetails(ctx context.Context, id uint64) (*experimentation.ExperimentRunDetails, error) {
	return nil, nil
}

func (fs *mockStorer) RegisterTransformation(transformation experimentstore.Transformation) error {
	return fmt.Errorf("Not implemented")
}

func (fs *mockStorer) Close() {
}

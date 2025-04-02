package experimentstoremock

import (
	"context"
	"fmt"

	"github.com/uber-go/tally/v4"
	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/anypb"

	experimentation "github.com/lyft/clutch/backend/api/chaos/experimentation/v1"
	"github.com/lyft/clutch/backend/service"
	"github.com/lyft/clutch/backend/service/chaos/experimentation/experimentstore"
)

func NewMock(_ *anypb.Any, _ *zap.Logger, _ tally.Scope) (service.Service, error) {
	return &MockStorer{}, nil
}

type MockStorer struct {
	GetExperimentArguments getExperimentArguments
}

var _ experimentstore.Storer = (*MockStorer)(nil)

type getExperimentArguments struct {
	ConfigType string
}

func (fs *MockStorer) CreateExperiment(ctx context.Context, es *experimentstore.ExperimentSpecification) (*experimentstore.Experiment, error) {
	return nil, nil
}

func (fs *MockStorer) CreateOrGetExperiment(ctx context.Context, es *experimentstore.ExperimentSpecification) (*experimentstore.CreateOrGetExperimentResult, error) {
	return nil, nil
}

func (fs *MockStorer) CancelExperimentRun(ctx context.Context, runId string, reason string) error {
	return nil
}

func (fs *MockStorer) GetExperiments(ctx context.Context, configType string, status experimentation.GetExperimentsRequest_Status) ([]*experimentstore.Experiment, error) {
	fs.GetExperimentArguments = getExperimentArguments{ConfigType: configType}
	return nil, nil
}

func (fs *MockStorer) GetListView(ctx context.Context) ([]*experimentation.ListViewItem, error) {
	return nil, nil
}

func (fs *MockStorer) GetExperimentRunDetails(ctx context.Context, runId string) (*experimentation.ExperimentRunDetails, error) {
	return nil, nil
}

func (fs *MockStorer) RegisterTransformation(transformation experimentstore.Transformation) error {
	return fmt.Errorf("Not implemented")
}

func (fs *MockStorer) Close() {
}

var _ experimentstore.Storer = &MockStorer{}

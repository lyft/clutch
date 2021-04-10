package experimentstore

import (
	"time"

	experimentation "github.com/lyft/clutch/backend/api/chaos/experimentation/v1"
)

func NewRunDetails(run *ExperimentRun, config *ExperimentConfig, transformer *Transformer, now time.Time) (*experimentation.ExperimentRunDetails, error) {
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
		Status:     run.Status(now),
		Properties: &experimentation.PropertiesList{Items: properties},
		Config:     config.Config,
	}, nil
}

package experimentstore

import (
	"time"

	experimentation "github.com/lyft/clutch/backend/api/chaos/experimentation/v1"
)

func NewRunListView(run *ExperimentRun, config *ExperimentConfig, transformer *Transformer, now time.Time) (*experimentation.ListViewItem, error) {
	runProperties, err := run.CreateProperties(now)
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
	propertiesMapItems := make(map[string]*experimentation.Property)
	for _, p := range properties {
		propertiesMapItems[p.Id] = p
	}

	return &experimentation.ListViewItem{
		Id:         run.Id,
		Properties: &experimentation.PropertiesMap{Items: propertiesMapItems},
	}, nil
}

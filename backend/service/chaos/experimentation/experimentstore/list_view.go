package experimentstore

import (
	experimentation "github.com/lyft/clutch/backend/api/chaos/experimentation/v1"
	"time"
)

func NewRunListView(run *ExperimentRun, config *ExperimentConfig, transformer *Transformer, now time.Time) (*experimentation.ListViewItem, error) {
	runProperties, err := run.CreateProperties(time.Now())
	if err != nil {
		return nil, err
	}

	configProperties, err := config.CreateProperties(transformer)
	if err != nil {
		return nil, err
	}

	properties := append(runProperties, configProperties...)
	propertiesMapItems := make(map[string]*experimentation.Property)
	for _, p := range properties {
		propertiesMapItems[p.Id] = p
	}

	return &experimentation.ListViewItem{
		Identifier: run.id,
		Properties: &experimentation.PropertiesMap{Items: propertiesMapItems},
	}, nil
}

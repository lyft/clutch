package experimentstore

import (
	experimentation "github.com/lyft/clutch/backend/api/chaos/experimentation/v1"
)

type Transformer struct {
	nameToConfigTransformMap map[string]func(*ExperimentConfig)([]*experimentation.Property, error)
}

func (c *Transformer) CreateProperties(config *ExperimentConfig) ([]*experimentation.Property, error) {
	transform, exists := c.nameToConfigTransformMap[config.Config.TypeUrl]
	if !exists {
		return []*experimentation.Property{}, nil
	}

	return transform(config)
}

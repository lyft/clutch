package experimentstore

import (
	experimentation "github.com/lyft/clutch/backend/api/chaos/experimentation/v1"
)

type Transformation struct {
	ConfigTypeUrl   string
	ConfigTransform func(config *ExperimentConfig) ([]*experimentation.Property, error)
}

type Transformer struct {
	nameToConfigTransformMap map[string]*Transformation
}

func NewTransformer() Transformer {
	t := Transformer{map[string]*Transformation{}}
	t.nameToConfigTransformMap = map[string]*Transformation{}
	return t
}

func (t *Transformer) Register(transformation Transformation) error {
	t.nameToConfigTransformMap[transformation.ConfigTypeUrl] = &transformation
	return nil
}

func (t *Transformer) CreateProperties(config *ExperimentConfig) ([]*experimentation.Property, error) {
	transform, exists := t.nameToConfigTransformMap[config.Config.TypeUrl]
	if !exists {
		return []*experimentation.Property{}, nil
	}

	return transform.ConfigTransform(config)
}

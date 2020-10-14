package experimentstore

import (
	"go.uber.org/zap"

	experimentation "github.com/lyft/clutch/backend/api/chaos/experimentation/v1"
)

type Transformation struct {
	ConfigTypeUrl   string
	ConfigTransform func(config *ExperimentConfig) ([]*experimentation.Property, error)
}

type Transformer struct {
	logger                   *zap.SugaredLogger
	nameToConfigTransformMap map[string]*Transformation
}

func NewTransformer(logger *zap.SugaredLogger) Transformer {
	t := Transformer{logger, map[string]*Transformation{}}
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

	properties, err := transform.ConfigTransform(config)
	if err != nil {
		t.logger.Errorw("error while creating properties", "error", err, "config", config)
	}

	return properties, err
}

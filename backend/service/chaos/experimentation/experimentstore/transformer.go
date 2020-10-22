package experimentstore

import (
	"go.uber.org/zap"

	experimentation "github.com/lyft/clutch/backend/api/chaos/experimentation/v1"
)

type Transformation struct {
	ConfigTypeUrl   string
	ConfigTransform func(config *ExperimentConfig) ([]*experimentation.Property, error)
	RunTransform    func(run *ExperimentRun, config *ExperimentConfig) ([]*experimentation.Property, error)
}

type Transformer struct {
	logger             *zap.SugaredLogger
	nameToTransformMap map[string][]*Transformation
}

func NewTransformer(logger *zap.SugaredLogger) Transformer {
	t := Transformer{logger, map[string][]*Transformation{}}
	t.nameToTransformMap = map[string][]*Transformation{}
	return t
}

func (tr *Transformer) Register(transformation Transformation) error {
	tr.nameToTransformMap[transformation.ConfigTypeUrl] = append(tr.nameToTransformMap[transformation.ConfigTypeUrl], &transformation)
	return nil
}

func (tr *Transformer) CreateProperties(run *ExperimentRun, config *ExperimentConfig) ([]*experimentation.Property, error) {
	transformations, exists := tr.nameToTransformMap[config.Config.TypeUrl]
	if !exists {
		return []*experimentation.Property{}, nil
	}

	var properties = []*experimentation.Property{}
	for _, t := range transformations {
		if t.RunTransform == nil {
			continue
		}

		currentTransformationProperties, err := t.RunTransform(run, config)
		if err != nil {
			tr.logger.Errorw("error while creating properties out of config",
				"error", err, "config", config, "run", run)
			return nil, err
		}

		properties = append(properties, currentTransformationProperties...)
	}

	return properties, nil
}

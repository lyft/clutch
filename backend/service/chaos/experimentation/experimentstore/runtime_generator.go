package experimentstore

import (
	"github.com/lyft/clutch/backend/module/chaos/serverexperimentation/xds"
	"go.uber.org/zap"

	experimentation "github.com/lyft/clutch/backend/api/chaos/experimentation/v1"
)

type RuntimeGeneration struct {
	ConfigTypeUrl   string
	GetEnforcingCluster func(experiment *experimentation.Experiment, logger *zap.SugaredLogger) (string, error)
	RuntimeKeysGeneration func(experiment *experimentation.Experiment, rtdsConfig *xds.RTDSConfig, logger *zap.SugaredLogger) ([]*RuntimeKeyValue, error)
}

type RuntimeGenerator struct {
	logger             *zap.SugaredLogger
	nameToGenerationMap map[string]*RuntimeGeneration
}

type RuntimeKeyValue struct {
	Key string
	Value uint32
}

func NewRuntimeGenerator(logger *zap.SugaredLogger) RuntimeGenerator {
	runtimeGenerator := RuntimeGenerator{logger, map[string]*RuntimeGeneration{}}
	runtimeGenerator.nameToGenerationMap = map[string]*RuntimeGeneration{}
	return runtimeGenerator
}

func (rg *RuntimeGenerator) Register(runtimeGeneration RuntimeGeneration) error {
	rg.nameToGenerationMap[runtimeGeneration.ConfigTypeUrl] = &runtimeGeneration
	return nil
}


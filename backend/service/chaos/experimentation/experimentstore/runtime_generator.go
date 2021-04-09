package experimentstore

import (
	"go.uber.org/zap"

	experimentation "github.com/lyft/clutch/backend/api/chaos/experimentation/v1"
)

type RuntimeGeneration struct {
	ConfigTypeUrl   string
	GetEnforcingCluster func(experiment *experimentation.Experiment, logger *zap.SugaredLogger) (string, error)
	RuntimeKeysGeneration func(experiment *experimentation.Experiment, runtimePrefixes *RuntimePrefixes, logger *zap.SugaredLogger) ([]*RuntimeKeyValue, error)
}

type RuntimeGenerator struct {
	logger             *zap.SugaredLogger
	nameToGenerationMap map[string]*RuntimeGeneration
}

type RuntimePrefixes struct {
	IngressPrefix string
	EgressPrefix  string
	RedisPrefix   string
}

type RuntimeKeyValue struct {
	Key   string
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


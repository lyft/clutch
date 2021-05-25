package xds

import (
	"fmt"

	redisexperimentation "github.com/lyft/clutch/backend/api/chaos/redisexperimentation/v1"
	"github.com/lyft/clutch/backend/module/chaos/experimentation/xds"
	"github.com/lyft/clutch/backend/service/chaos/experimentation/experimentstore"
)

const (
	redisErrorPercentage   = `%s.%s.error.error_percent`
	redisLatencyPercentage = `%s.%s.delay.fixed_delay_percent`
)

type RTDSFaultsGenerator struct {
	FaultRuntimePrefix string
}

func (g RTDSFaultsGenerator) GenerateResource(experiment *experimentstore.Experiment) (*xds.RTDSResource, error) {
	redisFaultConfig, ok := experiment.Config.Message.(*redisexperimentation.FaultConfig)
	if !ok {
		return nil, fmt.Errorf("RTDS redis faults generator couldn't generate faults due to config type casting failure (experiment run ID: %s, experiment config ID %s)", experiment.Run.Id, experiment.Config.Id)
	}

	downstreamCluster := redisFaultConfig.GetFaultTargeting().GetDownstreamCluster().GetName()
	upstreamCluster := redisFaultConfig.GetFaultTargeting().GetUpstreamCluster().GetName()

	percentageKey, percentageValue, err := g.createRedisRuntimeKeys(upstreamCluster, redisFaultConfig)
	if err != nil {
		return nil, err
	}

	return xds.NewRTDSResource(downstreamCluster, []*xds.RuntimeKeyValue{
		{
			Key:   percentageKey,
			Value: percentageValue,
		},
	})
}

func (g RTDSFaultsGenerator) createRedisRuntimeKeys(upstreamCluster string, redisFaultConfig *redisexperimentation.FaultConfig) (string, uint32, error) {
	var percentageKey string
	var percentageValue uint32

	switch redisFaultConfig.GetFault().(type) {
	case *redisexperimentation.FaultConfig_ErrorFault:
		percentageValue = redisFaultConfig.GetErrorFault().GetPercentage().GetPercentage()
		percentageKey = fmt.Sprintf(redisErrorPercentage, g.FaultRuntimePrefix, upstreamCluster)
	case *redisexperimentation.FaultConfig_LatencyFault:
		percentageValue = redisFaultConfig.GetLatencyFault().GetPercentage().GetPercentage()
		percentageKey = fmt.Sprintf(redisLatencyPercentage, g.FaultRuntimePrefix, upstreamCluster)
	default:
		return "", 0, fmt.Errorf("unknown fault type %v", redisFaultConfig)
	}

	return percentageKey, percentageValue, nil
}

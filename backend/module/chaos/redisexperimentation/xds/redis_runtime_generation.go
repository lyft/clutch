package xds

import (
	"fmt"

	"go.uber.org/zap"

	experimentation "github.com/lyft/clutch/backend/api/chaos/experimentation/v1"
	redisexperimentation "github.com/lyft/clutch/backend/api/chaos/redisexperimentation/v1"
	"github.com/lyft/clutch/backend/service/chaos/experimentation/experimentstore"
)

const (
	// REDIS FAULT
	// a given downstream service to the upstream redis cluster faults
	RedisLatencyPercentageWithDownstream = `%s.%s.delay.fixed_delay_percent`
	RedisErrorPercentageWithDownstream   = `%s.%s.delay.error_percent`
)

func GetEnforcingCluster(experiment *experimentation.Experiment, logger *zap.SugaredLogger) (string, error) {
	redisFaultConfig := &redisexperimentation.FaultConfig{}
	if !maybeUnmarshalRedisFaultTest(experiment, redisFaultConfig) {
		logger.Errorw("Invalid redis fault config", "config", redisFaultConfig)
		return "", nil
	}

	upstreamCluster, _ := getRedisClusterPair(redisFaultConfig)
	return upstreamCluster, nil
}

func RuntimeKeysGeneration(experiment *experimentation.Experiment, runtimePrefixes *experimentstore.RuntimePrefixes, logger *zap.SugaredLogger) ([]*experimentstore.RuntimeKeyValue, error) {
	redisFaultConfig := &redisexperimentation.FaultConfig{}
	if !maybeUnmarshalRedisFaultTest(experiment, redisFaultConfig) {
		logger.Errorw("Invalid redis fault config", "config", redisFaultConfig)
		return nil, nil
	}

	upstreamCluster, downstreamCluster := getRedisClusterPair(redisFaultConfig)

	percentageKey, percentageValue, err := createRedisRuntimeKeys(upstreamCluster, redisFaultConfig, runtimePrefixes)
	if err != nil {
		logger.Errorw("Unable to create runtime keys", "config", redisFaultConfig)
		return nil, nil
	}

	logger.Debugw("Fault details",
		"upstream_cluster", upstreamCluster,
		"downstream_cluster", downstreamCluster,
		"fault_type", percentageKey,
		"percentage", percentageValue)

	return []*experimentstore.RuntimeKeyValue{
		{
			Key:   percentageKey,
			Value: percentageValue,
		},
	}, nil
}

func createRedisRuntimeKeys(upstreamCluster string, redisFaultConfig *redisexperimentation.FaultConfig, runtimePrefixes *experimentstore.RuntimePrefixes) (string, uint32, error) {
	var percentageKey string
	var percentageValue uint32

	redisPrefix := runtimePrefixes.RedisPrefix

	switch redisFaultConfig.GetFault().(type) {
	case *redisexperimentation.FaultConfig_ErrorFault:
		percentageValue = redisFaultConfig.GetErrorFault().GetPercentage().GetPercentage()

		percentageKey = fmt.Sprintf(RedisErrorPercentageWithDownstream, redisPrefix, upstreamCluster)

	case *redisexperimentation.FaultConfig_LatencyFault:
		percentageValue = redisFaultConfig.GetLatencyFault().GetPercentage().GetPercentage()
		percentageKey = fmt.Sprintf(RedisLatencyPercentageWithDownstream, redisPrefix, upstreamCluster)

	default:
		return "", 0, fmt.Errorf("unknown fault type %v", redisFaultConfig)
	}

	return percentageKey, percentageValue, nil
}

func maybeUnmarshalRedisFaultTest(experiment *experimentation.Experiment, redisFaultConfig *redisexperimentation.FaultConfig) bool {
	err := experiment.GetConfig().UnmarshalTo(redisFaultConfig)
	if err != nil {
		return false
	}

	switch redisFaultConfig.GetFault().(type) {
	case *redisexperimentation.FaultConfig_ErrorFault:
		return true
	case *redisexperimentation.FaultConfig_LatencyFault:
		return true
	default:
		return false
	}
}

func getRedisClusterPair(redisFaultConfig *redisexperimentation.FaultConfig) (string, string) {
	var downstream, upstream string

	downstream = redisFaultConfig.GetFaultTargeting().GetDownstreamCluster().GetName()
	upstream = redisFaultConfig.GetFaultTargeting().GetUpstreamCluster().GetName()

	return upstream, downstream
}

package xds

import (
	"time"

	gcpRuntimeServiceV3 "github.com/envoyproxy/go-control-plane/envoy/service/runtime/v3"
	gcpTypes "github.com/envoyproxy/go-control-plane/pkg/cache/types"
	pstruct "github.com/golang/protobuf/ptypes/struct"
	"go.uber.org/zap"

	experimentation "github.com/lyft/clutch/backend/api/chaos/experimentation/v1"
	"github.com/lyft/clutch/backend/service/chaos/experimentation/experimentstore"
)

type RTDSConfig struct {
	// Name of the RTDS layer in Envoy config i.e. envoy.yaml
	layerName string

	// Runtime prefix for ingress faults
	ingressPrefix string

	// Runtime prefix for egress faults
	egressPrefix string

	// Runtime prefix for redis faults
	RedisPrefix string
}

func generateRTDSResource(experiments []*experimentation.Experiment, storer experimentstore.Storer, rtdsConfig *RTDSConfig, ttl *time.Duration, logger *zap.SugaredLogger) []gcpTypes.ResourceWithTtl {
	var fieldMap = map[string]*pstruct.Value{}

	// No experiments meaning clear all experiments for the given upstream cluster
	for _, experiment := range experiments {
		// get the runtime generation for each experiment
		runtimeGeneration, err := storer.GetRuntimeGenerator(experiment.GetConfig().GetTypeUrl())
		if err != nil {
			logger.Errorw("No runtime generation registered to config", "config", experiment.GetConfig())
			continue
		}

		// get runtime values from experiment types runtime generation logic
		runtimeKeyValues, err := runtimeGeneration.RuntimeKeysGeneration(experiment, rtdsConfig, logger)
		if err != nil {
			logger.Errorw("Unable to create RTDS runtime keys", "config", experiment)
			continue
		}

		// apply list of runtime key values to fieldMap
		for _, runtimeKeyValue := range runtimeKeyValues {
			fieldMap[runtimeKeyValue.Key] = &pstruct.Value{
				Kind: &pstruct.Value_NumberValue{
					NumberValue: float64(runtimeKeyValue.Value),
				},
			}
		}
	}

	runtimeLayer := &pstruct.Struct{
		Fields: fieldMap,
	}

	// We only want to set a TTL for non-empty layers. This ensures that we don't spam clients with heartbeat responses
	// unless they have an active runtime override.
	var resourceTTL *time.Duration
	if len(runtimeLayer.Fields) > 0 {
		resourceTTL = ttl
	}

	return []gcpTypes.ResourceWithTtl{{
		Resource: &gcpRuntimeServiceV3.Runtime{
			Name:  rtdsConfig.layerName,
			Layer: runtimeLayer,
		},
		Ttl: resourceTTL,
	}}
}

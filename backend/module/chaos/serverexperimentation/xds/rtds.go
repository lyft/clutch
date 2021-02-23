package xds

import (
	"fmt"
	"time"

	gcpRuntimeServiceV3 "github.com/envoyproxy/go-control-plane/envoy/service/runtime/v3"
	gcpTypes "github.com/envoyproxy/go-control-plane/pkg/cache/types"
	pstruct "github.com/golang/protobuf/ptypes/struct"
	"go.uber.org/zap"

	experimentation "github.com/lyft/clutch/backend/api/chaos/experimentation/v1"
	serverexperimentation "github.com/lyft/clutch/backend/api/chaos/serverexperimentation/v1"
)

const (
	// INTERNAL FAULT
	// a given downstream service to a given upstream service faults
	LatencyPercentageWithDownstream = `%s.%s.delay.fixed_delay_percent`
	LatencyDurationWithDownstream   = `%s.%s.delay.fixed_duration_ms`
	HTTPPercentageWithDownstream    = `%s.%s.abort.abort_percent`
	HTTPStatusWithDownstream        = `%s.%s.abort.http_status`

	// all downstream service to a given upstream faults
	LatencyPercentageWithoutDownstream = `%s.delay.fixed_delay_percent`
	LatencyDurationWithoutDownstream   = `%s.delay.fixed_duration_ms`
	HTTPPercentageWithoutDownstream    = `%s.abort.abort_percent`
	HTTPStatusWithoutDownstream        = `%s.abort.http_status`

	// EXTERNAL FAULT
	// a given downstream service to a given external upstream faults
	LatencyPercentageForExternal = `%s.%s.delay.fixed_delay_percent`
	LatencyDurationForExternal   = `%s.%s.delay.fixed_duration_ms`
	HTTPPercentageForExternal    = `%s.%s.abort.abort_percent`
	HTTPStatusForExternal        = `%s.%s.abort.http_status`
)

type RTDSConfig struct {
	// Name of the RTDS layer in Envoy config i.e. envoy.yaml
	layerName string

	// Runtime prefix for ingress faults
	ingressPrefix string

	// Runtime prefix for egress faults
	egressPrefix string
}

func generateRTDSResource(experiments []*experimentation.Experiment, rtdsConfig *RTDSConfig, ttl *time.Duration, logger *zap.SugaredLogger) []gcpTypes.ResourceWithTtl {
	var fieldMap = map[string]*pstruct.Value{}

	// No experiments meaning clear all experiments for the given upstream cluster
	for _, experiment := range experiments {
		httpFaultConfig := &serverexperimentation.HTTPFaultConfig{}
		if !maybeUnmarshalFaultTest(experiment, httpFaultConfig) {
			continue
		}

		upstreamCluster, downstreamCluster, err := getClusterPair(httpFaultConfig)
		if err != nil {
			logger.Errorw("Invalid http fault config", "config", httpFaultConfig)
			continue
		}

		percentageKey, percentageValue, faultKey, faultValue, err := createRuntimeKeys(upstreamCluster, downstreamCluster, httpFaultConfig, rtdsConfig)
		if err != nil {
			logger.Errorw("Unable to create RTDS runtime keys", "config", httpFaultConfig)
			continue
		}

		fieldMap[percentageKey] = &pstruct.Value{
			Kind: &pstruct.Value_NumberValue{
				NumberValue: float64(percentageValue),
			},
		}

		fieldMap[faultKey] = &pstruct.Value{
			Kind: &pstruct.Value_NumberValue{
				NumberValue: float64(faultValue),
			},
		}

		logger.Debugw("Fault details",
			"upstream_cluster", upstreamCluster,
			"downstream_cluster", downstreamCluster,
			"fault_type", faultKey,
			"percentage", percentageValue,
			"value", faultValue,
			"fault_enforcer", getEnforcer(httpFaultConfig),
		)
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

func createRuntimeKeys(upstreamCluster string, downstreamCluster string, httpFaultConfig *serverexperimentation.HTTPFaultConfig, rtdsConfig *RTDSConfig) (string, uint32, string, uint32, error) {
	var percentageKey string
	var percentageValue uint32
	var faultKey string
	var faultValue uint32

	ingressPrefix := rtdsConfig.ingressPrefix
	egressPrefix := rtdsConfig.egressPrefix

	switch httpFaultConfig.GetFault().(type) {
	case *serverexperimentation.HTTPFaultConfig_AbortFault:
		abort := httpFaultConfig.GetAbortFault()
		percentageValue = abort.GetPercentage().GetPercentage()
		faultValue = abort.GetAbortStatus().GetHttpStatusCode()

		switch httpFaultConfig.GetFaultTargeting().GetEnforcer().(type) {
		case *serverexperimentation.FaultTargeting_DownstreamEnforcing:
			// Abort External Fault
			percentageKey = fmt.Sprintf(HTTPPercentageForExternal, egressPrefix, upstreamCluster)
			faultKey = fmt.Sprintf(HTTPStatusForExternal, egressPrefix, upstreamCluster)

		case *serverexperimentation.FaultTargeting_UpstreamEnforcing:
			// Abort Internal Fault for all downstream services
			if downstreamCluster == "" {
				percentageKey = fmt.Sprintf(HTTPPercentageWithoutDownstream, ingressPrefix)
				faultKey = fmt.Sprintf(HTTPStatusWithoutDownstream, ingressPrefix)
			} else {
				// Abort Internal Fault for a given downstream services
				percentageKey = fmt.Sprintf(HTTPPercentageWithDownstream, ingressPrefix, downstreamCluster)
				faultKey = fmt.Sprintf(HTTPStatusWithDownstream, ingressPrefix, downstreamCluster)
			}

		default:
			return "", 0, "", 0, fmt.Errorf("unknown enforcer %v", httpFaultConfig)
		}

	case *serverexperimentation.HTTPFaultConfig_LatencyFault:
		latency := httpFaultConfig.GetLatencyFault()
		percentageValue = latency.GetPercentage().GetPercentage()
		faultValue = latency.GetLatencyDuration().GetFixedDurationMs()

		switch httpFaultConfig.GetFaultTargeting().GetEnforcer().(type) {
		case *serverexperimentation.FaultTargeting_DownstreamEnforcing:
			// Latency External Fault
			percentageKey = fmt.Sprintf(LatencyPercentageForExternal, egressPrefix, upstreamCluster)
			faultKey = fmt.Sprintf(LatencyDurationForExternal, egressPrefix, upstreamCluster)

		case *serverexperimentation.FaultTargeting_UpstreamEnforcing:
			// Latency Internal Fault for all downstream services
			if downstreamCluster == "" {
				percentageKey = fmt.Sprintf(LatencyPercentageWithoutDownstream, ingressPrefix)
				faultKey = fmt.Sprintf(LatencyDurationWithoutDownstream, ingressPrefix)
			} else {
				// Latency Internal Fault for a given downstream services
				percentageKey = fmt.Sprintf(LatencyPercentageWithDownstream, ingressPrefix, downstreamCluster)
				faultKey = fmt.Sprintf(LatencyDurationWithDownstream, ingressPrefix, downstreamCluster)
			}

		default:
			return "", 0, "", 0, fmt.Errorf("unknown enforcer %v", httpFaultConfig)
		}

	default:
		return "", 0, "", 0, fmt.Errorf("unknown fault type %v", httpFaultConfig)
	}

	return percentageKey, percentageValue, faultKey, faultValue, nil
}

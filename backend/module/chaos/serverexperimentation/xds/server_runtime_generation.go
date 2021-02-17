package xds

import (
	"fmt"
	"go.uber.org/zap"

	experimentation "github.com/lyft/clutch/backend/api/chaos/experimentation/v1"
	serverexperimentation "github.com/lyft/clutch/backend/api/chaos/serverexperimentation/v1"
	"github.com/lyft/clutch/backend/service/chaos/experimentation/experimentstore"
)

const (
	// INGRESS FAULT
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

	// EGRESS FAULT
	LatencyPercentageForEgress = `%s.%s.delay.fixed_delay_percent`
	LatencyDurationForEgress   = `%s.%s.delay.fixed_duration_ms`
	HTTPPercentageForEgress    = `%s.%s.abort.abort_percent`
	HTTPStatusForEgress        = `%s.%s.abort.http_status`
)

func GetEnforcingCluster(experiment *experimentation.Experiment, logger *zap.SugaredLogger) (string, error) {
	httpFaultConfig := &serverexperimentation.HTTPFaultConfig{}
	if !maybeUnmarshalFaultTest(experiment, httpFaultConfig) {
		return "", nil
	}

	upstreamCluster, downstreamCluster, err := getClusterPair(httpFaultConfig)
	if err != nil {
		logger.Errorw("Invalid http fault config", "config", httpFaultConfig)
		return "", err
	}

	switch httpFaultConfig.GetFaultTargeting().GetEnforcer().(type) {
	case *serverexperimentation.FaultTargeting_DownstreamEnforcing:
		return downstreamCluster, nil
	case *serverexperimentation.FaultTargeting_UpstreamEnforcing:
		return upstreamCluster, nil
	default:
		logger.Errorw("unknown enforcer %v", httpFaultConfig)
		return "", err
	}
}

func RuntimeKeysGeneration(experiment *experimentation.Experiment, rtdsConfig *RTDSConfig, logger *zap.SugaredLogger) ([]*experimentstore.RuntimeKeyValue, error) {
	httpFaultConfig := &serverexperimentation.HTTPFaultConfig{}
	if !maybeUnmarshalFaultTest(experiment, httpFaultConfig) {
		return nil, nil
	}

	upstreamCluster, downstreamCluster, err := getClusterPair(httpFaultConfig)
	if err != nil {
		logger.Errorw("Invalid http fault config", "config", httpFaultConfig)
		return nil, nil
	}

	percentageKey, percentageValue, faultKey, faultValue, err := createRuntimeKeys(upstreamCluster, downstreamCluster, httpFaultConfig, rtdsConfig)
	if err != nil {
		logger.Errorw("Unable to create RTDS runtime keys", "config", httpFaultConfig)
		return nil, nil
	}

	logger.Debugw("Fault details",
		"upstream_cluster", upstreamCluster,
		"downstream_cluster", downstreamCluster,
		"fault_type", faultKey,
		"percentage", percentageValue,
		"value", faultValue,
		"fault_enforcer", getEnforcer(httpFaultConfig),
	)

	return []*experimentstore.RuntimeKeyValue{
		{
			Key:   percentageKey,
			Value: percentageValue,
		},
		{
			Key:   faultKey,
			Value: faultValue,
		},
	}, nil
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
			// Abort Egress Fault
			percentageKey = fmt.Sprintf(HTTPPercentageForEgress, egressPrefix, upstreamCluster)
			faultKey = fmt.Sprintf(HTTPStatusForEgress, egressPrefix, upstreamCluster)

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
			// Latency Egress Fault
			percentageKey = fmt.Sprintf(LatencyPercentageForEgress, egressPrefix, upstreamCluster)
			faultKey = fmt.Sprintf(LatencyDurationForEgress, egressPrefix, upstreamCluster)

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

func maybeUnmarshalFaultTest(experiment *experimentation.Experiment, httpFaultConfig *serverexperimentation.HTTPFaultConfig) bool {
	err := experiment.GetConfig().UnmarshalTo(httpFaultConfig)
	if err != nil {
		return false
	}

	switch httpFaultConfig.GetFault().(type) {
	case *serverexperimentation.HTTPFaultConfig_AbortFault:
		return true
	case *serverexperimentation.HTTPFaultConfig_LatencyFault:
		return true
	default:
		return false
	}
}

func getClusterPair(httpFaultConfig *serverexperimentation.HTTPFaultConfig) (string, string, error) {
	var downstream, upstream string

	switch httpFaultConfig.GetFaultTargeting().GetEnforcer().(type) {
	case *serverexperimentation.FaultTargeting_DownstreamEnforcing:
		downstreamEnforcing := httpFaultConfig.GetFaultTargeting().GetDownstreamEnforcing()

		switch downstreamEnforcing.GetDownstreamType().(type) {
		case *serverexperimentation.DownstreamEnforcing_DownstreamCluster:
			downstream = downstreamEnforcing.GetDownstreamCluster().GetName()
		default:
			return "", "", fmt.Errorf("unknown downstream type of downstream enforcing %v", downstreamEnforcing.GetDownstreamType())
		}

		switch downstreamEnforcing.GetUpstreamType().(type) {
		case *serverexperimentation.DownstreamEnforcing_UpstreamCluster:
			upstream = downstreamEnforcing.GetUpstreamCluster().GetName()
		default:
			return "", "", fmt.Errorf("unknown upstream type of downstream enforcing %v", downstreamEnforcing.GetUpstreamType())
		}

	case *serverexperimentation.FaultTargeting_UpstreamEnforcing:
		upstreamEnforcing := httpFaultConfig.GetFaultTargeting().GetUpstreamEnforcing()

		switch upstreamEnforcing.GetDownstreamType().(type) {
		case *serverexperimentation.UpstreamEnforcing_DownstreamCluster:
			downstream = upstreamEnforcing.GetDownstreamCluster().GetName()
		default:
			return "", "", fmt.Errorf("unknown downstream type of upstream enforcing %v", upstreamEnforcing.GetDownstreamType())
		}

		switch upstreamEnforcing.GetUpstreamType().(type) {
		case *serverexperimentation.UpstreamEnforcing_UpstreamCluster:
			upstream = upstreamEnforcing.GetUpstreamCluster().GetName()
		case *serverexperimentation.UpstreamEnforcing_UpstreamPartialSingleCluster:
			upstream = upstreamEnforcing.GetUpstreamPartialSingleCluster().GetName()
		default:
			return "", "", fmt.Errorf("unknown upstream type of upstream enforcing %v", upstreamEnforcing.GetUpstreamType())
		}

	default:
		return "", "", fmt.Errorf("unknown enforcer %v", httpFaultConfig.GetFaultTargeting())
	}

	return upstream, downstream, nil
}

func getEnforcer(httpFaultConfig *serverexperimentation.HTTPFaultConfig) string {
	switch httpFaultConfig.GetFaultTargeting().GetEnforcer().(type) {
	case *serverexperimentation.FaultTargeting_DownstreamEnforcing:
		return "downstreamEnforcing"
	case *serverexperimentation.FaultTargeting_UpstreamEnforcing:
		return "upstreamEnforcing"
	default:
		return "unknown"
	}
}


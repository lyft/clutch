package xds

import (
	"fmt"

	serverexperimentationv1 "github.com/lyft/clutch/backend/api/chaos/serverexperimentation/v1"
	"github.com/lyft/clutch/backend/module/chaos/experimentation/xds"
	"github.com/lyft/clutch/backend/service/chaos/experimentation/experimentstore"
)

const (
	// INGRESS FAULT
	// a given downstream service to a given upstream service faults
	latencyPercentageWithDownstream = `%s.%s.delay.fixed_delay_percent`
	latencyDurationWithDownstream   = `%s.%s.delay.fixed_duration_ms`
	httpPercentageWithDownstream    = `%s.%s.abort.abort_percent`
	httpStatusWithDownstream        = `%s.%s.abort.http_status`

	// all downstream service to a given upstream faults
	latencyPercentageWithoutDownstream = `%s.delay.fixed_delay_percent`
	latencyDurationWithoutDownstream   = `%s.delay.fixed_duration_ms`
	httpPercentageWithoutDownstream    = `%s.abort.abort_percent`
	httpStatusWithoutDownstream        = `%s.abort.http_status`

	// EGRESS FAULT
	// a given downstream service to a given external upstream faults
	latencyPercentageForEgress = `%s.%s.delay.fixed_delay_percent`
	latencyDurationForEgress   = `%s.%s.delay.fixed_duration_ms`
	httpPercentageForEgress    = `%s.%s.abort.abort_percent`
	httpStatusForEgress        = `%s.%s.abort.http_status`
)

type RTDSFaultsGenerator struct {
	// Runtime prefix for ingress faults
	IngressFaultRuntimePrefix string
	// Runtime prefix for egress faults
	EgressFaultRuntimePrefix string
}

func (g RTDSFaultsGenerator) GenerateResource(experiment *experimentstore.Experiment) (*xds.RTDSResource, error) {
	httpFaultConfig, ok := experiment.Config.Message.(*serverexperimentationv1.HTTPFaultConfig)
	if !ok {
		return nil, fmt.Errorf("RTDS server faults generator couldn't generate faults due to config type casting failure (experiment run ID: %s, experiment config ID %s)", experiment.Run.Id, experiment.Config.Id)
	}

	cluster, err := GetEnforcingCluster(httpFaultConfig)
	if err != nil {
		return nil, err
	}

	upstreamCluster, downstreamCluster, err := GetClusterPair(httpFaultConfig)
	if err != nil {
		return nil, err
	}

	percentageKey, percentageValue, faultKey, faultValue, err := g.createRuntimeKeys(upstreamCluster, downstreamCluster, httpFaultConfig)
	if err != nil {
		return nil, err
	}

	keyValuePairs := make([]*xds.RuntimeKeyValue, 2)
	keyValuePairs[0] = &xds.RuntimeKeyValue{Key: percentageKey, Value: percentageValue}
	keyValuePairs[1] = &xds.RuntimeKeyValue{Key: faultKey, Value: faultValue}

	return xds.NewRTDSResource(cluster, keyValuePairs)
}

func (g RTDSFaultsGenerator) createRuntimeKeys(upstreamCluster string, downstreamCluster string, httpFaultConfig *serverexperimentationv1.HTTPFaultConfig) (string, uint32, string, uint32, error) {
	var percentageKey string
	var percentageValue uint32
	var faultKey string
	var faultValue uint32

	switch httpFaultConfig.GetFault().(type) {
	case *serverexperimentationv1.HTTPFaultConfig_AbortFault:
		abort := httpFaultConfig.GetAbortFault()
		percentageValue = abort.GetPercentage().GetPercentage()
		faultValue = abort.GetAbortStatus().GetHttpStatusCode()

		switch httpFaultConfig.GetFaultTargeting().GetEnforcer().(type) {
		case *serverexperimentationv1.FaultTargeting_DownstreamEnforcing:
			// Abort Egress Fault
			percentageKey = fmt.Sprintf(httpPercentageForEgress, g.EgressFaultRuntimePrefix, upstreamCluster)
			faultKey = fmt.Sprintf(httpStatusForEgress, g.EgressFaultRuntimePrefix, upstreamCluster)

		case *serverexperimentationv1.FaultTargeting_UpstreamEnforcing:
			// Abort Internal Fault for all downstream services
			if downstreamCluster == "" {
				percentageKey = fmt.Sprintf(httpPercentageWithoutDownstream, g.IngressFaultRuntimePrefix)
				faultKey = fmt.Sprintf(httpStatusWithoutDownstream, g.IngressFaultRuntimePrefix)
			} else {
				// Abort Internal Fault for a given downstream services
				percentageKey = fmt.Sprintf(httpPercentageWithDownstream, g.IngressFaultRuntimePrefix, downstreamCluster)
				faultKey = fmt.Sprintf(httpStatusWithDownstream, g.IngressFaultRuntimePrefix, downstreamCluster)
			}

		default:
			return "", 0, "", 0, fmt.Errorf("unknown enforcer %v", httpFaultConfig)
		}

	case *serverexperimentationv1.HTTPFaultConfig_LatencyFault:
		latency := httpFaultConfig.GetLatencyFault()
		percentageValue = latency.GetPercentage().GetPercentage()
		faultValue = latency.GetLatencyDuration().GetFixedDurationMs()

		switch httpFaultConfig.GetFaultTargeting().GetEnforcer().(type) {
		case *serverexperimentationv1.FaultTargeting_DownstreamEnforcing:
			// Latency Egress Fault
			percentageKey = fmt.Sprintf(latencyPercentageForEgress, g.EgressFaultRuntimePrefix, upstreamCluster)
			faultKey = fmt.Sprintf(latencyDurationForEgress, g.EgressFaultRuntimePrefix, upstreamCluster)

		case *serverexperimentationv1.FaultTargeting_UpstreamEnforcing:
			// Latency Internal Fault for all downstream services
			if downstreamCluster == "" {
				percentageKey = fmt.Sprintf(latencyPercentageWithoutDownstream, g.IngressFaultRuntimePrefix)
				faultKey = fmt.Sprintf(latencyDurationWithoutDownstream, g.IngressFaultRuntimePrefix)
			} else {
				// Latency Internal Fault for a given downstream services
				percentageKey = fmt.Sprintf(latencyPercentageWithDownstream, g.IngressFaultRuntimePrefix, downstreamCluster)
				faultKey = fmt.Sprintf(latencyDurationWithDownstream, g.IngressFaultRuntimePrefix, downstreamCluster)
			}

		default:
			return "", 0, "", 0, fmt.Errorf("unknown enforcer %v", httpFaultConfig)
		}

	default:
		return "", 0, "", 0, fmt.Errorf("unknown fault type %v", httpFaultConfig)
	}

	return percentageKey, percentageValue, faultKey, faultValue, nil
}

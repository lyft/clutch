package xds

import (
	"fmt"

	"github.com/golang/protobuf/ptypes/any"

	serverexperimentationv1 "github.com/lyft/clutch/backend/api/chaos/serverexperimentation/v1"
	rtdsv1 "github.com/lyft/clutch/backend/api/config/module/chaos/serverexperimentation/rtds/v1"
	"github.com/lyft/clutch/backend/module/chaos/experimentation/xds"
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
	// a given downstream service to a given external upstream faults
	LatencyPercentageForEgress = `%s.%s.delay.fixed_delay_percent`
	LatencyDurationForEgress   = `%s.%s.delay.fixed_duration_ms`
	HTTPPercentageForEgress    = `%s.%s.abort.abort_percent`
	HTTPStatusForEgress        = `%s.%s.abort.http_status`
)

type RTDSServerFaultsResourceGenerator struct {
	// Runtime prefix for ingress faults
	IngressFaultRuntimePrefix string
	// Runtime prefix for egress faults
	EgressFaultRuntimePrefix string
}

type RTDSServerFaultsResourceGeneratorFactory struct{}

func (RTDSServerFaultsResourceGeneratorFactory) Create(cfg *any.Any) (xds.RTDSResourceGenerator, error) {
	typedConfig := &rtdsv1.ServerFaultsRTDSResourceGeneratorConfig{}
	err := cfg.UnmarshalTo(typedConfig)
	if err != nil {
		return nil, err
	}

	return &RTDSServerFaultsResourceGenerator{
		IngressFaultRuntimePrefix: typedConfig.IngressFaultRuntimePrefix,
		EgressFaultRuntimePrefix:  typedConfig.EgressFaultRuntimePrefix,
	}, nil
}

func (g *RTDSServerFaultsResourceGenerator) GenerateResource(experiment *experimentstore.Experiment) (*xds.RTDSResource, error) {
	httpFaultConfig, ok := experiment.Config.Message.(*serverexperimentationv1.HTTPFaultConfig)
	if !ok {
		return xds.NewEmptyRTDSResource(""), nil
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

func (g *RTDSServerFaultsResourceGenerator) createRuntimeKeys(upstreamCluster string, downstreamCluster string, httpFaultConfig *serverexperimentationv1.HTTPFaultConfig) (string, uint32, string, uint32, error) {
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
			percentageKey = fmt.Sprintf(HTTPPercentageForEgress, g.EgressFaultRuntimePrefix, upstreamCluster)
			faultKey = fmt.Sprintf(HTTPStatusForEgress, g.EgressFaultRuntimePrefix, upstreamCluster)

		case *serverexperimentationv1.FaultTargeting_UpstreamEnforcing:
			// Abort Internal Fault for all downstream services
			if downstreamCluster == "" {
				percentageKey = fmt.Sprintf(HTTPPercentageWithoutDownstream, g.IngressFaultRuntimePrefix)
				faultKey = fmt.Sprintf(HTTPStatusWithoutDownstream, g.IngressFaultRuntimePrefix)
			} else {
				// Abort Internal Fault for a given downstream services
				percentageKey = fmt.Sprintf(HTTPPercentageWithDownstream, g.IngressFaultRuntimePrefix, downstreamCluster)
				faultKey = fmt.Sprintf(HTTPStatusWithDownstream, g.IngressFaultRuntimePrefix, downstreamCluster)
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
			percentageKey = fmt.Sprintf(LatencyPercentageForEgress, g.EgressFaultRuntimePrefix, upstreamCluster)
			faultKey = fmt.Sprintf(LatencyDurationForEgress, g.EgressFaultRuntimePrefix, upstreamCluster)

		case *serverexperimentationv1.FaultTargeting_UpstreamEnforcing:
			// Latency Internal Fault for all downstream services
			if downstreamCluster == "" {
				percentageKey = fmt.Sprintf(LatencyPercentageWithoutDownstream, g.IngressFaultRuntimePrefix)
				faultKey = fmt.Sprintf(LatencyDurationWithoutDownstream, g.IngressFaultRuntimePrefix)
			} else {
				// Latency Internal Fault for a given downstream services
				percentageKey = fmt.Sprintf(LatencyPercentageWithDownstream, g.IngressFaultRuntimePrefix, downstreamCluster)
				faultKey = fmt.Sprintf(LatencyDurationWithDownstream, g.IngressFaultRuntimePrefix, downstreamCluster)
			}

		default:
			return "", 0, "", 0, fmt.Errorf("unknown enforcer %v", httpFaultConfig)
		}

	default:
		return "", 0, "", 0, fmt.Errorf("unknown fault type %v", httpFaultConfig)
	}

	return percentageKey, percentageValue, faultKey, faultValue, nil
}

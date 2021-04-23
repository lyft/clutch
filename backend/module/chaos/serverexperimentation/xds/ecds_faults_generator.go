package xds

import (
	"fmt"
	"strings"

	gcpCoreV3 "github.com/envoyproxy/go-control-plane/envoy/config/core/v3"
	gcpRoute "github.com/envoyproxy/go-control-plane/envoy/config/route/v3"
	gcpFilterCommon "github.com/envoyproxy/go-control-plane/envoy/extensions/filters/common/fault/v3"
	gcpFilterFault "github.com/envoyproxy/go-control-plane/envoy/extensions/filters/http/fault/v3"
	gcpType "github.com/envoyproxy/go-control-plane/envoy/type/v3"
	proto "github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes/any"
	"github.com/golang/protobuf/ptypes/duration"

	serverexperimentation "github.com/lyft/clutch/backend/api/chaos/serverexperimentation/v1"
	ecdsv1 "github.com/lyft/clutch/backend/api/config/module/chaos/serverexperimentation/ecds/v1"
	"github.com/lyft/clutch/backend/module/chaos/experimentation/xds"
	"github.com/lyft/clutch/backend/service/chaos/experimentation/experimentstore"
)

const (
	envoyDownstreamServiceClusterHeader  = `x-envoy-downstream-service-cluster`
	faultFilterConfigNameForEgressFault  = `envoy.egress.extension_config.%s`
	faultFilterConfigNameForIngressFault = `envoy.extension_config`
	faultFilterTypeURL                   = `type.googleapis.com/envoy.extensions.filters.http.fault.v3.HTTPFault`

	delayPercentRuntime    = `ecds_runtime_override_do_not_use.http.delay.percentage`
	delayDurationRuntime   = `ecds_runtime_override_do_not_use.http.delay.fixed_duration_ms`
	abortHttpStatusRuntime = `ecds_runtime_override_do_not_use.http.abort.http_status`
	abortPercentRuntime    = `ecds_runtime_override_do_not_use.http.abort.abort_percent`
)

// Default abort and delay fault configs
var DefaultAbortFaultConfig = &gcpFilterFault.FaultAbort{
	ErrorType: &gcpFilterFault.FaultAbort_HttpStatus{
		HttpStatus: 503,
	},
	Percentage: &gcpType.FractionalPercent{
		Numerator:   0,
		Denominator: gcpType.FractionalPercent_HUNDRED,
	},
}

var DefaultDelayFaultConfig = &gcpFilterCommon.FaultDelay{
	FaultDelaySecifier: &gcpFilterCommon.FaultDelay_FixedDelay{
		FixedDelay: &duration.Duration{
			Nanos: 1000000, // 0.001 second
		},
	},
	Percentage: &gcpType.FractionalPercent{
		Numerator:   0,
		Denominator: gcpType.FractionalPercent_HUNDRED,
	},
}

type ECDSServerFaultsResourceGeneratorFactory struct{}

func (ECDSServerFaultsResourceGeneratorFactory) Create(cfg *any.Any) (xds.ECDSResourceGenerator, error) {
	typedConfig := &ecdsv1.ServerFaultsECDSResourceGeneratorConfig{}
	err := cfg.UnmarshalTo(typedConfig)
	if err != nil {
		return nil, err
	}

	return &ECDSServerFaultsResourceGenerator{}, nil
}

type ECDSServerFaultsResourceGenerator struct {
	SerializedDefaultFaultFilter []byte
}

func (g *ECDSServerFaultsResourceGenerator) GenerateResource(experiment *experimentstore.Experiment) (*xds.ECDSResource, error) {
	httpFaultConfig, ok := experiment.Config.Message.(*serverexperimentation.HTTPFaultConfig)
	if !ok {
		return nil, fmt.Errorf("ECDS server faults generator cannot generate faults for a given experiment (run ID: %s, config ID %s)", experiment.Run.Id, experiment.Config.Id)
	}

	upstreamCluster, downstreamCluster, err := GetClusterPair(httpFaultConfig)
	if err != nil {
		return nil, err
	}
	enforcingCluster, err := GetEnforcingCluster(httpFaultConfig)
	if err != nil {
		return nil, err
	}

	egressFault, abortFault, delayFault, err := g.createAbortDelayConfig(httpFaultConfig)
	if err != nil {
		return nil, err
	}

	abort := DefaultAbortFaultConfig
	delay := DefaultDelayFaultConfig
	isEgressFault := egressFault
	if abortFault != nil {
		abort = abortFault
	}
	if delayFault != nil {
		delay = delayFault
	}

	faultFilter := &gcpFilterFault.HTTPFault{
		Delay: delay,
		Abort: abort,

		// override runtimes so that default runtime is not used.
		DelayPercentRuntime:    delayPercentRuntime,
		DelayDurationRuntime:   delayDurationRuntime,
		AbortHttpStatusRuntime: abortHttpStatusRuntime,
		AbortPercentRuntime:    abortPercentRuntime,
	}

	var faultFilterName string
	if isEgressFault {
		faultFilterName = fmt.Sprintf(faultFilterConfigNameForEgressFault, upstreamCluster)
	} else {
		faultFilterName = faultFilterConfigNameForIngressFault
		// We match against the EnvoyDownstreamServiceCluster header to only apply these faults to the specific downstream cluster
		faultFilter.Headers = []*gcpRoute.HeaderMatcher{
			{
				Name: envoyDownstreamServiceClusterHeader,
				HeaderMatchSpecifier: &gcpRoute.HeaderMatcher_ExactMatch{
					ExactMatch: downstreamCluster,
				},
			},
		}
	}

	serializedFaultFilter, err := proto.Marshal(faultFilter)
	if err != nil {
		return nil, err
	}

	config := &gcpCoreV3.TypedExtensionConfig{
		Name: faultFilterName,
		TypedConfig: &any.Any{
			TypeUrl: faultFilterTypeURL,
			Value:   serializedFaultFilter,
		},
	}

	return xds.NewECDSResource(enforcingCluster, config)
}

func (g *ECDSServerFaultsResourceGenerator) GenerateDefaultResource(cluster string, resourceName string) (*xds.ECDSResource, error) {
	if resourceName != faultFilterConfigNameForIngressFault && !strings.HasPrefix(resourceName, fmt.Sprintf(faultFilterConfigNameForEgressFault, "")) {
		return xds.NewEmptyECDSResource(), nil
	}

	if g.SerializedDefaultFaultFilter == nil {
		faultFilter := &gcpFilterFault.HTTPFault{
			Delay: DefaultDelayFaultConfig,
			Abort: DefaultAbortFaultConfig,
			// override runtimes so that default runtime is not used.
			DelayPercentRuntime:    delayPercentRuntime,
			DelayDurationRuntime:   delayDurationRuntime,
			AbortHttpStatusRuntime: abortHttpStatusRuntime,
			AbortPercentRuntime:    abortPercentRuntime,
		}
		marshaledFilter, err := proto.Marshal(faultFilter)
		if err != nil {
			return nil, err
		}

		g.SerializedDefaultFaultFilter = marshaledFilter
	}

	config := &gcpCoreV3.TypedExtensionConfig{
		Name: resourceName,
		TypedConfig: &any.Any{
			TypeUrl: faultFilterTypeURL,
			Value:   g.SerializedDefaultFaultFilter,
		},
	}

	return xds.NewECDSResource(cluster, config)
}

func (g *ECDSServerFaultsResourceGenerator) createAbortDelayConfig(httpFaultConfig *serverexperimentation.HTTPFaultConfig) (bool, *gcpFilterFault.FaultAbort, *gcpFilterCommon.FaultDelay, error) {
	var isEgressFault bool
	var abort *gcpFilterFault.FaultAbort
	var delay *gcpFilterCommon.FaultDelay

	switch httpFaultConfig.GetFault().(type) {
	case *serverexperimentation.HTTPFaultConfig_AbortFault:
		switch httpFaultConfig.GetFaultTargeting().GetEnforcer().(type) {
		case *serverexperimentation.FaultTargeting_DownstreamEnforcing:
			// Abort Egress Fault
			isEgressFault = true

		case *serverexperimentation.FaultTargeting_UpstreamEnforcing:
			// Abort Internal Fault
			isEgressFault = false

		default:
			return false, nil, nil, fmt.Errorf("unknown enforcer %v", httpFaultConfig)
		}

		abort = &gcpFilterFault.FaultAbort{
			ErrorType: &gcpFilterFault.FaultAbort_HttpStatus{
				HttpStatus: httpFaultConfig.GetAbortFault().GetAbortStatus().GetHttpStatusCode(),
			},
			Percentage: &gcpType.FractionalPercent{
				Numerator:   httpFaultConfig.GetAbortFault().GetPercentage().GetPercentage(),
				Denominator: gcpType.FractionalPercent_HUNDRED,
			},
		}
	case *serverexperimentation.HTTPFaultConfig_LatencyFault:
		switch httpFaultConfig.GetFaultTargeting().GetEnforcer().(type) {
		case *serverexperimentation.FaultTargeting_DownstreamEnforcing:
			// Latency Egress Fault
			isEgressFault = true

		case *serverexperimentation.FaultTargeting_UpstreamEnforcing:
			// Latency Internal Fault for all downstream services
			isEgressFault = false

		default:
			return false, nil, nil, fmt.Errorf("unknown enforcer %v", httpFaultConfig)
		}

		delay = &gcpFilterCommon.FaultDelay{
			FaultDelaySecifier: &gcpFilterCommon.FaultDelay_FixedDelay{
				FixedDelay: &duration.Duration{
					Nanos: int32(httpFaultConfig.GetLatencyFault().GetLatencyDuration().GetFixedDurationMs() * 1000000),
				},
			},
			Percentage: &gcpType.FractionalPercent{
				Numerator:   httpFaultConfig.GetLatencyFault().GetPercentage().GetPercentage(),
				Denominator: gcpType.FractionalPercent_HUNDRED,
			},
		}
	default:
		return false, nil, nil, fmt.Errorf("unknown fault type %v", httpFaultConfig)
	}

	return isEgressFault, abort, delay, nil
}

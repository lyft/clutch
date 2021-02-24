package xds

import (
	"fmt"
	"github.com/golang/protobuf/ptypes/any"
	"github.com/golang/protobuf/ptypes/duration"
	"sync"
	"time"

	gcpCoreV3 "github.com/envoyproxy/go-control-plane/envoy/config/core/v3"
	gcpRoute "github.com/envoyproxy/go-control-plane/envoy/config/route/v3"
	gcpFilterCommon "github.com/envoyproxy/go-control-plane/envoy/extensions/filters/common/fault/v3"
	gcpFilterFault "github.com/envoyproxy/go-control-plane/envoy/extensions/filters/http/fault/v3"
	gcpType "github.com/envoyproxy/go-control-plane/envoy/type/v3"
	gcpTypes "github.com/envoyproxy/go-control-plane/pkg/cache/types"
	"go.uber.org/zap"

	proto "github.com/golang/protobuf/proto"

	experimentation "github.com/lyft/clutch/backend/api/chaos/experimentation/v1"
	serverexperimentation "github.com/lyft/clutch/backend/api/chaos/serverexperimentation/v1"
)

const (
	HeaderXEnvoyDownstreamServiceCluster = `x-envoy-downstream-service-cluster`
	ExternalFaultFilterConfigName        = `envoy.egress.extension_config.%s`
	InternalFaultFilterConfigName        = `envoy.extension_config`
	FaultFilterTypeURL                   = `type.googleapis.com/envoy.extensions.filters.http.fault.v3.HTTPFault`

	DelayPercentRuntime    = "%s.http.delay.percentage"
	DelayDurationRuntime   = "%s.http.delay.fixed_duration_ms"
	AbortHttpStatusRuntime = "%s.http.abort.http_status"
	AbortPercentRuntime    = "%s.http.abort.abort_percent"
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
			Seconds: 0,
		},
	},
	Percentage: &gcpType.FractionalPercent{
		Numerator:   0,
		Denominator: gcpType.FractionalPercent_HUNDRED,
	},
}

type ECDSConfig struct {
	enabledClusters []string

	faultRuntimePrefix string

	requestedResourcesMap sync.Map
}

func generateECDSResource(experiments []*experimentation.Experiment, ecdsConfig *ECDSConfig, ttl *time.Duration, logger *zap.SugaredLogger) []gcpTypes.ResourceWithTtl {
	var downstreamCluster string
	var upstreamCluster string
	var isExternalFault = false

	abort := DefaultAbortFaultConfig
	delay := DefaultDelayFaultConfig

	for _, experiment := range experiments {
		httpFaultConfig := &serverexperimentation.HTTPFaultConfig{}
		if !maybeUnmarshalFaultTest(experiment, httpFaultConfig) {
			continue
		}

		var err error
		upstreamCluster, downstreamCluster, err = getClusterPair(httpFaultConfig)
		if err != nil {
			logger.Errorw("Invalid http fault config", "config", httpFaultConfig)
			continue
		}

		isExternalFault, abort, delay, err = createAbortDelayConfig(httpFaultConfig, downstreamCluster)
		if err != nil {
			logger.Errorw("Unable to create ECDS filter fault", "config", httpFaultConfig)
			continue
		}

		// We can only perform one fault test per cluster with ECDS. Hence we break.
		break
	}

	faultFilter := &gcpFilterFault.HTTPFault{
		Delay: delay,
		Abort: abort,

		// override runtimes so that default runtime is not used.
		DelayPercentRuntime:    fmt.Sprintf(DelayPercentRuntime, ecdsConfig.faultRuntimePrefix),
		DelayDurationRuntime:   fmt.Sprintf(DelayDurationRuntime, ecdsConfig.faultRuntimePrefix),
		AbortHttpStatusRuntime: fmt.Sprintf(AbortHttpStatusRuntime, ecdsConfig.faultRuntimePrefix),
		AbortPercentRuntime:    fmt.Sprintf(AbortPercentRuntime, ecdsConfig.faultRuntimePrefix),
	}

	var faultFilterName string
	if isExternalFault {
		faultFilterName = fmt.Sprintf(ExternalFaultFilterConfigName, upstreamCluster)
	} else {
		faultFilterName = InternalFaultFilterConfigName

		// match downstream cluster with request header else fault will be applied to requests coming from all downstream
		faultFilter.Headers = []*gcpRoute.HeaderMatcher{
			{
				Name: HeaderXEnvoyDownstreamServiceCluster,
				HeaderMatchSpecifier: &gcpRoute.HeaderMatcher_ExactMatch{
					ExactMatch: downstreamCluster,
				},
			},
		}
	}

	serializedFaultFilter, err := proto.Marshal(faultFilter)
	if err != nil {
		logger.Warnw("Unable to unmarshal fault filter", "faultFilter", faultFilter)
		return nil
	}

	return []gcpTypes.ResourceWithTtl{{
		Resource: &gcpCoreV3.TypedExtensionConfig{
			Name: faultFilterName,
			TypedConfig: &any.Any{
				TypeUrl: FaultFilterTypeURL,
				Value:   serializedFaultFilter,
			},
		},
		Ttl: ttl,
	}}
}

func generateEmptyECDSResource(cluster string, ecdsConfig *ECDSConfig, logger *zap.SugaredLogger) []gcpTypes.ResourceWithTtl {
	faultFilter := &gcpFilterFault.HTTPFault{
		Delay: DefaultDelayFaultConfig,
		Abort: DefaultAbortFaultConfig,

		// override runtimes so that default runtime is not used.
		DelayPercentRuntime:    fmt.Sprintf(DelayPercentRuntime, ecdsConfig.faultRuntimePrefix),
		DelayDurationRuntime:   fmt.Sprintf(DelayDurationRuntime, ecdsConfig.faultRuntimePrefix),
		AbortHttpStatusRuntime: fmt.Sprintf(AbortHttpStatusRuntime, ecdsConfig.faultRuntimePrefix),
		AbortPercentRuntime:    fmt.Sprintf(AbortPercentRuntime, ecdsConfig.faultRuntimePrefix),
	}

	serializedFaultFilter, err := proto.Marshal(faultFilter)
	if err != nil {
		logger.Warnw("Unable to unmarshal fault filter", "faultFilter", faultFilter)
		return nil
	}

	clusterResources, isPresent := ecdsConfig.requestedResourcesMap.Load(cluster)
	if !isPresent {
		return nil
	}

	var resources []gcpTypes.ResourceWithTtl
	for _, resourceName := range clusterResources.([]string) {
		resource := gcpTypes.ResourceWithTtl{
			Resource: &gcpCoreV3.TypedExtensionConfig{
				Name: resourceName,
				TypedConfig: &any.Any{
					TypeUrl: FaultFilterTypeURL,
					Value:   serializedFaultFilter,
				},
			},
			Ttl: nil,
		}
		resources = append(resources, resource)
	}

	return resources
}

func createAbortDelayConfig(httpFaultConfig *serverexperimentation.HTTPFaultConfig, downstreamCluster string) (bool, *gcpFilterFault.FaultAbort, *gcpFilterCommon.FaultDelay, error) {
	var isExternalFault = false
	abort := DefaultAbortFaultConfig
	delay := DefaultDelayFaultConfig

	switch httpFaultConfig.GetFault().(type) {
	case *serverexperimentation.HTTPFaultConfig_AbortFault:
		switch httpFaultConfig.GetFaultTargeting().GetEnforcer().(type) {
		case *serverexperimentation.FaultTargeting_DownstreamEnforcing:
			// Abort External Fault
			isExternalFault = true

		case *serverexperimentation.FaultTargeting_UpstreamEnforcing:
			// Abort Internal Fault
			isExternalFault = false

		default:
			// If error, return default abort and default delay configs
			return false, abort, delay, fmt.Errorf("unknown enforcer %v", httpFaultConfig)
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
			// Latency External Fault
			isExternalFault = true

		case *serverexperimentation.FaultTargeting_UpstreamEnforcing:
			// Latency Internal Fault for all downstream services
			isExternalFault = false

		default:
			// If error, return default abort and default delay configs
			return false, abort, delay, fmt.Errorf("unknown enforcer %v", httpFaultConfig)
		}

		delay = &gcpFilterCommon.FaultDelay{
			FaultDelaySecifier: &gcpFilterCommon.FaultDelay_FixedDelay{
				FixedDelay: &duration.Duration{
					Seconds: int64(httpFaultConfig.GetLatencyFault().GetLatencyDuration().GetFixedDurationMs()),
				},
			},
			Percentage: &gcpType.FractionalPercent{
				Numerator:   httpFaultConfig.GetLatencyFault().GetPercentage().GetPercentage(),
				Denominator: gcpType.FractionalPercent_HUNDRED,
			},
		}
	default:
		return false, abort, delay, fmt.Errorf("unknown fault type %v", httpFaultConfig)
	}

	return isExternalFault, abort, delay, nil
}

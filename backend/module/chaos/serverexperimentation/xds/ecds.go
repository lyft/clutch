package xds

import (
	"fmt"
	"sync"
	"time"

	gcpCoreV3 "github.com/envoyproxy/go-control-plane/envoy/config/core/v3"
	gcpRoute "github.com/envoyproxy/go-control-plane/envoy/config/route/v3"
	gcpFilterCommon "github.com/envoyproxy/go-control-plane/envoy/extensions/filters/common/fault/v3"
	gcpFilterFault "github.com/envoyproxy/go-control-plane/envoy/extensions/filters/http/fault/v3"
	gcpType "github.com/envoyproxy/go-control-plane/envoy/type/v3"
	gcpTypes "github.com/envoyproxy/go-control-plane/pkg/cache/types"
	proto "github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes/any"
	"github.com/golang/protobuf/ptypes/duration"
	"go.uber.org/zap"

	experimentation "github.com/lyft/clutch/backend/api/chaos/experimentation/v1"
	serverexperimentation "github.com/lyft/clutch/backend/api/chaos/serverexperimentation/v1"
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

// Serialized default fault filter
var SerializedDefaultFaultFilter []byte

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

type SafeEcdsResourceMap struct {
	mu                    sync.Mutex
	requestedResourcesMap map[string]map[string]struct{}
}

func (safeResource *SafeEcdsResourceMap) getResourcesFromCluster(cluster string) []string {
	var resourceNames []string

	safeResource.mu.Lock()
	if _, exists := safeResource.requestedResourcesMap[cluster]; exists {
		for resourceName := range safeResource.requestedResourcesMap[cluster] {
			resourceNames = append(resourceNames, resourceName)
		}
	}
	safeResource.mu.Unlock()

	return resourceNames
}

func (safeResource *SafeEcdsResourceMap) setResourcesForCluster(cluster string, newResources []string) {
	safeResource.mu.Lock()
	defer safeResource.mu.Unlock()

	resources := make(map[string]struct{})
	if _, exists := safeResource.requestedResourcesMap[cluster]; exists {
		resources = safeResource.requestedResourcesMap[cluster]
	}

	for _, newResource := range newResources {
		resources[newResource] = struct{}{}
	}

	safeResource.requestedResourcesMap[cluster] = resources
}

type ECDSConfig struct {
	enabledClusters map[string]struct{}

	ecdsResourceMap *SafeEcdsResourceMap
}

func generateECDSResource(experiments []*experimentation.Experiment, cluster string, ttl *time.Duration, logger *zap.SugaredLogger) []gcpTypes.ResourceWithTtl {
	var downstreamCluster string
	var upstreamCluster string
	var isEgressFault = false

	abort := DefaultAbortFaultConfig
	delay := DefaultDelayFaultConfig

	if len(experiments) > 1 {
		// We can only perform one fault test per cluster with ECDS
		logger.Infow("Found multiple experiments per cluster. Only first one will be executed", "cluster", cluster, "experiments", experiments)
	}

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

		egressFault, abortFault, delayFault, err := createAbortDelayConfig(httpFaultConfig, downstreamCluster)
		if err != nil {
			logger.Errorw("Unable to create ECDS filter fault", "config", httpFaultConfig)
			continue
		}

		isEgressFault = egressFault
		if abortFault != nil {
			abort = abortFault
		}

		if delayFault != nil {
			delay = delayFault
		}

		// We can only perform one fault test per cluster with ECDS. Hence we break
		break
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

	logger.Debugw("ECDS Fault filter", "faultFilter", faultFilter, "upstreamcluster", upstreamCluster)

	serializedFaultFilter, err := proto.Marshal(faultFilter)
	if err != nil {
		logger.Warnw("Unable to unmarshal fault filter", "faultFilter", faultFilter)
		return []gcpTypes.ResourceWithTtl{}
	}

	return []gcpTypes.ResourceWithTtl{{
		Resource: &gcpCoreV3.TypedExtensionConfig{
			Name: faultFilterName,
			TypedConfig: &any.Any{
				TypeUrl: faultFilterTypeURL,
				Value:   serializedFaultFilter,
			},
		},
		Ttl: ttl,
	}}
}

func generateEmptyECDSResource(cluster string, ecdsConfig *ECDSConfig, logger *zap.SugaredLogger) []gcpTypes.ResourceWithTtl {
	// cache serialized
	if SerializedDefaultFaultFilter == nil {
		faultFilter := &gcpFilterFault.HTTPFault{
			Delay: DefaultDelayFaultConfig,
			Abort: DefaultAbortFaultConfig,

			// override runtimes so that default runtime is not used.
			DelayPercentRuntime:    delayPercentRuntime,
			DelayDurationRuntime:   delayDurationRuntime,
			AbortHttpStatusRuntime: abortHttpStatusRuntime,
			AbortPercentRuntime:    abortPercentRuntime,
		}

		SerializedDefaultFaultFilter, _ = proto.Marshal(faultFilter)
	}

	if _, exists := ecdsConfig.ecdsResourceMap.requestedResourcesMap[cluster]; !exists {
		logger.Debugw("Cluster not found in ECDS resource map", "cluster", cluster)
		return []gcpTypes.ResourceWithTtl{}
	}

	resourceNames := ecdsConfig.ecdsResourceMap.getResourcesFromCluster(cluster)

	var resources []gcpTypes.ResourceWithTtl
	for _, resourceName := range resourceNames {
		resource := gcpTypes.ResourceWithTtl{
			Resource: &gcpCoreV3.TypedExtensionConfig{
				Name: resourceName,
				TypedConfig: &any.Any{
					TypeUrl: faultFilterTypeURL,
					Value:   SerializedDefaultFaultFilter,
				},
			},
			Ttl: nil,
		}
		resources = append(resources, resource)
	}

	return resources
}

func createAbortDelayConfig(httpFaultConfig *serverexperimentation.HTTPFaultConfig, downstreamCluster string) (bool, *gcpFilterFault.FaultAbort, *gcpFilterCommon.FaultDelay, error) {
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

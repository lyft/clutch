package rtds

import (
	"context"
	"fmt"
	"time"

	gcpDiscovery "github.com/envoyproxy/go-control-plane/envoy/service/discovery/v2"
	gcpRuntimeServiceV3 "github.com/envoyproxy/go-control-plane/envoy/service/runtime/v3"
	gcpTypes "github.com/envoyproxy/go-control-plane/pkg/cache/types"
	gcpCacheV2 "github.com/envoyproxy/go-control-plane/pkg/cache/v2"
	gcpCacheV3 "github.com/envoyproxy/go-control-plane/pkg/cache/v3"
	gcpResourceV2 "github.com/envoyproxy/go-control-plane/pkg/resource/v2"
	gcpResourceV3 "github.com/envoyproxy/go-control-plane/pkg/resource/v3"
	"github.com/golang/protobuf/ptypes"
	pstruct "github.com/golang/protobuf/ptypes/struct"
	"github.com/mitchellh/hashstructure/v2"
	"go.uber.org/zap"

	experimentation "github.com/lyft/clutch/backend/api/chaos/experimentation/v1"
	serverexperimentation "github.com/lyft/clutch/backend/api/chaos/serverexperimentation/v1"
	"github.com/lyft/clutch/backend/service/chaos/experimentation/experimentstore"
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

// cacheWrapper is a wrapper interface that abstracts away the cache operations to make it easier
// to reuse as much of the code as possible for V2/V3.
type cacheWrapper interface {
	GetStatusKeys() []string
	GetSnapshotVersion(key string) (string, error)
	SetRuntimeLayer(nodeName string, layerName string, layer *pstruct.Struct, version string) error
}

type cacheWrapperV2 struct {
	gcpCacheV2.SnapshotCache
}

func (c *cacheWrapperV2) GetSnapshotVersion(key string) (string, error) {
	snapshot, err := c.GetSnapshot(key)
	if err != nil {
		return "", nil
	}

	return snapshot.GetVersion(gcpResourceV2.RuntimeType), nil
}

func (c *cacheWrapperV2) SetRuntimeLayer(nodeName string, layerName string, layer *pstruct.Struct, version string) error {
	runtimes := []gcpTypes.Resource{
		&gcpDiscovery.Runtime{
			Name:  layerName,
			Layer: layer,
		},
	}
	snapshot := gcpCacheV2.NewSnapshot(version, nil, nil, nil, nil, runtimes, nil)
	err := c.SetSnapshot(nodeName, snapshot)
	if err != nil {
		return err
	}

	return nil
}

type cacheWrapperV3 struct {
	gcpCacheV3.SnapshotCache
	ttl *time.Duration
}

func (c *cacheWrapperV3) GetSnapshotVersion(key string) (string, error) {
	snapshot, err := c.GetSnapshot(key)
	if err != nil {
		return "", nil
	}

	return snapshot.GetVersion(gcpResourceV3.RuntimeType), nil
}

func (c *cacheWrapperV3) SetRuntimeLayer(nodeName string, layerName string, layer *pstruct.Struct, version string) error {
	// We only want to set a TTL for non-empty layers. This ensures that we don't spam clients with heartbeat responses
	// unless they have an active runtime override.
	var resourceTTL *time.Duration
	if len(layer.Fields) > 0 {
		resourceTTL = c.ttl
	}

	runtimes := []gcpTypes.ResourceWithTtl{{
		Resource: &gcpRuntimeServiceV3.Runtime{
			Name:  layerName,
			Layer: layer,
		},
		Ttl: resourceTTL,
	}}

	snapshot := gcpCacheV3.NewSnapshotWithTtls(version, nil, nil, nil, nil, runtimes, nil)
	err := c.SetSnapshot(nodeName, snapshot)
	if err != nil {
		return err
	}

	return nil
}

func PeriodicallyRefreshCache(s *Server) {
	ticker := time.NewTicker(s.cacheRefreshInterval)
	go func() {
		for range ticker.C {
			s.logger.Info("Refreshing RTDS cache")
			refreshCache(s.ctx, s.storer, &cacheWrapperV2{s.snapshotCacheV2}, s.rtdsLayerName, s.ingressPrefix, s.egressPrefix, s.logger)
			refreshCache(s.ctx, s.storer, &cacheWrapperV3{s.snapshotCacheV3, s.resourceTTL}, s.rtdsLayerName, s.ingressPrefix, s.egressPrefix, s.logger)
		}
	}()
}

func refreshCache(ctx context.Context, storer experimentstore.Storer, snapshotCache cacheWrapper, rtdsLayerName string,
	ingressPrefix string, egressPrefix string, logger *zap.SugaredLogger) {
	allRunningExperiments, err := storer.GetExperiments(ctx, "type.googleapis.com/clutch.chaos.serverexperimentation.v1.HTTPFaultConfig", experimentation.GetExperimentsRequest_STATUS_RUNNING)
	if err != nil {
		logger.Errorw("Failed to get data from experiments store", "error", err)

		// If failed to get data from DB, stop all ongoing faults.
		allRunningExperiments = []*experimentation.Experiment{}
	}

	clusterFaultMap := make(map[string][]*experimentation.Experiment)
	for _, experiment := range allRunningExperiments {
		httpFaultConfig := &serverexperimentation.HTTPFaultConfig{}
		if !maybeUnmarshalFaultTest(experiment, httpFaultConfig) {
			continue
		}

		upstreamCluster, downstreamCluster, err := getClusterPair(httpFaultConfig)
		if err != nil {
			logger.Errorw("Invalid http fault config", "config", httpFaultConfig)
			continue
		}

		switch httpFaultConfig.GetFaultTargeting().GetEnforcer().(type) {
		case *serverexperimentation.FaultTargeting_DownstreamEnforcing:
			clusterFaultMap[upstreamCluster] = append(clusterFaultMap[upstreamCluster], experiment)
		case *serverexperimentation.FaultTargeting_UpstreamEnforcing:
			clusterFaultMap[downstreamCluster] = append(clusterFaultMap[downstreamCluster], experiment)
		default:
			logger.Errorw("unknown enforcer %v", httpFaultConfig)
			continue
		}
	}

	// Settings snapshot with empty faults to remove the faults
	for _, cluster := range snapshotCache.GetStatusKeys() {
		if _, exist := clusterFaultMap[cluster]; !exist {
			logger.Debugw("Removing faults for cluster", "cluster", cluster)
			err = setSnapshot(snapshotCache, rtdsLayerName, cluster, ingressPrefix, egressPrefix, []*experimentation.Experiment{}, logger)
			if err != nil {
				logger.Errorw("Unable to unset the fault for cluster", "cluster", cluster,
					"error", err)
			}
		}
	}

	// Create/Update faults
	for cluster, faults := range clusterFaultMap {
		logger.Infow("Injecting fault for cluster", "cluster", cluster)
		err := setSnapshot(snapshotCache, rtdsLayerName, cluster, ingressPrefix, egressPrefix, faults, logger)
		if err != nil {
			logger.Errorw("Unable to set the fault for cluster", "cluster", cluster,
				"error", err)
		}
	}
}

func setSnapshot(snapshotCache cacheWrapper, rtdsLayerName string, cluster string,
	ingressPrefix string, egressPrefix string, experiments []*experimentation.Experiment, logger *zap.SugaredLogger) error {
	// TODO(snowp): This code runs perioidcally and will compute a layer for consideration on every loop, which is somewhat wasteful
	// for empty layers. We should short circuit this logic to avoid wasting cycles/generating garbage for clusters with no active faults.
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

		percentageKey, percentageValue, faultKey, faultValue, err := createRuntimeKeys(upstreamCluster, downstreamCluster, httpFaultConfig, ingressPrefix, egressPrefix, logger)
		if err != nil {
			logger.Errorw("Unable to create runtime keys", "config", httpFaultConfig)
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
			"fault_enforcer", getEnforcer(httpFaultConfig))
	}

	runtimeLayer := &pstruct.Struct{
		Fields: fieldMap,
	}

	computedVersion, err := computeChecksum(runtimeLayer)
	if err != nil {
		logger.Errorw("Error computing version", "error", err)
		return err
	}

	currentSnapshotVersion, err := snapshotCache.GetSnapshotVersion(cluster)
	if err == nil {
		if currentSnapshotVersion == computedVersion {
			// No change in snapshot of this cluster
			logger.Debugw("Fault exists for cluster", "cluster", cluster)
			return nil
		}
	}

	err = snapshotCache.SetRuntimeLayer(cluster, rtdsLayerName, runtimeLayer, computedVersion)
	if err != nil {
		logger.Errorw("Error setting snapshot", "error", err)
		return err
	}

	return nil
}

func createRuntimeKeys(upstreamCluster string, downstreamCluster string, httpFaultConfig *serverexperimentation.HTTPFaultConfig, ingressPrefix string, egressPrefix string, logger *zap.SugaredLogger) (string, uint32, string, uint32, error) {
	var percentageKey string
	var percentageValue uint32
	var faultKey string
	var faultValue uint32

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
			logger.Errorw("unknown enforcer %v", httpFaultConfig)
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
			logger.Errorw("unknown enforcer %v", httpFaultConfig)
			return "", 0, "", 0, fmt.Errorf("unknown enforcer %v", httpFaultConfig)
		}

	default:
		logger.Errorw("Unknown fault type %v", httpFaultConfig)
		return "", 0, "", 0, fmt.Errorf("unknown fault type %v", httpFaultConfig)
	}

	return percentageKey, percentageValue, faultKey, faultValue, nil
}

func maybeUnmarshalFaultTest(experiment *experimentation.Experiment, httpFaultConfig *serverexperimentation.HTTPFaultConfig) bool {
	err := ptypes.UnmarshalAny(experiment.GetConfig(), httpFaultConfig)
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

func computeChecksum(item interface{}) (string, error) {
	hash, err := hashstructure.Hash(item, hashstructure.FormatV1, &hashstructure.HashOptions{TagName: "json"})
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%d", hash), nil
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

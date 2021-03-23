package xds

import (
	"context"
	"fmt"
	"time"

	gcpTypes "github.com/envoyproxy/go-control-plane/pkg/cache/types"
	gcpCacheV3 "github.com/envoyproxy/go-control-plane/pkg/cache/v3"
	"github.com/golang/protobuf/ptypes"
	"github.com/mitchellh/hashstructure/v2"
	"github.com/uber-go/tally"
	"go.uber.org/zap"

	experimentation "github.com/lyft/clutch/backend/api/chaos/experimentation/v1"
	serverexperimentation "github.com/lyft/clutch/backend/api/chaos/serverexperimentation/v1"
	"github.com/lyft/clutch/backend/service/chaos/experimentation/experimentstore"
)

func PeriodicallyRefreshCache(s *Server) {
	ticker := time.NewTicker(s.cacheRefreshInterval)
	go func() {
		for range ticker.C {
			s.logger.Info("Refreshing xDS cache")
			refreshCache(s.ctx, s.storer, s.snapshotCacheV3, s.resourceTTL, s.rtdsConfig, s.ecdsConfig, s.xdsScope, s.logger)
		}
	}()
}

func refreshCache(ctx context.Context, storer experimentstore.Storer, snapshotCache gcpCacheV3.SnapshotCache, ttl *time.Duration, rtdsConfig *RTDSConfig, ecdsConfig *ECDSConfig, scope tally.Scope, logger *zap.SugaredLogger) {
	allRunningExperiments, err := storer.GetExperiments(ctx, "type.googleapis.com/clutch.chaos.serverexperimentation.v1.HTTPFaultConfig", experimentation.GetExperimentsRequest_STATUS_RUNNING)
	if err != nil {
		logger.Errorw("Failed to get data from experiments store", "error", err)

		// If failed to get data from DB, stop all ongoing experiments.
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
			clusterFaultMap[downstreamCluster] = append(clusterFaultMap[downstreamCluster], experiment)
		case *serverexperimentation.FaultTargeting_UpstreamEnforcing:
			clusterFaultMap[upstreamCluster] = append(clusterFaultMap[upstreamCluster], experiment)
		default:
			logger.Errorw("unknown enforcer %v", httpFaultConfig)
			continue
		}
	}

	// Settings snapshot with empty experiments to remove the experiments
	for _, cluster := range snapshotCache.GetStatusKeys() {
		if _, exist := clusterFaultMap[cluster]; !exist {
			emptyResources := make(map[gcpTypes.ResponseType][]gcpTypes.ResourceWithTtl)

			logger.Debugw("Removing experiments for cluster", "cluster", cluster)

			// in order to remove fault, we need to set the snapshot with default ecds config and default runtime resource
			emptyResources[gcpTypes.ExtensionConfig] = generateEmptyECDSResource(cluster, ecdsConfig, logger)
			emptyResources[gcpTypes.Runtime] = generateRTDSResource([]*experimentation.Experiment{}, rtdsConfig, ttl, logger)

			err := setSnapshot(emptyResources, cluster, snapshotCache, true, logger, scope)
			if err != nil {
				logger.Errorw("Unable to unset the fault for cluster", "cluster", cluster,
					"error", err)
			}
		}
	}

	// Create/Update experiments
	for cluster, experiments := range clusterFaultMap {
		resources := make(map[gcpTypes.ResponseType][]gcpTypes.ResourceWithTtl)

		logger.Debugw("Injecting fault for cluster", "cluster", cluster)

		// Enable ECDS if cluster is ECDS enabled else enable RTDS
		if _, exists := ecdsConfig.enabledClusters[cluster]; exists {
			resources[gcpTypes.ExtensionConfig] = generateECDSResource(experiments, cluster, ttl, logger)
			resources[gcpTypes.Runtime] = generateRTDSResource([]*experimentation.Experiment{}, rtdsConfig, ttl, logger)
		} else {
			resources[gcpTypes.Runtime] = generateRTDSResource(experiments, rtdsConfig, ttl, logger)
			resources[gcpTypes.ExtensionConfig] = generateEmptyECDSResource(cluster, ecdsConfig, logger)
		}

		err := setSnapshot(resources, cluster, snapshotCache, false, logger, scope)
		if err != nil {
			logger.Errorw("Unable to set the fault for cluster", "cluster", cluster,
				"error", err)
		}
	}
}

func setSnapshot(resourceMap map[gcpTypes.ResponseType][]gcpTypes.ResourceWithTtl, cluster string, snapshotCache gcpCacheV3.SnapshotCache, isSnapshotEmpty bool, logger *zap.SugaredLogger, scope tally.Scope) error {
	currentSnapshot, _ := snapshotCache.GetSnapshot(cluster)

	snapshot := gcpCacheV3.Snapshot{}
	isNewSnapshotSame := true
	for resourceType, resources := range resourceMap {
		newComputedVersion, err := computeChecksum(resources)
		if err != nil {
			logger.Errorw("Unable to compute checksum", "cluster", cluster, "resourceType", resourceType, "resources", resources)
			return nil
		}

		if newComputedVersion != currentSnapshot.Resources[resourceType].Version {
			isNewSnapshotSame = false
		}

		snapshot.Resources[resourceType] = gcpCacheV3.NewResourcesWithTtl(newComputedVersion, resources)
	}

	if isNewSnapshotSame {
		// No change in snapshot of this cluster
		logger.Debugw("Not setting snapshot as its same as old one", "cluster", cluster, "snapshot", currentSnapshot, "resourceMap", resourceMap)
		return nil
	}

	logger.Infow("Setting snapshot", "cluster", cluster, "resources", resourceMap)
	err := snapshotCache.SetSnapshot(cluster, snapshot)
	if err != nil {
		return err
	}

	if isSnapshotEmpty {
		scope.Tagged(map[string]string{
			"cluster": cluster,
		}).Gauge("active_faults").Update(-1)
	} else {
		scope.Tagged(map[string]string{
			"cluster": cluster,
		}).Gauge("active_faults").Update(1)
	}

	return nil
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

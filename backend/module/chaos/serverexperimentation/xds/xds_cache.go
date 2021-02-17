package xds

import (
	"context"
	"fmt"
	"time"

	gcpTypes "github.com/envoyproxy/go-control-plane/pkg/cache/types"
	gcpCacheV3 "github.com/envoyproxy/go-control-plane/pkg/cache/v3"
	"github.com/mitchellh/hashstructure/v2"
	"github.com/uber-go/tally"
	"go.uber.org/zap"

	experimentation "github.com/lyft/clutch/backend/api/chaos/experimentation/v1"
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
	allRunningExperiments, err := storer.GetExperiments(ctx, "", experimentation.GetExperimentsRequest_STATUS_RUNNING)
	if err != nil {
		logger.Errorw("Failed to get data from experiments store", "error", err)

		// If failed to get data from DB, stop all ongoing experiments.
		allRunningExperiments = []*experimentation.Experiment{}
	}

	clusterFaultMap := make(map[string][]*experimentation.Experiment)
	for _, experiment := range allRunningExperiments {
		// get the runtime generation for each experiment
		runtimeGeneration, err := storer.GetRuntimeGenerator(experiment.GetConfig().GetTypeUrl())
		if err != nil {
			logger.Errorw("No runtime generation registered to config", "config", experiment.GetConfig())
			continue
		}

		// get cluster name that will enforce the fault and append experiment
		clusterName, err := runtimeGeneration.GetEnforcingCluster(experiment, logger)
		clusterFaultMap[clusterName] = append(clusterFaultMap[clusterName], experiment)
	}

	// Settings snapshot with empty experiments to remove the experiments
	for _, cluster := range snapshotCache.GetStatusKeys() {
		if _, exist := clusterFaultMap[cluster]; !exist {
			emptyResources := make(map[gcpTypes.ResponseType][]gcpTypes.ResourceWithTtl)

			logger.Debugw("Removing experiments for cluster", "cluster", cluster)

			// in order to remove fault, we need to set the snapshot with default ecds config and default runtime resource
			emptyResources[gcpTypes.ExtensionConfig] = generateEmptyECDSResource(cluster, ecdsConfig, logger)
			emptyResources[gcpTypes.Runtime] = generateRTDSResource([]*experimentation.Experiment{}, storer, rtdsConfig, ttl, logger)

			err := setSnapshot(emptyResources, cluster, snapshotCache, logger)
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
			resources[gcpTypes.Runtime] = generateRTDSResource([]*experimentation.Experiment{}, storer, rtdsConfig, ttl, logger)
			resources[gcpTypes.ExtensionConfig] = generateECDSResource(experiments, cluster, ttl, logger)
		} else {
			resources[gcpTypes.Runtime] = generateRTDSResource(experiments, storer, rtdsConfig, ttl, logger)
			resources[gcpTypes.ExtensionConfig] = generateEmptyECDSResource(cluster, ecdsConfig, logger)
		}

		err := setSnapshot(resources, cluster, snapshotCache, logger)
		if err != nil {
			logger.Errorw("Unable to set the fault for cluster", "cluster", cluster,
				"error", err)
		}
	}

	scope.Gauge("active_faults").Update(float64(len(clusterFaultMap)))
}

func setSnapshot(resourceMap map[gcpTypes.ResponseType][]gcpTypes.ResourceWithTtl, cluster string, snapshotCache gcpCacheV3.SnapshotCache, logger *zap.SugaredLogger) error {
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

	return nil
}

func computeChecksum(item interface{}) (string, error) {
	hash, err := hashstructure.Hash(item, hashstructure.FormatV1, &hashstructure.HashOptions{TagName: "json"})
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%d", hash), nil
}

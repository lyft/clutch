package xds

import (
	"context"
	"fmt"
	"time"

	gcpTypes "github.com/envoyproxy/go-control-plane/pkg/cache/types"
	gcpCacheV3 "github.com/envoyproxy/go-control-plane/pkg/cache/v3"
	"github.com/golang/protobuf/ptypes"
	"github.com/mitchellh/hashstructure/v2"
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
			refreshCache(s.ctx, s.storer, s.snapshotCacheV3, s.resourceTTL, s.rtdsConfig, s.logger)
		}
	}()
}

func refreshCache(ctx context.Context, storer experimentstore.Storer, snapshotCache gcpCacheV3.SnapshotCache, ttl *time.Duration, rtdsConfig *RTDSConfig, logger *zap.SugaredLogger) {
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
			clusterFaultMap[upstreamCluster] = append(clusterFaultMap[upstreamCluster], experiment)
		case *serverexperimentation.FaultTargeting_UpstreamEnforcing:
			clusterFaultMap[downstreamCluster] = append(clusterFaultMap[downstreamCluster], experiment)
		default:
			logger.Errorw("unknown enforcer %v", httpFaultConfig)
			continue
		}
	}

	// Settings snapshot with empty experiments to remove the experiments
	for _, cluster := range snapshotCache.GetStatusKeys() {
		if _, exist := clusterFaultMap[cluster]; !exist {
			logger.Debugw("Removing experiments for cluster", "cluster", cluster)

			// in order to remove fault, we need to set the snapshot with empty runtime resource
			emptyRuntimeResource := generateRTDSResource([]*experimentation.Experiment{}, rtdsConfig, ttl, logger)

			err := setSnapshot(emptyRuntimeResource, cluster, snapshotCache)
			if err != nil {
				logger.Errorw("Unable to unset the fault for cluster", "cluster", cluster,
					"error", err)
			}
		}
	}

	// Create/Update experiments
	for cluster, experiments := range clusterFaultMap {
		logger.Infow("Injecting fault for cluster", "cluster", cluster)

		runtimeResource := generateRTDSResource(experiments, rtdsConfig, ttl, logger)
		err := setSnapshot(runtimeResource, cluster, snapshotCache)
		if err != nil {
			logger.Errorw("Unable to set the fault for cluster", "cluster", cluster,
				"error", err)
		}
	}
}

func setSnapshot(resource []gcpTypes.ResourceWithTtl, cluster string, snapshotCache gcpCacheV3.SnapshotCache) error {
	computedVersion, err := computeChecksum(resource)
	if err != nil {
		return err
	}

	currentSnapshotVersion := ""
	currentSnapshot, err := snapshotCache.GetSnapshot(cluster)
	if err == nil {
		currentSnapshotVersion = currentSnapshot.GetVersion(cluster)
	}

	if currentSnapshotVersion == computedVersion {
		// No change in snapshot of this cluster
		return nil
	}

	snapshot := gcpCacheV3.NewSnapshotWithTtls(computedVersion, nil, nil, nil, nil, resource, nil)
	err = snapshotCache.SetSnapshot(cluster, snapshot)
	if err != nil {
		return err
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

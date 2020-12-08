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
}

func (c *cacheWrapperV3) GetSnapshotVersion(key string) (string, error) {
	snapshot, err := c.GetSnapshot(key)
	if err != nil {
		return "", nil
	}

	return snapshot.GetVersion(gcpResourceV3.RuntimeType), nil
}

func (c *cacheWrapperV3) SetRuntimeLayer(nodeName string, layerName string, layer *pstruct.Struct, version string) error {
	runtimes := []gcpTypes.Resource{
		&gcpRuntimeServiceV3.Runtime{
			Name:  layerName,
			Layer: layer,
		},
	}
	snapshot := gcpCacheV3.NewSnapshot(version, nil, nil, nil, nil, runtimes, nil)
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
			refreshCache(s.ctx, s.storer, &cacheWrapperV3{s.snapshotCacheV3}, s.rtdsLayerName, s.ingressPrefix, s.egressPrefix, s.logger)
		}
	}()
}

func refreshCache(ctx context.Context, storer experimentstore.Storer, snapshotCache cacheWrapper, rtdsLayerName string,
	ingressPrefix string, egressPrefix string, logger *zap.SugaredLogger) {
	allRunningExperiments, err := storer.GetExperiments(ctx, "type.googleapis.com/clutch.chaos.serverexperimentation.v1.TestConfig", experimentation.GetExperimentsRequest_STATUS_RUNNING)
	if err != nil {
		logger.Errorw("Failed to get data from experiments store", "error", err)

		// If failed to get data from DB, stop all ongoing faults.
		allRunningExperiments = []*experimentation.Experiment{}
	}

	clusterFaultMap := make(map[string][]*experimentation.Experiment)
	for _, experiment := range allRunningExperiments {
		testConfig := &serverexperimentation.TestConfig{}
		if !isFaultTest(experiment, testConfig) {
			continue
		}

		upstreamCluster := testConfig.GetClusterPair().GetUpstreamCluster()
		downstreamCluster := testConfig.GetClusterPair().GetDownstreamCluster()
		faultInjectionCluster := testConfig.GetClusterPair().GetFaultInjectionCluster()

		switch faultInjectionCluster {
		case serverexperimentation.FaultInjectionCluster_FAULTINJECTIONCLUSTER_UPSTREAM:
			clusterFaultMap[upstreamCluster] = append(clusterFaultMap[upstreamCluster], experiment)

		case serverexperimentation.FaultInjectionCluster_FAULTINJECTIONCLUSTER_DOWNSTREAM:
			clusterFaultMap[downstreamCluster] = append(clusterFaultMap[downstreamCluster], experiment)

		default:
			logger.Errorw("Invalid fault injection cluster found", "upstream", upstreamCluster, "downstream", downstreamCluster, "faultInjectionCluster", faultInjectionCluster)
		}
	}

	// Settings snapshot with empty faults to remove the faults
	for _, cluster := range snapshotCache.GetStatusKeys() {
		if _, exist := clusterFaultMap[cluster]; !exist {
			logger.Infow("Removing faults for cluster", "cluster", cluster)
			err = setSnapshot(snapshotCache, rtdsLayerName, cluster, ingressPrefix, egressPrefix, []*experimentation.Experiment{}, logger)
			if err != nil {
				logger.Errorw("Unable to unset the fault for cluster", "cluster", cluster,
					"error", err)
				panic(err)
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
			panic(err)
		}
	}
}

func setSnapshot(snapshotCache cacheWrapper, rtdsLayerName string, cluster string,
	ingressPrefix string, egressPrefix string, experiments []*experimentation.Experiment, logger *zap.SugaredLogger) error {
	var fieldMap = map[string]*pstruct.Value{}

	// No experiments meaning clear all experiments for the given upstream cluster
	if len(experiments) != 0 {
		for _, experiment := range experiments {
			testConfig := &serverexperimentation.TestConfig{}
			if !isFaultTest(experiment, testConfig) {
				continue
			}
			percentageKey, percentageValue, faultKey, faultValue := createRuntimeKeys(testConfig, ingressPrefix, egressPrefix, logger)

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
			clusterPair := testConfig.GetClusterPair()
			logger.Debugw("Fault details",
				"upstream_cluster", clusterPair.GetUpstreamCluster(),
				"downstream_cluster", clusterPair.GetDownstreamCluster(),
				"fault_type", faultKey,
				"percentage", percentageValue,
				"value", faultValue,
				"fault_injection_type", testConfig.GetClusterPair().GetFaultInjectionCluster())
		}
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

func createRuntimeKeys(testConfig *serverexperimentation.TestConfig, ingressPrefix string, egressPrefix string, logger *zap.SugaredLogger) (string, float32, string, int32) {
	var percentageKey string
	var percentageValue float32
	var faultKey string
	var faultValue int32

	target := testConfig.GetClusterPair()
	faultInjectionCluster := testConfig.GetClusterPair().GetFaultInjectionCluster()

	switch testConfig.GetFault().(type) {
	case *serverexperimentation.TestConfig_Abort:
		abort := testConfig.GetAbort()
		percentageValue = abort.Percent
		faultValue = abort.HttpStatus

		switch faultInjectionCluster {
		case serverexperimentation.FaultInjectionCluster_FAULTINJECTIONCLUSTER_DOWNSTREAM:
			// Abort External Fault
			percentageKey = fmt.Sprintf(HTTPPercentageForExternal, egressPrefix, target.GetUpstreamCluster())
			faultKey = fmt.Sprintf(HTTPStatusForExternal, egressPrefix, target.GetUpstreamCluster())

		case serverexperimentation.FaultInjectionCluster_FAULTINJECTIONCLUSTER_UPSTREAM:
			// Abort Internal Fault for all downstream services
			if target.GetDownstreamCluster() == "" {
				percentageKey = fmt.Sprintf(HTTPPercentageWithoutDownstream, ingressPrefix)
				faultKey = fmt.Sprintf(HTTPStatusWithoutDownstream, ingressPrefix)
			} else {
				// Abort Internal Fault for a given downstream services
				percentageKey = fmt.Sprintf(HTTPPercentageWithDownstream, ingressPrefix, target.GetDownstreamCluster())
				faultKey = fmt.Sprintf(HTTPStatusWithDownstream, ingressPrefix, target.GetDownstreamCluster())
			}

		default:
			logger.Errorw("Invalid fault injection cluster found", "upstream", target.GetUpstreamCluster(), "downstream", target.GetDownstreamCluster(), "faultInjectionCluster", faultInjectionCluster)
		}

	case *serverexperimentation.TestConfig_Latency:
		latency := testConfig.GetLatency()
		percentageValue = latency.Percent
		faultValue = latency.DurationMs

		switch faultInjectionCluster {
		case serverexperimentation.FaultInjectionCluster_FAULTINJECTIONCLUSTER_DOWNSTREAM:
			// Latency External Fault
			percentageKey = fmt.Sprintf(LatencyPercentageForExternal, egressPrefix, target.GetUpstreamCluster())
			faultKey = fmt.Sprintf(LatencyDurationForExternal, egressPrefix, target.GetUpstreamCluster())

		case serverexperimentation.FaultInjectionCluster_FAULTINJECTIONCLUSTER_UPSTREAM:
			// Latency Internal Fault for all downstream services
			if target.GetDownstreamCluster() == "" {
				percentageKey = fmt.Sprintf(LatencyPercentageWithoutDownstream, ingressPrefix)
				faultKey = fmt.Sprintf(LatencyDurationWithoutDownstream, ingressPrefix)
			} else {
				// Latency Internal Fault for a given downstream services
				percentageKey = fmt.Sprintf(LatencyPercentageWithDownstream, ingressPrefix, target.GetDownstreamCluster())
				faultKey = fmt.Sprintf(LatencyDurationWithDownstream, ingressPrefix, target.GetDownstreamCluster())
			}

		default:
			logger.Errorw("Invalid fault injection cluster found", "upstream", target.GetUpstreamCluster(), "downstream", target.GetDownstreamCluster(), "faultInjectionCluster", faultInjectionCluster)
		}

	default:
		logger.Errorw("Unknown fault type: %t", testConfig)
		panic("Unknown fault type")
	}

	return percentageKey, percentageValue, faultKey, faultValue
}

func isFaultTest(experiment *experimentation.Experiment, testConfig *serverexperimentation.TestConfig) bool {
	err := ptypes.UnmarshalAny(experiment.GetConfig(), testConfig)
	if err != nil {
		return false
	}

	switch testConfig.GetFault().(type) {
	case *serverexperimentation.TestConfig_Abort:
		return true
	case *serverexperimentation.TestConfig_Latency:
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

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
	"github.com/mitchellh/hashstructure"
	"go.uber.org/zap"

	experimentation "github.com/lyft/clutch/backend/api/chaos/experimentation/v1"
	serverexperimentation "github.com/lyft/clutch/backend/api/chaos/serverexperimentation/v1"
	"github.com/lyft/clutch/backend/service/chaos/experimentation/experimentstore"
)

const (
	LatencyPercentageWithoutDownstream = `fault.http.delay.fixed_delay_percent`
	LatencyPercentageWithDownstream    = `fault.http.%s.delay.fixed_delay_percent`
	LatencyDurationWithoutDownstream   = `fault.http.delay.fixed_duration_ms`
	LatencyDurationWithDownstream      = `fault.http.%s.delay.fixed_duration_ms`
	HTTPPercentageWithoutDownstream    = `fault.http.abort.abort_percent`
	HTTPPercentageWithDownstream       = `fault.http.%s.abort.abort_percent`
	HTTPStatusWithoutDownstream        = `fault.http.abort.http_status`
	HTTPStatusWithDownstream           = `fault.http.%s.abort.http_status`
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
			refreshCache(s.ctx, s.storer, &cacheWrapperV2{s.snapshotCacheV2}, s.rtdsLayerName, s.logger)
			refreshCache(s.ctx, s.storer, &cacheWrapperV3{s.snapshotCacheV3}, s.rtdsLayerName, s.logger)
		}
	}()
}

func refreshCache(ctx context.Context, storer experimentstore.Storer, snapshotCache cacheWrapper, rtdsLayerName string,
	logger *zap.SugaredLogger) {
	allExperiments, err := storer.GetExperiments(ctx, "type.googleapis.com/clutch.chaos.serverexperimentation.v1.TestConfig", experimentation.GetExperimentsRequest_RUNNING)
	if err != nil {
		logger.Errorw("Failed to get data from experiments store", "error", err)

		// If failed to get data from DB, stop all ongoing faults.
		allExperiments = []*experimentation.Experiment{}
	}

	// Group faults by upstream cluster
	upstreamClusterFaultMap := make(map[string][]*experimentation.Experiment)
	for _, experiment := range allExperiments {
		testConfig := &serverexperimentation.TestConfig{}
		if !isFaultTest(experiment, testConfig) {
			continue
		}
		clusterPair := testConfig.GetClusterPair()
		upstreamClusterFaultMap[clusterPair.UpstreamCluster] =
			append(upstreamClusterFaultMap[clusterPair.UpstreamCluster], experiment)
	}

	// Settings snapshot with empty faults to remove the faults
	for _, upstreamCluster := range snapshotCache.GetStatusKeys() {
		if _, exist := upstreamClusterFaultMap[upstreamCluster]; !exist {
			logger.Infow("Removing faults for upstream cluster", "cluster", upstreamCluster)
			err = setSnapshot(snapshotCache, rtdsLayerName, upstreamCluster, []*experimentation.Experiment{}, logger)
			if err != nil {
				logger.Errorw("Unable to unset the fault for upstream cluster", "cluster", upstreamCluster,
					"error", err)
				panic(err)
			}
		}
	}

	// Create/Update faults
	for upstreamCluster, faults := range upstreamClusterFaultMap {
		logger.Infow("Injecting fault for upstream cluster", "cluster", upstreamCluster)
		err := setSnapshot(snapshotCache, rtdsLayerName, upstreamCluster, faults, logger)
		if err != nil {
			logger.Errorw("Unable to set the fault for upstream cluster", "cluster", upstreamCluster,
				"error", err)
			panic(err)
		}
	}
}

func setSnapshot(snapshotCache cacheWrapper, rtdsLayerName string, upstreamCluster string,
	experiments []*experimentation.Experiment, logger *zap.SugaredLogger) error {
	var fieldMap = map[string]*pstruct.Value{}

	// No experiments meaning clear all experiments for the given upstream cluster
	if len(experiments) != 0 {
		for _, experiment := range experiments {
			testConfig := &serverexperimentation.TestConfig{}
			if !isFaultTest(experiment, testConfig) {
				continue
			}
			percentageKey, percentageValue, faultKey, faultValue := createRuntimeKeys(testConfig, logger)

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
				"upstream_cluster", clusterPair.UpstreamCluster,
				"downstream_cluster", clusterPair.DownstreamCluster,
				"fault_type", faultKey,
				"percentage", percentageValue,
				"value", faultValue)
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

	currentSnapshotVersion, err := snapshotCache.GetSnapshotVersion(upstreamCluster)
	if err == nil {
		if currentSnapshotVersion == computedVersion {
			// No change in snapshot of this upstream cluster
			logger.Debugw("Fault exists for upstream cluster", "cluster", upstreamCluster)
			return nil
		}
	}

	err = snapshotCache.SetRuntimeLayer(upstreamCluster, rtdsLayerName, runtimeLayer, computedVersion)
	if err != nil {
		logger.Errorw("Error setting snapshot", "error", err)
		return err
	}

	return nil
}

func createRuntimeKeys(testConfig *serverexperimentation.TestConfig, logger *zap.SugaredLogger) (string, float32, string, int32) {
	var percentageKey string
	var percentageValue float32
	var faultKey string
	var faultValue int32

	target := testConfig.GetClusterPair()
	switch testConfig.GetFault().(type) {
	case *serverexperimentation.TestConfig_Abort:
		abort := testConfig.GetAbort()
		percentageValue = abort.Percent
		faultValue = abort.HttpStatus

		if target.DownstreamCluster == "" {
			percentageKey = HTTPPercentageWithoutDownstream
			faultKey = HTTPStatusWithoutDownstream
		} else {
			percentageKey = fmt.Sprintf(HTTPPercentageWithDownstream, target.DownstreamCluster)
			faultKey = fmt.Sprintf(HTTPStatusWithDownstream, target.DownstreamCluster)
		}
	case *serverexperimentation.TestConfig_Latency:
		latency := testConfig.GetLatency()
		percentageValue = latency.Percent
		faultValue = latency.DurationMs

		if target.DownstreamCluster == "" {
			percentageKey = LatencyPercentageWithoutDownstream
			faultKey = LatencyDurationWithoutDownstream
		} else {
			percentageKey = fmt.Sprintf(LatencyPercentageWithDownstream, target.DownstreamCluster)
			faultKey = fmt.Sprintf(LatencyDurationWithDownstream, target.DownstreamCluster)
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
	hash, err := hashstructure.Hash(item, &hashstructure.HashOptions{TagName: "json"})
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%d", hash), nil
}

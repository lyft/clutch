package rtds

import (
	"context"
	"fmt"
	"reflect"
	"time"

	gcpDiscovery "github.com/envoyproxy/go-control-plane/envoy/service/discovery/v2"
	gcpCache "github.com/envoyproxy/go-control-plane/pkg/cache/v2"
	pstruct "github.com/golang/protobuf/ptypes/struct"
	"github.com/mitchellh/hashstructure"
	"go.uber.org/zap"

	experimentation "github.com/lyft/clutch/backend/api/chaos/experimentation/v1"
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

func PeriodicallyRefreshCache(s *Server) {
	ticker := time.NewTicker(s.cacheRefreshInterval)
	go func() {
		for range ticker.C {
			s.logger.Info("Refreshing RTDS cache")
			refreshCache(s.ctx, s.experimentStore, s.snapshotCache, s.rtdsLayerName, s.logger)
		}
	}()
}

func refreshCache(ctx context.Context, store experimentstore.ExperimentStore, snapshotCache gcpCache.SnapshotCache, rtdsLayerName string,
	logger *zap.SugaredLogger) {
	allExperiments, err := store.GetExperiments(ctx)
	if err != nil {
		logger.Errorw("Failed to get data from experiments store", "error", err)
		return
	}

	// Group faults by upstream cluster
	upstreamClusterFaultMap := make(map[string][]*experimentation.Experiment)
	for _, experiment := range allExperiments {
		if !isFaultTest(experiment) {
			continue
		}
		clusterPair := getClusterPair(experiment)
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

func setSnapshot(snapshotCache gcpCache.SnapshotCache, rtdsLayerName string, upstreamCluster string,
	experiments []*experimentation.Experiment, logger *zap.SugaredLogger) error {
	var fieldMap = map[string]*pstruct.Value{}

	// No experiments meaning clear all experiments for the given upstream cluster
	if len(experiments) != 0 {
		for _, experiment := range experiments {
			percentageKey, percentageValue, faultKey, faultValue := createRuntimeKeys(experiment, logger)

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
			clusterPair := getClusterPair(experiment)
			logger.Debugw("Fault details",
				"upstream_cluster", clusterPair.UpstreamCluster,
				"downstream_cluster", clusterPair.DownstreamCluster,
				"fault_type", faultKey,
				"percentage", percentageValue,
				"value", faultValue)
		}
	}

	runtimes := []gcpCache.Resource{
		&gcpDiscovery.Runtime{
			Name: rtdsLayerName,
			Layer: &pstruct.Struct{
				Fields: fieldMap,
			},
		},
	}

	computedVersion, err := computeChecksum(runtimes)
	if err != nil {
		logger.Errorw("Error computing version", "error", err)
		return err
	}

	currentSnapshot, _ := snapshotCache.GetSnapshot(upstreamCluster)
	if !reflect.DeepEqual(currentSnapshot, gcpCache.Snapshot{}) {
		currentVersion := currentSnapshot.GetVersion(gcpCache.RuntimeType)
		if currentVersion == computedVersion {
			// No change in snapshot of this upstream cluster
			logger.Debugw("Fault exists for upstream cluster", "cluster", upstreamCluster)
			return nil
		}
	}

	snapshot := gcpCache.NewSnapshot(computedVersion, nil, nil, nil, nil, runtimes)
	err = snapshotCache.SetSnapshot(upstreamCluster, snapshot)
	if err != nil {
		logger.Errorw("Error setting snapshot", "error", err)
		return err
	}

	return nil
}

func createRuntimeKeys(experiment *experimentation.Experiment, logger *zap.SugaredLogger) (string, float32, string, int32) {
	var percentageKey string
	var percentageValue float32
	var faultKey string
	var faultValue int32

	switch experiment.GetTestSpecification().GetConfig().(type) {
	case *experimentation.TestSpecification_Abort:
		abort := experiment.GetTestSpecification().GetAbort()
		percentageValue = abort.Percent
		faultValue = abort.HttpStatus

		target := abort.GetClusterPair()
		if target.DownstreamCluster == "" {
			percentageKey = HTTPPercentageWithoutDownstream
			faultKey = HTTPStatusWithoutDownstream
		} else {
			percentageKey = fmt.Sprintf(HTTPPercentageWithDownstream, target.DownstreamCluster)
			faultKey = fmt.Sprintf(HTTPStatusWithDownstream, target.DownstreamCluster)
		}
	case *experimentation.TestSpecification_Latency:
		latency := experiment.GetTestSpecification().GetLatency()
		percentageValue = latency.Percent
		faultValue = latency.DurationMs

		target := latency.GetClusterPair()
		if target.DownstreamCluster == "" {
			percentageKey = LatencyPercentageWithoutDownstream
			faultKey = LatencyDurationWithoutDownstream
		} else {
			percentageKey = fmt.Sprintf(LatencyPercentageWithDownstream, target.DownstreamCluster)
			faultKey = fmt.Sprintf(LatencyDurationWithDownstream, target.DownstreamCluster)
		}
	default:
		logger.Errorw("Unknown fault type: %t", experiment)
		panic("Unknown fault type")
	}

	return percentageKey, percentageValue, faultKey, faultValue
}

func isFaultTest(experiment *experimentation.Experiment) bool {
	switch experiment.GetTestSpecification().GetConfig().(type) {
	case *experimentation.TestSpecification_Abort:
		return true
	case *experimentation.TestSpecification_Latency:
		return true
	default:
		return false
	}
}

func getClusterPair(experiment *experimentation.Experiment) *experimentation.ClusterPairTarget {
	switch experiment.GetTestSpecification().GetConfig().(type) {
	case *experimentation.TestSpecification_Abort:
		abort := experiment.GetTestSpecification().GetAbort()
		return abort.GetClusterPair()
	case *experimentation.TestSpecification_Latency:
		latency := experiment.GetTestSpecification().GetLatency()
		return latency.GetClusterPair()
	}
	panic("unknown fault type")
}

func computeChecksum(item interface{}) (string, error) {
	hash, err := hashstructure.Hash(item, nil)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%d", hash), nil
}

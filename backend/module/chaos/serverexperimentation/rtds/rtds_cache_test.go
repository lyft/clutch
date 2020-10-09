package rtds

import (
	"context"
	"fmt"
	"testing"

	gcpDiscovery "github.com/envoyproxy/go-control-plane/envoy/service/discovery/v2"
	gcpRuntimeServiceV3 "github.com/envoyproxy/go-control-plane/envoy/service/runtime/v3"
	gcpTypes "github.com/envoyproxy/go-control-plane/pkg/cache/types"
	gcpCache "github.com/envoyproxy/go-control-plane/pkg/cache/v2"
	gcpCacheV3 "github.com/envoyproxy/go-control-plane/pkg/cache/v3"
	gcpResource "github.com/envoyproxy/go-control-plane/pkg/resource/v2"
	gcpResourceV3 "github.com/envoyproxy/go-control-plane/pkg/resource/v3"
	"github.com/golang/protobuf/ptypes"
	pstruct "github.com/golang/protobuf/ptypes/struct"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"

	experimentation "github.com/lyft/clutch/backend/api/chaos/experimentation/v1"
	serverexperimentation "github.com/lyft/clutch/backend/api/chaos/serverexperimentation/v1"
)

func createAbortExperiment(t *testing.T, upstreamCluster string, downstreamCluster string, faultPercent float32, httpStatus int32) *experimentation.Experiment {
	config := &serverexperimentation.TestConfig{
		ClusterPair: &serverexperimentation.ClusterPairTarget{
			DownstreamCluster: downstreamCluster,
			UpstreamCluster:   upstreamCluster,
		},
		Fault: &serverexperimentation.TestConfig_Abort{
			Abort: &serverexperimentation.AbortFaultConfig{
				Percent:    faultPercent,
				HttpStatus: httpStatus,
			},
		},
	}

	anyConfig, err := ptypes.MarshalAny(config)
	if err != nil {
		t.Errorf("marshalAny failed: %v", err)
	}

	return &experimentation.Experiment{Config: anyConfig}
}

func createLatencyExperiment(t *testing.T, upstreamCluster string, downstreamCluster string, latencyPercent float32, duration int32) *experimentation.Experiment {
	config := &serverexperimentation.TestConfig{
		ClusterPair: &serverexperimentation.ClusterPairTarget{
			DownstreamCluster: downstreamCluster,
			UpstreamCluster:   upstreamCluster,
		},
		Fault: &serverexperimentation.TestConfig_Latency{
			Latency: &serverexperimentation.LatencyFaultConfig{
				Percent:    latencyPercent,
				DurationMs: duration,
			},
		},
	}

	anyConfig, err := ptypes.MarshalAny(config)
	if err != nil {
		t.Errorf("setSnapshot failed %v", err)
	}

	return &experimentation.Experiment{Config: anyConfig}
}

func mockGenerateFaultData(t *testing.T) []*experimentation.Experiment {
	return []*experimentation.Experiment{
		createAbortExperiment(t, "serviceA", "serviceB", 10, 404),
		createAbortExperiment(t, "serviceC", "", 20, 504),
		createLatencyExperiment(t, "serviceA", "serviceD", 30, 100),
		createLatencyExperiment(t, "serviceE", "", 40, 200),
	}
}

func TestSetSnapshotV2(t *testing.T) {
	testCache := gcpCache.NewSnapshotCache(false, gcpCache.IDHash{}, nil)
	testRtdsLayerName := "testRtdsLayerName"
	testUpstreamCluster := "serviceA"
	mockExperimentList := mockGenerateFaultData(t)

	var testUpstreamClusterFaults []*experimentation.Experiment
	for _, experiment := range mockExperimentList {
		config := &serverexperimentation.TestConfig{}
		err := ptypes.UnmarshalAny(experiment.GetConfig(), config)
		if err != nil {
			t.Errorf("unmarshalAny failed %v", err)
		}

		target := config.GetClusterPair()
		if target.GetUpstreamCluster() == testUpstreamCluster {
			testUpstreamClusterFaults = append(testUpstreamClusterFaults, experiment)
		}
	}

	err := setSnapshot(&cacheWrapperV2{testCache}, testRtdsLayerName, testUpstreamCluster, testUpstreamClusterFaults, zap.NewNop().Sugar())
	if err != nil {
		t.Errorf("setSnapshot failed %v", err)
	}

	snapshot, err := testCache.GetSnapshot(testUpstreamCluster)
	if err != nil {
		t.Errorf("Snapshot not found for cluster %s", testUpstreamCluster)
	}

	resources := snapshot.GetResources(gcpResource.RuntimeType)
	if resources == nil {
		t.Errorf("no resources")
	}

	resource := resources[testRtdsLayerName]
	if resources == nil {
		t.Errorf("no RTDS resources")
	}

	runtime := resource.(*gcpDiscovery.Runtime)
	fields := runtime.GetLayer().GetFields()
	assert.Equal(t, 4, len(fields))
	assert.EqualValues(t, 10, fields["fault.http.serviceB.abort.abort_percent"].GetNumberValue())
	assert.EqualValues(t, 404, fields["fault.http.serviceB.abort.http_status"].GetNumberValue())
	assert.EqualValues(t, 30, fields["fault.http.serviceD.delay.fixed_delay_percent"].GetNumberValue())
	assert.EqualValues(t, 100, fields["fault.http.serviceD.delay.fixed_duration_ms"].GetNumberValue())
}

func TestSetSnapshotV3(t *testing.T) {
	testCache := gcpCacheV3.NewSnapshotCache(false, gcpCacheV3.IDHash{}, nil)
	testRtdsLayerName := "testRtdsLayerName"
	testUpstreamCluster := "serviceA"
	mockExperimentList := mockGenerateFaultData(t)

	var testUpstreamClusterFaults []*experimentation.Experiment
	for _, experiment := range mockExperimentList {
		config := &serverexperimentation.TestConfig{}
		err := ptypes.UnmarshalAny(experiment.GetConfig(), config)
		if err != nil {
			t.Errorf("unmarshalAny failed %v", err)
		}

		target := config.GetClusterPair()
		if target.GetUpstreamCluster() == testUpstreamCluster {
			testUpstreamClusterFaults = append(testUpstreamClusterFaults, experiment)
		}
	}

	err := setSnapshot(&cacheWrapperV3{testCache}, testRtdsLayerName, testUpstreamCluster, testUpstreamClusterFaults, zap.NewNop().Sugar())
	if err != nil {
		t.Errorf("setSnapshot failed %v", err)
	}

	snapshot, err := testCache.GetSnapshot(testUpstreamCluster)
	if err != nil {
		t.Errorf("Snapshot not found for cluster %s", testUpstreamCluster)
	}

	resources := snapshot.GetResources(gcpResourceV3.RuntimeType)
	if resources == nil {
		t.Errorf("no resources")
	}

	resource := resources[testRtdsLayerName]
	if resources == nil {
		t.Errorf("no RTDS resources")
	}

	runtime := resource.(*gcpRuntimeServiceV3.Runtime)
	fields := runtime.GetLayer().GetFields()
	assert.Equal(t, 4, len(fields))
	assert.EqualValues(t, 10, fields["fault.http.serviceB.abort.abort_percent"].GetNumberValue())
	assert.EqualValues(t, 404, fields["fault.http.serviceB.abort.http_status"].GetNumberValue())
	assert.EqualValues(t, 30, fields["fault.http.serviceD.delay.fixed_delay_percent"].GetNumberValue())
	assert.EqualValues(t, 100, fields["fault.http.serviceD.delay.fixed_duration_ms"].GetNumberValue())
}

func TestRefreshCache(t *testing.T) {
	es := &mockExperimentStore{}
	testCache := gcpCache.NewSnapshotCache(false, gcpCache.IDHash{}, nil)
	refreshCache(context.Background(), es, &cacheWrapperV2{testCache}, "test_layer", nil)
	assert.Equal(t, es.getExperimentArguments.configType, "type.googleapis.com/clutch.chaos.serverexperimentation.v1.TestConfig")
}

func TestCreateRuntimeKeys(t *testing.T) {
	testDataList := mockGenerateFaultData(t)
	for _, testExperiment := range testDataList {
		var expectedPercentageKey string
		var expectedPercentageValue float32
		var expectedFaultKey string
		var expectedFaultValue int32

		config := &serverexperimentation.TestConfig{}
		err := ptypes.UnmarshalAny(testExperiment.GetConfig(), config)
		if err != nil {
			t.Errorf("unmarshalAny failed %v", err)
		}

		target := config.GetClusterPair()
		switch config.GetFault().(type) {
		case *serverexperimentation.TestConfig_Abort:
			abort := config.GetAbort()
			expectedFaultValue = abort.HttpStatus
			expectedPercentageValue = abort.Percent

			if target.DownstreamCluster == "" {
				expectedPercentageKey = HTTPPercentageWithoutDownstream
				expectedFaultKey = HTTPStatusWithoutDownstream
			} else {
				expectedPercentageKey = fmt.Sprintf(HTTPPercentageWithDownstream, target.DownstreamCluster)
				expectedFaultKey = fmt.Sprintf(HTTPStatusWithDownstream, target.DownstreamCluster)
			}
		case *serverexperimentation.TestConfig_Latency:
			latency := config.GetLatency()
			expectedFaultValue = latency.DurationMs
			expectedPercentageValue = latency.Percent
			if target.DownstreamCluster == "" {
				expectedPercentageKey = LatencyPercentageWithoutDownstream
				expectedFaultKey = LatencyDurationWithoutDownstream
			} else {
				expectedPercentageKey = fmt.Sprintf(LatencyPercentageWithDownstream, target.DownstreamCluster)
				expectedFaultKey = fmt.Sprintf(LatencyDurationWithDownstream, target.DownstreamCluster)
			}
		}

		percentageKey, percentageValue, faultKey, faultValue := createRuntimeKeys(config, zap.NewNop().Sugar())

		assert.Equal(t, expectedPercentageKey, percentageKey)
		assert.Equal(t, expectedPercentageValue, percentageValue)
		assert.Equal(t, expectedFaultKey, faultKey)
		assert.Equal(t, expectedFaultValue, faultValue)
	}
}

func TestComputeVersionReturnValue(t *testing.T) {
	rtdsLayerName := "TestRtdsLayer"

	mockRuntime := []gcpTypes.Resource{
		&gcpDiscovery.Runtime{
			Name: rtdsLayerName,
			Layer: &pstruct.Struct{
				Fields: map[string]*pstruct.Value{},
			},
		},
	}

	checksum, _ := computeChecksum(mockRuntime)
	assert.Equal(t, "4464232557941748628", checksum)

	mockRuntime2 := []gcpTypes.Resource{
		&gcpDiscovery.Runtime{
			Name: rtdsLayerName + "foo",
			Layer: &pstruct.Struct{
				Fields: map[string]*pstruct.Value{},
			},
		},
	}

	checksum2, _ := computeChecksum(mockRuntime2)
	assert.NotEqual(t, checksum, checksum2)
}

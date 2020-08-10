package rtds

import (
	"fmt"
	"testing"

	gcpDiscovery "github.com/envoyproxy/go-control-plane/envoy/service/discovery/v2"
	gcpCache "github.com/envoyproxy/go-control-plane/pkg/cache/v2"
	pstruct "github.com/golang/protobuf/ptypes/struct"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"

	experimentation "github.com/lyft/clutch/backend/api/chaos/experimentation/v1"
)

func createAbortExperiment(upstreamCluster string, downstreamCluster string, faultPercent float32, httpStatus int32) *experimentation.Experiment {
	return &experimentation.Experiment{
		TestSpecification: &experimentation.TestSpecification{
			Config: &experimentation.TestSpecification_Abort{
				Abort: &experimentation.AbortFault{
					Target: &experimentation.AbortFault_ClusterPair{
						ClusterPair: &experimentation.ClusterPairTarget{
							DownstreamCluster: downstreamCluster,
							UpstreamCluster:   upstreamCluster,
						},
					},
					Percent:    faultPercent,
					HttpStatus: httpStatus},
			},
		},
	}
}

func createLatencyExperiment(upstreamCluster string, downstreamCluster string, latencyPercent float32, duration int32) *experimentation.Experiment {
	return &experimentation.Experiment{
		TestSpecification: &experimentation.TestSpecification{
			Config: &experimentation.TestSpecification_Latency{
				Latency: &experimentation.LatencyFault{
					Target: &experimentation.LatencyFault_ClusterPair{
						ClusterPair: &experimentation.ClusterPairTarget{
							DownstreamCluster: downstreamCluster,
							UpstreamCluster:   upstreamCluster,
						},
					},
					Percent:    latencyPercent,
					DurationMs: duration},
			},
		},
	}
}

func mockGenerateFaultData() []*experimentation.Experiment {
	return []*experimentation.Experiment{
		createAbortExperiment("serviceA", "serviceB", 10, 404),
		createAbortExperiment("serviceC", "", 20, 504),
		createLatencyExperiment("serviceA", "serviceD", 30, 100),
		createLatencyExperiment("serviceE", "", 40, 200),
	}
}

func TestSetSnapshot(t *testing.T) {
	testCache := gcpCache.NewSnapshotCache(false, gcpCache.IDHash{}, nil)
	testRtdsLayerName := "testRtdsLayerName"
	testUpstreamCluster := "serviceA"
	mockExperimentList := mockGenerateFaultData()

	var testUpstreamClusterFaults []*experimentation.Experiment
	for _, experiment := range mockExperimentList {
		switch experiment.GetTestSpecification().GetConfig().(type) {
		case *experimentation.TestSpecification_Abort:
			target := experiment.GetTestSpecification().GetAbort().GetClusterPair()
			if target.GetUpstreamCluster() == testUpstreamCluster {
				testUpstreamClusterFaults = append(testUpstreamClusterFaults, experiment)
			}
		case *experimentation.TestSpecification_Latency:
			target := experiment.GetTestSpecification().GetLatency().GetClusterPair()
			if target.GetUpstreamCluster() == testUpstreamCluster {
				testUpstreamClusterFaults = append(testUpstreamClusterFaults, experiment)
			}
		}
	}

	err := setSnapshot(testCache, testRtdsLayerName, testUpstreamCluster, testUpstreamClusterFaults, zap.NewNop().Sugar())
	if err != nil {
		t.Errorf("setSnapshot failed %v", err)
	}

	snapshot, err := testCache.GetSnapshot(testUpstreamCluster)
	if err != nil {
		t.Errorf("Snapshot not found for cluster %s", testUpstreamCluster)
	}

	resources := snapshot.GetResources(gcpCache.RuntimeType)
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

func TestCreateRuntimeKeys(t *testing.T) {
	testDataList := mockGenerateFaultData()
	for _, testExperiment := range testDataList {
		var expectedPercentageKey string
		var expectedPercentageValue float32
		var expectedFaultKey string
		var expectedFaultValue int32

		switch testExperiment.GetTestSpecification().GetConfig().(type) {
		case *experimentation.TestSpecification_Abort:
			abort := testExperiment.GetTestSpecification().GetAbort()
			expectedFaultValue = abort.HttpStatus
			expectedPercentageValue = abort.Percent
			target := abort.GetClusterPair()
			if target.DownstreamCluster == "" {
				expectedPercentageKey = HTTPPercentageWithoutDownstream
				expectedFaultKey = HTTPStatusWithoutDownstream
			} else {
				expectedPercentageKey = fmt.Sprintf(HTTPPercentageWithDownstream, target.DownstreamCluster)
				expectedFaultKey = fmt.Sprintf(HTTPStatusWithDownstream, target.DownstreamCluster)
			}
		case *experimentation.TestSpecification_Latency:
			latency := testExperiment.GetTestSpecification().GetLatency()
			expectedFaultValue = latency.DurationMs
			expectedPercentageValue = latency.Percent
			target := latency.GetClusterPair()
			if target.DownstreamCluster == "" {
				expectedPercentageKey = LatencyPercentageWithoutDownstream
				expectedFaultKey = LatencyDurationWithoutDownstream
			} else {
				expectedPercentageKey = fmt.Sprintf(LatencyPercentageWithDownstream, target.DownstreamCluster)
				expectedFaultKey = fmt.Sprintf(LatencyDurationWithDownstream, target.DownstreamCluster)
			}
		}

		percentageKey, percentageValue, faultKey, faultValue := createRuntimeKeys(testExperiment, zap.NewNop().Sugar())

		assert.Equal(t, expectedPercentageKey, percentageKey)
		assert.Equal(t, expectedPercentageValue, percentageValue)
		assert.Equal(t, expectedFaultKey, faultKey)
		assert.Equal(t, expectedFaultValue, faultValue)
	}
}

func TestComputeVersionReturnValue(t *testing.T) {
	rtdsLayerName := "TestRtdsLayer"

	mockRuntime := []gcpCache.Resource{
		&gcpDiscovery.Runtime{
			Name: rtdsLayerName,
			Layer: &pstruct.Struct{
				Fields: map[string]*pstruct.Value{},
			},
		},
	}

	checksum, _ := computeChecksum(mockRuntime)
	assert.Equal(t, "18363920902738959582", checksum)
}

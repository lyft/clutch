package xds

import (
	"context"
	"testing"
	"time"

	gcpDiscovery "github.com/envoyproxy/go-control-plane/envoy/service/discovery/v2"
	gcpRuntimeServiceV3 "github.com/envoyproxy/go-control-plane/envoy/service/runtime/v3"
	gcpTypes "github.com/envoyproxy/go-control-plane/pkg/cache/types"
	gcpCacheV3 "github.com/envoyproxy/go-control-plane/pkg/cache/v3"
	gcpResourceV3 "github.com/envoyproxy/go-control-plane/pkg/resource/v3"
	"github.com/golang/protobuf/ptypes"
	pstruct "github.com/golang/protobuf/ptypes/struct"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"

	experimentation "github.com/lyft/clutch/backend/api/chaos/experimentation/v1"
	serverexperimentation "github.com/lyft/clutch/backend/api/chaos/serverexperimentation/v1"
	xds_testing "github.com/lyft/clutch/backend/module/chaos/serverexperimentation/xds/testing"
)

const (
	INTERNAL = "internal"
	EXTERNAL = "external"
	ABORT    = "abort"
	LATENCY  = "latency"
)

func createExperiment(t *testing.T, upstreamCluster string, downstreamCluster string, faultPercent uint32, faultValue uint32, faultInjectorEnforcing string, faultType string) *experimentation.Experiment {
	config := &serverexperimentation.HTTPFaultConfig{}

	upstreamSingleCluster := &serverexperimentation.SingleCluster{
		Name: upstreamCluster,
	}

	downstreamSingleCluster := &serverexperimentation.SingleCluster{
		Name: downstreamCluster,
	}

	switch faultType {
	case ABORT:
		config.Fault = &serverexperimentation.HTTPFaultConfig_AbortFault{
			AbortFault: &serverexperimentation.AbortFault{
				Percentage:  &serverexperimentation.FaultPercentage{Percentage: faultPercent},
				AbortStatus: &serverexperimentation.FaultAbortStatus{HttpStatusCode: faultValue},
			},
		}
	case LATENCY:
		config.Fault = &serverexperimentation.HTTPFaultConfig_LatencyFault{
			LatencyFault: &serverexperimentation.LatencyFault{
				Percentage:      &serverexperimentation.FaultPercentage{Percentage: faultPercent},
				LatencyDuration: &serverexperimentation.FaultLatencyDuration{FixedDurationMs: faultValue},
			},
		}
	}

	switch faultInjectorEnforcing {
	case INTERNAL:
		config.FaultTargeting = &serverexperimentation.FaultTargeting{
			Enforcer: &serverexperimentation.FaultTargeting_UpstreamEnforcing{
				UpstreamEnforcing: &serverexperimentation.UpstreamEnforcing{
					UpstreamType: &serverexperimentation.UpstreamEnforcing_UpstreamCluster{
						UpstreamCluster: upstreamSingleCluster,
					},
					DownstreamType: &serverexperimentation.UpstreamEnforcing_DownstreamCluster{
						DownstreamCluster: downstreamSingleCluster,
					},
				},
			},
		}
	case EXTERNAL:
		config.FaultTargeting = &serverexperimentation.FaultTargeting{
			Enforcer: &serverexperimentation.FaultTargeting_DownstreamEnforcing{
				DownstreamEnforcing: &serverexperimentation.DownstreamEnforcing{
					UpstreamType: &serverexperimentation.DownstreamEnforcing_UpstreamCluster{
						UpstreamCluster: upstreamSingleCluster,
					},
					DownstreamType: &serverexperimentation.DownstreamEnforcing_DownstreamCluster{
						DownstreamCluster: downstreamSingleCluster,
					},
				},
			},
		}
	}

	anyConfig, err := ptypes.MarshalAny(config)
	if err != nil {
		t.Errorf("marshalAny failed: %v", err)
	}

	return &experimentation.Experiment{Config: anyConfig}
}

func mockGenerateFaultData(t *testing.T) []*experimentation.Experiment {
	return []*experimentation.Experiment{
		// Abort - Service B -> Service A (Internal)
		createExperiment(t, "serviceA", "serviceB", 10, 404, INTERNAL, ABORT),

		// Abort - All downstream -> Service C (Internal)
		createExperiment(t, "serviceC", "", 20, 504, INTERNAL, ABORT),

		// Latency - Service D -> Service A (Internal)
		createExperiment(t, "serviceA", "serviceD", 30, 100, INTERNAL, LATENCY),

		// Latency - All downstream -> Service E (Internal)
		createExperiment(t, "serviceE", "", 40, 200, INTERNAL, LATENCY),

		// Abort - Service A -> Service X (External)
		createExperiment(t, "serviceX", "serviceA", 65, 400, EXTERNAL, ABORT),

		// Latency - Service F -> Service Y (External)
		createExperiment(t, "serviceY", "serviceF", 40, 200, EXTERNAL, LATENCY),
	}
}

func TestSetSnapshotV3(t *testing.T) {
	testCache := gcpCacheV3.NewSnapshotCache(false, gcpCacheV3.IDHash{}, nil)
	testRtdsLayerName := "testRtdsLayerName"
	ingressPrefix := "ingress"
	egressPrefix := "egress"
	testCluster := "serviceA"
	mockExperimentList := mockGenerateFaultData(t)

	rtdsConfig := RTDSConfig{
		layerName:     testRtdsLayerName,
		ingressPrefix: ingressPrefix,
		egressPrefix:  egressPrefix,
	}

	var testClusterFaults []*experimentation.Experiment
	for _, experiment := range mockExperimentList {
		config := &serverexperimentation.HTTPFaultConfig{}
		err := ptypes.UnmarshalAny(experiment.GetConfig(), config)
		if err != nil {
			t.Errorf("unmarshalAny failed %v", err)
		}

		upstream, downstream, err := getClusterPair(config)
		assert.NoError(t, err)
		if upstream == testCluster || downstream == testCluster {
			testClusterFaults = append(testClusterFaults, experiment)
		}
	}

	runtimeResource := generateRTDSResource(testClusterFaults, &rtdsConfig, nil, zap.NewNop().Sugar())
	assert.Nil(t, runtimeResource[0].Ttl)
	err := setSnapshot(runtimeResource, testCluster, testCache, zap.NewNop().Sugar())
	if err != nil {
		t.Errorf("setSnapshot failed %v", err)
	}

	snapshot, err := testCache.GetSnapshot(testCluster)
	if err != nil {
		t.Errorf("Snapshot not found for cluster %s", testCluster)
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
	assert.Equal(t, 6, len(fields))
	assert.EqualValues(t, 10, fields["ingress.serviceB.abort.abort_percent"].GetNumberValue())
	assert.EqualValues(t, 404, fields["ingress.serviceB.abort.http_status"].GetNumberValue())
	assert.EqualValues(t, 30, fields["ingress.serviceD.delay.fixed_delay_percent"].GetNumberValue())
	assert.EqualValues(t, 100, fields["ingress.serviceD.delay.fixed_duration_ms"].GetNumberValue())
	assert.EqualValues(t, 65, fields["egress.serviceX.abort.abort_percent"].GetNumberValue())
	assert.EqualValues(t, 400, fields["egress.serviceX.abort.http_status"].GetNumberValue())
}

func TestSetSnapshotV3WithTTL(t *testing.T) {
	testCache := gcpCacheV3.NewSnapshotCache(false, gcpCacheV3.IDHash{}, nil)
	testRtdsLayerName := "testRtdsLayerName"
	ingressPrefix := "ingress"
	egressPrefix := "egress"
	testCluster := "serviceA"
	ttl := time.Duration(2 * 1000 * 1000 * 1000)
	mockExperimentList := mockGenerateFaultData(t)

	rtdsConfig := RTDSConfig{
		layerName:     testRtdsLayerName,
		ingressPrefix: ingressPrefix,
		egressPrefix:  egressPrefix,
	}

	var testClusterFaults []*experimentation.Experiment
	for _, experiment := range mockExperimentList {
		config := &serverexperimentation.HTTPFaultConfig{}
		err := ptypes.UnmarshalAny(experiment.GetConfig(), config)
		if err != nil {
			t.Errorf("unmarshalAny failed %v", err)
		}

		upstream, downstream, err := getClusterPair(config)
		assert.NoError(t, err)
		if upstream == testCluster || downstream == testCluster {
			testClusterFaults = append(testClusterFaults, experiment)
		}
	}

	runtimeResource := generateRTDSResource(testClusterFaults, &rtdsConfig, &ttl, zap.NewNop().Sugar())
	err := setSnapshot(runtimeResource, testCluster, testCache, zap.NewNop().Sugar())
	if err != nil {
		t.Errorf("setSnapshot failed %v", err)
	}

	snapshot, err := testCache.GetSnapshot(testCluster)
	if err != nil {
		t.Errorf("Snapshot not found for cluster %s", testCluster)
	}

	resources := snapshot.GetResourcesAndTtl(gcpResourceV3.RuntimeType)
	if resources == nil {
		t.Errorf("no resources")
	}

	resource := resources[testRtdsLayerName]
	if resources == nil {
		t.Errorf("no RTDS resources")
	}

	assert.Equal(t, ttl, *resource.Ttl)

	runtime := resource.Resource.(*gcpRuntimeServiceV3.Runtime)
	fields := runtime.GetLayer().GetFields()
	assert.Equal(t, 6, len(fields))
	assert.EqualValues(t, 10, fields["ingress.serviceB.abort.abort_percent"].GetNumberValue())
	assert.EqualValues(t, 404, fields["ingress.serviceB.abort.http_status"].GetNumberValue())
	assert.EqualValues(t, 30, fields["ingress.serviceD.delay.fixed_delay_percent"].GetNumberValue())
	assert.EqualValues(t, 100, fields["ingress.serviceD.delay.fixed_duration_ms"].GetNumberValue())
	assert.EqualValues(t, 65, fields["egress.serviceX.abort.abort_percent"].GetNumberValue())
	assert.EqualValues(t, 400, fields["egress.serviceX.abort.http_status"].GetNumberValue())
}

func TestRefreshCache(t *testing.T) {
	s := xds_testing.MockStorer{}
	rtdsConfig := RTDSConfig{
		layerName:     "test_layer",
		ingressPrefix: "ingress",
		egressPrefix:  "egress",
	}
	testCache := gcpCacheV3.NewSnapshotCache(false, gcpCacheV3.IDHash{}, nil)
	refreshCache(context.Background(), &s, testCache, nil, &rtdsConfig, nil)
	assert.Equal(t, s.GetExperimentArguments.ConfigType, "type.googleapis.com/clutch.chaos.serverexperimentation.v1.HTTPFaultConfig")
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

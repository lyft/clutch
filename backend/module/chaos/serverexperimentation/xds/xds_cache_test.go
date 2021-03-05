package xds

import (
	"context"
	"fmt"
	"testing"
	"time"

	gcpCoreV3 "github.com/envoyproxy/go-control-plane/envoy/config/core/v3"
	gcpFilterFault "github.com/envoyproxy/go-control-plane/envoy/extensions/filters/http/fault/v3"
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

func TestSetSnapshotRTDS(t *testing.T) {
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

	resources := make(map[gcpTypes.ResponseType][]gcpTypes.ResourceWithTtl)
	resources[gcpTypes.Runtime] = generateRTDSResource(testClusterFaults, &rtdsConfig, nil, zap.NewNop().Sugar())
	assert.Nil(t, resources[gcpTypes.Runtime][0].Ttl)
	err := setSnapshot(resources, testCluster, testCache, zap.NewNop().Sugar())
	if err != nil {
		t.Errorf("setSnapshot failed %v", err)
	}

	snapshot, err := testCache.GetSnapshot(testCluster)
	if err != nil {
		t.Errorf("Snapshot not found for cluster %s", testCluster)
	}

	rtdsResources := snapshot.GetResources(gcpResourceV3.RuntimeType)
	if rtdsResources == nil {
		t.Errorf("no resources")
	}

	resource := rtdsResources[testRtdsLayerName]
	if rtdsResources == nil {
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

	resources := make(map[gcpTypes.ResponseType][]gcpTypes.ResourceWithTtl)
	resources[gcpTypes.Runtime] = generateRTDSResource(testClusterFaults, &rtdsConfig, &ttl, zap.NewNop().Sugar())
	err := setSnapshot(resources, testCluster, testCache, zap.NewNop().Sugar())
	if err != nil {
		t.Errorf("setSnapshot failed %v", err)
	}

	snapshot, err := testCache.GetSnapshot(testCluster)
	if err != nil {
		t.Errorf("Snapshot not found for cluster %s", testCluster)
	}

	rtdsResources := snapshot.GetResourcesAndTtl(gcpResourceV3.RuntimeType)
	if rtdsResources == nil {
		t.Errorf("no resources")
	}

	resource := rtdsResources[testRtdsLayerName]
	if rtdsResources == nil {
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
	testCluster := "clusterA"
	s := xds_testing.MockStorer{}
	rtdsConfig := RTDSConfig{
		layerName:     "test_layer",
		ingressPrefix: "ingress",
		egressPrefix:  "egress",
	}

	ecdsConfig := ECDSConfig{
		ecdsResourceMap: &SafeEcdsResourceMap{},
		enabledClusters: map[string]struct{}{testCluster: struct{}{}},
	}

	testCache := gcpCacheV3.NewSnapshotCache(false, gcpCacheV3.IDHash{}, nil)
	refreshCache(context.Background(), &s, testCache, nil, &rtdsConfig, &ecdsConfig, nil)
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

func TestSetSnapshotECDSInternalFault(t *testing.T) {
	testCache := gcpCacheV3.NewSnapshotCache(false, gcpCacheV3.IDHash{}, nil)
	testCluster := "serviceA"
	var downstreamCluster string
	mockExperimentList := mockGenerateFaultData(t)

	var testClusterFaults []*experimentation.Experiment
	for _, experiment := range mockExperimentList {
		config := &serverexperimentation.HTTPFaultConfig{}
		err := ptypes.UnmarshalAny(experiment.GetConfig(), config)
		if err != nil {
			t.Errorf("unmarshalAny failed %v", err)
		}

		upstream, downstream, err := getClusterPair(config)
		assert.NoError(t, err)

		downstreamCluster = downstream

		// Test Internal Fault
		if upstream == testCluster {
			testClusterFaults = append(testClusterFaults, experiment)

			// Since ECDS can only perform one experiment per cluster
			break
		}
	}

	resources := make(map[gcpTypes.ResponseType][]gcpTypes.ResourceWithTtl)
	resources[gcpTypes.ExtensionConfig] = generateECDSResource(testClusterFaults, nil, zap.NewNop().Sugar())
	assert.Nil(t, resources[gcpTypes.ExtensionConfig][0].Ttl)
	err := setSnapshot(resources, testCluster, testCache, zap.NewNop().Sugar())
	if err != nil {
		t.Errorf("setSnapshot failed %v", err)
	}

	snapshot, err := testCache.GetSnapshot(testCluster)
	if err != nil {
		t.Errorf("Snapshot not found for cluster %s", testCluster)
	}

	ecdsResources := snapshot.Resources[gcpTypes.ExtensionConfig].Items
	if ecdsResources == nil {
		t.Errorf("no resources")
	}

	assert.Contains(t, ecdsResources, InternalFaultFilterConfigName)

	extensionConfig := ecdsResources[InternalFaultFilterConfigName].Resource.(*gcpCoreV3.TypedExtensionConfig)

	httpFaultFilter := &gcpFilterFault.HTTPFault{}
	err = ptypes.UnmarshalAny(extensionConfig.TypedConfig, httpFaultFilter)
	if err != nil {
		t.Errorf("Unable to unmarshall typed config")
	}

	assert.EqualValues(t, 10, httpFaultFilter.GetAbort().GetPercentage().GetNumerator())
	assert.EqualValues(t, 404, httpFaultFilter.GetAbort().GetHttpStatus())
	assert.EqualValues(t, 0, httpFaultFilter.GetDelay().GetPercentage().GetNumerator())
	assert.EqualValues(t, 0, httpFaultFilter.GetDelay().GetFixedDelay().GetSeconds())
	assert.EqualValues(t, 1, len(httpFaultFilter.GetHeaders()))
	assert.EqualValues(t, EnvoyDownstreamServiceClusterHeader, httpFaultFilter.GetHeaders()[0].GetName())
	assert.EqualValues(t, downstreamCluster, httpFaultFilter.GetHeaders()[0].GetExactMatch())
}

func TestSetSnapshotECDSExternalFault(t *testing.T) {
	testCache := gcpCacheV3.NewSnapshotCache(false, gcpCacheV3.IDHash{}, nil)
	testCluster := "serviceF"
	var upstreamCluster string
	mockExperimentList := mockGenerateFaultData(t)

	var testClusterFaults []*experimentation.Experiment
	for _, experiment := range mockExperimentList {
		config := &serverexperimentation.HTTPFaultConfig{}
		err := ptypes.UnmarshalAny(experiment.GetConfig(), config)
		if err != nil {
			t.Errorf("unmarshalAny failed %v", err)
		}

		upstream, downstream, err := getClusterPair(config)
		assert.NoError(t, err)

		upstreamCluster = upstream

		// Test External Fault
		if downstream == testCluster {
			testClusterFaults = append(testClusterFaults, experiment)

			// Since ECDS can only perform one experiment per cluster
			break
		}
	}

	resources := make(map[gcpTypes.ResponseType][]gcpTypes.ResourceWithTtl)
	resources[gcpTypes.ExtensionConfig] = generateECDSResource(testClusterFaults, nil, zap.NewNop().Sugar())
	assert.Nil(t, resources[gcpTypes.ExtensionConfig][0].Ttl)
	err := setSnapshot(resources, testCluster, testCache, zap.NewNop().Sugar())
	if err != nil {
		t.Errorf("setSnapshot failed %v", err)
	}

	snapshot, err := testCache.GetSnapshot(testCluster)
	if err != nil {
		t.Errorf("Snapshot not found for cluster %s", testCluster)
	}

	ecdsResources := snapshot.Resources[gcpTypes.ExtensionConfig].Items
	if ecdsResources == nil {
		t.Errorf("no resources")
	}

	assert.Contains(t, ecdsResources, fmt.Sprintf(ExternalFaultFilterConfigName, upstreamCluster))

	extensionConfig := ecdsResources[fmt.Sprintf(ExternalFaultFilterConfigName, upstreamCluster)].Resource.(*gcpCoreV3.TypedExtensionConfig)

	httpFaultFilter := &gcpFilterFault.HTTPFault{}
	err = ptypes.UnmarshalAny(extensionConfig.TypedConfig, httpFaultFilter)
	if err != nil {
		t.Errorf("Unable to unmarshall typed config")
	}

	assert.EqualValues(t, 0, httpFaultFilter.GetAbort().GetPercentage().GetNumerator())
	assert.EqualValues(t, 503, httpFaultFilter.GetAbort().GetHttpStatus())
	assert.EqualValues(t, 40, httpFaultFilter.GetDelay().GetPercentage().GetNumerator())
	assert.EqualValues(t, 200, httpFaultFilter.GetDelay().GetFixedDelay().GetSeconds())
	assert.EqualValues(t, 0, len(httpFaultFilter.GetHeaders()))
}

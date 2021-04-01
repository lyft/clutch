package xds

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"
	"time"

	gcpCoreV3 "github.com/envoyproxy/go-control-plane/envoy/config/core/v3"
	gcpRoute "github.com/envoyproxy/go-control-plane/envoy/config/route/v3"
	gcpFilterCommon "github.com/envoyproxy/go-control-plane/envoy/extensions/filters/common/fault/v3"
	gcpFilterFault "github.com/envoyproxy/go-control-plane/envoy/extensions/filters/http/fault/v3"
	gcpDiscovery "github.com/envoyproxy/go-control-plane/envoy/service/discovery/v2"
	gcpRuntimeServiceV3 "github.com/envoyproxy/go-control-plane/envoy/service/runtime/v3"
	gcpType "github.com/envoyproxy/go-control-plane/envoy/type/v3"
	gcpTypes "github.com/envoyproxy/go-control-plane/pkg/cache/types"
	gcpCacheV3 "github.com/envoyproxy/go-control-plane/pkg/cache/v3"
	gcpResourceV3 "github.com/envoyproxy/go-control-plane/pkg/resource/v3"
	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/duration"
	pstruct "github.com/golang/protobuf/ptypes/struct"
	"github.com/stretchr/testify/assert"
	"github.com/uber-go/tally"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"

	experimentation "github.com/lyft/clutch/backend/api/chaos/experimentation/v1"
	serverexperimentation "github.com/lyft/clutch/backend/api/chaos/serverexperimentation/v1"
	"github.com/lyft/clutch/backend/mock/service/chaos/experimentation/experimentstoremock"
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

	anyConfig, err := anypb.New(config)
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

func ProtoEqual(t *testing.T, a, b proto.Message) {
	t.Helper()

	jsonA, _ := json.Marshal(a)
	jsonB, _ := json.Marshal(b)
	assert.True(t, proto.Equal(a, b), "not equal: \n"+string(jsonA)+"\n"+string(jsonB))
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
	s := experimentstoremock.SimpleStorer{}

	scope := tally.NewTestScope("", nil)

	rtdsConfig := RTDSConfig{
		layerName:     "test_layer",
		ingressPrefix: "ingress",
		egressPrefix:  "egress",
	}

	ecdsConfig := ECDSConfig{
		ecdsResourceMap: &SafeEcdsResourceMap{},
		enabledClusters: map[string]struct{}{testCluster: {}},
	}

	now := time.Now()
	config := serverexperimentation.HTTPFaultConfig{
		Fault: &serverexperimentation.HTTPFaultConfig_AbortFault{
			AbortFault: &serverexperimentation.AbortFault{
				Percentage: &serverexperimentation.FaultPercentage{
					Percentage: 10,
				},
				AbortStatus: &serverexperimentation.FaultAbortStatus{
					HttpStatusCode: 503,
				},
			},
		},
		FaultTargeting: &serverexperimentation.FaultTargeting{
			Enforcer: &serverexperimentation.FaultTargeting_UpstreamEnforcing{
				UpstreamEnforcing: &serverexperimentation.UpstreamEnforcing{
					UpstreamType: &serverexperimentation.UpstreamEnforcing_UpstreamCluster{
						UpstreamCluster: &serverexperimentation.SingleCluster{
							Name: "cluster",
						},
					},
					DownstreamType: &serverexperimentation.UpstreamEnforcing_DownstreamCluster{
						DownstreamCluster: &serverexperimentation.SingleCluster{
							Name: "cluster",
						},
					},
				},
			},
		},
	}
	a, err := anypb.New(&config)
	assert.NoError(t, err)

	_, err = s.CreateExperiment(context.Background(), a, &now, &now)
	assert.NoError(t, err)

	testCache := gcpCacheV3.NewSnapshotCache(false, gcpCacheV3.IDHash{}, nil)
	refreshCache(context.Background(), &s, testCache, nil, &rtdsConfig, &ecdsConfig, scope, zap.NewNop().Sugar())

	experiments, err := s.GetExperiments(context.Background(), "type.googleapis.com/clutch.chaos.serverexperimentation.v1.HTTPFaultConfig", experimentation.GetExperimentsRequest_STATUS_RUNNING)
	assert.NoError(t, err)

	assert.Equal(t, 1, len(experiments))
	assert.Equal(t, float64(1), scope.Snapshot().Gauges()["active_faults+"].Value())
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
		assert.NoError(t, err)

		upstream, downstream, err := getClusterPair(config)
		assert.NoError(t, err)

		downstreamCluster = downstream

		// Test Internal Fault
		if upstream == testCluster {
			testClusterFaults = append(testClusterFaults, experiment)

			// Currently ECDS can only perform one experiment per cluster. We pick the first experiment in the list.
			break
		}
	}

	resources := make(map[gcpTypes.ResponseType][]gcpTypes.ResourceWithTtl)
	resources[gcpTypes.ExtensionConfig] = generateECDSResource(testClusterFaults, testCluster, nil, zap.NewNop().Sugar())
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

	assert.Contains(t, ecdsResources, faultFilterConfigNameForIngressFault)

	extensionConfig := ecdsResources[faultFilterConfigNameForIngressFault].Resource.(*gcpCoreV3.TypedExtensionConfig)

	httpFaultFilter := &gcpFilterFault.HTTPFault{}
	err = ptypes.UnmarshalAny(extensionConfig.TypedConfig, httpFaultFilter)
	if err != nil {
		t.Errorf("Unable to unmarshall typed config")
	}

	expectedHTTPFaultFilter := &gcpFilterFault.HTTPFault{
		Delay: &gcpFilterCommon.FaultDelay{
			FaultDelaySecifier: &gcpFilterCommon.FaultDelay_FixedDelay{
				FixedDelay: &duration.Duration{
					Nanos: 1000000, // 0.001 second
				},
			},
			Percentage: &gcpType.FractionalPercent{
				Numerator:   0,
				Denominator: gcpType.FractionalPercent_HUNDRED,
			},
		},
		Abort: &gcpFilterFault.FaultAbort{
			ErrorType: &gcpFilterFault.FaultAbort_HttpStatus{
				HttpStatus: 404,
			},
			Percentage: &gcpType.FractionalPercent{
				Numerator:   10,
				Denominator: gcpType.FractionalPercent_HUNDRED,
			},
		},
		Headers: []*gcpRoute.HeaderMatcher{
			{
				Name: envoyDownstreamServiceClusterHeader,
				HeaderMatchSpecifier: &gcpRoute.HeaderMatcher_ExactMatch{
					ExactMatch: downstreamCluster,
				},
			},
		},

		// override runtimes so that default runtime is not used.
		DelayPercentRuntime:    delayPercentRuntime,
		DelayDurationRuntime:   delayDurationRuntime,
		AbortHttpStatusRuntime: abortHttpStatusRuntime,
		AbortPercentRuntime:    abortPercentRuntime,
	}

	ProtoEqual(t, httpFaultFilter, expectedHTTPFaultFilter)
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
	resources[gcpTypes.ExtensionConfig] = generateECDSResource(testClusterFaults, testCluster, nil, zap.NewNop().Sugar())
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

	assert.Contains(t, ecdsResources, fmt.Sprintf(faultFilterConfigNameForEgressFault, upstreamCluster))

	extensionConfig := ecdsResources[fmt.Sprintf(faultFilterConfigNameForEgressFault, upstreamCluster)].Resource.(*gcpCoreV3.TypedExtensionConfig)

	httpFaultFilter := &gcpFilterFault.HTTPFault{}
	err = ptypes.UnmarshalAny(extensionConfig.TypedConfig, httpFaultFilter)
	if err != nil {
		t.Errorf("Unable to unmarshall typed config")
	}

	expectedHTTPFaultFilter := &gcpFilterFault.HTTPFault{
		Delay: &gcpFilterCommon.FaultDelay{
			FaultDelaySecifier: &gcpFilterCommon.FaultDelay_FixedDelay{
				FixedDelay: &duration.Duration{
					Nanos: 200000000, // 200 ms
				},
			},
			Percentage: &gcpType.FractionalPercent{
				Numerator:   40,
				Denominator: gcpType.FractionalPercent_HUNDRED,
			},
		},
		Abort: &gcpFilterFault.FaultAbort{
			ErrorType: &gcpFilterFault.FaultAbort_HttpStatus{
				HttpStatus: 503,
			},
			Percentage: &gcpType.FractionalPercent{
				Numerator:   0,
				Denominator: gcpType.FractionalPercent_HUNDRED,
			},
		},

		// override runtimes so that default runtime is not used.
		DelayPercentRuntime:    delayPercentRuntime,
		DelayDurationRuntime:   delayDurationRuntime,
		AbortHttpStatusRuntime: abortHttpStatusRuntime,
		AbortPercentRuntime:    abortPercentRuntime,
	}

	ProtoEqual(t, httpFaultFilter, expectedHTTPFaultFilter)
}

package xds

import (
	"context"
	"errors"
	"testing"
	"time"

	gcpCoreV3 "github.com/envoyproxy/go-control-plane/envoy/config/core/v3"
	gcpDiscovery "github.com/envoyproxy/go-control-plane/envoy/service/discovery/v2"
	gcpRuntimeServiceV3 "github.com/envoyproxy/go-control-plane/envoy/service/runtime/v3"
	gcpTypes "github.com/envoyproxy/go-control-plane/pkg/cache/types"
	gcpCacheV3 "github.com/envoyproxy/go-control-plane/pkg/cache/v3"
	gcpResourceV3 "github.com/envoyproxy/go-control-plane/pkg/resource/v3"
	"github.com/golang/protobuf/ptypes/any"
	pstruct "github.com/golang/protobuf/ptypes/struct"
	"github.com/stretchr/testify/assert"
	"github.com/uber-go/tally"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"
	"google.golang.org/protobuf/types/known/timestamppb"
	"google.golang.org/protobuf/types/known/wrapperspb"

	experimentation "github.com/lyft/clutch/backend/api/chaos/experimentation/v1"
	"github.com/lyft/clutch/backend/mock/service/chaos/experimentation/experimentstoremock"
	"github.com/lyft/clutch/backend/service/chaos/experimentation/experimentstore"
)

func TestRTDSResourcesRefresh(t *testing.T) {
	rtdsConfig := RTDSConfig{
		layerName: "test_layer",
	}
	ecdsConfig := ECDSConfig{
		ecdsResourceMap: &SafeEcdsResourceMap{},
		enabledClusters: make(map[string]struct{}),
	}

	configString := &wrapperspb.StringValue{}
	configInt := &wrapperspb.Int64Value{}
	configDouble := &wrapperspb.DoubleValue{}
	configBool := &wrapperspb.BoolValue{}
	configBytes := &wrapperspb.BytesValue{}

	input := []struct {
		data proto.Message
	}{
		{data: configString},
		{data: configInt},
		{data: configDouble},
		{data: configBool},
		{data: configBytes},
	}

	now := time.Now()
	futureTimestamp := timestamppb.New(now.Add(time.Hour))

	s := experimentstoremock.SimpleStorer{}
	for _, i := range input {
		any, err := anypb.New(i.data)
		assert.NoError(t, err)
		es, err := experimentstore.NewExperimentSpecification(
			&experimentation.CreateExperimentData{
				EndTime: futureTimestamp,
				Config:  any,
			}, now)
		assert.NoError(t, err)
		_, err = s.CreateExperiment(context.Background(), es)
		assert.NoError(t, err)
	}

	scope := tally.NewTestScope("", nil)
	l, err := zap.NewDevelopment()
	assert.NoError(t, err)

	rtdsGeneratorsByTypeUrl := map[string]RTDSResourceGenerator{
		// The error should be ignored and the process of generating resources
		// should continue.
		TypeUrl(configString): &MockRTDSResourceGenerator{
			err: errors.New("foo"),
		},
		// This resource should work just fine.
		TypeUrl(configDouble): &MockRTDSResourceGenerator{
			resource: &RTDSResource{
				Cluster:          "foo",
				RuntimeKeyValues: []*RuntimeKeyValue{{Key: "foo1", Value: 1}},
			},
		},
		// Even though both of these resources target the same cluster,
		// they can coexist since they specify different runtime keys.
		TypeUrl(configBool): &MockRTDSResourceGenerator{
			resource: &RTDSResource{
				Cluster:          "foo",
				RuntimeKeyValues: []*RuntimeKeyValue{{Key: "foo2", Value: 2}},
			},
		},
		// This resource cannot coexist with similar one from the above
		// since it targets the same cluster and uses the same runtime key.
		TypeUrl(configInt): &MockRTDSResourceGenerator{
			resource: &RTDSResource{
				Cluster:          "foo",
				RuntimeKeyValues: []*RuntimeKeyValue{{Key: "foo1", Value: 1}},
			},
		},
		// This resource can coexist with other resources that use `foo1` key
		// since it targets a cluster that's different compared to what other
		// resources use.
		TypeUrl(configBytes): &MockRTDSResourceGenerator{
			resource: &RTDSResource{
				Cluster:          "bar",
				RuntimeKeyValues: []*RuntimeKeyValue{{Key: "foo1", Value: 1}},
			},
		},
	}

	cache := gcpCacheV3.NewSnapshotCache(false, gcpCacheV3.IDHash{}, nil)
	p := Poller{
		ctx:                                context.Background(),
		storer:                             &s,
		snapshotCache:                      cache,
		resourceTtl:                        time.Duration(time.Second),
		cacheRefreshInterval:               time.Duration(time.Second),
		rtdsConfig:                         &rtdsConfig,
		ecdsConfig:                         &ecdsConfig,
		rtdsGeneratorsByTypeUrl:            rtdsGeneratorsByTypeUrl,
		ecdsGeneratorsByTypeUrl:            make(map[string]ECDSResourceGenerator),
		rtdsResourceGenerationFailureCount: scope.Counter("generation_failures"),
		ecdsResourceGenerationFailureCount: scope.Counter("c2"),
		ecdsDefaultResourceGenerationFailureCount: scope.Counter("c1"),
		setCacheSnapshotSuccessCount:              scope.Counter("c4"),
		setCacheSnapshotFailureCount:              scope.Counter("c5"),
		activeFaultsGauge:                         scope.Gauge("active_faults"),
		logger:                                    l.Sugar(),
	}
	defer p.Stop()
	p.Start()

	awaitGaugeEquals(t, scope, "active_faults+", 3)
	assert.Equal(t, int64(1), scope.Snapshot().Counters()["generation_failures+"].Value())

	// Verify runtime of cluster "foo"
	snapshotFoo, err := cache.GetSnapshot("foo")
	assert.NoError(t, err)
	rtdsResourcesClusterFoo := snapshotFoo.GetResourcesAndTtl(gcpResourceV3.RuntimeType)
	resourcesClusterFoo := rtdsResourcesClusterFoo[rtdsConfig.layerName]
	assert.NotNil(t, resourcesClusterFoo)
	fieldsFoo := resourcesClusterFoo.Resource.(*gcpRuntimeServiceV3.Runtime).GetLayer().GetFields()
	assert.EqualValues(t, 1, fieldsFoo["foo1"].GetNumberValue())
	assert.EqualValues(t, 2, fieldsFoo["foo2"].GetNumberValue())

	// Verify runtime of cluster "foo"
	snapshotBar, err := cache.GetSnapshot("bar")
	assert.NoError(t, err)
	rtdsResourcesClusterBar := snapshotBar.GetResourcesAndTtl(gcpResourceV3.RuntimeType)
	resourcesClusterBar := rtdsResourcesClusterBar[rtdsConfig.layerName]
	assert.NotNil(t, resourcesClusterBar)
	fieldsBar := resourcesClusterBar.Resource.(*gcpRuntimeServiceV3.Runtime).GetLayer().GetFields()
	assert.EqualValues(t, 1, fieldsBar["foo1"].GetNumberValue())
}

func TestECDSResourcesRefresh(t *testing.T) {
	rtdsConfig := RTDSConfig{
		layerName: "test_layer",
	}
	ecdsConfig := ECDSConfig{
		ecdsResourceMap: &SafeEcdsResourceMap{},
		enabledClusters: map[string]struct{}{
			"foo": struct{}{},
			"bar": struct{}{},
		},
	}

	configString := &wrapperspb.StringValue{}
	anyConfigString, err := anypb.New(configString)
	assert.NoError(t, err)
	configInt := &wrapperspb.Int64Value{}
	configDouble := &wrapperspb.DoubleValue{}
	configBool := &wrapperspb.BoolValue{}
	configBytes := &wrapperspb.BytesValue{}

	input := []struct {
		data proto.Message
	}{
		{data: configString},
		{data: configInt},
		{data: configDouble},
		{data: configBool},
		{data: configBytes},
	}

	now := time.Now()
	futureTimestamp := timestamppb.New(now.Add(time.Hour))

	s := experimentstoremock.SimpleStorer{}
	for _, i := range input {
		any, err := anypb.New(i.data)
		assert.NoError(t, err)
		es, err := experimentstore.NewExperimentSpecification(
			&experimentation.CreateExperimentData{
				EndTime: futureTimestamp,
				Config:  any,
			}, now)
		assert.NoError(t, err)
		_, err = s.CreateExperiment(context.Background(), es)
		assert.NoError(t, err)
	}

	scope := tally.NewTestScope("", nil)
	l, err := zap.NewDevelopment()
	assert.NoError(t, err)

	serializedFilter, err := proto.Marshal(anyConfigString)
	assert.NoError(t, err)

	ecdsGeneratorsByTypeUrl := map[string]ECDSResourceGenerator{
		// The error should be ignored and the process of generating resources
		// should continue.
		TypeUrl(configString): &MockECDSResourceGenerator{
			err: errors.New("foo"),
		},
		// Even though both of these resources target the same cluster,
		// they can coexist since they specify different runtime keys.
		TypeUrl(configDouble): &MockECDSResourceGenerator{
			resource: &ECDSResource{
				Cluster: "foo",
				ExtensionConfig: &gcpCoreV3.TypedExtensionConfig{
					Name: "filter1",
				},
			},
		},
		// This one is in conflict with the previous one since both of them
		// target the same cluster. Only one of them will be applied.
		TypeUrl(configBool): &MockECDSResourceGenerator{
			resource: &ECDSResource{
				Cluster: "foo",
				ExtensionConfig: &gcpCoreV3.TypedExtensionConfig{
					Name: "filter1",
				},
			},
		},
		// This resource can coexist with similar one from the above
		// since it targets a unique cluster.
		TypeUrl(configInt): &MockECDSResourceGenerator{
			resource: &ECDSResource{
				Cluster: "bar",
				ExtensionConfig: &gcpCoreV3.TypedExtensionConfig{
					Name: "filter2",
					TypedConfig: &any.Any{
						TypeUrl: TypeUrl(configString),
						Value:   serializedFilter,
					},
				},
			},
		},
		// This resource is ignored since ECDSConfig passed to poller doesn't
		// specify cluster "zar" as ECDS enabled one.
		TypeUrl(configBytes): &MockECDSResourceGenerator{
			resource: &ECDSResource{
				Cluster: "zar",
				ExtensionConfig: &gcpCoreV3.TypedExtensionConfig{
					Name: "filter3",
				},
			},
		},
	}

	cache := gcpCacheV3.NewSnapshotCache(false, gcpCacheV3.IDHash{}, nil)
	p := Poller{
		ctx:                                context.Background(),
		storer:                             &s,
		snapshotCache:                      cache,
		resourceTtl:                        time.Duration(time.Second),
		cacheRefreshInterval:               time.Duration(time.Second),
		rtdsConfig:                         &rtdsConfig,
		ecdsConfig:                         &ecdsConfig,
		rtdsGeneratorsByTypeUrl:            make(map[string]RTDSResourceGenerator),
		ecdsGeneratorsByTypeUrl:            ecdsGeneratorsByTypeUrl,
		rtdsResourceGenerationFailureCount: scope.Counter("c1"),
		ecdsResourceGenerationFailureCount: scope.Counter("c2"),
		ecdsDefaultResourceGenerationFailureCount: scope.Counter("c1"),
		setCacheSnapshotSuccessCount:              scope.Counter("c4"),
		setCacheSnapshotFailureCount:              scope.Counter("c5"),
		activeFaultsGauge:                         scope.Gauge("active_faults"),
		logger:                                    l.Sugar(),
	}
	defer p.Stop()
	p.Start()

	awaitGaugeEquals(t, scope, "active_faults+", 2)

	snapshotClusterFoo, err := cache.GetSnapshot("foo")
	assert.NoError(t, err)
	assert.NotNil(t, snapshotClusterFoo.Resources[gcpTypes.ExtensionConfig].Items["filter1"].Resource.(*gcpCoreV3.TypedExtensionConfig))

	snapshotClusterBar, err := cache.GetSnapshot("bar")
	assert.NoError(t, err)
	serializedFilterResource := snapshotClusterBar.Resources[gcpTypes.ExtensionConfig].Items["filter2"].Resource.(*gcpCoreV3.TypedExtensionConfig)
	unserializedResourceFilter := &wrapperspb.StringValue{}
	err = serializedFilterResource.TypedConfig.UnmarshalTo(unserializedResourceFilter)
	assert.NoError(t, err)

	_, err = cache.GetSnapshot("zar")
	assert.Error(t, err)
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

func awaitGaugeEquals(t *testing.T, scope tally.TestScope, gauge string, value float64) {
	t.Helper()

	timeout := time.NewTimer(5 * time.Second)
	for {
		v := float64(0)
		c, ok := scope.Snapshot().Gauges()[gauge]
		if ok {
			v = c.Value()
		}

		if value == v {
			return
		}

		select {
		case <-timeout.C:
			t.Errorf("timed out waiting for %s to become %f, last value %f", gauge, value, v)
			return
		case <-time.After(100 * time.Millisecond):
			break
		}
	}
}

type MockRTDSResourceGenerator struct {
	resource *RTDSResource
	err      error
}

func (g *MockRTDSResourceGenerator) GenerateResource(experiment *experimentstore.Experiment) (*RTDSResource, error) {
	if g.err != nil {
		return nil, g.err
	}
	return g.resource, nil
}

type MockECDSResourceGenerator struct {
	resource *ECDSResource
	err      error
}

func (g *MockECDSResourceGenerator) GenerateDefaultResource(cluster string, resourceName string) (*ECDSResource, error) {
	if g.err != nil {
		return nil, g.err
	}
	return g.resource, nil
}
func (g *MockECDSResourceGenerator) GenerateResource(experiment *experimentstore.Experiment) (*ECDSResource, error) {
	if g.err != nil {
		return nil, g.err
	}
	return g.resource, nil
}

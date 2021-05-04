package xds

import (
	"context"
	"errors"
	"testing"
	"time"

	envoy_config_core_v3 "github.com/envoyproxy/go-control-plane/envoy/config/core/v3"
	gcpDiscoveryV3 "github.com/envoyproxy/go-control-plane/envoy/service/discovery/v3"
	gcpRuntimeServiceV3 "github.com/envoyproxy/go-control-plane/envoy/service/runtime/v3"
	"github.com/stretchr/testify/assert"
	"github.com/uber-go/tally"
	rpc_status "google.golang.org/genproto/googleapis/rpc/status"
	"google.golang.org/protobuf/types/known/anypb"
	"google.golang.org/protobuf/types/known/durationpb"

	xdsv1 "github.com/lyft/clutch/backend/api/config/module/chaos/experimentation/xds/v1"
	"github.com/lyft/clutch/backend/module/chaos/experimentation/xds/internal/xdstest"
	"github.com/lyft/clutch/backend/service/chaos/experimentation/experimentstore"
)

func TestServerStats(t *testing.T) {
	xdsConfig := &xdsv1.Config{
		RtdsLayerName:     "rtds",
		HeartbeatInterval: durationpb.New(time.Second),
	}

	testServer, err := xdstest.NewTestModuleServer(New, false, xdsConfig)
	assert.NoError(t, err)
	defer testServer.Stop()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Verify V3 stats.
	v3Stream, err := newV3Stream(ctx, "rtds", "cluster", testServer)
	assert.NoError(t, err)
	defer v3Stream.close(t)

	// Regular flow.
	_, err = v3Stream.sendV3RequestAndAwaitResponse("", "")
	assert.NoError(t, err)

	assert.Equal(t, int64(1), testServer.Scope.Snapshot().Counters()["test.rtds.totalResourcesServed+"].Value())
	assert.Equal(t, int64(0), testServer.Scope.Snapshot().Counters()["test.rtds.totalErrorsReceived+"].Value())

	// Error response from xDS client.
	err = v3Stream.stream.Send(&gcpDiscoveryV3.DiscoveryRequest{ErrorDetail: &rpc_status.Status{}})
	assert.NoError(t, err)

	assert.Equal(t, int64(1), testServer.Scope.Snapshot().Counters()["test.rtds.totalResourcesServed+"].Value())
	// Async verification here since it appears that we don't get a response back in this case, so we
	// aren't able to synchronize on the response.
	awaitCounterEquals(t, testServer.Scope, "test.rtds.totalErrorsReceived+", 1)
}

// Verifies that TTL and heartbeating is done when configured to do so.
func TestResourceTTL(t *testing.T) {
	RTDSGeneratorsByTypeUrl[TypeUrl(&durationpb.Duration{})] = &MockRTDSResourceGenerator{
		resource: &RTDSResource{
			Cluster:          "cluster",
			RuntimeKeyValues: []*RuntimeKeyValue{{Key: "foo", Value: 1}},
		},
	}

	xdsConfig := &xdsv1.Config{
		RtdsLayerName: "rtds",
	}

	testServer, err := xdstest.NewTestModuleServer(New, true, xdsConfig)
	assert.NoError(t, err)
	defer testServer.Stop()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	now := time.Now()
	config := durationpb.Duration{}
	a, err := anypb.New(&config)
	assert.NoError(t, err)

	es := experimentstore.ExperimentSpecification{StartTime: now, EndTime: &now, Config: a}
	_, err = testServer.Storer.CreateExperiment(context.Background(), &es)
	assert.NoError(t, err)

	// First we look at a stream for a cluster that has an active fault. This should result in a TTL'd
	// resource that sends heartbeats.
	s, err := newV3Stream(ctx, "rtds", "cluster", testServer)
	assert.NoError(t, err)
	defer s.close(t)

	// Flow:
	// Send initial request
	// Receive response
	// Send request responding to response
	// Receive heartbeat response
	r, err := s.sendV3RequestAndAwaitResponse("", "")
	assert.NoError(t, err)

	resource := &gcpDiscoveryV3.Resource{}
	assert.NoError(t, r.Resources[0].UnmarshalTo(resource))

	assert.Equal(t, int64(2), resource.Ttl.Seconds)

	r, err = s.sendV3RequestAndAwaitResponse(r.VersionInfo, r.Nonce)
	assert.NoError(t, err)

	resource = &gcpDiscoveryV3.Resource{}
	assert.NoError(t, r.Resources[0].UnmarshalTo(resource))

	assert.Equal(t, int64(2), resource.Ttl.Seconds)
	assert.Nil(t, resource.Resource)

	// Second we look at a stream for a cluster that should not receive any faults.
	// Here we do not expect to see TTLs or heartbeats.
	noFaultStream, err := newV3Stream(ctx, "rtds", "other-cluster", testServer)
	assert.NoError(t, err)
	defer s.close(t)

	// Flow:
	// Send initial request
	// Receive response
	// Send request responding to response
	// <no response>
	r, err = noFaultStream.sendV3RequestAndAwaitResponse("", "")
	assert.NoError(t, err)

	runtime := &gcpRuntimeServiceV3.Runtime{}
	assert.NoError(t, r.Resources[0].UnmarshalTo(runtime))

	assert.Empty(t, runtime.Layer.Fields)

	assert.NoError(t, noFaultStream.sendV3RequestWithoutResponse(r.VersionInfo, r.Nonce))
}

func awaitCounterEquals(t *testing.T, scope tally.TestScope, counter string, value int64) {
	t.Helper()

	timeout := time.NewTimer(time.Second)
	for {
		v := int64(0)
		c, ok := scope.Snapshot().Counters()[counter]
		if ok {
			v = c.Value()
		}

		if value == v {
			return
		}

		select {
		case <-timeout.C:
			t.Errorf("timed out waiting for %s to become %d, last value %d", counter, value, v)
			return
		case <-time.After(100 * time.Millisecond):
			break
		}
	}
}

// Helper class for testing calls over a single RTDS stream.
type v3StreamWrapper struct {
	stream  gcpRuntimeServiceV3.RuntimeDiscoveryService_StreamRuntimeClient
	layer   string
	cluster string
}

func newV3Stream(ctx context.Context, layer string, cluster string, t *xdstest.TestModuleServer) (*v3StreamWrapper, error) {
	conn, err := t.ClientConn()
	if err != nil {
		return nil, err
	}

	client := gcpRuntimeServiceV3.NewRuntimeDiscoveryServiceClient(conn)
	s, err := client.StreamRuntime(ctx)
	if err != nil {
		return nil, err
	}

	return &v3StreamWrapper{stream: s, layer: layer, cluster: cluster}, err
}

func (v *v3StreamWrapper) close(t *testing.T) {
	assert.NoError(t, v.stream.CloseSend())
}

func (v *v3StreamWrapper) sendV3RequestWithoutResponse(version string, nonce string) error {
	err := v.stream.Send(&gcpDiscoveryV3.DiscoveryRequest{
		Node:          &envoy_config_core_v3.Node{Cluster: v.cluster},
		ResourceNames: []string{v.layer},
	})
	if err != nil {
		return err
	}

	// Capture the error if we get one, otherwise we pass nil over the channel which signals a response.
	// TODO(snowp): Output the response during test failures if necessary.
	responseCh := make(chan error)
	go func() {
		_, err := v.stream.Recv()
		responseCh <- err
	}()

	// Wait up to two seconds for a response.
	select {
	case err := <-responseCh:
		if err == nil {
			return errors.New("unexpected response")
		}
		return err
	case <-time.After(2 * time.Second):
	}

	return nil
}

func (v *v3StreamWrapper) sendV3RequestAndAwaitResponse(version string, nonce string) (*gcpDiscoveryV3.DiscoveryResponse, error) {
	err := v.stream.Send(&gcpDiscoveryV3.DiscoveryRequest{
		VersionInfo:   version,
		ResponseNonce: nonce,
		Node:          &envoy_config_core_v3.Node{Cluster: v.cluster},
		ResourceNames: []string{v.layer},
	})
	if err != nil {
		return nil, err
	}

	response, err := v.stream.Recv()
	if err != nil {
		return nil, err
	}

	return response, nil
}

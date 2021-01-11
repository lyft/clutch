package rtds

import (
	"context"
	"errors"
	"net"
	"testing"
	"time"

	envoy_api_v2 "github.com/envoyproxy/go-control-plane/envoy/api/v2"
	envoy_config_core_v3 "github.com/envoyproxy/go-control-plane/envoy/config/core/v3"
	gcpDiscoveryV2 "github.com/envoyproxy/go-control-plane/envoy/service/discovery/v2"
	gcpDiscoveryV3 "github.com/envoyproxy/go-control-plane/envoy/service/discovery/v3"
	gcpRuntimeServiceV3 "github.com/envoyproxy/go-control-plane/envoy/service/runtime/v3"
	"github.com/golang/protobuf/ptypes"
	"github.com/stretchr/testify/assert"
	"github.com/uber-go/tally"
	"go.uber.org/zap"
	rpc_status "google.golang.org/genproto/googleapis/rpc/status"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/durationpb"

	serverexperimentation "github.com/lyft/clutch/backend/api/chaos/serverexperimentation/v1"
	rtdsconfigv1 "github.com/lyft/clutch/backend/api/config/module/chaos/experimentation/rtds/v1"
	"github.com/lyft/clutch/backend/module/moduletest"
	"github.com/lyft/clutch/backend/service"
	"github.com/lyft/clutch/backend/service/chaos/experimentation/experimentstore"
)

func TestServerStats(t *testing.T) {
	testServer := newTestServer(t, false)
	defer testServer.stop()

	// Connect to the test server.
	conn, err := testServer.clientConn()
	assert.NoError(t, err)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Verify V3 stats.
	v3Stream, err := newV3Stream(ctx, "tests", "cluster", testServer)
	assert.NoError(t, err)
	defer v3Stream.close(t)

	// Regular flow.
	_, err = v3Stream.sendV3RequestAndAwaitResponse("", "")
	assert.NoError(t, err)

	assert.Equal(t, int64(1), testServer.scope.Snapshot().Counters()["test.v3.totalResourcesServed+"].Value())
	assert.Equal(t, int64(0), testServer.scope.Snapshot().Counters()["test.v3.totalErrorsReceived+"].Value())

	// Error response from xDS client.
	err = v3Stream.stream.Send(&gcpDiscoveryV3.DiscoveryRequest{ErrorDetail: &rpc_status.Status{}})
	assert.NoError(t, err)

	assert.Equal(t, int64(1), testServer.scope.Snapshot().Counters()["test.v3.totalResourcesServed+"].Value())
	// Async verification here since it appears that we don't get a response back in this case, so we
	// aren't able to synchronize on the response.
	awaitCounterEquals(t, testServer.scope, "test.v3.totalErrorsReceived+", 1)

	// Verify V2 stats.
	v2Client := gcpDiscoveryV2.NewRuntimeDiscoveryServiceClient(conn)
	v2Stream, err := v2Client.StreamRuntime(ctx)
	assert.NoError(t, err)
	defer func() {
		err := v2Stream.CloseSend()
		assert.NoError(t, err)
	}()

	// Regular flow.
	err = v2Stream.Send(&envoy_api_v2.DiscoveryRequest{})
	assert.NoError(t, err)

	_, err = v2Stream.Recv()
	assert.NoError(t, err)

	assert.Equal(t, int64(1), testServer.scope.Snapshot().Counters()["test.v2.totalResourcesServed+"].Value())

	// Error response from xDS client.
	err = v2Stream.Send(&envoy_api_v2.DiscoveryRequest{ErrorDetail: &rpc_status.Status{}})
	assert.NoError(t, err)

	assert.Equal(t, int64(1), testServer.scope.Snapshot().Counters()["test.v2.totalResourcesServed+"].Value())
	// Async verification here since it appears that we don't get a response back in this case, so we
	// aren't able to synchronize on the response.
	awaitCounterEquals(t, testServer.scope, "test.v2.totalErrorsReceived+", 1)
}

// Verifies that TTL and heartbeating is done when configured to do so.
func TestResourceTTL(t *testing.T) {
	testServer := newTestServer(t, true)
	defer testServer.stop()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	now := time.Now()
	config := serverexperimentation.HTTPFaultConfig{
		Fault: &serverexperimentation.HTTPFaultConfig_AbortFault{
			AbortFault: &serverexperimentation.AbortFault{
				Percentage:    &serverexperimentation.FaultPercentage{
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
	a, err := ptypes.MarshalAny(&config)
	assert.NoError(t, err)

	_, err = testServer.storer.CreateExperiment(context.Background(), a, &now, &now)
	assert.NoError(t, err)

	// First we look at a stream for a cluster that has an active fault. This should result in a TTL'd
	// resource that sends heartbeats.
	s, err := newV3Stream(ctx, "tests", "cluster", testServer)
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
	assert.NoError(t, ptypes.UnmarshalAny(r.Resources[0], resource))

	assert.Equal(t, int64(1), resource.Ttl.Seconds)

	r, err = s.sendV3RequestAndAwaitResponse(r.VersionInfo, r.Nonce)
	assert.NoError(t, err)

	resource = &gcpDiscoveryV3.Resource{}
	assert.NoError(t, ptypes.UnmarshalAny(r.Resources[0], resource))

	assert.Equal(t, int64(1), resource.Ttl.Seconds)
	assert.Nil(t, resource.Resource)

	// Second we look at a stream for a cluster that should not receive any faults.
	// Here we do not expect to see TTLs or heartbeats.
	noFaultStream, err := newV3Stream(ctx, "tests", "other-cluster", testServer)
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
	assert.NoError(t, ptypes.UnmarshalAny(r.Resources[0], runtime))

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

// Helper class for initializing a test RTDS server.
type testServer struct {
	registrar *moduletest.TestRegistrar
	scope     tally.TestScope
	storer    *simpleStorer
}

func newTestServer(t *testing.T, ttl bool) testServer {
	t.Helper()
	server := testServer{}

	server.storer = &simpleStorer{}
	service.Registry[experimentstore.Name] = server.storer

	// Set up a test server listening to :9000.
	config := &rtdsconfigv1.Config{
		RtdsLayerName:             "tests",
		CacheRefreshInterval:      ptypes.DurationProto(time.Second),
		IngressFaultRuntimePrefix: "ingress",
		EgressFaultRuntimePrefix:  "egress",
	}

	if ttl {
		config.ResourceTtl = &durationpb.Duration{
			Seconds: 1,
		}
		config.HeartbeatInterval = &durationpb.Duration{
			Seconds: 1,
		}
	}

	any, err := ptypes.MarshalAny(config)
	assert.NoError(t, err)
	server.scope = tally.NewTestScope("test", nil)

	logger, err := zap.NewDevelopment()
	assert.NoError(t, err)

	m, err := New(any, logger, server.scope)
	assert.NoError(t, err)

	server.registrar = moduletest.NewRegisterChecker()

	err = m.Register(server.registrar)
	assert.NoError(t, err)

	l, err := net.Listen("tcp", "localhost:9000")
	assert.NoError(t, err)

	go func() {
		err := server.registrar.GRPCServer().Serve(l)
		assert.NoError(t, err)
	}()

	return server
}

func (t *testServer) stop() {
	t.registrar.GRPCServer().Stop()
}

func (t *testServer) clientConn() (*grpc.ClientConn, error) {
	return grpc.Dial("localhost:9000", grpc.WithInsecure())
}

// Helper class for testing calls over a single RTDS stream.
type v3StreamWrapper struct {
	stream  gcpRuntimeServiceV3.RuntimeDiscoveryService_StreamRuntimeClient
	layer   string
	cluster string
}

func newV3Stream(ctx context.Context, layer string, cluster string, t testServer) (*v3StreamWrapper, error) {
	conn, err := t.clientConn()
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

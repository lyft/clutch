package rtds

import (
	"net"
	"testing"
	"time"

	"github.com/golang/protobuf/ptypes"
	rtdsconfigv1 "github.com/lyft/clutch/backend/api/config/module/chaos/experimentation/rtds/v1"
	rtds_testing "github.com/lyft/clutch/backend/module/chaos/serverexperimentation/rtds/testing"
	"github.com/lyft/clutch/backend/module/moduletest"
	"github.com/lyft/clutch/backend/service"
	"github.com/lyft/clutch/backend/service/chaos/experimentation/experimentstore"
	"github.com/stretchr/testify/assert"
	"github.com/uber-go/tally"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/durationpb"
)

// Helper class for initializing a test RTDS server.
type testServer struct {
	registrar *moduletest.TestRegistrar
	scope     tally.TestScope
	storer    *rtds_testing.SimpleStorer
}

func newTestServer(t *testing.T, ttl bool) testServer {
	t.Helper()
	server := testServer{}

	server.storer = &rtds_testing.SimpleStorer{}
	service.Registry[experimentstore.Name] = server.storer

	// Set up a test server listening to :9000.
	config := &rtdsconfigv1.Config{
		RtdsLayerName:             "rtds",
		CacheRefreshInterval:      ptypes.DurationProto(time.Second),
		IngressFaultRuntimePrefix: "fault.http",
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

	go func() {
		l, err := net.Listen("tcp", "localhost:9000")
		assert.NoError(t, err)

		err = server.registrar.GRPCServer().Serve(l)
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

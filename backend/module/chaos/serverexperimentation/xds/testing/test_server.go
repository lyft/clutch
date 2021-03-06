package testing

import (
	"github.com/lyft/clutch/backend/mock/service/chaos/experimentation/experimentstoremock"
	"net"
	"time"

	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/any"
	"github.com/uber-go/tally"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/durationpb"

	rtdsconfigv1 "github.com/lyft/clutch/backend/api/config/module/chaos/experimentation/xds/v1"
	"github.com/lyft/clutch/backend/module"
	"github.com/lyft/clutch/backend/module/moduletest"
	"github.com/lyft/clutch/backend/service"
	"github.com/lyft/clutch/backend/service/chaos/experimentation/experimentstore"
)

// Helper class for initializing a test RTDS server.
type TestServer struct {
	registrar *moduletest.TestRegistrar
	Scope     tally.TestScope
	Storer    *experimentstoremock.SimpleStorer
}

func NewTestServer(c func(cfg *any.Any, logger *zap.Logger, scope tally.Scope) (module.Module, error), ttl bool) TestServer {
	server := TestServer{}

	server.Storer = &experimentstoremock.SimpleStorer{}
	service.Registry[experimentstore.Name] = server.Storer

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

	server.registrar = moduletest.NewRegisterChecker()

	cfg, err := ptypes.MarshalAny(config)
	if err != nil {
		panic(err)
	}
	server.Scope = tally.NewTestScope("test", nil)

	logger, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}
	m, err := c(cfg, logger, server.Scope)
	if err != nil {
		panic(err)
	}

	err = m.Register(server.registrar)
	if err != nil {
		panic(err)
	}

	go func() {
		//nolint:gosec
		l, err := net.Listen("tcp", "0.0.0.0:9000")
		if err != nil {
			panic(err)
		}

		err = server.registrar.GRPCServer().Serve(l)
		if err != nil {
			panic(err)
		}
	}()

	return server
}

func (t *TestServer) Stop() {
	t.registrar.GRPCServer().Stop()
}

func (t *TestServer) ClientConn() (*grpc.ClientConn, error) {
	return grpc.Dial("localhost:9000", grpc.WithInsecure())
}

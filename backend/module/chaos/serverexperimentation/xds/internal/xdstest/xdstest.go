package xdstest

import (
	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/any"
	rtdsconfigv1 "github.com/lyft/clutch/backend/api/config/module/chaos/experimentation/xds/v1"
	"github.com/lyft/clutch/backend/internal/test/integration/helper/testserver"
	"github.com/lyft/clutch/backend/mock/service/chaos/experimentation/experimentstoremock"
	"github.com/lyft/clutch/backend/module"
	"github.com/lyft/clutch/backend/service"
	"github.com/lyft/clutch/backend/service/chaos/experimentation/experimentstore"
	"github.com/uber-go/tally"
	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/anypb"
	"google.golang.org/protobuf/types/known/durationpb"
	"time"
)

// Helper class for initializing a test RTDS server.
type TestModuleServer struct {
	*xdstest.TestServer

	Storer *experimentstoremock.SimpleStorer
}

func NewTestModuleServer(c func(cfg *any.Any, logger *zap.Logger, scope tally.Scope) (module.Module, error), ttl bool) *TestModuleServer {
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

	cfg, err := anypb.New(config)
	if err != nil {
		panic(err)
	}

	ms := &TestModuleServer{
		TestServer: xdstest.New(),
		Storer:     &experimentstoremock.SimpleStorer{},
	}
	service.Registry[experimentstore.Name] = ms.Storer

	m, err := c(cfg, ms.Logger, ms.Scope)
	if err != nil {
		panic(err)
	}

	ms.Register(m)

	ms.Run()

	return ms

}


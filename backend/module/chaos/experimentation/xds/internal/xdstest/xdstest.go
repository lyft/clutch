package xdstest

import (
	"time"

	"github.com/golang/protobuf/ptypes/any"
	"github.com/uber-go/tally"
	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/anypb"
	"google.golang.org/protobuf/types/known/durationpb"

	xdsconfigv1 "github.com/lyft/clutch/backend/api/config/module/chaos/experimentation/xds/v1"
	xdstest "github.com/lyft/clutch/backend/internal/test/integration/helper/testserver"
	"github.com/lyft/clutch/backend/mock/service/chaos/experimentation/experimentstoremock"
	"github.com/lyft/clutch/backend/module"
	"github.com/lyft/clutch/backend/service"
	"github.com/lyft/clutch/backend/service/chaos/experimentation/experimentstore"
)

// Helper class for initializing a test RTDS server.
type TestModuleServer struct {
	*xdstest.TestServer

	Storer *experimentstoremock.SimpleStorer
}

func NewTestModuleServer(c func(cfg *any.Any, logger *zap.Logger, scope tally.Scope) (module.Module, error), ttl bool, config *xdsconfigv1.Config) (*TestModuleServer, error) {
	// Set up a test server listening to :9000.
	if ttl {
		config.ResourceTtl = &durationpb.Duration{
			Seconds: 2,
		}
		config.HeartbeatInterval = &durationpb.Duration{
			Seconds: 1,
		}
	}

	if config.CacheRefreshInterval == nil {
		config.CacheRefreshInterval = durationpb.New(time.Second * 1)
	}

	cfg, err := anypb.New(config)
	if err != nil {
		return nil, err
	}

	ms := &TestModuleServer{
		TestServer: xdstest.New(),
		Storer:     &experimentstoremock.SimpleStorer{},
	}
	service.Registry[experimentstore.Name] = ms.Storer

	m, err := c(cfg, ms.Logger, ms.Scope)
	if err != nil {
		return nil, err
	}

	ms.Register(m)

	err = ms.Run()
	if err != nil {
		return nil, err
	}

	return ms, nil
}

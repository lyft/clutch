package topology

import (
	"database/sql"
	"errors"
	"os"
	"os/signal"
	"syscall"

	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/any"
	"github.com/uber-go/tally"
	"go.uber.org/zap"
	"golang.org/x/net/context"

	topologyv1cfg "github.com/lyft/clutch/backend/api/config/service/topology/v1"
	topologyv1 "github.com/lyft/clutch/backend/api/topology/v1"
	"github.com/lyft/clutch/backend/service"
	pgservice "github.com/lyft/clutch/backend/service/db/postgres"
)

const Name = "clutch.service.topology"

type Service interface{}

type client struct {
	config *topologyv1cfg.Config

	db    *sql.DB
	log   *zap.Logger
	scope tally.Scope
}

// CacheableTopology is implemented by a service that wishes to enable the topology API feature set
//
// By implementing this interface the topology service will automatically setup all services which implement it.
// Automatically ingesting Resource objects via the `GetTopologyObjectChannel()` function.
// This enables users to make use of the Topology APIs with these new Topology Resources.
//
type CacheableTopology interface {
	CacheEnabled() bool
	StartTopologyCaching(ctx context.Context) (<-chan *topologyv1.UpdateCacheRequest, error)
}

func New(cfg *any.Any, logger *zap.Logger, scope tally.Scope) (service.Service, error) {
	topologyConfig := &topologyv1cfg.Config{}
	err := ptypes.UnmarshalAny(cfg, topologyConfig)
	if err != nil {
		return nil, err
	}

	p, ok := service.Registry[pgservice.Name]
	if !ok {
		return nil, errors.New("Please config the datastore [clutch.service.db.postgres] to use the topology service")
	}

	dbClient, ok := p.(pgservice.Client)
	if !ok {
		return nil, errors.New("Unable to get the datastore client")
	}

	c := &client{
		config: topologyConfig,
		db:     dbClient.DB(),
		log:    logger,
		scope:  scope,
	}

	if topologyConfig.Cache == nil {
		c.log.Info("Topology caching is disabled")
		return c, nil
	}

	ctx, ctxCancelFunc := context.WithCancel(context.Background())
	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	go func() {
		<-sigc
		c.log.Info("Caught shutdown signal, shutting down topology caching and releasing advisory lock")
		ctxCancelFunc()
	}()

	c.log.Info("Topology caching is enabled")
	go c.acquireTopologyCacheLock(ctx)

	return c, nil
}

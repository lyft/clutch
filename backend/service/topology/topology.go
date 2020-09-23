package topology

import (
	"database/sql"
	"errors"

	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/any"
	"github.com/uber-go/tally"
	"go.uber.org/zap"

	topologyconfigv1 "github.com/lyft/clutch/backend/api/config/service/topology/v1"
	topologyv1 "github.com/lyft/clutch/backend/api/topology/v1"
	"github.com/lyft/clutch/backend/service"
	pgservice "github.com/lyft/clutch/backend/service/db/postgres"
)

const Name = "clutch.service.topology"

type Service interface {
	// GetByID(ctx context.Context, key string, resolverTypeUrl string)
	// GetByLabel(ctx context.Context, labels map[string]string, resolverTypeUrl string)

	SetCache(obj topologyv1.TopologyObject)
	DeleteCache(obj topologyv1.TopologyObject)

	// LeaderElect()
	// deleteExpiredCache()
}

type client struct {
	config *topologyconfigv1.Config

	db    *sql.DB
	log   *zap.Logger
	scope tally.Scope
}

// CacheableTopology is implemented by a service that wishes to enable the topology API feature set
//
// By implmenting this interface the topology service will automatically setup all services which
// implement this interface. Automatically ingesting TopologyObjects via the `GetTopologyObjectChannel()` function.
// This enables users to make use of the Topology APIs with these new TopologyObjects.
//
type CacheableTopology interface {
	CacheEnabled() bool
	GetTopologyObjectChannel() chan topologyv1.TopologyCacheObject
}

func New(cfg *any.Any, logger *zap.Logger, scope tally.Scope) (service.Service, error) {
	p, ok := service.Registry[pgservice.Name]
	if !ok {
		return nil, errors.New("could not find database service")
	}

	dbClient, ok := p.(pgservice.Client)
	if !ok {
		return nil, errors.New("Unable to get dbClient")
	}

	topologyConfig := &topologyconfigv1.Config{}
	err := ptypes.UnmarshalAny(cfg, topologyConfig)
	if err != nil {
		return nil, err
	}

	c, err := &client{
		config: topologyConfig,
		db:     dbClient.DB(),
		log:    logger,
		scope:  scope,
	}, nil

	c.startCaching()

	return c, err
}

func (c *client) LeaderElect() (bool, error) {
	var lock bool
	if err := c.db.QueryRow("SELECT pg_try_advisory_lock(100)").Scan(&lock); err != nil {
		return lock, err
	}

	return lock, nil
}

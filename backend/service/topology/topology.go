package topology

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"time"

	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/any"
	"github.com/uber-go/tally"
	"go.uber.org/zap"

	topologyv1 "github.com/lyft/clutch/backend/api/config/service/topology/v1"
	"github.com/lyft/clutch/backend/service"
	pgservice "github.com/lyft/clutch/backend/service/db/postgres"
)

const Name = "clutch.service.topology"

type Service interface {
	aquireTopologyCacheLock()
	startTopologyCache()
}

type client struct {
	config *topologyv1.Config

	db    *sql.DB
	log   *zap.Logger
	scope tally.Scope
}

func New(cfg *any.Any, logger *zap.Logger, scope tally.Scope) (service.Service, error) {
	topologyConfig := &topologyv1.Config{}
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

	c, err := &client{
		config: topologyConfig,
		db:     dbClient.DB(),
		log:    logger,
		scope:  scope,
	}, nil

	go c.aquireTopologyCacheLock()

	return c, err
}

// Perfroms leader election by aquiring an postgres advisory lock
func (c *client) aquireTopologyCacheLock() {
	var lock bool

	// We create our own connection to use for aquirting the advisory lock
	// If the connection is severed for any reason the advisory lock will automatically unlock
	conn, err := c.db.Conn(context.Background())
	if err != nil {
		log.Printf("err: %v", err)
	}

	// Infinately try to aquire the advisory lock
	// Once the lock is aquired we start caching, this is a blocking operation
	for {
		c.log.Debug("trying to aquire advisory lock")

		// Advisory locks only take a arbatrary bigint value
		// In this case 100 is the lock for topology caching
		_ = conn.QueryRowContext(context.Background(), "SELECT pg_try_advisory_lock(100)").Scan(&lock)
		if err != nil {
			log.Printf("err: %v", err)
		}

		if lock {
			c.log.Debug("aquired the advisory lock, starting to cache topology now...")
			c.startTopologyCache()
		}

		time.Sleep(time.Second * 10)
	}
}

func (c *client) startTopologyCache() {
	c.log.Debug("implement me")
}

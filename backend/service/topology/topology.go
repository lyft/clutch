package topology

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/apex/log"
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

func (c *client) aquireTopologyCacheLock() {
	var lock bool

	conn, err := c.db.Conn(context.Background())
	if err != nil {
		log.Errorf("err: %v", err)
	}

	for {
		_ = conn.QueryRowContext(context.Background(), "SELECT pg_try_advisory_lock(100)").Scan(&lock)
		if err != nil {
			log.Errorf("err: %v", err)
		}

		if lock {
			c.startTopologyCache()
		}

		time.Sleep(time.Second * 10)
	}
}

func (c *client) startTopologyCache() {
	time.Sleep(time.Hour * 2)
}

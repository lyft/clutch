package topology

import (
	"database/sql"
	"errors"

	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/any"
	"github.com/uber-go/tally"
	"go.uber.org/zap"

	topologyv1 "github.com/lyft/clutch/backend/api/config/service/topology/v1"
	"github.com/lyft/clutch/backend/service"
	pgservice "github.com/lyft/clutch/backend/service/db/postgres"
)

const Name = "clutch.service.topology"

type Service interface{}

type client struct {
	config *topologyv1.Config

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
	GetTopologyObjectChannel() chan TopologyObject
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

	return c, err
}

package topology

import (
	"context"
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

type Service interface {
	GetByID(ctx context.Context, key string, resolverTypeUrl string)
	GetByLabel(ctx context.Context, labels map[string]string, resolverTypeUrl string)

	SetCache(ctx context.Context, key string, resolverTypeUrl string, data []byte, metadata map[string]string)
	DeleteCache(ctx context.Context, key string)

	deleteExpiredCache()
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

	return c, err
}

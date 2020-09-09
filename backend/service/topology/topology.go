package topology

import (
	"context"
	"database/sql"
	"errors"
	"log"

	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/any"
	"github.com/uber-go/tally"
	"go.uber.org/zap"

	topologyv1 "github.com/lyft/clutch/backend/api/config/service/topology/v1"
	"github.com/lyft/clutch/backend/service"
	pgservice "github.com/lyft/clutch/backend/service/db/postgres"
	k8sservice "github.com/lyft/clutch/backend/service/k8s"
)

const Name = "clutch.service.topology"

type Service interface {
	// GetByID(ctx context.Context, key string, resolverTypeUrl string)
	// GetByLabel(ctx context.Context, labels map[string]string, resolverTypeUrl string)

	SetCache(key string, resolverTypeUrl string, data []byte)
	DeleteCache(key string)

	// DeleteExpiredCache()
	// LeaderElect()
	// ManageCache()
}

type Client struct {
	config *topologyv1.Config

	db    *sql.DB
	log   *zap.Logger
	scope tally.Scope
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

	topologyConfig := &topologyv1.Config{}
	err := ptypes.UnmarshalAny(cfg, topologyConfig)
	if err != nil {
		return nil, err
	}

	c, err := &Client{
		config: topologyConfig,
		db:     dbClient.DB(),
		log:    logger,
		scope:  scope,
	}, nil

	c.startCaching()

	return c, err
}

func (c *Client) startCaching() {
	// if k8s is enabled then cache it
	k8sClient, _ := service.Registry["clutch.service.k8s"]
	svc, _ := k8sClient.(k8sservice.Service)
	svc.ManageCache(c)

}

func (c *Client) DeleteCache(id string) {
	const deleteQuery = `
		DELETE FROM topology_cache WHERE id = $1
	`
	_, err := c.db.ExecContext(context.Background(), deleteQuery, id)
	if err != nil {
		log.Printf("%v", err)
		return
	}
}

func (c *Client) SetCache(id string, resolver_type_url string, data []byte) {
	const upsertQuery = `
		INSERT INTO topology_cache (id, data, resolver_type_url)
		VALUES ($1, $2, $3)
		ON CONFLICT (id) DO UPDATE SET
			id = EXCLUDED.id,
			data = EXCLUDED.data,
			resolver_type_url = EXCLUDED.resolver_type_url
	`

	_, err := c.db.ExecContext(context.Background(), upsertQuery, id, data, resolver_type_url)
	if err != nil {
		log.Printf("%v", err)
		return
	}
}

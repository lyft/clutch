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

type Client interface {
	GetByID(ctx context.Context, key string, resolverTypeUrl string)
	GetByLabel(ctx context.Context, labels map[string]string, resolverTypeUrl string)

	SetCache(ctx context.Context, key string, resolverTypeUrl string, data any.Any)

	DeleteExpiredCache()
	LeaderElect()
	ManageCache()
}

type client struct {
	config *topologyv1.Config

	isLeader bool
	db       *sql.DB
	log      *zap.Logger
	scope    tally.Scope
}

const Name = "clutch.service.topology"

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

	c, err := &client{
		config: topologyConfig,
		db:     dbClient.DB(),
		log:    logger,
		scope:  scope,
	}, nil

	c.PopulateCacheFromKubernetes()

	return c, err
}

// pretend we are the leader
func (c *client) PopulateCacheFromKubernetes() {
	log.Print("topology is enabled and starting k8s cache.")

	// if k8s is enabled then cache it
	client, ok := service.Registry["clutch.service.k8s"]
	if !ok {
		return
	}

	svc, ok := client.(k8sservice.Service)
	if !ok {
		return
	}

	svc.PopulateCache(c.db)
}

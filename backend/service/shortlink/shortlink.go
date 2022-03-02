package shortlink

import (
	"context"
	"database/sql"
	"errors"

	"github.com/golang/protobuf/ptypes/any"
	"github.com/uber-go/tally/v4"
	"go.uber.org/zap"

	shortlinkv1 "github.com/lyft/clutch/backend/api/shortlink/v1"
	"github.com/lyft/clutch/backend/service"
	pgservice "github.com/lyft/clutch/backend/service/db/postgres"
)

const (
	Name = "clutch.service.shortlink"
)

type Service interface {
	Create(context.Context, string, []*shortlinkv1.ShareableState) (string, error)
	Get(context.Context, string) (string, []*shortlinkv1.ShareableState, error)
}

type client struct {
	db    *sql.DB
	log   *zap.Logger
	scope tally.Scope
}

func New(cfg *any.Any, logger *zap.Logger, scope tally.Scope) (service.Service, error) {
	p, ok := service.Registry[pgservice.Name]
	if !ok {
		return nil, errors.New("Please config the datastore [clutch.service.db.postgres] to use the shortlink service")
	}

	dbClient, ok := p.(pgservice.Client)
	if !ok {
		return nil, errors.New("Unable to get the datastore client")
	}

	c := &client{
		db:    dbClient.DB(),
		log:   logger,
		scope: scope,
	}

	return c, nil
}

func (c *client) Create(ctx context.Context, path string, state []*shortlinkv1.ShareableState) (string, error) {
	return "", errors.New("not implemented")
}

func (c *client) Get(ctx context.Context, hash string) (string, []*shortlinkv1.ShareableState, error) {
	return "", nil, errors.New("not implemented")
}

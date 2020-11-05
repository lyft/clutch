package authn

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/lyft/clutch/backend/service"
	"github.com/lyft/clutch/backend/service/db/postgres"
)

type authnUser struct {
	id                   string
	providerRefreshToken []byte
}

type repository struct {
	db *sql.DB
}

func newRepository() (*repository, error) {
	svcName := postgres.Name
	svc, ok := service.Registry[svcName]
	if !ok {
		return nil, fmt.Errorf("database '%s' not registered", svcName)
	}

	pg, ok := svc.(postgres.Client)
	if !ok {
		return nil, fmt.Errorf("database does not implement the required interface")
	}

	return &repository{db: pg.DB()}, nil
}

const createOrUpdateUser = `
INSERT INTO authn_users (id, provider_refresh_token) VALUES ($1, $2)
ON CONFLICT (id)
DO UPDATE SET
    provider_refresh_token = EXCLUDED.provider_refresh_token
RETURNING id, provider_refresh_token
`

func (r *repository) createOrUpdateUser(ctx context.Context, id string, providerRefreshToken []byte) (*authnUser, error) {
	q := r.db.QueryRowContext(ctx, createOrUpdateUser, id, providerRefreshToken)

	ret := &authnUser{}
	err := q.Scan(&ret.id, &ret.providerRefreshToken)
	if err != nil {
		return nil, err
	}

	return ret, nil
}

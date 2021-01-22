package authn

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/lyft/clutch/backend/service"
	"github.com/lyft/clutch/backend/service/db/postgres"
)

type authnToken struct {
	userID   string
	provider string

	accessToken  []byte
	refreshToken []byte
	idToken      []byte
	expiry       time.Time
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

// #nosec G101
const createOrUpdateProviderToken = `
INSERT INTO authn_tokens (user_id, provider, access_token, refresh_token, id_token, expiry) VALUES ($1, $2, $3, $4, $5, $6)
ON CONFLICT (user_id, provider) DO UPDATE SET
    user_id = EXCLUDED.user_id,
    provider = EXCLUDED.provider,
    access_token = EXCLUDED.access_token,
    refresh_token = EXCLUDED.refresh_token,
    id_token = EXCLUDED.id_token,
	expiry = EXCLUDED.expiry
`

func (r *repository) createOrUpdateProviderToken(ctx context.Context, token *authnToken) error {
	_, err := r.db.ExecContext(ctx, createOrUpdateProviderToken,
		token.userID, token.provider, token.accessToken, token.refreshToken, token.idToken, token.expiry)

	return err
}

// #nosec G101
const readProviderToken = `
SELECT user_id, provider, access_token, refresh_token, id_token, expiry FROM authn_tokens WHERE user_id = $1 AND provider = $2
`

func (r *repository) readProviderToken(ctx context.Context, userID, provider string) (*authnToken, error) {
	t := &authnToken{}

	q := r.db.QueryRowContext(ctx, readProviderToken, userID, provider)
	err := q.Scan(&t.userID, &t.provider, &t.accessToken, &t.refreshToken, &t.idToken, &t.expiry)
	if err != nil {
		return nil, err
	}

	return t, nil
}

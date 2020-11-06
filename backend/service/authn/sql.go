package authn

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/lyft/clutch/backend/service"
	"github.com/lyft/clutch/backend/service/db/postgres"
)

type authnToken struct {
	userID   string
	provider string

	tokenType    string
	idToken      []byte
	accessToken  []byte
	refreshToken []byte
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
INSERT INTO authn_tokens (user_id, provider, token_type, id_token, access_token, refresh_token) VALUES ($1, $2, $3, $4, $5, $6)
ON CONFLICT DO UPDATE SET
    user_id = EXCLUDED.user_id,
    provider = EXCLUDED.provider,
    token_type = EXCLUDED.token_type,
    id_token = EXCLUDED.id_token,
    access_token = EXCLUDED.access_token,
    refresh_token = EXCLUDED.refresh_token
`

func (r *repository) createOrUpdateProviderToken(ctx context.Context, token *authnToken) error {
	_, err := r.db.ExecContext(ctx, createOrUpdateProviderToken,
		token.userID, token.provider, token.tokenType, token.idToken, token.accessToken, token.refreshToken)

	return err
}

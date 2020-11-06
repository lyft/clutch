package authn

import (
	"context"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestCreateOrUpdateUser(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	assert.NoError(t, err)
	defer db.Close()

	tok := &authnToken{
		userID:       "foo@example.com",
		provider:     "okta",
		tokenType:    "oidc",
		idToken:      []byte("id"),
		accessToken:  []byte("access"),
		refreshToken: []byte("refresh"),
	}

	mock.ExpectExec(createOrUpdateProviderToken).
		WithArgs(tok.userID, tok.provider, tok.tokenType, tok.idToken, tok.accessToken, tok.refreshToken).
		WillReturnResult(sqlmock.NewResult(0, 1))

	r := &repository{db: db}
	err = r.createOrUpdateProviderToken(context.Background(), tok)
	assert.NoError(t, err)
	//assert.NoError(t, mock.ExpectationsWereMet())
}

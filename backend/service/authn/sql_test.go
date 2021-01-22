package authn

import (
	"context"
	"testing"
	"time"

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
		idToken:      []byte("id"),
		accessToken:  []byte("access"),
		refreshToken: []byte("refresh"),
		expiry:       time.Now(),
	}

	mock.ExpectExec(createOrUpdateProviderToken).
		WithArgs(tok.userID, tok.provider, tok.accessToken, tok.refreshToken, tok.idToken, tok.expiry).
		WillReturnResult(sqlmock.NewResult(0, 1))

	r := &repository{db: db}
	err = r.createOrUpdateProviderToken(context.Background(), tok)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

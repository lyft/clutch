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

	id := "foo@example.com"
	refreshToken := []byte("abdcdefgzzz")

	r := &repository{db: db}
	mock.ExpectQuery(createOrUpdateUser).
		WithArgs(id, refreshToken).
		WillReturnRows(sqlmock.NewRows([]string{"id", "provider_refresh_token"}).AddRow(id, refreshToken))

	u, err := r.createOrUpdateUser(context.Background(), id, refreshToken)
	assert.NoError(t, err)
	assert.Equal(t, id, u.id)
	assert.Equal(t, refreshToken, u.providerRefreshToken)
}

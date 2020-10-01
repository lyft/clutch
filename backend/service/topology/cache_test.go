package topology

import (
	"context"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestTryAdvisoryLock(t *testing.T) {
	db, _, err := sqlmock.New()
	assert.NoError(t, err)

	conn, err := db.Conn(context.Background())
	assert.NoError(t, err)

	lock := tryAdvisoryLock(conn, 1)
	assert.True(t, lock)
}

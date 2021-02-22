package topology

import (
	"context"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestConvertLockIdToAdvisoryLockId(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		id     string
		input  string
		expect uint32
	}{
		{
			id:     "key with chars",
			input:  "topologycache",
			expect: 1953460335,
		},
		{
			id:     "key with special chars",
			input:  "*()#@&!*(#!@",
			expect: 707275043,
		},
	}

	for _, tt := range testCases {
		tt := tt
		t.Run(tt.id, func(t *testing.T) {
			t.Parallel()
			id := convertLockIdToAdvisoryLockId(tt.input)
			assert.Equal(t, tt.expect, id)
		})
	}
}

func TestGetCacheLockConn(t *testing.T) {
	db, _, err := sqlmock.New()
	assert.NoError(t, err)

	topo := &client{
		db: db,
	}

	// Assert that the cacheLockConn get created
	err = topo.getCacheLockConn(context.Background())
	assert.NoError(t, err)
	assert.NotNil(t, topo.cacheLockConn)
}

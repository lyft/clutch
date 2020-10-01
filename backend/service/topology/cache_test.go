package topology

import (
	"log"
	"testing"

	topologyv1 "github.com/lyft/clutch/backend/api/topology/v1"

	"database/sql/driver"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/golang/protobuf/ptypes"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap/zaptest"
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
			log.Printf("%v", id)
			assert.Equal(t, tt.expect, id)
		})
	}
}

var topologyCacheColumns = []string{
	"id",
	"data",
	"resolver_type_url",
	"metadata",
	"updated_at",
}

func TestSetCache(t *testing.T) {
	log := zaptest.NewLogger(t)
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)

	topologyClient := &client{
		log: log,
		db:  db,
	}

	pb, err := ptypes.MarshalAny(&topologyv1.Resource{})
	assert.NoError(t, err)

	cacheItem := &topologyv1.Resource{
		Id: "pod-name",
		Pb: pb,
		Metadata: map[string]string{
			"key": "value",
		},
	}

	mock.ExpectExec("INSERT INTO topology_cache (id, resolver_type_url, data, metadata) VALUES ($1, $2, $3, $4)").WithArgs([]driver.Value{
		sqlmock.AnyArg(),
	})

	topologyClient.SetCache(cacheItem)

	// we make sure that all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

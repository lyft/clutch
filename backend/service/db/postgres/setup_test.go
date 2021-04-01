package postgres

import (
	"testing"

	"github.com/golang/protobuf/ptypes/any"
	"github.com/stretchr/testify/assert"
	"github.com/uber-go/tally"
	"go.uber.org/zap/zaptest"
	"google.golang.org/protobuf/types/known/anypb"

	postgresv1 "github.com/lyft/clutch/backend/api/config/service/db/postgres/v1"
)

func TestNew(t *testing.T) {
	t.Parallel()

	scope := tally.NewTestScope("", nil)
	logger := zaptest.NewLogger(t)
	cfg := &any.Any{TypeUrl: "type.googleapis.com/clutch.config.service.db.postgres.v1.Config"}

	svc, err := New(cfg, logger, scope)
	assert.Error(t, err)
	assert.Nil(t, svc)

	pgconfig := &postgresv1.Config{
		Connection: &postgresv1.Connection{
			Host:   "localhost",
			Port:   5432,
			User:   "clutch",
			Dbname: "clutch",
		},
	}

	cfg, err = anypb.New(pgconfig)
	assert.NoError(t, err)

	svc, err = New(cfg, logger, scope)
	assert.NoError(t, err)
	assert.NotNil(t, svc)

	client, ok := svc.(Client)
	assert.True(t, ok)
	assert.NotNil(t, client.DB())
}

func TestConnString(t *testing.T) {
	t.Parallel()

	tests := []struct {
		cfg      *postgresv1.Connection
		expected string
		err      bool
	}{
		{
			cfg: &postgresv1.Connection{
				Host:    "localhost",
				Port:    5432,
				User:    "clutch",
				Dbname:  "clutch",
				SslMode: 6,
			},
			expected: "host=localhost port=5432 dbname=clutch user=clutch sslmode=verify-full",
		},
		{
			cfg: &postgresv1.Connection{
				Host:   "localhost",
				Port:   5432,
				User:   "clutch",
				Dbname: "clutch",
				Authn:  &postgresv1.Connection_Password{Password: "password"},
			},
			expected: "host=localhost port=5432 dbname=clutch user=clutch password=password",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.expected, func(t *testing.T) {
			t.Parallel()

			result, err := connString(tt.cfg)
			assert.Equal(t, tt.err, err != nil)
			assert.Equal(t, tt.expected, result)
		})
	}
}

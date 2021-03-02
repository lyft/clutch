package postgres

// <!-- START clutchdoc -->
// description: Provides a connection to the configured PostgreSQL database.
// <!-- END clutchdoc -->

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/any"
	_ "github.com/lib/pq"
	"github.com/uber-go/tally"
	"go.uber.org/zap"

	postgresv1 "github.com/lyft/clutch/backend/api/config/service/db/postgres/v1"
	"github.com/lyft/clutch/backend/service"
)

const Name = "clutch.service.db.postgres"

type client struct {
	sqlDB  *sql.DB
	logger *zap.Logger
	scope  tally.Scope
}

type Client interface {
	DB() *sql.DB
}

func (c *client) DB() *sql.DB { return c.sqlDB }

func New(cfg *any.Any, logger *zap.Logger, scope tally.Scope) (service.Service, error) {
	pgcfg := &postgresv1.Config{}
	err := ptypes.UnmarshalAny(cfg, pgcfg)
	if err != nil {
		return nil, err
	}

	connection, err := connString(pgcfg.Connection)
	if err != nil {
		return nil, err
	}

	sqlDB, err := sql.Open("postgres", connection)
	if err != nil {
		return nil, err
	}

	// This pointer  for MaxIdleConnections so we can check the exsistance of this variable.
	// Ints are not capable of being nil the default is the value 0, which is a valid input
	// in this case and we want to fallback to the drivers default if nothing is set.
	maxIdleConnections := &pgcfg.MaxIdleConnections
	if maxIdleConnections != nil {
		sqlDB.SetMaxIdleConns(int(pgcfg.MaxIdleConnections))
	}

	return &client{logger: logger, scope: scope, sqlDB: sqlDB}, nil
}

func connString(cfg *postgresv1.Connection) (string, error) {
	if cfg == nil {
		return "", errors.New("no connection information")
	}

	connection := fmt.Sprintf(
		"host=%s port=%d dbname=%s user=%s",
		cfg.Host, cfg.Port, cfg.Dbname, cfg.User,
	)

	switch cfg.GetSslMode() {
	case postgresv1.Connection_UNSPECIFIED:
		break
	default:
		mode := strings.ReplaceAll(strings.ToLower(cfg.SslMode.String()), "_", "-")
		connection += fmt.Sprintf(" sslmode=%s", mode)
	}

	switch cfg.GetAuthn().(type) {
	case *postgresv1.Connection_Password:
		connection += fmt.Sprintf(" password=%s", cfg.GetPassword())
	default:
		break
	}

	return connection, nil
}

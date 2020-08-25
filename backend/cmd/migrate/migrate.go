package main

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
	"github.com/uber-go/tally"
	"go.uber.org/zap"

	"github.com/lyft/clutch/backend/gateway"
	clutchpg "github.com/lyft/clutch/backend/service/db/postgres"
)

func run() error {
	f := gateway.ParseFlags()
	cfg := gateway.MustReadOrValidateConfig(f)

	logger, _ := zap.NewProduction()

	var sqlDB *sql.DB
	for _, s := range cfg.Services {
		if s.Name == clutchpg.Name {
			pgdb, err := clutchpg.New(s.TypedConfig, logger, tally.NoopScope)
			if err != nil {
				logger.Fatal("error creating db", zap.Error(err))
			}
			sqlDB = pgdb.(clutchpg.Client).DB()
			break
		}
	}

	if sqlDB == nil {
		logger.Fatal("no database found in config", zap.String("file", f.ConfigPath))
	}

	if err := sqlDB.Ping(); err != nil {
		return fmt.Errorf("error pinging db: %w", err)
	}

	driver, err := postgres.WithInstance(sqlDB, &postgres.Config{})
	if err != nil {
		return fmt.Errorf("error creating pg driver: %w", err)
	}

	migrationDir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("could not get working dir: %w", err)
	}
	migrationDir = filepath.Join(migrationDir, "migrations")

	m, err := migrate.NewWithDatabaseInstance(
		fmt.Sprintf("file://%s", migrationDir),
		"postgres", driver)
	if err != nil {
		return fmt.Errorf("error creating migrator: %w", err)
	}

	err = m.Up()
	if err != nil {
		return fmt.Errorf("failed running migrations: %w", err)
	}

	return nil
}

func main() {
	if err := run(); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(-1)
	}
}

package main

import (
	"database/sql"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/golang/protobuf/ptypes"
	_ "github.com/lib/pq"
	"github.com/uber-go/tally"
	"go.uber.org/zap"

	postgresv1 "github.com/lyft/clutch/backend/api/config/service/db/postgres/v1"
	"github.com/lyft/clutch/backend/gateway"
	clutchpg "github.com/lyft/clutch/backend/service/db/postgres"
)

type MigrateFlags struct {
	Force     bool
	BaseFlags *gateway.Flags
}

func (m *MigrateFlags) Link() {
	m.BaseFlags = &gateway.Flags{}
	m.BaseFlags.Link()

	flag.BoolVar(&m.Force, "f", false, "do not ask user for confirmation")
}

func run() error {
	f := &MigrateFlags{}
	f.Link()
	flag.Parse()

	cfg := gateway.MustReadOrValidateConfig(f.BaseFlags)

	logger, _ := zap.NewDevelopment()
	logger = logger.WithOptions(zap.AddStacktrace(zap.FatalLevel + 1))

	var sqlDB *sql.DB
	var hostInfo string
	for _, s := range cfg.Services {
		if s.Name == clutchpg.Name {
			pgdb, err := clutchpg.New(s.TypedConfig, logger, tally.NoopScope)
			if err != nil {
				logger.Fatal("error creating db", zap.Error(err))
			}

			cfg := &postgresv1.Config{}
			if err := ptypes.UnmarshalAny(s.TypedConfig, cfg); err != nil {
				logger.Fatal("could not convert config", zap.Error(err))
			}

			sqlDB = pgdb.(clutchpg.Client).DB()
			hostInfo = fmt.Sprintf("%s@%s:%d", cfg.Connection.User, cfg.Connection.Host, cfg.Connection.Port)

			break
		}
	}

	if sqlDB == nil {
		logger.Fatal("no database found in config", zap.String("file", f.BaseFlags.ConfigPath))
	}
	fmt.Printf("connecting to '%s' for migration\n", hostInfo)
	if !f.Force {
		fmt.Printf("this could cause irrevocable data loss, is this okay? [Y/n] ")
		var answer string
		if _, err := fmt.Scanln(&answer); err != nil && err.Error() != "unexpected newline" {
			logger.Fatal("could not read user input", zap.Error(err))
		}
		if answer != "" || strings.ToLower(answer) == "y" {
			fmt.Println("aborting")
			os.Exit(1)
		}
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

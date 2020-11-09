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

type migrateLogger struct {
	logger *zap.SugaredLogger
}

func (m *migrateLogger) Printf(format string, v ...interface{}) {
	m.logger.Infof(strings.TrimRight(format, "\n"), v...)
}

func (m *migrateLogger) Verbose() bool {
	return true
}

func Run() {
	// Read flags and config.
	f := &MigrateFlags{}
	f.Link()
	flag.Parse()

	cfg := gateway.MustReadOrValidateConfig(f.BaseFlags)

	logger, _ := zap.NewDevelopment()
	logger = logger.WithOptions(zap.AddStacktrace(zap.FatalLevel + 1))

	// Find the database in config and instantiate the service.
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

	// Verify that user wants to continue (unless -f for force is passed as a flag).
	logger.Info("using database", zap.String("hostInfo", hostInfo))
	if !f.Force {
		logger.Warn("migration has the potential to cause irrevocable data loss, verify host information above")

		fmt.Printf("\n*** Continue with migration? [y/N] ")
		var answer string
		if _, err := fmt.Scanln(&answer); err != nil && err.Error() != "unexpected newline" {
			logger.Fatal("could not read user input", zap.Error(err))
		}
		if strings.ToLower(answer) != "y" {
			logger.Fatal("aborting, enter 'y' to continue or use the '-f' (force) option")
		}
		fmt.Println()
	}

	// Ping database and bring up driver.
	if err := sqlDB.Ping(); err != nil {
		logger.Fatal("error pinging db", zap.Error(err))
	}

	driver, err := postgres.WithInstance(sqlDB, &postgres.Config{})
	if err != nil {
		logger.Fatal("error creating pg driver", zap.Error(err))
	}

	// Create migrator.
	migrationDir, err := os.Getwd()
	if err != nil {
		logger.Fatal("could not get working dir", zap.Error(err))
	}
	migrationDir = filepath.Join(migrationDir, "migrations")

	m, err := migrate.NewWithDatabaseInstance(
		fmt.Sprintf("file://%s", migrationDir),
		"postgres", driver)
	if err != nil {
		logger.Fatal("error creating migrator", zap.Error(err))
	}

	m.Log = &migrateLogger{
		logger: logger.Sugar(),
	}

	// Apply migrations!
	logger.Info("applying migrations", zap.String("migrationDir", migrationDir))
	err = m.Up()
	if err != nil && err != migrate.ErrNoChange {
		logger.Fatal("failed running migrations", zap.Error(err))
	}
}

func main() {
	Run()
}

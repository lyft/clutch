package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
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

	gatewayv1 "github.com/lyft/clutch/backend/api/config/gateway/v1"
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

// Migrator handles all migration operations both up and down.
type Migrator struct {
	log   *zap.Logger
	flags *MigrateFlags
	cfg   *gatewayv1.Config
}

func (m *Migrator) setupSqlClient() (*sql.DB, string) {
	// Find the database in config and instantiate the service.
	var sqlDB *sql.DB
	var hostInfo string
	for _, s := range m.cfg.Services {
		if s.Name == clutchpg.Name {
			pgdb, err := clutchpg.New(s.TypedConfig, m.log, tally.NoopScope)
			if err != nil {
				log.Fatal("error creating db", zap.Error(err))
			}

			cfg := &postgresv1.Config{}
			if err := ptypes.UnmarshalAny(s.TypedConfig, cfg); err != nil {
				log.Fatal("could not convert config", zap.Error(err))
			}

			sqlDB = pgdb.(clutchpg.Client).DB()
			hostInfo = fmt.Sprintf("%s@%s:%d", cfg.Connection.User, cfg.Connection.Host, cfg.Connection.Port)

			break
		}
	}
	if sqlDB == nil {
		log.Fatal("no database found in config", zap.String("file", m.flags.BaseFlags.ConfigPath))
	}

	return sqlDB, hostInfo
}

// Asks the user to confrim an action, this can be skipped by using the force flag.
// If the input is not 'y' we log fatal and exit.
func (m *Migrator) confirmWithUser(msg string, hostInfo string) {
	// Verify that user wants to continue (unless -f for force is passed as a flag).
	m.log.Info("using database", zap.String("hostInfo", hostInfo))
	if !m.flags.Force {
		m.log.Warn(msg)

		fmt.Printf("\n*** Continue with migration? [y/N] ")
		var answer string
		if _, err := fmt.Scanln(&answer); err != nil && err.Error() != "unexpected newline" {
			m.log.Fatal("could not read user input", zap.Error(err))
		}
		if strings.ToLower(answer) != "y" {
			m.log.Fatal("aborting, enter 'y' to continue or use the '-f' (force) option")
		}
		fmt.Println()
	}
}

func (m *Migrator) Up() {
	sqlDB, hostInfo := m.setupSqlClient()

	msg := "migration has the potential to cause irrevocable data loss, verify host information above"
	m.confirmWithUser(msg, hostInfo)

	// Ping database and bring up driver.
	if err := sqlDB.Ping(); err != nil {
		m.log.Fatal("error pinging db", zap.Error(err))
	}

	driver, err := postgres.WithInstance(sqlDB, &postgres.Config{})
	if err != nil {
		m.log.Fatal("error creating pg driver", zap.Error(err))
	}

	// Create migrator.
	migrationDir, err := os.Getwd()
	if err != nil {
		m.log.Fatal("could not get working dir", zap.Error(err))
	}
	migrationDir = filepath.Join(migrationDir, "migrations")

	sqlMigrate, err := migrate.NewWithDatabaseInstance(
		fmt.Sprintf("file://%s", migrationDir),
		"postgres", driver)
	if err != nil {
		m.log.Fatal("error creating migrator", zap.Error(err))
	}

	sqlMigrate.Log = &migrateLogger{
		logger: m.log.Sugar(),
	}

	// Apply migrations!
	m.log.Info("applying migrations", zap.String("migrationDir", migrationDir))
	err = sqlMigrate.Up()
	if err != nil && err != migrate.ErrNoChange {
		m.log.Fatal("failed running migrations", zap.Error(err))
	}
}

func main() {
	f := &MigrateFlags{}
	f.Link()
	flag.Parse()

	cfg := gateway.MustReadOrValidateConfig(f.BaseFlags)

	logger, _ := zap.NewDevelopment()
	logger = logger.WithOptions(zap.AddStacktrace(zap.FatalLevel + 1))

	migrator := &Migrator{
		log:   logger,
		flags: f,
		cfg:   cfg,
	}

	migrator.Up()
}

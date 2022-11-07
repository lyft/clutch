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
	_ "github.com/lib/pq"
	gatewayv1 "github.com/lyft/clutch/backend/api/config/gateway/v1"
	postgresv1 "github.com/lyft/clutch/backend/api/config/service/db/postgres/v1"
	"github.com/lyft/clutch/backend/gateway"
	clutchpg "github.com/lyft/clutch/backend/service/db/postgres"
	"github.com/uber-go/tally/v4"
	"go.uber.org/zap"
)

type MigrateFlags struct {
	Force        bool
	Down         bool
	Namespace    string
	MigrationDir string
	BaseFlags    *gateway.Flags
}

func (m *MigrateFlags) Link() {
	m.BaseFlags = &gateway.Flags{}
	m.BaseFlags.Link()

	flag.BoolVar(&m.Force, "f", false, "do not ask user for confirmation")
	flag.BoolVar(&m.Down, "down", false, "migrates down by one version")

	flag.StringVar(&m.Namespace, "namespace", "", "when overriding the migrations directory, a namespace must be supplied for independent versioning")
	flag.StringVar(&m.MigrationDir, "migrationDir", "", "override the migrations directory (used when managing private or federated modules)")
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
			if err := s.TypedConfig.UnmarshalTo(cfg); err != nil {
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
func (m *Migrator) confirmWithUser(msg string) {
	_, hostInfo := m.setupSqlClient()
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

// Sets up the sql migrator while also perfomring some pre flight checks such as a db ping.
func (m *Migrator) setupSqlMigrator() *migrate.Migrate {
	sqlDB, _ := m.setupSqlClient()

	// Ping database and bring up driver.
	if err := sqlDB.Ping(); err != nil {
		m.log.Fatal("error pinging db", zap.Error(err))
	}

	// Determine version storage table.
	cfg := &postgres.Config{
		MigrationsTable: postgres.DefaultMigrationsTable,
	}
	if m.flags.Namespace != "" {
		cfg.MigrationsTable = fmt.Sprintf("%s_%s", postgres.DefaultMigrationsTable, m.flags.Namespace)
	}
	m.log.Info("using migration table", zap.String("migrationTable", cfg.MigrationsTable))

	// Create driver.
	driver, err := postgres.WithInstance(sqlDB, cfg)
	if err != nil {
		m.log.Fatal("error creating pg driver", zap.Error(err))
	}

	// Create migrator.
	migrationDir := m.flags.MigrationDir
	if m.flags.MigrationDir == "" {
		var err error
		migrationDir, err = os.Getwd()
		if err != nil {
			m.log.Fatal("could not get working dir", zap.Error(err))
		}
		migrationDir = filepath.Join(migrationDir, "migrations")
	}

	absPath, _ := filepath.Abs(migrationDir)
	m.log.Info("using migration directory", zap.String("migrationDir", migrationDir), zap.String("absoluteMigrationDir", absPath))

	sqlMigrate, err := migrate.NewWithDatabaseInstance(
		fmt.Sprintf("file://%s", migrationDir),
		"postgres", driver)
	if err != nil {
		m.log.Fatal("error creating migrator", zap.Error(err))
	}

	sqlMigrate.Log = &migrateLogger{
		logger: m.log.Sugar(),
	}

	return sqlMigrate
}

func (m *Migrator) Up() {
	sqlMigrate := m.setupSqlMigrator()

	msg := "migration has the potential to cause irrevocable data loss, verify information above"
	m.confirmWithUser(msg)

	// Apply migrations!
	m.log.Info("applying up migrations")
	err := sqlMigrate.Up()
	if err != nil && err != migrate.ErrNoChange {
		m.log.Fatal("failed running migrations", zap.Error(err))
	}
}

func (m *Migrator) Down() {
	sqlMigrate := m.setupSqlMigrator()
	version, _, err := sqlMigrate.Version()
	if err != nil {
		m.log.Fatal("failed to aquire migration version", zap.Error(err))
	}

	msg := fmt.Sprintf(
		"Migrating DOWN by ONE version from (%d -> %d) this migration has the potential to cause irrevocable data loss, verify host information above",
		version, (version - 1))

	m.confirmWithUser(msg)

	// Migrate back by 1
	m.log.Info("applying migrations down")
	err = sqlMigrate.Steps(-1)
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

	// Check that both namespace and migration dir are either empty or filled.
	if (f.Namespace == "") != (f.MigrationDir == "") {
		logger.Fatal("namespace and migration dir must both be provided to use an alternate migration path")
		os.Exit(1)
	}

	migrator := &Migrator{
		log:   logger,
		flags: f,
		cfg:   cfg,
	}

	if f.Down {
		migrator.Down()
	} else {
		migrator.Up()
	}
}

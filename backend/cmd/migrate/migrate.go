package main

import (
	"database/sql"
	"errors"
	"flag"
	"fmt"
	"os"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
)

func run(args []string) error {
	flags := flag.NewFlagSet(args[0], flag.ExitOnError)
	var (
		migrationDir = flags.String("migrations", "./migrations", "Directory with sql migrations.")
		// TODO: Read from the gateway config to create the connection.
		//config       = flags.String("config", "../backend/clutch-config.yaml", "config yaml to parse for db connection")
		connection = flags.String("connection", "", "Connection string for the database.")
	)

	if err := flags.Parse(os.Args[1:]); err != nil {
		flags.PrintDefaults()
		return err
	}

	if *connection == "" {
		flags.PrintDefaults()
		return errors.New("connection string cannot be empty")
	}

	sqlDB, err := sql.Open("postgres", *connection)
	if err != nil {
		return fmt.Errorf("cannot open db connection: %w", err)
	}

	err = sqlDB.Ping()
	if err != nil {
		return fmt.Errorf("error pinging db: %w", err)
	}

	driver, err := postgres.WithInstance(sqlDB, &postgres.Config{})
	if err != nil {
		return fmt.Errorf("error creating pg driver: %w", err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		fmt.Sprintf("file://%s", *migrationDir),
		"postgres", driver)
	if err != nil {
		return fmt.Errorf("err with migration: %w", err)
	}

	err = m.Up()
	if err != nil {
		return fmt.Errorf("failed running migrations: %w", err)
	}

	return nil
}

func main() {
	if err := run(os.Args); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(-1)
	}
}

package database

import (
	"errors"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq" // Postgresql driver
)

func Migrate(db sqlx.DB) error {
	sourceDriver, err := iofs.New(fs, "migrations")
	if err != nil {
		return fmt.Errorf("cannot create source driver to migrate sql schema: %w", err)
	}

	driver, err := postgres.WithInstance(db.DB, &postgres.Config{})
	if err != nil {
		return fmt.Errorf("cannot create postgresql instance: %w", err)
	}

	migration, err := migrate.NewWithInstance("iofs", sourceDriver, "recorder", driver)
	if err != nil {
		return fmt.Errorf("cannot initialize migration: %w", err)
	}

	if err := migration.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return fmt.Errorf("cannot migrate sql schema: %w", err)
	}

	return nil
}

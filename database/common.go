package database

import (
	"context"
	"embed"
	"fmt"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq" // Postgresql driver

	"jazzine/m3u8recorder/libs/logging"
)

//go:embed migrations/*.sql
var fs embed.FS

const connStr = "postgres://postgres:%20%20%20%20@localhost:5433/recorder?sslmode=disable"

func NewConnection() (*sqlx.DB, error) {
	db, err := sqlx.Open("postgres", connStr)

	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}

// SaveOrUpdate saves or updates an element (based on the given query) in postgresql database
func SaveOrUpdate[Entity any](db *sqlx.DB, ctx context.Context, query string, entity Entity) error {
	tx, err := db.Beginx()
	if err != nil {
		return fmt.Errorf("cannot begin transaction: %w", err)
	}
	rollback := func() {
		if err := tx.Rollback(); err != nil {
			logging.Log().
				WithError(err).
				Warn("failed to rollback transaction")
		}
	}
	if _, err := tx.NamedExecContext(ctx, query, entity); err != nil {
		rollback()
		return fmt.Errorf("cannot create query to save/update element: %w", err)
	}
	if err := tx.Commit(); err != nil {
		rollback()
		return fmt.Errorf("cannot commit transaction to save/update element: %w", err)
	}
	return nil
}

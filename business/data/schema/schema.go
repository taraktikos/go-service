package schema

import (
	"context"
	_ "embed"
	"fmt"

	"github.com/ardanlabs/darwin"
	"github.com/jmoiron/sqlx"
)

var (
	//go:embed sql/schema.sql
	schema string

	//go:embed sql/seed.sql
	seed string

	//go:embed sql/delete.sql
	delete string
)

func Migrate(ctx context.Context, db *sqlx.DB) error {
	driver, err := darwin.NewGenericDriver(db.DB, darwin.PostgresDialect{})
	if err != nil {
		return fmt.Errorf("failed to construct darwin driver: %w", err)
	}
	d := darwin.New(driver, darwin.ParseMigrations(schema))
	return d.Migrate()
}

func Seed(ctx context.Context, db *sqlx.DB) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	if _, err := tx.Exec(seed); err != nil {
		if err := tx.Rollback(); err != nil {
			return err
		}
		return err
	}
	return tx.Commit()
}

func DeleteAll(db *sqlx.DB) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	if _, err := tx.Exec(delete); err != nil {
		if err := tx.Rollback(); err != nil {
			return err
		}
		return err
	}
	return tx.Commit()
}

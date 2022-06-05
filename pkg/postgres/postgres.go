// Package postgres provides functions for interacting with postgres-compatible databases and includes
// migration scripts.
package postgres

import (
	"context"
	"database/sql"
	"embed"
	"errors"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	migratepostgres "github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"gocloud.dev/postgres"
	_ "gocloud.dev/postgres/awspostgres"
	_ "gocloud.dev/postgres/gcppostgres"
)

//go:embed migrations/*.sql
var migrations embed.FS

// Open a connection to the database instance described in the URL, performing migrations on success.
func Open(ctx context.Context, url string) (*sql.DB, error) {
	db, err := postgres.Open(ctx, url)
	if err != nil {
		return nil, err
	}

	if err = db.PingContext(ctx); err != nil {
		return nil, err
	}

	if err = performMigrations(db); err != nil {
		return nil, fmt.Errorf("failed to migrate database: %w", err)
	}

	return db, nil
}

func performMigrations(db *sql.DB) error {
	source, err := iofs.New(migrations, "migrations")
	if err != nil {
		return err
	}

	driver, err := migratepostgres.WithInstance(db, &migratepostgres.Config{})
	if err != nil {
		return err
	}

	migration, err := migrate.NewWithInstance("iofs", source, "postgres", driver)
	if err != nil {
		return err
	}

	err = migration.Up()
	switch {
	case errors.Is(err, migrate.ErrNoChange):
		break
	case err != nil:
		return err
	}

	return nil
}

// The TransactionFunc is a function invoked when calling WithinTransaction.
type TransactionFunc func(ctx context.Context, tx *sql.Tx) error

// WithinTransaction invokes the TransactionFunc, providing it an sql.Tx transaction to perform SQL operations
// against. If the TransactionFunc returns a non-nil error, the transaction is rolled back. Otherwise, it is
// committed. If the transaction is intended to be read-only, use WithinReadOnlyTransaction.
func WithinTransaction(ctx context.Context, db *sql.DB, fn TransactionFunc) error {
	return withinTransaction(ctx, db, &sql.TxOptions{}, fn)
}

// WithinReadOnlyTransaction invokes the TransactionFunc, providing it an sql.Tx transaction to perform SQL operations
// against. If the TransactionFunc returns a non-nil error, the transaction is rolled back. Otherwise, it is
// committed. If the transaction is intended to include write operations, use WithinTransaction.
func WithinReadOnlyTransaction(ctx context.Context, db *sql.DB, fn TransactionFunc) error {
	return withinTransaction(ctx, db, &sql.TxOptions{ReadOnly: true}, fn)
}

func withinTransaction(ctx context.Context, db *sql.DB, opts *sql.TxOptions, fn TransactionFunc) error {
	tx, err := db.BeginTx(ctx, opts)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	if err = fn(ctx, tx); err != nil {
		if err = tx.Rollback(); err != nil {
			return fmt.Errorf("failed to roll back transaction: %w", err)
		}

		return err
	}

	return tx.Commit()
}

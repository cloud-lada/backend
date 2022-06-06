package status

import (
	"context"
	"database/sql"
	"errors"

	"github.com/cloud-lada/backend/pkg/postgres"
)

type (
	// The PostgresRepository is a Repository implementation that queries status data from a postgres-compatible
	// database.
	PostgresRepository struct {
		db *sql.DB
	}
)

// NewPostgresRepository returns a new instance of the PostgresRepository type that will perform queries against
// the provided sql.DB instance.
func NewPostgresRepository(db *sql.DB) *PostgresRepository {
	return &PostgresRepository{db: db}
}

// Status returns the current status of data within the database.
func (r *PostgresRepository) Status(ctx context.Context) (Status, error) {
	var status Status

	err := postgres.WithinReadOnlyTransaction(ctx, r.db, func(ctx context.Context, tx *sql.Tx) error {
		const q = `SELECT MAX(timestamp) FROM reading`

		row := tx.QueryRowContext(ctx, q)
		err := row.Scan(&status.LastIngestTimestamp)
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil
		case err != nil:
			return err
		default:
			return row.Err()
		}
	})

	return status, err
}

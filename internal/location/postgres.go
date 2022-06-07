package location

import (
	"context"
	"database/sql"
	"errors"

	"github.com/cloud-lada/backend/pkg/postgres"
)

type (
	// The PostgresRepository is a Repository implementation that queries location data from a postgres-compatible
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

// Latest returns the most recent location data stored within the database.
func (r *PostgresRepository) Latest(ctx context.Context) (Location, error) {
	var location Location
	var err error

	err = postgres.WithinReadOnlyTransaction(ctx, r.db, func(ctx context.Context, tx *sql.Tx) error {
		location.Latitude, err = r.latestLatitude(ctx, tx)
		if err != nil {
			return err
		}

		location.Longitude, err = r.latestLongitude(ctx, tx)
		return err
	})

	return location, err
}

func (r *PostgresRepository) latestLatitude(ctx context.Context, tx *sql.Tx) (float64, error) {
	const q = `
		SELECT value FROM reading
		WHERE sensor = 'location_latitude'
		ORDER BY timestamp DESC
		FETCH FIRST ROW ONLY
	`

	var value float64
	row := tx.QueryRowContext(ctx, q)

	err := row.Scan(&value)
	switch {
	case errors.Is(err, sql.ErrNoRows):
		return 0, nil
	case err != nil:
		return 0, err
	default:
		return value, nil
	}
}

func (r *PostgresRepository) latestLongitude(ctx context.Context, tx *sql.Tx) (float64, error) {
	const q = `
		SELECT value FROM reading
		WHERE sensor = 'location_longitude'
		ORDER BY timestamp DESC
		FETCH FIRST ROW ONLY
	`

	var value float64
	row := tx.QueryRowContext(ctx, q)

	err := row.Scan(&value)
	switch {
	case errors.Is(err, sql.ErrNoRows):
		return 0, nil
	case err != nil:
		return 0, err
	default:
		return value, nil
	}
}

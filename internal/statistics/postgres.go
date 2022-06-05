package statistics

import (
	"context"
	"database/sql"
	"errors"

	"github.com/cloud-lada/backend/pkg/postgres"
)

type (
	// The PostgresRepository is a Repository implementation that queries statistical data from a postgres-compatible
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

// Latest returns a Statistics type whose fields will be populated with the most recent data available within
// the database.
func (r *PostgresRepository) Latest(ctx context.Context) (Statistics, error) {
	var stats Statistics
	var err error

	err = postgres.WithinReadOnlyTransaction(ctx, r.db, func(ctx context.Context, tx *sql.Tx) error {
		stats.Speed, err = r.latestReading(ctx, tx, "speed")
		if err != nil {
			return err
		}

		stats.Fuel, err = r.latestReading(ctx, tx, "fuel")
		if err != nil {
			return err
		}

		stats.EngineTemperature, err = r.latestReading(ctx, tx, "engine_temperature")
		if err != nil {
			return err
		}

		stats.Revolutions, err = r.latestReading(ctx, tx, "revolution")
		return err
	})

	return stats, err
}

func (r *PostgresRepository) latestReading(ctx context.Context, tx *sql.Tx, sensor string) (float64, error) {
	const q = `
		SELECT value FROM reading
		WHERE sensor = $1
		ORDER BY timestamp DESC 
		FETCH FIRST ROW ONLY
	`

	var value float64
	row := tx.QueryRowContext(ctx, q, sensor)

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

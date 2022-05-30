// Package reading provides persistence logic for sensor readings.
package reading

import (
	"context"
	"database/sql"
	"time"

	"github.com/cloud-lada/backend/pkg/postgres"
)

type (
	// The PostgresRepository type is used to persist Reading data into a PostgreSQL instance.
	PostgresRepository struct {
		db *sql.DB
	}

	ForEachFunc func(ctx context.Context, reading Reading) error
)

// NewPostgresRepository returns a new instance of the PostgresRepository type that will perform queries against
// the provided sql.DB instance.
func NewPostgresRepository(db *sql.DB) *PostgresRepository {
	return &PostgresRepository{db: db}
}

// Save the reading to the database. If a reading already exists for the sensor at the given timestamp, do nothing.
func (pr *PostgresRepository) Save(ctx context.Context, reading Reading) error {
	return postgres.WithinTransaction(ctx, pr.db, func(ctx context.Context, tx *sql.Tx) error {
		const q = `
			INSERT INTO reading (sensor, value, timestamp) VALUES ($1, $2, $3)
			ON CONFLICT (sensor, timestamp) DO NOTHING
		`

		_, err := tx.ExecContext(ctx, q, reading.Sensor, reading.Value, reading.Timestamp)
		return err
	})
}

// ForEachOnDate iterates through all readings stored in the database on the date component of the given time. For
// each record, the ForEachFunc is invoked. Iteration will stop when there are no more records, the context is
// cancelled or the ForEachFunc returns an error. Readings are processed in batches of 100 at the time.
func (pr *PostgresRepository) ForEachOnDate(ctx context.Context, date time.Time, fn ForEachFunc) error {
	return postgres.WithinTransaction(ctx, pr.db, func(ctx context.Context, tx *sql.Tx) error {
		const cursorQuery = `
			DECLARE reading_cursor CURSOR FOR 
			    SELECT sensor, value, timestamp FROM reading
				WHERE timestamp >= $1::DATE AND timestamp < ($1::DATE + INTERVAL '1 day')
		`

		if _, err := tx.ExecContext(ctx, cursorQuery, date); err != nil {
			return err
		}

		const fetchQuery = "FETCH 100 FROM reading_cursor"

		for {
			rows, err := tx.QueryContext(ctx, fetchQuery)
			if err != nil {
				return err
			}

			readings := make([]Reading, 0)
			for rows.Next() {
				if err = rows.Err(); err != nil {
					return err
				}

				var reading Reading
				if err = rows.Scan(&reading.Sensor, &reading.Value, &reading.Timestamp); err != nil {
					return err
				}

				readings = append(readings, reading)
			}

			if len(readings) == 0 {
				return rows.Err()
			}

			for _, reading := range readings {
				if err = fn(ctx, reading); err != nil {
					return err
				}
			}
		}
	})
}

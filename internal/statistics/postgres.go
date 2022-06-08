package statistics

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/cloud-lada/backend/internal/reading"
	"github.com/cloud-lada/backend/pkg/closers"
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
		stats.Speed, err = r.latestReading(ctx, tx, reading.SensorTypeSpeed)
		if err != nil {
			return err
		}

		stats.Fuel, err = r.latestReading(ctx, tx, reading.SensorTypeFuel)
		if err != nil {
			return err
		}

		stats.EngineTemperature, err = r.latestReading(ctx, tx, reading.SensorTypeEngineTemperature)
		if err != nil {
			return err
		}

		stats.Revolutions, err = r.latestReading(ctx, tx, reading.SensorTypeRevolution)
		return err
	})

	return stats, err
}

func (r *PostgresRepository) latestReading(ctx context.Context, tx *sql.Tx, sensor reading.SensorType) (float64, error) {
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

// ForDate queries the database for time-bucketed statistics on a given date for a sensor type. Statistics are averaged
// in 15 minute intervals.
func (r *PostgresRepository) ForDate(ctx context.Context, date time.Time, sensor reading.SensorType) ([]Statistic, error) {
	out := make([]Statistic, 0)
	err := postgres.WithinReadOnlyTransaction(ctx, r.db, func(ctx context.Context, tx *sql.Tx) error {
		// This query returns the average value of the sensor over 15 minute increments. It uses the
		// time_bucket_gapfill function to automatically fill in times for the rest of the day when
		// the dataset is incomplete. It uses locf to set the values for the missing timestamps to the
		// previous known value. The time range is inferred by the WHERE clause of the query.
		//
		// https://docs.timescale.com/api/latest/hyperfunctions/gapfilling-interpolation/time_bucket_gapfill
		// https://docs.timescale.com/api/latest/hyperfunctions/gapfilling-interpolation/locf
		const q = `
			SELECT 
				sensor,
				COALESCE(AVG(value), 0),
				locf(AVG(value)),
				time_bucket_gapfill('15 minutes', timestamp) AS bucket
			FROM reading 
			WHERE 
				sensor = $1
				AND timestamp >= $2::DATE 
				AND timestamp < ($2::DATE + INTERVAL '1 day')
			GROUP BY bucket, sensor
			ORDER BY bucket ASC
		`

		rows, err := tx.QueryContext(ctx, q, sensor, date)
		if err != nil {
			return err
		}
		defer closers.Close(rows)

		var locf sql.NullFloat64
		for rows.Next() {
			var stat Statistic

			// We don't actually care about the locf value, but it allows timescale to fill in all the missing
			// intervals with a NULL value that we coalesce to zero. So we scan it and throw it away.
			if err = rows.Scan(&stat.Sensor, &stat.Value, &locf, &stat.Timestamp); err != nil {
				return err
			}

			out = append(out, stat)
		}

		if err = rows.Err(); err != nil {
			return err
		}

		return rows.Close()
	})

	return out, err
}

package statistics_test

import (
	"testing"
	"time"

	"github.com/cloud-lada/backend/internal/reading"
	"github.com/cloud-lada/backend/internal/statistics"
	"github.com/cloud-lada/backend/pkg/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPostgresRepository_Latest(t *testing.T) {
	if testing.Short() {
		t.Skip()
		return
	}

	ctx := testutil.Context(t)
	db := testutil.Postgres(t, ctx)

	readings := reading.NewPostgresRepository(db)
	stats := statistics.NewPostgresRepository(db)

	// Insert readings that we can query
	seed := []reading.Reading{
		{
			Sensor:    reading.SensorTypeSpeed,
			Value:     10,
			Timestamp: time.Now(),
		},
		{
			Sensor:    reading.SensorTypeFuel,
			Value:     50,
			Timestamp: time.Now(),
		},
		{
			Sensor:    reading.SensorTypeEngineTemperature,
			Value:     30,
			Timestamp: time.Now(),
		},
		{
			Sensor:    reading.SensorTypeRevolution,
			Value:     1000,
			Timestamp: time.Now(),
		},
	}

	for _, s := range seed {
		require.NoError(t, readings.Save(ctx, s))
	}

	t.Run("It should return the latest statistics", func(t *testing.T) {
		expected := statistics.Statistics{
			Speed:             10,
			Fuel:              50,
			EngineTemperature: 30,
			Revolutions:       1000,
		}

		actual, err := stats.Latest(ctx)
		require.NoError(t, err)
		assert.EqualValues(t, expected, actual)
	})
}

func TestPostgresRepository_ForDate(t *testing.T) {
	if testing.Short() {
		t.Skip()
		return
	}

	ctx := testutil.Context(t)
	db := testutil.Postgres(t, ctx)

	readings := reading.NewPostgresRepository(db)
	stats := statistics.NewPostgresRepository(db)

	var seed []reading.Reading

	// Create readings for 1AM, this allows us to test the gapfilling for the hours before and
	// after.
	for i := 60; i < 120; i++ {
		ts := time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC)
		seed = append(seed, reading.Reading{
			Sensor:    reading.SensorTypeSpeed,
			Value:     10,
			Timestamp: ts.Add(time.Minute * time.Duration(i)),
		})
	}

	for _, s := range seed {
		require.NoError(t, readings.Save(ctx, s))
	}

	t.Run("It should return time bucketed statistics", func(t *testing.T) {
		date := time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC)

		actual, err := stats.ForDate(ctx, date, reading.SensorTypeSpeed)
		require.NoError(t, err)

		// We expect 96 readings for 15 minute increments over 24 hours.
		require.Len(t, actual, 96)

		for _, elem := range actual {
			assert.EqualValues(t, reading.SensorTypeSpeed, elem.Sensor)

			// We only set values for the 1am hour, so anything before and after that should have a value of zero.
			if elem.Timestamp.Hour() == 1 {
				assert.EqualValues(t, float64(10), elem.Value)
			} else {
				assert.EqualValues(t, float64(0), elem.Value)
			}

			assert.NotZero(t, elem.Timestamp)
		}
	})
}

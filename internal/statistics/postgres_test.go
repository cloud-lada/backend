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
			Sensor:    "speed",
			Value:     10,
			Timestamp: time.Now(),
		},
		{
			Sensor:    "fuel",
			Value:     50,
			Timestamp: time.Now(),
		},
		{
			Sensor:    "engine_temperature",
			Value:     30,
			Timestamp: time.Now(),
		},
		{
			Sensor:    "revolution",
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

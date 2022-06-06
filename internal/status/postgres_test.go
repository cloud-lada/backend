package status_test

import (
	"testing"
	"time"

	"github.com/cloud-lada/backend/internal/reading"
	"github.com/cloud-lada/backend/internal/status"
	"github.com/cloud-lada/backend/pkg/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPostgresRepository_Status(t *testing.T) {
	if testing.Short() {
		t.Skip()
		return
	}

	ctx := testutil.Context(t)
	db := testutil.Postgres(t, ctx)

	readings := reading.NewPostgresRepository(db)
	statuses := status.NewPostgresRepository(db)

	// Insert readings that we can query
	seed := []reading.Reading{
		{
			Sensor:    reading.SensorTypeSpeed,
			Value:     10,
			Timestamp: time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
		},
		{
			Sensor:    reading.SensorTypeSpeed,
			Value:     10,
			Timestamp: time.Date(2021, 1, 1, 1, 0, 0, 0, time.UTC),
		},
		{
			Sensor:    reading.SensorTypeSpeed,
			Value:     10,
			Timestamp: time.Date(2021, 1, 1, 2, 0, 0, 0, time.UTC),
		},
	}

	for _, s := range seed {
		require.NoError(t, readings.Save(ctx, s))
	}

	t.Run("It should return the latest status", func(t *testing.T) {
		expected := status.Status{
			LastIngestTimestamp: time.Date(2021, 1, 1, 2, 0, 0, 0, time.UTC),
		}

		actual, err := statuses.Status(ctx)
		require.NoError(t, err)
		assert.EqualValues(t, expected, actual)
	})
}

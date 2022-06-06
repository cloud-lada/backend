package location_test

import (
	"testing"
	"time"

	"github.com/cloud-lada/backend/internal/location"
	"github.com/cloud-lada/backend/internal/reading"
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
	locations := location.NewPostgresRepository(db)

	// Insert readings that we can query
	seed := []reading.Reading{
		{
			Sensor:    reading.SensorTypeLocationLatitude,
			Value:     50,
			Timestamp: time.Now(),
		},
		{
			Sensor:    reading.SensorTypeLocationLongitude,
			Value:     51,
			Timestamp: time.Now(),
		},
	}

	for _, s := range seed {
		require.NoError(t, readings.Save(ctx, s))
	}

	t.Run("It should return the latest location", func(t *testing.T) {
		expected := location.Location{
			Longitude: 51,
			Latitude:  50,
		}

		actual, err := locations.Latest(ctx)
		require.NoError(t, err)
		assert.EqualValues(t, expected, actual)
	})
}

package reading_test

import (
	"context"
	"testing"
	"time"

	"github.com/cloud-lada/backend/internal/reading"
	"github.com/cloud-lada/backend/pkg/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPostgresRepository_Save(t *testing.T) {
	if testing.Short() {
		t.Skip()
		return
	}

	ctx := testutil.Context(t)
	db := testutil.Postgres(t, ctx)
	repo := reading.NewPostgresRepository(db)

	t.Run("It should store a reading", func(t *testing.T) {
		assert.NoError(t, repo.Save(ctx, reading.Reading{
			Sensor:    reading.SensorTypeSpeed,
			Value:     100,
			Timestamp: time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC),
		}))
	})

	t.Run("It should not error for a duplicate reading", func(t *testing.T) {
		assert.NoError(t, repo.Save(ctx, reading.Reading{
			Sensor:    reading.SensorTypeSpeed,
			Value:     100,
			Timestamp: time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC),
		}))
	})
}

func TestPostgresRepository_ForEachOnDate(t *testing.T) {
	if testing.Short() {
		t.Skip()
		return
	}

	ctx := testutil.Context(t)
	db := testutil.Postgres(t, ctx)
	repo := reading.NewPostgresRepository(db)

	now := time.Now().UTC()
	readings := []reading.Reading{
		{
			Sensor:    reading.SensorTypeSpeed,
			Value:     50,
			Timestamp: now.Add(-time.Hour * 24),
		},
		{
			Sensor:    reading.SensorTypeSpeed,
			Value:     100,
			Timestamp: now,
		},
	}

	for _, r := range readings {
		require.NoError(t, repo.Save(ctx, r))
	}

	assert.NoError(t, repo.ForEachOnDate(ctx, now, func(ctx context.Context, reading reading.Reading) error {
		assert.EqualValues(t, readings[1].Sensor, reading.Sensor)
		assert.EqualValues(t, readings[1].Value, reading.Value)
		assert.EqualValues(t, readings[1].Timestamp.Unix(), reading.Timestamp.Unix())
		return nil
	}))
}

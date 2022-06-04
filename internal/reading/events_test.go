package reading_test

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"testing"
	"time"

	"github.com/cloud-lada/backend/internal/reading"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEventHandler_HandleEvent(t *testing.T) {
	t.Parallel()

	tt := []struct {
		Name         string
		Data         json.RawMessage
		Error        error
		ExpectsError bool
		Expected     reading.Reading
	}{
		{
			Name: "It should store a reading",
			Data: marshal(t, reading.Reading{
				Sensor:    "speed",
				Value:     100,
				Timestamp: time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC),
			}),
			Expected: reading.Reading{
				Sensor:    "speed",
				Value:     100,
				Timestamp: time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC),
			},
		},
		{
			Name:         "It should return an error for an invalid reading",
			Data:         []byte("this will not unmarshal"),
			ExpectsError: true,
		},
		{
			Name: "It should return repository errors",
			Data: marshal(t, reading.Reading{
				Sensor:    "speed",
				Value:     100,
				Timestamp: time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC),
			}),
			Error:        io.EOF,
			ExpectsError: true,
		},
	}

	for _, tc := range tt {
		t.Run(tc.Name, func(t *testing.T) {
			ctx := context.Background()
			readings := &MockRepository{err: tc.Error}

			handler := reading.NewEventHandler(readings, log.New(io.Discard, "", log.Flags()))

			err := handler.HandleEvent(ctx, tc.Data)
			if tc.ExpectsError {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			assert.EqualValues(t, tc.Expected, readings.saved)
		})
	}
}

func marshal(t *testing.T, value interface{}) []byte {
	t.Helper()

	data, err := json.Marshal(value)
	require.NoError(t, err)
	return data
}

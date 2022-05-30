package dump_test

import (
	"bytes"
	"compress/gzip"
	"context"
	"encoding/json"
	"io"
	"log"
	"testing"
	"time"

	"github.com/cloud-lada/backend/internal/dump"
	"github.com/cloud-lada/backend/internal/reading"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDumper_Dump(t *testing.T) {
	t.Parallel()

	tt := []struct {
		Name         string
		Date         time.Time
		Seed         []reading.Reading
		Error        error
		ExpectsError bool
	}{
		{
			Name: "It should write all readings from the repository to the sink",
			Date: time.Now(),
			Seed: []reading.Reading{
				{
					Sensor:    "speed",
					Value:     100,
					Timestamp: time.Now().UTC(),
				},
				{
					Sensor:    "speed",
					Value:     200,
					Timestamp: time.Now().UTC(),
				},
				{
					Sensor:    "speed",
					Value:     300,
					Timestamp: time.Now().UTC(),
				},
			},
		},
		{
			Name:         "It should propagate errors from the repository",
			Date:         time.Now(),
			Error:        io.EOF,
			ExpectsError: true,
		},
	}

	for _, tc := range tt {
		t.Run(tc.Name, func(t *testing.T) {
			ctx := context.Background()
			readings := &MockRepository{readings: tc.Seed, err: tc.Error}
			blobs := &MockSink{buffer: bytes.NewBuffer([]byte{})}

			config := dump.Config{
				Date:     tc.Date,
				Readings: readings,
				Blobs:    blobs,
				Logger:   log.New(io.Discard, "", log.Flags()),
			}

			err := dump.New(config).Dump(ctx)
			if tc.ExpectsError {
				assert.Error(t, err)
				return
			}

			assert.EqualValues(t, tc.Date.Format("2006-02-01.json.gz"), blobs.name)

			reader, err := gzip.NewReader(blobs.buffer)
			require.NoError(t, err)

			decoder := json.NewDecoder(reader)
			for _, expected := range tc.Seed {
				var actual reading.Reading
				require.NoError(t, decoder.Decode(&actual))
				assert.EqualValues(t, expected, actual)
			}

			require.NoError(t, reader.Close())
		})
	}
}

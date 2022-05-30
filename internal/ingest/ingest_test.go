package ingest_test

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/cloud-lada/backend/internal/ingest"
	"github.com/cloud-lada/backend/internal/reading"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestIngestor_Ingest(t *testing.T) {
	t.Parallel()

	tt := []struct {
		Name           string
		Readings       []reading.Reading
		PublishError   error
		ExpectedCode   int
		ExpectedAPIKey string
	}{
		{
			Name:           "It should accept valid readings and publish them",
			ExpectedCode:   http.StatusOK,
			ExpectedAPIKey: "ingestor-key",
			Readings: []reading.Reading{
				{
					Sensor:    "speed",
					Value:     65,
					Timestamp: time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC),
				},
				{
					Sensor:    "speed",
					Value:     70,
					Timestamp: time.Date(2022, 1, 1, 0, 1, 0, 0, time.UTC),
				},
			},
		},
		{
			Name:         "It should return unauthorized for an invalid API key",
			ExpectedCode: http.StatusUnauthorized,
			Readings: []reading.Reading{
				{
					Sensor:    "speed",
					Value:     65,
					Timestamp: time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC),
				},
			},
		},
		{
			Name:           "It should return internal server error for publishing errors",
			ExpectedCode:   http.StatusInternalServerError,
			ExpectedAPIKey: "ingestor-key",
			PublishError:   io.EOF,
			Readings: []reading.Reading{
				{
					Sensor:    "speed",
					Value:     65,
					Timestamp: time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC),
				},
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.Name, func(t *testing.T) {
			sink := &MockEventWriter{err: tc.PublishError}
			config := ingest.Config{
				APIKey: tc.ExpectedAPIKey,
				Writer: sink,
				Logger: log.New(io.Discard, "", log.Flags()),
			}

			h := ingest.New(config)

			router := mux.NewRouter()
			h.Register(router)

			body := bytes.NewBuffer([]byte{})
			encoder := json.NewEncoder(body)
			for _, r := range tc.Readings {
				require.NoError(t, encoder.Encode(r))
			}

			w := httptest.NewRecorder()
			r := httptest.NewRequest(http.MethodPost, "http://example.org/ingest", body)
			r.SetBasicAuth(tc.ExpectedAPIKey, "")
			r.Header.Set("Content-Type", "application/stream+json")

			router.ServeHTTP(w, r)
			assert.EqualValues(t, tc.ExpectedCode, w.Code)

			if tc.ExpectedCode >= http.StatusMultipleChoices {
				return
			}

			require.Len(t, sink.messages, len(tc.Readings))

			for i, message := range sink.messages {
				var r reading.Reading

				require.NoError(t, json.Unmarshal(message.Data, &r))
				assert.EqualValues(t, tc.Readings[i], r)
			}
		})
	}
}

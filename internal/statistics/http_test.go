package statistics_test

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"path"
	"testing"
	"time"

	"github.com/cloud-lada/backend/internal/reading"
	"github.com/cloud-lada/backend/internal/statistics"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHTTP_Latest(t *testing.T) {
	t.Parallel()

	tt := []struct {
		Name         string
		Expected     statistics.Statistics
		Error        error
		ExpectsError bool
		ExpectedCode int
	}{
		{
			Name: "It should return the latest statistics",
			Expected: statistics.Statistics{
				Speed:             10,
				Fuel:              11,
				EngineTemperature: 12,
				Revolutions:       13,
			},
			ExpectedCode: http.StatusOK,
		},
		{
			Name:         "It should return return errors from the repository",
			Error:        io.EOF,
			ExpectsError: true,
			ExpectedCode: http.StatusInternalServerError,
		},
	}

	for _, tc := range tt {
		t.Run(tc.Name, func(t *testing.T) {
			repo := &MockRepository{latest: tc.Expected, err: tc.Error}
			api := statistics.NewHTTP(repo)

			router := mux.NewRouter()
			api.Register(router)

			w := httptest.NewRecorder()
			r := httptest.NewRequest(http.MethodGet, "/statistics/latest", nil)

			router.ServeHTTP(w, r)
			assert.EqualValues(t, tc.ExpectedCode, w.Code)
			if tc.ExpectsError {
				return
			}

			var actual statistics.Statistics
			require.NoError(t, json.NewDecoder(w.Body).Decode(&actual))
			assert.EqualValues(t, tc.Expected, actual)
		})
	}
}

func TestHTTP_ForDate(t *testing.T) {
	t.Parallel()

	tt := []struct {
		Name         string
		Expected     []statistics.Statistic
		Date         time.Time
		Sensor       reading.SensorType
		Error        error
		ExpectsError bool
		ExpectedCode int
	}{
		{
			Name:         "It should return statistics for the date",
			Date:         time.Now(),
			Sensor:       reading.SensorTypeSpeed,
			ExpectedCode: http.StatusOK,
			Expected: []statistics.Statistic{
				{
					Sensor:    reading.SensorTypeSpeed,
					Value:     10,
					Timestamp: time.Now(),
				},
			},
		},
		{
			Name:         "It should return return errors from the repository",
			Error:        io.EOF,
			ExpectsError: true,
			ExpectedCode: http.StatusInternalServerError,
			Date:         time.Now(),
			Sensor:       reading.SensorTypeSpeed,
		},
		{
			Name:         "It should return bad request for an invalid sensor",
			ExpectsError: true,
			ExpectedCode: http.StatusBadRequest,
			Date:         time.Now(),
			Sensor:       "invalid",
		},
	}

	for _, tc := range tt {
		t.Run(tc.Name, func(t *testing.T) {
			repo := &MockRepository{stats: tc.Expected, err: tc.Error}
			api := statistics.NewHTTP(repo)

			router := mux.NewRouter()
			api.Register(router)

			w := httptest.NewRecorder()

			uri := path.Join("/statistics", "sensor", string(tc.Sensor), "date", tc.Date.Format("2006-02-01"))
			r := httptest.NewRequest(http.MethodGet, uri, nil)

			router.ServeHTTP(w, r)
			assert.EqualValues(t, tc.ExpectedCode, w.Code)
			if tc.ExpectsError {
				return
			}

			var actuals []statistics.Statistic
			require.NoError(t, json.NewDecoder(w.Body).Decode(&actuals))
			require.Len(t, actuals, len(tc.Expected))
		})
	}
}

package statistics_test

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

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
			repo := &MockRepository{stats: tc.Expected, err: tc.Error}
			api := statistics.NewHTTP(repo)

			router := mux.NewRouter()
			api.Register(router)

			w := httptest.NewRecorder()
			r := httptest.NewRequest(http.MethodGet, "/api/statistics/latest", nil)

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

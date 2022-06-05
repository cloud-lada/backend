package location_test

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/cloud-lada/backend/internal/location"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHTTP_Latest(t *testing.T) {
	t.Parallel()

	tt := []struct {
		Name         string
		Expected     location.Location
		Error        error
		ExpectsError bool
		ExpectedCode int
	}{
		{
			Name: "It should return the latest location",
			Expected: location.Location{
				Longitude: 50,
				Latitude:  51,
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
			repo := &MockRepository{location: tc.Expected, err: tc.Error}
			api := location.NewHTTP(repo)

			router := mux.NewRouter()
			api.Register(router)

			w := httptest.NewRecorder()
			r := httptest.NewRequest(http.MethodGet, "/location/latest", nil)

			router.ServeHTTP(w, r)
			assert.EqualValues(t, tc.ExpectedCode, w.Code)
			if tc.ExpectsError {
				return
			}

			var actual location.Location
			require.NoError(t, json.NewDecoder(w.Body).Decode(&actual))
			assert.EqualValues(t, tc.Expected, actual)
		})
	}
}

package status_test

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/cloud-lada/backend/internal/status"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHTTP_Status(t *testing.T) {
	t.Parallel()

	tt := []struct {
		Name         string
		Expected     status.Status
		ExpectedCode int
		ExpectsError bool
		Error        error
	}{
		{
			Name: "It should return the current status",
			Expected: status.Status{
				LastIngestTimestamp: time.Now(),
			},
			ExpectedCode: http.StatusOK,
		},
		{
			Name:         "It should propagate repository errors",
			Error:        io.EOF,
			ExpectsError: true,
			ExpectedCode: http.StatusInternalServerError,
		},
	}

	for _, tc := range tt {
		t.Run(tc.Name, func(t *testing.T) {
			repo := &MockRepository{status: tc.Expected, err: tc.Error}

			router := mux.NewRouter()
			status.NewHTTP(repo).Register(router)

			w := httptest.NewRecorder()
			r := httptest.NewRequest(http.MethodGet, "/status", nil)

			router.ServeHTTP(w, r)
			assert.EqualValues(t, tc.ExpectedCode, w.Code)
			if tc.ExpectsError {
				return
			}

			var actual status.Status
			require.NoError(t, json.NewDecoder(w.Body).Decode(&actual))
			assert.EqualValues(t, tc.Expected.LastIngestTimestamp.Unix(), actual.LastIngestTimestamp.Unix())
		})
	}
}

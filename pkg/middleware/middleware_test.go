package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/cloud-lada/backend/pkg/middleware"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

func TestAPIKey(t *testing.T) {
	t.Parallel()

	tt := []struct {
		Name              string
		APIKey            string
		BasicAuthUsername string
		ExpectedCode      int
	}{
		{
			Name:              "It should allow requests that have the correct API key",
			APIKey:            "example",
			BasicAuthUsername: "example",
			ExpectedCode:      http.StatusOK,
		},
		{
			Name:              "It should reject requests that have an incorrect API key",
			APIKey:            "example",
			BasicAuthUsername: "wrong",
			ExpectedCode:      http.StatusUnauthorized,
		},
	}

	for _, tc := range tt {
		t.Run(tc.Name, func(t *testing.T) {
			router := mux.NewRouter()
			router.Use(middleware.APIKey(tc.APIKey))
			router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {})

			w := httptest.NewRecorder()
			r := httptest.NewRequest(http.MethodGet, "/", nil)
			r.SetBasicAuth(tc.BasicAuthUsername, "")

			router.ServeHTTP(w, r)
			assert.EqualValues(t, tc.ExpectedCode, w.Code)
		})
	}

}

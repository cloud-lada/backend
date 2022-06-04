// Package middleware contains HTTP middleware functions used by HTTP servers.
package middleware

import (
	"net/http"

	"github.com/gorilla/mux"
)

// APIKey returns a mux.MiddlewareFunc implementation that will check the basic authentication credentials for a
// username matching the provided key. It will return a 401 response if there are no basic authentication credentials
// or if the username does not match the api key.
func APIKey(key string) mux.MiddlewareFunc {
	return func(handler http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			username, _, ok := r.BasicAuth()
			switch {
			case !ok:
				http.Error(w, "no credentials provided", http.StatusUnauthorized)
				return
			case username != key || username == "":
				http.Error(w, "invalid api key", http.StatusUnauthorized)
				return
			default:
				handler.ServeHTTP(w, r)
			}
		})
	}
}

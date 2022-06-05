package location

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

type (
	// The HTTP type contains HTTP request handlers that serve location data.
	HTTP struct {
		location Repository
	}

	// The Repository interface describes types that can query location database from persistent
	// storage.
	Repository interface {
		Latest(ctx context.Context) (Location, error)
	}
)

// NewHTTP returns a new instance of the HTTP type that will serve location data queried from the
// Repository implementation.
func NewHTTP(location Repository) *HTTP {
	return &HTTP{location: location}
}

// Latest handles an inbound HTTP GET request that returns the latest location stored within the Repository.
func (h *HTTP) Latest(w http.ResponseWriter, r *http.Request) {
	stats, err := h.location.Latest(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err = json.NewEncoder(w).Encode(stats); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// Register the HTTP routes into the given router.
func (h *HTTP) Register(router *mux.Router) {
	router.HandleFunc("/location/latest", h.Latest).Methods(http.MethodGet)
}

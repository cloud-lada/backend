package statistics

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

type (
	// The HTTP type contains HTTP request handlers that serve statistical data.
	HTTP struct {
		statistics Repository
	}

	// The Repository interface describes types that can query statistical database from persistent
	// storage.
	Repository interface {
		Latest(ctx context.Context) (Statistics, error)
	}
)

// NewHTTP returns a new instance of the HTTP type that will serve statistical data queried from the
// Repository implementation.
func NewHTTP(statistics Repository) *HTTP {
	return &HTTP{statistics: statistics}
}

// Latest handles an inbound HTTP GET request that returns the latest statistics stored within the Repository.
func (h *HTTP) Latest(w http.ResponseWriter, r *http.Request) {
	stats, err := h.statistics.Latest(r.Context())
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
	router.HandleFunc("/statistics/latest", h.Latest).Methods(http.MethodGet)
}

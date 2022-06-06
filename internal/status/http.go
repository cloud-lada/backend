package status

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

type (
	// The HTTP type contains HTTP request handlers that return the status of data ingestion.
	HTTP struct {
		statuses Repository
	}

	// The Repository interface describes types that can query the status of the data within the persistent
	// storage.
	Repository interface {
		Status(ctx context.Context) (Status, error)
	}
)

// NewHTTP returns a new instance of the HTTP type that will serve status data queried from the
// Repository implementation.
func NewHTTP(statuses Repository) *HTTP {
	return &HTTP{statuses: statuses}
}

// Status handles an inbound HTTP GET request that returns information on the status of data ingestion from the
// Repository.
func (h *HTTP) Status(w http.ResponseWriter, r *http.Request) {
	status, err := h.statuses.Status(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err = json.NewEncoder(w).Encode(status); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// Register the HTTP routes into the given router.
func (h *HTTP) Register(router *mux.Router) {
	router.HandleFunc("/status", h.Status).Methods(http.MethodGet)
}

package statistics

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/cloud-lada/backend/internal/reading"
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
		ForDate(ctx context.Context, date time.Time, sensor reading.SensorType) ([]Statistic, error)
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

// ForDate handles an inbound HTTP GET request that returns an array of sensor statistics for a specific date
// from the Repository.
func (h *HTTP) ForDate(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	sensor := reading.SensorType(vars["sensor"])
	if !sensor.Valid() {
		http.Error(w, "invalid sensor type", http.StatusBadRequest)
		return
	}

	date, err := time.Parse("2006-01-02", vars["date"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	stats, err := h.statistics.ForDate(r.Context(), date, sensor)
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
	router.HandleFunc("/statistics/sensor/{sensor}/date/{date}", h.ForDate).Methods(http.MethodGet)
}

// Package ingest provides the HTTP handling methods for inbound sensor readings.
package ingest

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"

	"github.com/cloud-lada/backend/internal/reading"
	"github.com/gorilla/mux"
)

type (
	// The Ingestor type is responsible for responding to HTTP requests to publish events onto an event sink and
	// probing for liveness/readiness.
	Ingestor struct {
		writer EventWriter
		apiKey string
		logger *log.Logger
	}

	// The EventWriter interface describes types that can publish messages onto an event stream.
	EventWriter interface {
		Write(ctx context.Context, message json.RawMessage) error
	}

	// The Config type contains fields used to configure the Ingestor.
	Config struct {
		Writer EventWriter
		APIKey string
		Logger *log.Logger
	}
)

// New returns a new instance of the Ingestor type that will publish Reading events onto the provided
// EventWriter implementation to the desired subject. The Ingestor.Register method should be used to register
// the handling methods onto an HTTP router.
func New(config Config) *Ingestor {
	return &Ingestor{
		writer: config.Writer,
		apiKey: config.APIKey,
		logger: config.Logger,
	}
}

// Ingest readings from the request body, publishing each onto the configured EventWriter. This method expects
// the request body to contain a JSON stream of individual readings. Each reading is validated then published.
// It expects basic authentication on the inbound request where the password matches the configured API key.
func (h *Ingestor) Ingest(w http.ResponseWriter, r *http.Request) {
	apiKey, _, ok := r.BasicAuth()
	switch {
	case !ok:
		http.Error(w, "no credentials provided", http.StatusForbidden)
		return
	case apiKey != h.apiKey || apiKey == "":
		http.Error(w, "invalid api key", http.StatusUnauthorized)
		return
	}

	ctx := r.Context()
	decoder := json.NewDecoder(r.Body)

	// For efficiency, read the contents of the stream one JSON-object at a time. This will allow the server
	// to publish readings without loading the entire payload in-memory. It could be that we go substantial
	// amounts of time without an internet connection so there's a possibility of large uploads. Several days
	// worth of readings could trigger an OOM.
	for {
		select {
		case <-ctx.Done():
			http.Error(w, ctx.Err().Error(), http.StatusRequestTimeout)
			return
		default:
			var request reading.Reading

			// We decode each reading one-by-one to ensure their format is correct.
			err := decoder.Decode(&request)
			switch {
			case errors.Is(err, io.EOF):
				return
			case err != nil:
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			case !request.Valid():
				h.logger.Println("Invalid reading:", request)
				continue
			}

			data, err := json.Marshal(request)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			if err = h.writer.Write(ctx, data); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			h.logger.Println("Ingested reading:", request)
		}
	}
}

// Register the Ingestor's routes onto the HTTP router.
func (h *Ingestor) Register(router *mux.Router) {
	router.HandleFunc("/ingest", h.Ingest).
		Methods(http.MethodPost).
		Headers("Content-Type", "application/stream+json")
}

package reading

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type (
	// The HTTP type is responsible for handling HTTP requests and publishing events onto an event sink.
	HTTP struct {
		writer EventWriter
		logger *log.Logger
	}

	// The EventWriter interface describes types that can publish messages onto an event stream.
	EventWriter interface {
		Write(ctx context.Context, message json.RawMessage) error
	}
)

// NewHTTP returns a new instance of the HTTP type that will publish Reading events onto the provided
// EventWriter implementation to the desired subject. The HTTP.Register method should be used to register
// the handling methods onto an HTTP router.
func NewHTTP(events EventWriter, logger *log.Logger) *HTTP {
	return &HTTP{
		writer: events,
		logger: logger,
	}
}

type (
	// The IngestResponse type is the response DTO when calling HTTP.Ingest. It contains an array of all readings
	// that failed validation.
	IngestResponse struct {
		Invalid []Reading `json:"invalid,omitempty"`
	}
)

// Ingest readings from the request body, publishing each onto the configured EventWriter. This method expects
// the request body to contain a JSON stream of individual readings. Each reading is validated then published.
func (h *HTTP) Ingest(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	decoder := json.NewDecoder(r.Body)

	var resp IngestResponse

	// For efficiency, read the contents of the stream one JSON-object at a time. This will allow the server
	// to publish readings without loading the entire payload in-memory. It could be that we go substantial
	// amounts of time without an internet connection so there's a possibility of large uploads. Several days
	// worth of readings could trigger an OOM.
ingest:
	for {
		select {
		case <-ctx.Done():
			http.Error(w, ctx.Err().Error(), http.StatusRequestTimeout)
			return
		default:
			var request Reading

			// We decode each reading one-by-one to ensure their format is correct.
			err := decoder.Decode(&request)
			switch {
			case errors.Is(err, io.EOF):
				break ingest
			case err != nil:
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			case !request.Valid():
				resp.Invalid = append(resp.Invalid, request)
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

	w.Header().Set("Content-Type", "application/json")
	if len(resp.Invalid) > 0 {
		w.WriteHeader(http.StatusBadRequest)
	}

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// Register the HTTP's routes onto the HTTP router.
func (h *HTTP) Register(router *mux.Router) {
	router.HandleFunc("/ingest", h.Ingest).
		Methods(http.MethodPost).
		Headers("Content-Type", "application/stream+json")
}

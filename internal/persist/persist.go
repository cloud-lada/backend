// Package persist provides logic for handling inbound readings and persisting them to a repository.
package persist

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/cloud-lada/backend/internal/reading"
)

type (
	// The Persistor type is used to handle inbound readings and persist them to the Repository implementation.
	Persistor struct {
		logger   *log.Logger
		readings Repository
	}

	// The Config type contains fields used to configure a Persistor.
	Config struct {
		Logger   *log.Logger
		Readings Repository
	}

	// The Repository interface describes types that can store readings.
	Repository interface {
		Save(ctx context.Context, reading reading.Reading) error
	}
)

// New returns a new instance of the Persistor type, configured using the provided Config.
func New(config Config) *Persistor {
	return &Persistor{logger: config.Logger, readings: config.Readings}
}

// HandleEvent handles an inbound JSON payload. It expects the payload to be unmarshallable into a storage.Reading
// type. Once decoded, the reading is persisted via the Repository implementation.
func (h *Persistor) HandleEvent(ctx context.Context, message json.RawMessage) error {
	var request reading.Reading
	if err := json.Unmarshal(message, &request); err != nil {
		return fmt.Errorf("failed to unmarshal reading: %w", err)
	}

	if err := h.readings.Save(ctx, request); err != nil {
		return fmt.Errorf("failed to store reading: %w", err)
	}

	h.logger.Println("Persisted reading:", request)
	return nil
}

package reading

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
)

type (
	// The EventHandler type is used to handle inbound readings and persist them to the Repository implementation.
	EventHandler struct {
		logger   *log.Logger
		readings Repository
	}

	// The Repository interface describes types that can store readings.
	Repository interface {
		Save(ctx context.Context, reading Reading) error
	}
)

// NewEventHandler returns a new instance of the EventHandler type. Inbound events will be stored using the provided
// Repository implementation.
func NewEventHandler(readings Repository, logger *log.Logger) *EventHandler {
	return &EventHandler{logger: logger, readings: readings}
}

// HandleEvent handles an inbound JSON payload. It expects the payload to be unmarshallable into a storage.Reading
// type. Once decoded, the reading is persisted via the Repository implementation.
func (h *EventHandler) HandleEvent(ctx context.Context, message json.RawMessage) error {
	var request Reading
	if err := json.Unmarshal(message, &request); err != nil {
		return fmt.Errorf("failed to unmarshal reading: %w", err)
	}

	if err := h.readings.Save(ctx, request); err != nil {
		return fmt.Errorf("failed to store reading: %w", err)
	}

	h.logger.Println("Persisted reading:", request)
	return nil
}

package ingest_test

import (
	"context"
	"encoding/json"

	"github.com/cloud-lada/backend/internal/ingest"
)

type (
	MockEventWriter struct {
		ingest.EventWriter

		messages []MockMessage
		err      error
	}

	MockMessage struct {
		Data []byte
	}
)

func (m *MockEventWriter) Write(_ context.Context, message json.RawMessage) error {
	m.messages = append(m.messages, MockMessage{
		Data: message,
	})

	return m.err
}

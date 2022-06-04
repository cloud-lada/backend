package reading_test

import (
	"context"
	"encoding/json"

	"github.com/cloud-lada/backend/internal/reading"
)

type (
	MockEventWriter struct {
		reading.EventWriter

		messages []MockMessage
		err      error
	}

	MockMessage struct {
		Data []byte
	}

	MockRepository struct {
		saved reading.Reading
		err   error
	}
)

func (m *MockEventWriter) Write(_ context.Context, message json.RawMessage) error {
	m.messages = append(m.messages, MockMessage{
		Data: message,
	})

	return m.err
}

func (m *MockRepository) Save(ctx context.Context, reading reading.Reading) error {
	m.saved = reading
	return m.err
}

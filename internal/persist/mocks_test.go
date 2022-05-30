package persist_test

import (
	"context"

	"github.com/cloud-lada/backend/internal/reading"
)

type (
	MockRepository struct {
		saved reading.Reading
		err   error
	}
)

func (m *MockRepository) Save(ctx context.Context, reading reading.Reading) error {
	m.saved = reading
	return m.err
}

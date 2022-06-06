package status_test

import (
	"context"

	"github.com/cloud-lada/backend/internal/status"
)

type (
	MockRepository struct {
		status status.Status
		err    error
	}
)

func (m *MockRepository) Status(ctx context.Context) (status.Status, error) {
	return m.status, m.err
}

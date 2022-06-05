package location_test

import (
	"context"

	"github.com/cloud-lada/backend/internal/location"
)

type (
	MockRepository struct {
		location location.Location
		err      error
	}
)

func (m *MockRepository) Latest(ctx context.Context) (location.Location, error) {
	return m.location, m.err
}

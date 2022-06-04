package statistics_test

import (
	"context"

	"github.com/cloud-lada/backend/internal/statistics"
)

type (
	MockRepository struct {
		stats statistics.Statistics
		err   error
	}
)

func (m *MockRepository) Latest(ctx context.Context) (statistics.Statistics, error) {
	return m.stats, m.err
}

package statistics_test

import (
	"context"
	"time"

	"github.com/cloud-lada/backend/internal/reading"
	"github.com/cloud-lada/backend/internal/statistics"
)

type (
	MockRepository struct {
		latest statistics.Statistics
		stats  []statistics.Statistic
		err    error
	}
)

func (m *MockRepository) ForDate(ctx context.Context, date time.Time, sensor reading.SensorType) ([]statistics.Statistic, error) {
	return m.stats, m.err
}

func (m *MockRepository) Latest(ctx context.Context) (statistics.Statistics, error) {
	return m.latest, m.err
}

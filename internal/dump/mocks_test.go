package dump_test

import (
	"bytes"
	"context"
	"io"
	"time"

	"github.com/cloud-lada/backend/internal/reading"
)

type (
	MockRepository struct {
		err      error
		readings []reading.Reading
	}

	MockSink struct {
		name   string
		buffer *bytes.Buffer
	}

	NoopCloser struct {
		io.Writer
	}
)

func (n *NoopCloser) Close() error {
	return nil
}

func (m *MockSink) NewWriter(ctx context.Context, name string) (io.WriteCloser, error) {
	m.name = name
	return &NoopCloser{Writer: m.buffer}, nil
}

func (m *MockRepository) ForEachOnDate(ctx context.Context, date time.Time, fn reading.ForEachFunc) error {
	for _, r := range m.readings {
		if err := fn(ctx, r); err != nil {
			return err
		}
	}

	return m.err
}

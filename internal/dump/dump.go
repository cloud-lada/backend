// Package dump provides types that write a days worth of sensor data into a blob storage service.
package dump

import (
	"compress/gzip"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"time"

	"github.com/cloud-lada/backend/internal/reading"
)

type (
	// The Config type contains fields used to configure the Dumper.
	Config struct {
		Date     time.Time
		Readings Repository
		Blobs    Sink
		Logger   *log.Logger
	}

	// The Repository interface describes types that can iterate over reading data in a database.
	Repository interface {
		ForEachOnDate(ctx context.Context, date time.Time, fn reading.ForEachFunc) error
	}

	// The Sink interface describes types that provide keyed blobs where sensor data can be written.
	Sink interface {
		NewWriter(ctx context.Context, name string) (io.WriteCloser, error)
	}

	// The Dumper type is used to perform a daily dump of sensor data into a blob storage provider.
	Dumper struct {
		readings Repository
		blobs    Sink
		logger   *log.Logger
		date     time.Time
	}
)

// New returns a new instance of the Dumper type for the provided configuration values. For safety, it will automatically
// set Config.Date to be the earliest possible time for that day.
func New(config Config) *Dumper {
	return &Dumper{
		readings: config.Readings,
		blobs:    config.Blobs,
		logger:   config.Logger,
		// We want to get all the data for a given day, so we need to start at 00:00:00 for that specific day.
		date: time.Date(config.Date.Year(), config.Date.Month(), config.Date.Day(), 0, 0, 0, 0, config.Date.Location()),
	}
}

// Dump JSON-encoded readings for the configured date into the blob storage provider. Dumps will be JSON streams
// similar to how they are originally presented to the ingestor.
func (d *Dumper) Dump(ctx context.Context) error {
	name := d.date.Format("2006-02-01.json.gz")
	blob, err := d.blobs.NewWriter(ctx, name)
	if err != nil {
		return fmt.Errorf("failed to open blob: %w", err)
	}
	defer d.close(blob)

	archive := gzip.NewWriter(blob)
	defer d.close(archive)

	encoder := json.NewEncoder(archive)
	return d.readings.ForEachOnDate(ctx, d.date, func(ctx context.Context, reading reading.Reading) error {
		return encoder.Encode(reading)
	})
}

func (d *Dumper) close(c io.Closer) {
	if err := c.Close(); err != nil {
		d.logger.Println(err.Error())
	}
}

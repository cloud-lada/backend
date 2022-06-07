package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/cloud-lada/backend/internal/dump"
	"github.com/cloud-lada/backend/internal/reading"
	"github.com/cloud-lada/backend/pkg/blob"
	"github.com/cloud-lada/backend/pkg/closers"
	"github.com/cloud-lada/backend/pkg/postgres"
	"github.com/spf13/cobra"
)

var version = "dev"

func main() {
	var (
		blobStoreURL string
		databaseURL  string
		dumpDate     string
	)

	cmd := &cobra.Command{
		Use:     "dumper",
		Short:   "Dumps an entire day's worth of sensor data into blob storage",
		Version: version,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			date, err := time.Parse("2006-02-01", dumpDate)
			if err != nil {
				return fmt.Errorf("invalid dump date: %w", err)
			}

			db, err := postgres.Open(ctx, databaseURL)
			if err != nil {
				return fmt.Errorf("failed to connect to database: %w", err)
			}
			defer closers.Close(db)

			blobs, err := blob.Open(ctx, blobStoreURL)
			if err != nil {
				return fmt.Errorf("failed to connect to blob storage: %w", err)
			}
			defer closers.Close(db)

			logger := log.Default()
			dumper := dump.New(dump.Config{
				Date:     date,
				Readings: reading.NewPostgresRepository(db),
				Blobs:    blobs,
			})

			logger.Println("Creating dump for", dumpDate)
			return dumper.Dump(ctx)
		},
	}

	flags := cmd.PersistentFlags()
	flags.StringVar(&blobStoreURL, "blob-store-url", "", "The URL of the blob store to persist dumps to")
	flags.StringVar(&databaseURL, "database-url", "", "The URL of the database to read data from")
	flags.StringVar(&dumpDate, "dump-date", time.Now().Add(-time.Hour*24).Format("2006-02-01"), "The date to dump data for")

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, os.Kill, syscall.SIGTERM)
	if err := cmd.ExecuteContext(ctx); err != nil {
		cancel()
		os.Exit(1)
	}
}

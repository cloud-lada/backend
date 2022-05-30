package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/cloud-lada/backend/internal/persist"
	"github.com/cloud-lada/backend/internal/reading"
	"github.com/cloud-lada/backend/pkg/event"
	"github.com/cloud-lada/backend/pkg/postgres"
	"github.com/spf13/cobra"
	"golang.org/x/sync/errgroup"
)

var version = "dev"

func main() {
	var (
		eventReaderURL string
		databaseURL    string
	)

	cmd := &cobra.Command{
		Use:     "persistor",
		Short:   "Listens for sensor readings from en event bus and persists them to TimescaleDB",
		Version: version,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()

			db, err := postgres.Open(ctx, databaseURL)
			if err != nil {
				return fmt.Errorf("failed to connect to database: %w", err)
			}

			logger := log.Default()

			reader, err := event.NewReader(ctx, eventReaderURL)
			if err != nil {
				return fmt.Errorf("failed to create reader: %w", err)
			}

			persistor := persist.New(persist.Config{
				Logger:   logger,
				Readings: reading.NewPostgresRepository(db),
			})

			grp, ctx := errgroup.WithContext(ctx)
			grp.Go(func() error {
				return reader.Read(ctx, persistor.HandleEvent)
			})
			grp.Go(func() error {
				<-ctx.Done()
				return reader.Close()
			})
			grp.Go(func() error {
				<-ctx.Done()
				return db.Close()
			})

			logger.Println("Listening for events from", eventReaderURL)
			return grp.Wait()
		},
	}

	flags := cmd.PersistentFlags()
	flags.StringVar(&eventReaderURL, "event-reader-url", "", "The URL of the event bus to read messages from")
	flags.StringVar(&databaseURL, "database-url", "", "The URL of the database to persist data to")

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, os.Kill, syscall.SIGTERM)
	if err := cmd.ExecuteContext(ctx); err != nil {
		cancel()
		os.Exit(1)
	}
}

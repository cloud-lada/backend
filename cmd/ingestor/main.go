package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/cloud-lada/backend/internal/reading"
	"github.com/cloud-lada/backend/pkg/event"
	"github.com/cloud-lada/backend/pkg/middleware"
	"github.com/gorilla/mux"
	"github.com/spf13/cobra"
	"golang.org/x/sync/errgroup"
)

var version = "dev"

func main() {
	var (
		eventWriterURL string
		apiKey         string
		port           int
	)

	cmd := &cobra.Command{
		Use:     "ingestor",
		Short:   "Accepts sensor readings over HTTP and publishes them to a NATS subject",
		Version: version,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()

			writer, err := event.NewWriter(ctx, eventWriterURL)
			if err != nil {
				return fmt.Errorf("failed to connect to event bus: %w", err)
			}

			logger := log.Default()

			router := mux.NewRouter()
			router.Use(middleware.APIKey(apiKey))

			handler := reading.NewHTTP(writer, logger)
			handler.Register(router)

			svr := &http.Server{
				Addr:    fmt.Sprint(":", port),
				Handler: router,
			}

			grp, ctx := errgroup.WithContext(ctx)
			grp.Go(func() error {
				return svr.ListenAndServe()
			})
			grp.Go(func() error {
				<-ctx.Done()
				return svr.Shutdown(context.Background())
			})
			grp.Go(func() error {
				<-ctx.Done()
				return writer.Close()
			})

			logger.Println("Server started on port", port)
			return grp.Wait()
		},
	}

	flags := cmd.PersistentFlags()
	flags.IntVar(&port, "port", 5000, "The port to listen for HTTP requests from")
	flags.StringVar(&eventWriterURL, "event-writer-url", "", "The URL of the event bus to send messages to")
	flags.StringVar(&apiKey, "api-key", "", "The API key to use for basic authentication")

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, os.Kill, syscall.SIGTERM)
	if err := cmd.ExecuteContext(ctx); err != nil {
		cancel()
		os.Exit(1)
	}
}

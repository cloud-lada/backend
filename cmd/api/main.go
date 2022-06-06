package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/cloud-lada/backend/internal/location"
	"github.com/cloud-lada/backend/internal/statistics"
	"github.com/cloud-lada/backend/internal/status"
	"github.com/cloud-lada/backend/pkg/postgres"
	"github.com/gorilla/mux"
	"github.com/spf13/cobra"
	"golang.org/x/sync/errgroup"
)

var version = "dev"

func main() {
	var (
		databaseURL string
		port        int
	)

	cmd := &cobra.Command{
		Use:     "api",
		Short:   "Serves statistical data from the database via HTTP",
		Version: version,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()

			db, err := postgres.Open(ctx, databaseURL)
			if err != nil {
				return fmt.Errorf("failed to connect to database: %w", err)
			}

			logger := log.Default()
			router := mux.NewRouter()
			api := router.PathPrefix("/api").Subrouter()

			statistics.NewHTTP(statistics.NewPostgresRepository(db)).Register(api)
			location.NewHTTP(location.NewPostgresRepository(db)).Register(api)
			status.NewHTTP(status.NewPostgresRepository(db)).Register(api)

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

			logger.Println("Server started on port", port)
			return grp.Wait()
		},
	}

	flags := cmd.PersistentFlags()
	flags.IntVar(&port, "port", 5000, "The port to listen for HTTP requests from")
	flags.StringVar(&databaseURL, "database-url", "", "The URL of the database to read data from")

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, os.Kill, syscall.SIGTERM)
	if err := cmd.ExecuteContext(ctx); err != nil {
		cancel()
		os.Exit(1)
	}
}

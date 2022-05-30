// Package testutil provides test helper functions.
package testutil

import (
	"context"
	"database/sql"
	"testing"

	"github.com/cloud-lada/backend/pkg/postgres"
	"github.com/stretchr/testify/require"
)

// Context returns a context.Context that is cancelled when the test finishes.
func Context(t *testing.T) context.Context {
	ctx, cancel := context.WithCancel(context.Background())

	t.Cleanup(cancel)
	return ctx
}

// Postgres returns an sql.DB instance assumed to be available on localhost:5432. Migrations are performed and all
// data is deleted once the test is complete.
func Postgres(t *testing.T, ctx context.Context) *sql.DB {
	t.Helper()

	db, err := postgres.Open(ctx, "postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable")
	require.NoError(t, err)

	t.Cleanup(func() {
		_, err = db.ExecContext(ctx, "DELETE FROM reading WHERE sensor IS NOT NULl")
		require.NoError(t, err)
		require.NoError(t, db.Close())
	})

	return db
}

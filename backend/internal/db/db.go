package db

import (
	"context"
	"embed"
	"fmt"
	"log"
	"sort"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"
)

//go:embed migrations/*.sql
var migrationFS embed.FS

func Connect(ctx context.Context, databaseURL string) (*pgxpool.Pool, error) {
	if databaseURL == "" {
		return nil, fmt.Errorf("DATABASE_URL is empty")
	}
	pool, err := pgxpool.New(ctx, databaseURL)
	if err != nil {
		return nil, err
	}
	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, err
	}
	return pool, nil
}

func Migrate(ctx context.Context, pool *pgxpool.Pool) error {
	_, err := pool.Exec(ctx, `
		CREATE TABLE IF NOT EXISTS schema_migrations (
			version TEXT PRIMARY KEY,
			applied_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
		)`)
	if err != nil {
		return fmt.Errorf("schema_migrations: %w", err)
	}

	entries, err := migrationFS.ReadDir("migrations")
	if err != nil {
		return err
	}

	names := make([]string, 0, len(entries))
	for _, e := range entries {
		if !e.IsDir() && strings.HasSuffix(e.Name(), ".sql") {
			names = append(names, e.Name())
		}
	}
	sort.Strings(names)

	var appExists bool
	if err := pool.QueryRow(ctx, `
		SELECT EXISTS(
			SELECT 1 FROM information_schema.tables
			WHERE table_schema = 'public' AND table_name = 'applications'
		)`).Scan(&appExists); err != nil {
		return err
	}
	var appliedCount int
	if err := pool.QueryRow(ctx, `SELECT COUNT(*) FROM schema_migrations`).Scan(&appliedCount); err != nil {
		return err
	}
	if appExists && appliedCount == 0 {
		for _, name := range names {
			if _, err := pool.Exec(ctx,
				`INSERT INTO schema_migrations (version) VALUES ($1) ON CONFLICT DO NOTHING`, name,
			); err != nil {
				return err
			}
			log.Printf("recorded existing migration %s (already on db)", name)
		}
	}

	for _, name := range names {
		var exists bool
		err := pool.QueryRow(ctx,
			`SELECT EXISTS(SELECT 1 FROM schema_migrations WHERE version=$1)`, name,
		).Scan(&exists)
		if err != nil {
			return err
		}
		if exists {
			continue
		}

		body, err := migrationFS.ReadFile("migrations/" + name)
		if err != nil {
			return err
		}

		tx, err := pool.Begin(ctx)
		if err != nil {
			return err
		}
		if _, err := tx.Exec(ctx, string(body)); err != nil {
			_ = tx.Rollback(ctx)
			return fmt.Errorf("migration %s: %w", name, err)
		}
		if _, err := tx.Exec(ctx,
			`INSERT INTO schema_migrations (version) VALUES ($1)`, name,
		); err != nil {
			_ = tx.Rollback(ctx)
			return err
		}
		if err := tx.Commit(ctx); err != nil {
			return err
		}
		log.Printf("applied migration %s", name)
	}
	return nil
}

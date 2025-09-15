package postgres

import (
	"context"
	"embed"
	"fmt"
	"io/fs"
	"sort"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

//go:embed migrations/*.sql
var migrationsFS embed.FS

// RunMigrations 执行嵌入的 SQL 迁移，按文件名排序顺序执行。已应用的 migration 会被跳过。
func RunMigrations(ctx context.Context, pool *pgxpool.Pool) error {
	entries, err := fs.ReadDir(migrationsFS, "migrations")
	if err != nil {
		return fmt.Errorf("read migrations: %w", err)
	}
	var names []string
	for _, e := range entries {
		if !e.IsDir() {
			names = append(names, e.Name())
		}
	}
	sort.Strings(names)

	// ensure schema_migrations table exists
	const createTableSQL = `CREATE TABLE IF NOT EXISTS schema_migrations (version TEXT PRIMARY KEY, applied_at TIMESTAMPTZ NOT NULL DEFAULT now());`
	if _, err := pool.Exec(ctx, createTableSQL); err != nil {
		return fmt.Errorf("ensure schema_migrations: %w", err)
	}

	// load applied migrations
	rows, err := pool.Query(ctx, `SELECT version FROM schema_migrations`)
	if err != nil {
		return fmt.Errorf("query applied migrations: %w", err)
	}
	defer rows.Close()
	applied := make(map[string]struct{})
	for rows.Next() {
		var v string
		if err := rows.Scan(&v); err != nil {
			return fmt.Errorf("scan applied migrations: %w", err)
		}
		applied[v] = struct{}{}
	}
	if rows.Err() != nil {
		return fmt.Errorf("read applied migrations: %w", rows.Err())
	}

	for _, name := range names {
		if _, ok := applied[name]; ok {
			// skip already applied
			continue
		}

		b, err := migrationsFS.ReadFile("migrations/" + name)
		if err != nil {
			return fmt.Errorf("read migration %s: %w", name, err)
		}
		sql := string(b)

		// run in transaction with timeout
		ctx2, cancel := context.WithTimeout(ctx, 30*time.Second)
		tx, err := pool.Begin(ctx2)
		if err != nil {
			cancel()
			return fmt.Errorf("begin tx for %s: %w", name, err)
		}

		if _, err := tx.Exec(ctx2, sql); err != nil {
			tx.Rollback(ctx2)
			cancel()
			return fmt.Errorf("apply migration %s: %w", name, err)
		}

		if _, err := tx.Exec(ctx2, `INSERT INTO schema_migrations(version, applied_at) VALUES($1, now())`, name); err != nil {
			tx.Rollback(ctx2)
			cancel()
			return fmt.Errorf("record migration %s: %w", name, err)
		}

		if err := tx.Commit(ctx2); err != nil {
			cancel()
			return fmt.Errorf("commit migration %s: %w", name, err)
		}
		cancel()
	}
	return nil
}

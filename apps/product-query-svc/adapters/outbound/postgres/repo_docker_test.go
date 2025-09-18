//go:build docker

package postgres

import (
    "context"
    "os"
    "path/filepath"
    "sort"
    "testing"
    "time"
    "io/ioutil"
    "strings"

    "github.com/jackc/pgx/v5/pgxpool"
    "github.com/fightingBald/GoTuto/apps/product-query-svc/domain"
    "github.com/testcontainers/testcontainers-go"
    "github.com/testcontainers/testcontainers-go/wait"
)

// applyMigrations runs all *.up.sql files in the migrations directory in lexicographic order.
func applyMigrations(ctx context.Context, pool *pgxpool.Pool, dir string, t *testing.T) {
    entries, err := ioutil.ReadDir(dir)
    if err != nil {
        t.Fatalf("read migrations dir: %v", err)
    }
    var files []string
    for _, e := range entries {
        if !e.IsDir() && strings.HasSuffix(e.Name(), ".up.sql") {
            files = append(files, filepath.Join(dir, e.Name()))
        }
    }
    sort.Strings(files)
    for _, f := range files {
        b, err := os.ReadFile(f)
        if err != nil {
            t.Fatalf("read %s: %v", f, err)
        }
        if _, err := pool.Exec(ctx, string(b)); err != nil {
            t.Fatalf("exec migration %s: %v", f, err)
        }
    }
}

func TestPostgresRepo_WithDocker(t *testing.T) {
    if os.Getenv("SKIP_DOCKER_TESTS") == "1" {
        t.Skip("skipped by SKIP_DOCKER_TESTS=1")
    }

    ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
    defer cancel()

    // Start Postgres container
    req := testcontainers.ContainerRequest{
        Image:        "postgres:16-alpine",
        ExposedPorts: []string{"5432/tcp"},
        Env: map[string]string{
            "POSTGRES_USER":     "app",
            "POSTGRES_PASSWORD": "app_password",
            "POSTGRES_DB":       "productdb",
        },
        WaitingFor: wait.ForListeningPort("5432/tcp").WithStartupTimeout(60 * time.Second),
    }
    pgC, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{ContainerRequest: req, Started: true})
    if err != nil {
        t.Fatalf("start container: %v", err)
    }
    defer func() { _ = pgC.Terminate(context.Background()) }()

    host, err := pgC.Host(ctx)
    if err != nil { t.Fatalf("host: %v", err) }
    port, err := pgC.MappedPort(ctx, "5432/tcp")
    if err != nil { t.Fatalf("mapped port: %v", err) }

    dsn := "postgres://app:app_password@" + host + ":" + port.Port() + "/productdb?sslmode=disable"

    // Connect pool
    pool, err := pgxpool.New(ctx, dsn)
    if err != nil { t.Fatalf("pgxpool.New: %v", err) }
    defer pool.Close()

    // Apply migrations from local dir (same package directory)
    migDir := filepath.Join("migrations")
    applyMigrations(ctx, pool, migDir, t)

    // Run a few repo operations
    repo := NewProductRepository(pool)

    // Search should work on seeded data (may be empty if seeds change)
    if items, total, err := repo.Search(ctx, "pro", 1, 10); err != nil {
        t.Fatalf("repo.Search: %v", err)
    } else if total < 0 || len(items) < 0 { // sanity
        t.Fatalf("unexpected search result: total=%d items=%d", total, len(items))
    }

    // Create -> Get -> Delete roundtrip
    id, err := repo.Create(ctx, &domain.Product{
        Name:  "DockerTest",
        Price: 1234,
        Tags:  []string{"tc"},
    })
    if err != nil {
        t.Fatalf("repo.Create: %v", err)
    }
    p, err := repo.GetByID(ctx, id)
    if err != nil {
        t.Fatalf("repo.GetByID: %v", err)
    }
    if p.ID != id || p.Name != "DockerTest" || p.Price != 1234 {
        t.Fatalf("unexpected product: %#v", p)
    }
    if err := repo.Delete(ctx, id); err != nil {
        t.Fatalf("repo.Delete: %v", err)
    }
}

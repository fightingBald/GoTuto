package testutil

import (
    "context"
    "os"
    "path/filepath"
    "sort"
    "strings"
    "testing"

    "github.com/jackc/pgx/v5/pgxpool"
    "github.com/testcontainers/testcontainers-go"
    "github.com/testcontainers/testcontainers-go/wait"
)

// DSNFromEnvOrDocker returns a Postgres DSN. If DATABASE_URL is set, it is used.
// Otherwise a temporary Docker Postgres is started. The boolean indicates
// whether a temp container was started (call cleanup when true).
func DSNFromEnvOrDocker(ctx context.Context, t testing.TB) (dsn string, isTemp bool, cleanup func()) {
    t.Helper()
    if v := os.Getenv("DATABASE_URL"); v != "" {
        return v, false, func() {}
    }
    dsn, cleanup = StartDockerPostgres(ctx, t)
    return dsn, true, cleanup
}

// StartDockerPostgres launches a temporary Postgres container and returns the DSN and cleanup.
func StartDockerPostgres(ctx context.Context, t testing.TB) (dsn string, cleanup func()) {
    t.Helper()
    req := testcontainers.ContainerRequest{
        Image:        "postgres:16-alpine",
        ExposedPorts: []string{"5432/tcp"},
        Env: map[string]string{
            "POSTGRES_USER":     "app",
            "POSTGRES_PASSWORD": "app_password",
            "POSTGRES_DB":       "productdb",
        },
        WaitingFor: wait.ForListeningPort("5432/tcp"),
    }
    pgC, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{ContainerRequest: req, Started: true})
    if err != nil { t.Fatalf("start container: %v", err) }
    host, err := pgC.Host(ctx)
    if err != nil { t.Fatalf("host: %v", err) }
    port, err := pgC.MappedPort(ctx, "5432/tcp")
    if err != nil { t.Fatalf("mapped port: %v", err) }
    cleanup = func() { _ = pgC.Terminate(context.Background()) }
    dsn = "postgres://app:app_password@" + host + ":" + port.Port() + "/productdb?sslmode=disable"
    return dsn, cleanup
}

// ApplyMigrations applies repository migrations in order to the given pool.
// It looks for migrations under apps/product-query-svc/adapters/outbound/postgres/migrations
// starting from the module root.
func ApplyMigrations(ctx context.Context, t testing.TB, pool *pgxpool.Pool) {
    t.Helper()
    root := moduleRoot(t)
    migDir := filepath.Join(root, "apps", "product-query-svc", "adapters", "outbound", "postgres", "migrations")
    entries, err := os.ReadDir(migDir)
    if err != nil { t.Fatalf("read dir: %v", err) }
    var files []string
    for _, e := range entries {
        if !e.IsDir() && strings.HasSuffix(e.Name(), ".up.sql") {
            files = append(files, filepath.Join(migDir, e.Name()))
        }
    }
    sort.Strings(files)
    for _, f := range files {
        b, err := os.ReadFile(f)
        if err != nil { t.Fatalf("read %s: %v", f, err) }
        if _, err := pool.Exec(ctx, string(b)); err != nil { t.Fatalf("exec %s: %v", f, err) }
    }
}

// NewPool creates a pgxpool.Pool and fails the test on error.
func NewPool(ctx context.Context, t testing.TB, dsn string) *pgxpool.Pool {
    t.Helper()
    pool, err := pgxpool.New(ctx, dsn)
    if err != nil { t.Fatalf("pgxpool.New: %v", err) }
    return pool
}

// moduleRoot returns the directory containing go.mod by walking up from CWD.
func moduleRoot(t testing.TB) string {
    t.Helper()
    dir, err := os.Getwd()
    if err != nil { t.Fatalf("getwd: %v", err) }
    for {
        if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
            return dir
        }
        parent := filepath.Dir(dir)
        if parent == dir { t.Fatalf("go.mod not found from %s", dir) }
        dir = parent
    }
}


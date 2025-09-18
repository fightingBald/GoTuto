package test

import (
    "context"
    "encoding/json"
    "net/http"
    "net/http/httptest"
    "os"
    "path/filepath"
    "sort"
    "strconv"
    "strings"
    "testing"
    "time"

    appshttp "github.com/fightingBald/GoTuto/apps/product-query-svc/adapters/inbound/http"
    appspg "github.com/fightingBald/GoTuto/apps/product-query-svc/adapters/outbound/postgres"
    appsvc "github.com/fightingBald/GoTuto/apps/product-query-svc/app"
    "github.com/go-chi/chi/v5"
    "github.com/jackc/pgx/v5/pgxpool"
    "github.com/testcontainers/testcontainers-go"
    "github.com/testcontainers/testcontainers-go/wait"
)

// startDockerPostgres launches a temporary Postgres and returns DSN and cleanup.
func startDockerPostgres(t *testing.T, ctx context.Context) (dsn string, cleanup func()) {
    t.Helper()
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
    if err != nil { t.Fatalf("start container: %v", err) }
    host, err := pgC.Host(ctx)
    if err != nil { t.Fatalf("host: %v", err) }
    port, err := pgC.MappedPort(ctx, "5432/tcp")
    if err != nil { t.Fatalf("mapped port: %v", err) }
    cleanup = func() { _ = pgC.Terminate(context.Background()) }
    dsn = "postgres://app:app_password@" + host + ":" + port.Port() + "/productdb?sslmode=disable"
    return dsn, cleanup
}

// applyMigrations applies all *.up.sql from repo migrations in order.
func applyMigrations(t *testing.T, ctx context.Context, pool *pgxpool.Pool) {
    t.Helper()
    migDir := filepath.Join("..", "apps", "product-query-svc", "adapters", "outbound", "postgres", "migrations")
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

// TestCreateProduct_Postgres validates POST /products on a real Postgres.
// Requires DATABASE_URL to be set; otherwise it is skipped.
func TestCreateProduct_Postgres(t *testing.T) {
    dsn := os.Getenv("DATABASE_URL")
    var cleanup func()

    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()
    if dsn == "" {
        // fallback to a temporary Docker Postgres if no env provided
        var c func()
        dsn, c = startDockerPostgres(t, ctx)
        cleanup = c
    }

    pool, err := pgxpool.New(ctx, dsn)
    if err != nil {
        t.Fatalf("pgxpool.New: %v", err)
    }
    defer pool.Close()

    // Apply migrations only for temporary container
    if cleanup != nil {
        applyMigrations(t, ctx, pool)
    }

    repo := appspg.NewProductRepository(pool)
    svc := appsvc.NewProductService(repo)
    server := appshttp.NewServer(svc)

    r := chi.NewRouter()
    h := appshttp.HandlerFromMux(server, r)

    ts := httptest.NewServer(h)
    defer ts.Close()

    // Create
    body := `{"name":"CI Test Item","price":12.34}`
    resp, err := http.Post(ts.URL+"/products", "application/json", strings.NewReader(body))
    if err != nil {
        t.Fatalf("http post: %v", err)
    }
    defer resp.Body.Close()
    if resp.StatusCode != http.StatusCreated {
        t.Fatalf("expected 201, got %d", resp.StatusCode)
    }
    var created appshttp.Product
    if err := json.NewDecoder(resp.Body).Decode(&created); err != nil {
        t.Fatalf("decode created: %v", err)
    }
    if created.Id == 0 || created.Name != "CI Test Item" {
        t.Fatalf("unexpected created product: %+v", created)
    }

    // Fetch the created product by ID
    resp2, err := http.Get(ts.URL + "/products/" + strconv.FormatInt(created.Id, 10))
    if err != nil {
        t.Fatalf("http get: %v", err)
    }
    defer resp2.Body.Close()
    if resp2.StatusCode != http.StatusOK {
        t.Fatalf("expected 200, got %d", resp2.StatusCode)
    }
    var fetched appshttp.Product
    if err := json.NewDecoder(resp2.Body).Decode(&fetched); err != nil {
        t.Fatalf("decode fetched: %v", err)
    }
    if fetched.Id != created.Id || fetched.Name != created.Name {
        t.Fatalf("mismatch fetched vs created: created=%+v fetched=%+v", created, fetched)
    }

    // Cleanup: delete and verify 404 afterwards
    req, _ := http.NewRequest(http.MethodDelete, ts.URL+"/products/"+strconv.FormatInt(created.Id, 10), nil)
    resp3, err := http.DefaultClient.Do(req)
    if err != nil {
        t.Fatalf("http delete: %v", err)
    }
    resp3.Body.Close()
    if resp3.StatusCode != http.StatusNoContent {
        t.Fatalf("expected 204, got %d", resp3.StatusCode)
    }

    resp4, err := http.Get(ts.URL + "/products/" + strconv.FormatInt(created.Id, 10))
    if err != nil {
        t.Fatalf("http get after delete: %v", err)
    }
    resp4.Body.Close()
    if resp4.StatusCode != http.StatusNotFound {
        t.Fatalf("expected 404 after delete, got %d", resp4.StatusCode)
    }
    if cleanup != nil { cleanup() }
}

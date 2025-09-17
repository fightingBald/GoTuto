package test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	appshttp "github.com/fightingBald/GoTuto/apps/product-query-svc/adapters/inbound/http"
	appspg "github.com/fightingBald/GoTuto/apps/product-query-svc/adapters/outbound/postgres"
	appsvc "github.com/fightingBald/GoTuto/apps/product-query-svc/app"
	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// TestSearchProducts_Postgres seeds are applied via migrations in dev/CI.
// This test requires DATABASE_URL to be set; otherwise it is skipped.
func TestSearchProducts_Postgres(t *testing.T) {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		t.Skip("DATABASE_URL not set; skipping Postgres integration test")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		t.Fatalf("pgxpool.New: %v", err)
	}
	defer pool.Close()

	repo := appspg.NewProductRepository(pool)
	svc := appsvc.NewProductService(repo)
	server := appshttp.NewServer(svc)

	r := chi.NewRouter()
	h := appshttp.HandlerFromMux(server, r)

	ts := httptest.NewServer(h)
	defer ts.Close()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, ts.URL+"/products/search?q=wid&page=1&pageSize=10", nil)
	if err != nil {
		t.Fatalf("new request: %v", err)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("http do: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("unexpected status: %d", resp.StatusCode)
	}

	var pl appshttp.ProductList
	if err := json.NewDecoder(resp.Body).Decode(&pl); err != nil {
		t.Fatalf("decode: %v", err)
	}

	if pl.Total < len(pl.Items) {
		t.Fatalf("expected total >= items length; got total=%d items=%d", pl.Total, len(pl.Items))
	}
}

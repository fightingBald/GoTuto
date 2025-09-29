package http_pg_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"
	"time"

	appshttp "github.com/fightingBald/GoTuto/apps/product-query-svc/adapters/inbound/http"
	appspg "github.com/fightingBald/GoTuto/apps/product-query-svc/adapters/outbound/postgres"
	productapp "github.com/fightingBald/GoTuto/apps/product-query-svc/application/product"
	userapp "github.com/fightingBald/GoTuto/apps/product-query-svc/application/user"
	"github.com/fightingBald/GoTuto/internal/testutil"
	"github.com/go-chi/chi/v5"
)

// TestCreateProduct_Postgres validates POST /products on a real Postgres.
func TestCreateProduct_Postgres(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	pool := testutil.NewPool(ctx, t, pgDSN)
	defer pool.Close()
	if pgTemp {
		testutil.ApplyMigrations(ctx, t, pool)
	}

	productRepo := appspg.NewProductRepository(pool)
	userRepo := appspg.NewUserRepository(pool)
	productSvc := productapp.NewService(productRepo)
	userSvc := userapp.NewService(userRepo)
	server := appshttp.NewServer(productSvc, userSvc)

	r := chi.NewRouter()
	strict := appshttp.NewStrictHTTPHandler(server, nil)
	h := appshttp.HandlerFromMux(strict, r)

	ts := httptest.NewServer(h)
	defer ts.Close()

	var created appshttp.Product
	t.Run("create 201", func(t *testing.T) {
		body := `{"name":"CI Test Item","price":12.34}`
		resp, err := http.Post(ts.URL+"/products", "application/json", strings.NewReader(body))
		if err != nil {
			t.Fatalf("http post: %v", err)
		}
		t.Cleanup(func() { resp.Body.Close() })
		if resp.StatusCode != http.StatusCreated {
			t.Fatalf("expected 201, got %d", resp.StatusCode)
		}
		if err := json.NewDecoder(resp.Body).Decode(&created); err != nil {
			t.Fatalf("decode created: %v", err)
		}
		if created.Id == 0 || created.Name != "CI Test Item" {
			t.Fatalf("unexpected created product: %+v", created)
		}
	})

	t.Run("get 200", func(t *testing.T) {
		resp2, err := http.Get(ts.URL + "/products/" + strconv.FormatInt(created.Id, 10))
		if err != nil {
			t.Fatalf("http get: %v", err)
		}
		t.Cleanup(func() { resp2.Body.Close() })
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
	})

	t.Run("update 200", func(t *testing.T) {
		body := `{"name":"CI Test Item Updated","price":15.67}`
		req, _ := http.NewRequest(http.MethodPut, ts.URL+"/products/"+strconv.FormatInt(created.Id, 10), strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Fatalf("http put: %v", err)
		}
		t.Cleanup(func() { resp.Body.Close() })
		if resp.StatusCode != http.StatusOK {
			t.Fatalf("expected 200, got %d", resp.StatusCode)
		}
		var updated appshttp.Product
		if err := json.NewDecoder(resp.Body).Decode(&updated); err != nil {
			t.Fatalf("decode updated: %v", err)
		}
		if updated.Name != "CI Test Item Updated" {
			t.Fatalf("expected updated name, got %+v", updated)
		}
		created = updated
	})

	t.Run("get after update 200", func(t *testing.T) {
		resp, err := http.Get(ts.URL + "/products/" + strconv.FormatInt(created.Id, 10))
		if err != nil {
			t.Fatalf("http get: %v", err)
		}
		t.Cleanup(func() { resp.Body.Close() })
		if resp.StatusCode != http.StatusOK {
			t.Fatalf("expected 200, got %d", resp.StatusCode)
		}
		var fetched appshttp.Product
		if err := json.NewDecoder(resp.Body).Decode(&fetched); err != nil {
			t.Fatalf("decode fetched: %v", err)
		}
		if fetched.Name != "CI Test Item Updated" {
			t.Fatalf("expected updated name, got %+v", fetched)
		}
	})

	t.Run("delete 204", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodDelete, ts.URL+"/products/"+strconv.FormatInt(created.Id, 10), nil)
		resp3, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Fatalf("http delete: %v", err)
		}
		resp3.Body.Close()
		if resp3.StatusCode != http.StatusNoContent {
			t.Fatalf("expected 204, got %d", resp3.StatusCode)
		}
	})

	t.Run("get after delete 404", func(t *testing.T) {
		resp4, err := http.Get(ts.URL + "/products/" + strconv.FormatInt(created.Id, 10))
		if err != nil {
			t.Fatalf("http get after delete: %v", err)
		}
		resp4.Body.Close()
		if resp4.StatusCode != http.StatusNotFound {
			t.Fatalf("expected 404 after delete, got %d", resp4.StatusCode)
		}
	})
}

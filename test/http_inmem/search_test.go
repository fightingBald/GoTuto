package http_inmem_test

import (
	"encoding/json"
	"net/http"
	"strings"
	"testing"

	appshttp "github.com/fightingBald/GoTuto/apps/product-query-svc/adapters/inbound/http"
	appsinmem "github.com/fightingBald/GoTuto/apps/product-query-svc/adapters/outbound/inmem"
	"github.com/fightingBald/GoTuto/internal/testutil"
)

func TestHTTP_InMem_Product(t *testing.T) {
	t.Run("search returns items", func(t *testing.T) {
		ts := testutil.NewHTTPServerWithRepo(appsinmem.NewInMemRepo())
		defer ts.Close()

		resp, err := http.Get(ts.URL + "/products/search?q=wid&page=1&pageSize=10")
		if err != nil {
			t.Fatalf("http get: %v", err)
		}
		defer resp.Body.Close()
		if resp.StatusCode != http.StatusOK {
			t.Fatalf("unexpected status: %d", resp.StatusCode)
		}
		var pl appshttp.ProductList
		if err := json.NewDecoder(resp.Body).Decode(&pl); err != nil {
			t.Fatalf("decode: %v", err)
		}
		if len(pl.Items) == 0 || pl.Total == 0 {
			t.Fatalf("expected seeded items in in-memory repo; got items=%d total=%d", len(pl.Items), pl.Total)
		}
	})

	t.Run("get id=1 returns product", func(t *testing.T) {
		ts := testutil.NewHTTPServerWithRepo(appsinmem.NewInMemRepo())
		defer ts.Close()

		resp, err := http.Get(ts.URL + "/products/1")
		if err != nil {
			t.Fatalf("http get: %v", err)
		}
		defer resp.Body.Close()
		if resp.StatusCode != http.StatusOK {
			t.Fatalf("unexpected status: %d", resp.StatusCode)
		}
		var p appshttp.Product
		if err := json.NewDecoder(resp.Body).Decode(&p); err != nil {
			t.Fatalf("decode: %v", err)
		}
		if p.Id != 1 || p.Name == "" {
			t.Fatalf("expected product id=1 with name; got id=%d name=%q", p.Id, p.Name)
		}
	})

	t.Run("update id=1 returns updated product", func(t *testing.T) {
		ts := testutil.NewHTTPServerWithRepo(appsinmem.NewInMemRepo())
		defer ts.Close()

		body := `{"name":"Updated Widget","price":15.25}`
		req, _ := http.NewRequest(http.MethodPut, ts.URL+"/products/1", strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Fatalf("http put: %v", err)
		}
		defer resp.Body.Close()
		if resp.StatusCode != http.StatusOK {
			t.Fatalf("expected 200, got %d", resp.StatusCode)
		}
		var updated appshttp.Product
		if err := json.NewDecoder(resp.Body).Decode(&updated); err != nil {
			t.Fatalf("decode: %v", err)
		}
		if updated.Id != 1 || updated.Name != "Updated Widget" {
			t.Fatalf("unexpected updated product: %+v", updated)
		}
	})

	t.Run("search with short q returns 400", func(t *testing.T) {
		ts := testutil.NewHTTPServerWithRepo(appsinmem.NewInMemRepo())
		defer ts.Close()

		resp, err := http.Get(ts.URL + "/products/search?q=ab")
		if err != nil {
			t.Fatalf("http get: %v", err)
		}
		defer resp.Body.Close()
		if resp.StatusCode != http.StatusBadRequest {
			t.Fatalf("expected 400 for short q; got %d", resp.StatusCode)
		}
	})
}

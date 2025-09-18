package http_inmem_test

import (
    "encoding/json"
    "net/http"
    "testing"

    appshttp "github.com/fightingBald/GoTuto/apps/product-query-svc/adapters/inbound/http"
    appsinmem "github.com/fightingBald/GoTuto/apps/product-query-svc/adapters/outbound/inmem"
    "github.com/fightingBald/GoTuto/internal/testutil"
)

func TestSearchProducts_InMem(t *testing.T) {
    ts := testutil.NewHTTPServerWithRepo(appsinmem.NewInMemRepo())
    defer ts.Close()

	// q must be at least 3 characters; use 'wid' to match 'Blue Widget'
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
}

func TestGetProduct_InMem(t *testing.T) {
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
}

func TestSearchProducts_QueryTooShort(t *testing.T) {
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
}

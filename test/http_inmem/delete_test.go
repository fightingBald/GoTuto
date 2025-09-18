package http_inmem_test

import (
    "net/http"
    "testing"

    appsinmem "github.com/fightingBald/GoTuto/apps/product-query-svc/adapters/outbound/inmem"
    "github.com/fightingBald/GoTuto/internal/testutil"
)

func TestDeleteProduct_InMem(t *testing.T) {
    ts := testutil.NewHTTPServerWithRepo(appsinmem.NewInMemRepo())
    defer ts.Close()

	// delete existing seeded id=1
	req, _ := http.NewRequest(http.MethodDelete, ts.URL+"/products/1", nil)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("http delete: %v", err)
	}
	resp.Body.Close()
	if resp.StatusCode != http.StatusNoContent {
		t.Fatalf("expected 204, got %d", resp.StatusCode)
	}

	// subsequent GET should be 404
	resp2, err := http.Get(ts.URL + "/products/1")
	if err != nil {
		t.Fatalf("http get: %v", err)
	}
	resp2.Body.Close()
	if resp2.StatusCode != http.StatusNotFound {
		t.Fatalf("expected 404 after delete, got %d", resp2.StatusCode)
	}
}

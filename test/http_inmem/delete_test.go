package http_inmem_test

import (
	"net/http"
	"testing"

	appsinmem "github.com/fightingBald/GoTuto/apps/product-query-svc/adapters/outbound/inmem"
	"github.com/fightingBald/GoTuto/internal/testutil"
)

func TestDeleteProduct_InMem(t *testing.T) {
	store := appsinmem.NewInMemRepo()
	ts := testutil.NewHTTPServer(store, store)
	defer ts.Close()

	t.Run("delete id=1 returns 204", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodDelete, ts.URL+"/products/1", nil)
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Fatalf("http delete: %v", err)
		}
		resp.Body.Close()
		if resp.StatusCode != http.StatusNoContent {
			t.Fatalf("expected 204, got %d", resp.StatusCode)
		}
	})

	t.Run("get after delete returns 404", func(t *testing.T) {
		resp2, err := http.Get(ts.URL + "/products/1")
		if err != nil {
			t.Fatalf("http get: %v", err)
		}
		resp2.Body.Close()
		if resp2.StatusCode != http.StatusNotFound {
			t.Fatalf("expected 404 after delete, got %d", resp2.StatusCode)
		}
	})
}

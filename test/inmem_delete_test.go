package test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	appshttp "github.com/fightingBald/GoTuto/apps/product-query-svc/adapters/http"
	appsinmem "github.com/fightingBald/GoTuto/apps/product-query-svc/adapters/inmem"
	appsvc "github.com/fightingBald/GoTuto/apps/product-query-svc/app"
	"github.com/go-chi/chi/v5"
)

func TestDeleteProduct_InMem(t *testing.T) {
	repo := appsinmem.NewInMemRepo()
	svc := appsvc.NewProductService(repo)
	server := appshttp.NewServer(svc)

	r := chi.NewRouter()
	// 仅挂载 OpenAPI 生成的路由（包含 DELETE）
	h := appshttp.HandlerFromMux(server, r)

	ts := httptest.NewServer(h)
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

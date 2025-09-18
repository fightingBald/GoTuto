package testutil

import (
    "net/http"
    "net/http/httptest"

    app "github.com/fightingBald/GoTuto/apps/product-query-svc/app"
    httpadapter "github.com/fightingBald/GoTuto/apps/product-query-svc/adapters/inbound/http"
    "github.com/fightingBald/GoTuto/apps/product-query-svc/ports"
    "github.com/go-chi/chi/v5"
)

// NewHTTPHandlerWithRepo wires repo -> service -> http server handler.
func NewHTTPHandlerWithRepo(repo ports.ProductRepo) http.Handler {
    svc := app.NewProductService(repo)
    server := httpadapter.NewServer(svc)
    r := chi.NewRouter()
    return httpadapter.HandlerFromMux(server, r)
}

// NewHTTPServerWithRepo starts an httptest.Server for convenience.
func NewHTTPServerWithRepo(repo ports.ProductRepo) *httptest.Server {
    h := NewHTTPHandlerWithRepo(repo)
    return httptest.NewServer(h)
}


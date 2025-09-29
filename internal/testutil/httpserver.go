package testutil

import (
	"net/http"
	"net/http/httptest"

	httpadapter "github.com/fightingBald/GoTuto/apps/product-query-svc/adapters/inbound/http"
	productapp "github.com/fightingBald/GoTuto/apps/product-query-svc/application/product"
	userapp "github.com/fightingBald/GoTuto/apps/product-query-svc/application/user"
	"github.com/fightingBald/GoTuto/apps/product-query-svc/ports/outbound"
	"github.com/go-chi/chi/v5"
)

// NewHTTPHandler wires repos -> services -> HTTP handler.
func NewHTTPHandler(productRepo outbound.ProductRepository, userRepo outbound.UserRepository) http.Handler {
	productSvc := productapp.NewService(productRepo)
	userSvc := userapp.NewService(userRepo)
	server := httpadapter.NewServer(productSvc, userSvc)
	r := chi.NewRouter()
	strict := httpadapter.NewStrictHTTPHandler(server, nil)
	return httpadapter.HandlerFromMux(strict, r)
}

// NewHTTPServer starts an httptest.Server for convenience.
func NewHTTPServer(productRepo outbound.ProductRepository, userRepo outbound.UserRepository) *httptest.Server {
	h := NewHTTPHandler(productRepo, userRepo)
	return httptest.NewServer(h)
}

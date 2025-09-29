package testutil

import (
	"net/http"
	"net/http/httptest"

	httpadapter "github.com/fightingBald/GoTuto/apps/product-query-svc/adapters/inbound/http"
	productapp "github.com/fightingBald/GoTuto/apps/product-query-svc/application/product"
	userapp "github.com/fightingBald/GoTuto/apps/product-query-svc/application/user"
	"github.com/fightingBald/GoTuto/apps/product-query-svc/ports/outbound"
)

// NewHTTPHandler wires repos -> services -> HTTP handler.
func NewHTTPHandler(productRepo outbound.ProductRepository, userRepo outbound.UserRepository) http.Handler {
	productSvc := productapp.NewService(productRepo)
	userSvc := userapp.NewService(userRepo)
	server := httpadapter.NewServer(productSvc, userSvc)
	h, err := httpadapter.NewAPIHandler(server, nil)
	if err != nil {
		panic(err)
	}
	return h
}

// NewHTTPServer starts an httptest.Server for convenience.
func NewHTTPServer(productRepo outbound.ProductRepository, userRepo outbound.UserRepository) *httptest.Server {
	h := NewHTTPHandler(productRepo, userRepo)
	return httptest.NewServer(h)
}

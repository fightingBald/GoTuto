package testutil

import (
	"net/http"
	"net/http/httptest"

	httpadapter "github.com/fightingBald/GoTuto/apps/product-query-svc/adapters/inbound/http"
	app "github.com/fightingBald/GoTuto/apps/product-query-svc/app"
	"github.com/fightingBald/GoTuto/apps/product-query-svc/ports"
	"github.com/go-chi/chi/v5"
)

// NewHTTPHandler wires repos -> services -> HTTP handler.
func NewHTTPHandler(productRepo ports.ProductRepo, userRepo ports.UserRepo) http.Handler {
	productSvc := app.NewProductService(productRepo)
	userSvc := app.NewUserService(userRepo)
	server := httpadapter.NewServer(productSvc, userSvc)
	r := chi.NewRouter()
	return httpadapter.HandlerFromMux(server, r)
}

// NewHTTPServer starts an httptest.Server for convenience.
func NewHTTPServer(productRepo ports.ProductRepo, userRepo ports.UserRepo) *httptest.Server {
	h := NewHTTPHandler(productRepo, userRepo)
	return httptest.NewServer(h)
}

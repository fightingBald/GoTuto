package testutil

import (
	"net/http"
	"net/http/httptest"

	httpadapter "github.com/fightingBald/GoTuto/apps/product-query-svc/adapters/inbound/http"
	commentapp "github.com/fightingBald/GoTuto/apps/product-query-svc/application/comment"
	productapp "github.com/fightingBald/GoTuto/apps/product-query-svc/application/product"
	userapp "github.com/fightingBald/GoTuto/apps/product-query-svc/application/user"
	"github.com/fightingBald/GoTuto/apps/product-query-svc/ports/outbound"
)

// NewHTTPHandler wires repos -> services -> HTTP handler.
func NewHTTPHandler(productRepo outbound.ProductRepository, userRepo outbound.UserRepository, commentRepo outbound.CommentRepository) http.Handler {
	productSvc := productapp.NewService(productRepo)
	userSvc := userapp.NewService(userRepo)
	commentSvc := commentapp.NewService(commentRepo, productRepo, userRepo)
	server := httpadapter.NewServer(productSvc, userSvc, commentSvc)
	h, err := httpadapter.NewAPIHandler(server, nil)
	if err != nil {
		panic(err)
	}
	return h
}

// NewHTTPServer starts an httptest.Server for convenience.
func NewHTTPServer(productRepo outbound.ProductRepository, userRepo outbound.UserRepository, commentRepo outbound.CommentRepository) *httptest.Server {
	h := NewHTTPHandler(productRepo, userRepo, commentRepo)
	return httptest.NewServer(h)
}

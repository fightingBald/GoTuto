package httpadapter

import (
	"net/http"

	"github.com/fightingBald/GoTuto/apps/product-query-svc/ports"
)

// Server wires product and user services to HTTP handlers generated from OpenAPI.
type Server struct {
	products ports.ProductService
	users    ports.UserService
}

func NewServer(products ports.ProductService, users ports.UserService) *Server {
	return &Server{products: products, users: users}
}

// Health 健康检查
func (s *Server) Health(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("ok"))
}

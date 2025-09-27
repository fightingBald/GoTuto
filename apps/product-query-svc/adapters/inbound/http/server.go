package httpadapter

import (
	"net/http"

	"github.com/fightingBald/GoTuto/apps/product-query-svc/ports/inbound"
)

// Server wires product and user services to HTTP handlers generated from OpenAPI.
type Server struct {
	products inbound.ProductUseCases
	users    inbound.UserQueries
}

func NewServer(products inbound.ProductUseCases, users inbound.UserQueries) *Server {
	return &Server{products: products, users: users}
}

// Health 健康检查
func (s *Server) Health(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("ok"))
}

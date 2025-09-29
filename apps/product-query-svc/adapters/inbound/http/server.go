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

var _ StrictServerInterface = (*Server)(nil)

// NewStrictHTTPHandler wraps the server with oapi-codegen strict adapter using
// JSON error envelopes for request/response failures.
func NewStrictHTTPHandler(server *Server, middlewares []StrictMiddlewareFunc) ServerInterface {
	options := StrictHTTPServerOptions{
		RequestErrorHandlerFunc: func(w http.ResponseWriter, r *http.Request, err error) {
			writeError(w, http.StatusBadRequest, "INVALID_REQUEST", err.Error())
		},
		ResponseErrorHandlerFunc: func(w http.ResponseWriter, r *http.Request, err error) {
			writeError(w, http.StatusInternalServerError, "INTERNAL", err.Error())
		},
	}
	return NewStrictHandlerWithOptions(server, middlewares, options)
}

// Health 健康检查
func (s *Server) Health(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("ok"))
}

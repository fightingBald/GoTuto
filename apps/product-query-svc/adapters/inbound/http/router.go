package httpadapter

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	nethttpmiddleware "github.com/oapi-codegen/nethttp-middleware"
)

// NewAPIHandler returns a chi-backed handler wired with the strict server and
// OpenAPI request validator.
func NewAPIHandler(server *Server, strictMiddlewares []StrictMiddlewareFunc, middlewares ...func(http.Handler) http.Handler) (http.Handler, error) {
	swagger, err := GetSwagger()
	if err != nil {
		return nil, fmt.Errorf("load swagger spec: %w", err)
	}

	r := chi.NewRouter()
	for _, mw := range middlewares {
		r.Use(mw)
	}
	r.Use(nethttpmiddleware.OapiRequestValidator(swagger))

	strict := NewStrictHTTPHandler(server, strictMiddlewares)
	return HandlerFromMux(strict, r), nil
}

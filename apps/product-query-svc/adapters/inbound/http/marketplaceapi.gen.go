// Package httpadapter provides primitives to interact with the openapi HTTP API.
// Minimal shim to satisfy our handlers and router without full codegen.
package httpadapter

import (
	"github.com/go-chi/chi/v5"
	"net/http"
	"strconv"
)

// Error defines model for Error.
type Error struct {
	Code    string `json:"code"`
	Details *[]struct {
		Field  *string `json:"field,omitempty"`
		Reason *string `json:"reason,omitempty"`
	} `json:"details,omitempty"`
	Message string `json:"message"`
}

// Product defines model for Product.
type Product struct {
	Id    int64   `json:"id"`
	Name  string  `json:"name"`
	Price float32 `json:"price"`
}

// ProductList defines model for ProductList.
type ProductList struct {
	Items    []Product `json:"items"`
	Page     int       `json:"page"`
	PageSize int       `json:"pageSize"`
	Total    int       `json:"total"`
}

// ErrorResponse alias.
type ErrorResponse = Error

// SearchProductsParams defines parameters for SearchProducts.
type SearchProductsParams struct {
	Q        *string `form:"q,omitempty" json:"q,omitempty"`
	Page     *int    `form:"page,omitempty" json:"page,omitempty"`
	PageSize *int    `form:"pageSize,omitempty" json:"pageSize,omitempty"`
}

// ServerInterface represents all server handlers.
type ServerInterface interface {
	SearchProducts(w http.ResponseWriter, r *http.Request, params SearchProductsParams)
	DeleteProductByID(w http.ResponseWriter, r *http.Request, id int64)
	GetProductByID(w http.ResponseWriter, r *http.Request, id int64)
}

// Unimplemented server returns 501 for each endpoint.
type Unimplemented struct{}

func (_ Unimplemented) SearchProducts(w http.ResponseWriter, r *http.Request, params SearchProductsParams) {
	w.WriteHeader(http.StatusNotImplemented)
}
func (_ Unimplemented) DeleteProductByID(w http.ResponseWriter, r *http.Request, id int64) {
	w.WriteHeader(http.StatusNotImplemented)
}
func (_ Unimplemented) GetProductByID(w http.ResponseWriter, r *http.Request, id int64) {
	w.WriteHeader(http.StatusNotImplemented)
}

// MiddlewareFunc is compatible with chi middlewares.
type MiddlewareFunc func(http.Handler) http.Handler

// ChiServerOptions configures the router helpers.
type ChiServerOptions struct {
	BaseURL     string
	BaseRouter  chi.Router
	Middlewares []MiddlewareFunc
}

// Handler creates http.Handler with routing matching OpenAPI spec (minimal).
func Handler(si ServerInterface) http.Handler { return HandlerWithOptions(si, ChiServerOptions{}) }

// HandlerFromMux binds onto an existing mux.
func HandlerFromMux(si ServerInterface, r chi.Router) http.Handler {
	return HandlerFromMuxWithOptions(si, r, ChiServerOptions{})
}

// HandlerWithOptions mounts routes on a new mux using options.
func HandlerWithOptions(si ServerInterface, options ChiServerOptions) http.Handler {
	var r chi.Router
	if options.BaseRouter == nil {
		r = chi.NewRouter()
	} else {
		r = options.BaseRouter
	}
	for _, m := range options.Middlewares {
		if m != nil {
			r.Use(m)
		}
	}

	r.Group(func(r chi.Router) {
		r.Get("/products/search", func(w http.ResponseWriter, req *http.Request) {
			q := req.URL.Query().Get("q")
			var qPtr *string
			if q != "" {
				qPtr = &q
			}
			var pagePtr, sizePtr *int
			if v := req.URL.Query().Get("page"); v != "" {
				if n, err := strconv.Atoi(v); err == nil {
					pagePtr = &n
				}
			}
			if v := req.URL.Query().Get("pageSize"); v != "" {
				if n, err := strconv.Atoi(v); err == nil {
					sizePtr = &n
				}
			}
			params := SearchProductsParams{Q: qPtr, Page: pagePtr, PageSize: sizePtr}
			si.SearchProducts(w, req, params)
		})
		r.Get("/products/{id}", func(w http.ResponseWriter, req *http.Request) {
			idStr := chi.URLParam(req, "id")
			if id, err := strconv.ParseInt(idStr, 10, 64); err == nil {
				si.GetProductByID(w, req, id)
			} else {
				w.WriteHeader(http.StatusBadRequest)
			}
		})
		r.Delete("/products/{id}", func(w http.ResponseWriter, req *http.Request) {
			idStr := chi.URLParam(req, "id")
			if id, err := strconv.ParseInt(idStr, 10, 64); err == nil {
				si.DeleteProductByID(w, req, id)
			} else {
				w.WriteHeader(http.StatusBadRequest)
			}
		})
	})

	return r
}

// HandlerFromMuxWithOptions mounts routes onto the provided mux with options.
func HandlerFromMuxWithOptions(si ServerInterface, r chi.Router, options ChiServerOptions) http.Handler {
	for _, m := range options.Middlewares {
		if m != nil {
			r.Use(m)
		}
	}
	r.Group(func(r chi.Router) {
		r.Get("/products/search", func(w http.ResponseWriter, req *http.Request) {
			q := req.URL.Query().Get("q")
			var qPtr *string
			if q != "" {
				qPtr = &q
			}
			var pagePtr, sizePtr *int
			if v := req.URL.Query().Get("page"); v != "" {
				if n, err := strconv.Atoi(v); err == nil {
					pagePtr = &n
				}
			}
			if v := req.URL.Query().Get("pageSize"); v != "" {
				if n, err := strconv.Atoi(v); err == nil {
					sizePtr = &n
				}
			}
			params := SearchProductsParams{Q: qPtr, Page: pagePtr, PageSize: sizePtr}
			si.SearchProducts(w, req, params)
		})
		r.Get("/products/{id}", func(w http.ResponseWriter, req *http.Request) {
			idStr := chi.URLParam(req, "id")
			if id, err := strconv.ParseInt(idStr, 10, 64); err == nil {
				si.GetProductByID(w, req, id)
			} else {
				w.WriteHeader(http.StatusBadRequest)
			}
		})
		r.Delete("/products/{id}", func(w http.ResponseWriter, req *http.Request) {
			idStr := chi.URLParam(req, "id")
			if id, err := strconv.ParseInt(idStr, 10, 64); err == nil {
				si.DeleteProductByID(w, req, id)
			} else {
				w.WriteHeader(http.StatusBadRequest)
			}
		})
	})
	return r
}

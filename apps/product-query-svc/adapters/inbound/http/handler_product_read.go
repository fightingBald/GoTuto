package httpadapter

import (
    "net/http"
    "errors"

    "github.com/fightingBald/GoTuto/apps/product-query-svc/domain"
    "github.com/fightingBald/GoTuto/apps/product-query-svc/ports"
)

type Server struct{ svc ports.ProductService }

func NewServer(s ports.ProductService) *Server { return &Server{svc: s} }

func (s *Server) GetProductByID(w http.ResponseWriter, r *http.Request, id int64) {
    p, err := s.svc.GetProduct(r.Context(), id)
    if err != nil {
        if errors.Is(err, domain.ErrNotFound) {
            writeError(w, http.StatusNotFound, "NOT_FOUND", err.Error())
        } else {
            writeError(w, http.StatusInternalServerError, "INTERNAL", err.Error())
        }
        return
    }
	out := Product{Id: p.ID, Name: p.Name, Price: float32(p.Price) / 100.0}
	writeJSON(w, http.StatusOK, out)
}

func (s *Server) SearchProducts(w http.ResponseWriter, r *http.Request, params SearchProductsParams) {
	q := ""
	if params.Q != nil {
		q = *params.Q
	}
	// Enforce OpenAPI minLength:3 for q when provided
	if q != "" && len(q) < 3 {
		writeError(w, http.StatusBadRequest, "INVALID_QUERY", "q must be at least 3 characters if provided")
		return
	}
	page := 1
	if params.Page != nil {
		page = *params.Page
	}
	pageSize := 20
	if params.PageSize != nil {
		pageSize = *params.PageSize
	}
	items, total, err := s.svc.SearchProducts(r.Context(), q, page, pageSize)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "INTERNAL", err.Error())
		return
	}
	var out []Product
	for _, it := range items {
		out = append(out, Product{Id: it.ID, Name: it.Name, Price: float32(it.Price) / 100.0})
	}
	resp := ProductList{Items: out, Page: page, PageSize: pageSize, Total: total}
	writeJSON(w, http.StatusOK, resp)
}

// Healthz 健康检查
func (s *Server) Health(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("ok"))
}

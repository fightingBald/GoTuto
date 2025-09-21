package httpadapter

import "net/http"

func (s *Server) GetProductByID(w http.ResponseWriter, r *http.Request, id int64) {
	if id <= 0 {
		writeError(w, http.StatusBadRequest, "INVALID_ID", "id must be a positive integer")
		return
	}
	p, err := s.products.GetProduct(r.Context(), id)
	if err != nil {
		writeDomainError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, presentProduct(p))
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
	items, total, err := s.products.SearchProducts(r.Context(), q, page, pageSize)
	if err != nil {
		writeDomainError(w, err)
		return
	}
	resp := ProductList{Items: presentProducts(items), Page: page, PageSize: pageSize, Total: total}
	writeJSON(w, http.StatusOK, resp)
}

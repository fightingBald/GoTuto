package httpadapter

import (
	"encoding/json"
	"net/http"

	"github.com/fightingBald/GoTuto/apps/product-query-svc/domain"
)

// DeleteProductByID implements OpenAPI operation: DELETE /products/{id}
func (s *Server) DeleteProductByID(w http.ResponseWriter, r *http.Request, id int64) {
	if id <= 0 {
		writeError(w, http.StatusBadRequest, "INVALID_ID", "id must be a positive integer")
		return
	}
	if err := s.products.Remove(r.Context(), id); err != nil {
		writeDomainError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// CreateProduct implements POST /products
func (s *Server) CreateProduct(w http.ResponseWriter, r *http.Request) {
	var in CreateProductJSONBody
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		writeError(w, http.StatusBadRequest, "INVALID_JSON", "invalid request body")
		return
	}
	cents := amountToCents(in.Price)
	p, err := domain.NewProduct(in.Name, cents, nil)
	if err != nil {
		writeDomainError(w, err)
		return
	}
	id, err := s.products.Create(r.Context(), p)
	if err != nil {
		writeDomainError(w, err)
		return
	}
	created := *p
	created.ID = id
	writeJSON(w, http.StatusCreated, presentProduct(&created))
}

// UpdateProduct implements OpenAPI operation: PUT /products/{id}
func (s *Server) UpdateProduct(w http.ResponseWriter, r *http.Request, id int64) {
	if id <= 0 {
		writeError(w, http.StatusBadRequest, "INVALID_ID", "id must be a positive integer")
		return
	}
	var in UpdateProductJSONBody
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		writeError(w, http.StatusBadRequest, "INVALID_JSON", "invalid request body")
		return
	}
	cents := amountToCents(in.Price)
	p, err := domain.NewProduct(in.Name, cents, nil)
	if err != nil {
		writeDomainError(w, err)
		return
	}
	p.ID = id
	updated, err := s.products.Update(r.Context(), p)
	if err != nil {
		writeDomainError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, presentProduct(updated))
}

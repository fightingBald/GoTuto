package httpadapter

import (
	"encoding/json"
	"errors"
	"math"
	"net/http"

	"github.com/fightingBald/GoTuto/apps/product-query-svc/domain"
)

// DeleteProductByID implements OpenAPI operation: DELETE /products/{id}
func (s *Server) DeleteProductByID(w http.ResponseWriter, r *http.Request, id int64) {
	if id <= 0 {
		writeError(w, http.StatusBadRequest, "INVALID_ID", "id must be a positive integer")
		return
	}
	if err := s.svc.DeleteProduct(r.Context(), id); err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			writeError(w, http.StatusNotFound, "NOT_FOUND", err.Error())
			return
		}
		writeError(w, http.StatusInternalServerError, "INTERNAL", err.Error())
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
	cents := int64(math.Round(float64(in.Price) * 100.0))
	p, err := domain.NewProduct(in.Name, cents, nil)
	if err != nil {
		writeError(w, http.StatusBadRequest, "VALIDATION", err.Error())
		return
	}
	id, err := s.svc.CreateProduct(r.Context(), p)
	if err != nil {
		if errors.Is(err, domain.ErrValidation) {
			writeError(w, http.StatusBadRequest, "VALIDATION", err.Error())
			return
		}
		writeError(w, http.StatusInternalServerError, "INTERNAL", err.Error())
		return
	}
	out := Product{Id: id, Name: p.Name, Price: float32(cents) / 100.0}
	writeJSON(w, http.StatusCreated, out)
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
	cents := int64(math.Round(float64(in.Price) * 100.0))
	p, err := domain.NewProduct(in.Name, cents, nil)
	if err != nil {
		writeError(w, http.StatusBadRequest, "VALIDATION", err.Error())
		return
	}
	p.ID = id
	updated, err := s.svc.UpdateProduct(r.Context(), p)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrValidation):
			writeError(w, http.StatusBadRequest, "VALIDATION", err.Error())
		case errors.Is(err, domain.ErrNotFound):
			writeError(w, http.StatusNotFound, "NOT_FOUND", err.Error())
		default:
			writeError(w, http.StatusInternalServerError, "INTERNAL", err.Error())
		}
		return
	}
	out := Product{Id: updated.ID, Name: updated.Name, Price: float32(updated.Price) / 100.0}
	writeJSON(w, http.StatusOK, out)
}

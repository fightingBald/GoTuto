package httpadapter

import (
	"net/http"
)

// DeleteProductByID implements OpenAPI operation: DELETE /products/{id}
func (s *Server) DeleteProductByID(w http.ResponseWriter, r *http.Request, id int64) {
	if id <= 0 {
		writeError(w, http.StatusBadRequest, "INVALID_ID", "id must be a positive integer")
		return
	}
	if err := s.svc.DeleteProduct(r.Context(), id); err != nil {
		if err.Error() == "not found" {
			writeError(w, http.StatusNotFound, "NOT_FOUND", err.Error())
			return
		}
		writeError(w, http.StatusInternalServerError, "INTERNAL", err.Error())
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

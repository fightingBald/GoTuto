package httpadapter

import "net/http"

func (s *Server) GetUserByID(w http.ResponseWriter, r *http.Request, id int64) {
	if id <= 0 {
		writeError(w, http.StatusBadRequest, "INVALID_ID", "id must be a positive integer")
		return
	}
	u, err := s.users.FetchByID(r.Context(), id)
	if err != nil {
		writeDomainError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, presentUser(u))
}

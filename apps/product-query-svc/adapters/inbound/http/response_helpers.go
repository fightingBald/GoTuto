package httpadapter

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/fightingBald/GoTuto/apps/product-query-svc/domain"
)

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

type errorBody struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Details *[]struct {
		Field  *string `json:"field,omitempty"`
		Reason *string `json:"reason,omitempty"`
	} `json:"details,omitempty"`
}

func writeError(w http.ResponseWriter, status int, code, message string) {
	writeJSON(w, status, errorBody{Code: code, Message: message})
}

func writeDomainError(w http.ResponseWriter, err error) {
	status, code := classifyDomainError(err)
	writeError(w, status, code, domainErrorMessage(status, err))
}

func classifyDomainError(err error) (int, string) {
	switch {
	case errors.Is(err, domain.ErrValidation):
		return http.StatusBadRequest, "VALIDATION"
	case errors.Is(err, domain.ErrNotFound):
		return http.StatusNotFound, "NOT_FOUND"
	default:
		return http.StatusInternalServerError, "INTERNAL"
	}
}

func domainErrorMessage(status int, err error) string {
	if status == http.StatusInternalServerError {
		return http.StatusText(status)
	}
	if errors.Is(err, domain.ErrValidation) {
		if split := strings.SplitN(err.Error(), "\n", 2); len(split) == 2 {
			return split[1]
		}
	}
	return err.Error()
}

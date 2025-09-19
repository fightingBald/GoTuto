package httpadapter

import (
    "encoding/json"
    "net/http"
)

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

type errorBody struct {
    Code    string       `json:"code"`
    Message string       `json:"message"`
    Details *[]struct{
        Field  *string `json:"field,omitempty"`
        Reason *string `json:"reason,omitempty"`
    } `json:"details,omitempty"`
}

func writeError(w http.ResponseWriter, status int, code, message string) {
    writeJSON(w, status, errorBody{Code: code, Message: message})
}

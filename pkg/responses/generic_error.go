package responses

import (
	"encoding/json"
	"net/http"
)

type GenericError struct {
	Error string `json:"error"`
}

func NewGenericError(msg string) *GenericError {
	return &GenericError{Error: msg}
}

func (e *GenericError) Encode(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Accept", "application/json")

	_ = json.NewEncoder(w).Encode(e)
}

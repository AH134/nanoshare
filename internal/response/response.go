package response

import (
	"encoding/json"
	"net/http"
)

type Envelope[T any] struct {
	Success bool      `json:"success"`
	Data    *T        `json:"data,omitzero"`
	Error   *APIError `json:"error,omitzero"`
}

type APIError struct {
	Code    string            `json:"code"`
	Message string            `json:"message"`
	Fields  map[string]string `json:"fields,omitzero"`
}

func writeJSON[T any](w http.ResponseWriter, statusCode int, v T) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	_ = json.NewEncoder(w).Encode(v)
}

func Success[T any](w http.ResponseWriter, statusCode int, data T) {
	writeJSON(w, statusCode, Envelope[T]{
		Success: true,
		Data:    new(T(data)),
	})
}

func Error(w http.ResponseWriter, statusCode int, apiErr APIError) {
	writeJSON(w, statusCode, Envelope[any]{
		Success: false,
		Error:   new(APIError(apiErr)),
	})
}

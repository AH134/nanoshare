package response

import (
	"encoding/json"
	"log/slog"
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
		Data:    new(data),
	})
}

func Error(w http.ResponseWriter, statusCode int, apiErr APIError) {
	writeJSON(w, statusCode, Envelope[any]{
		Success: false,
		Error:   new(apiErr),
	})
}

func InternalError(w http.ResponseWriter, logger *slog.Logger, msg string, err error) {
	logger.Error(msg, "err", err)
	Error(w, http.StatusInternalServerError, APIError{
		Code:    "INTERNAL_ERROR",
		Message: "An unexpected error occurred. Please try again later.",
	})
}

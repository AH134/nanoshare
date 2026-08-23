package handler

import (
	"net/http"

	"github.com/AH134/nanoshare/internal/response"
)

type HealthResponse struct {
	Status string `json:"status"`
}

func HealthCheck(w http.ResponseWriter, r *http.Request) {
	response.Success(w, http.StatusOK, HealthResponse{Status: "ok"})
}

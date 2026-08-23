package handler

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"time"

	"github.com/AH134/nanoshare/internal/database"
	"github.com/AH134/nanoshare/internal/response"
	"github.com/AH134/nanoshare/internal/session"
	"github.com/alexedwards/scs/v2"
	"golang.org/x/crypto/bcrypt"
)

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type AuthResponse struct {
	ID        int64     `json:"id"`
	Username  string    `json:"username"`
	CreatedAt time.Time `json:"createdAt"`
}

type AuthHandler struct {
	users          *database.UserRepository
	sessionManager *scs.SessionManager
	logger         *slog.Logger
}

func NewAuthHandler(users *database.UserRepository, sessionManager *scs.SessionManager, logger *slog.Logger) *AuthHandler {
	return &AuthHandler{
		users:          users,
		sessionManager: sessionManager,
		logger:         logger.With("handler", "auth"),
	}
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Warn("failed to decode login request body", "error", err)
		response.Error(w, http.StatusBadRequest, response.APIError{
			Code:    "BAD_REQUEST",
			Message: "Invalid JSON payload.",
		})
		return
	}

	user, err := h.users.GetByUsername(req.Username)
	if errors.Is(err, database.ErrNotFound) {
		h.logger.Warn("user not found", "username", req.Username, "error", err)
		response.Error(w, http.StatusUnauthorized, response.APIError{
			Code:    "INVALID_CREDENTIALS",
			Message: "Invalid username or password.",
		})
		return
	}
	if err != nil {
		response.InternalError(w, h.logger, "failed to get user", err)
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		h.logger.Warn("password mismatch", "username", req.Username, "error", err)
		response.Error(w, http.StatusUnauthorized, response.APIError{
			Code:    "INVALID_CREDENTIALS",
			Message: "Invalid username or password.",
		})
		return
	}

	if err := h.sessionManager.RenewToken(r.Context()); err != nil {
		response.InternalError(w, h.logger, "failed to renew session token", err)
		return
	}

	h.sessionManager.Put(r.Context(), session.DefaultUserIDKey, user.ID)

	resp := AuthResponse{
		ID:        user.ID,
		Username:  user.Username,
		CreatedAt: user.CreatedAt,
	}
	response.Success(w, http.StatusOK, resp)
}

func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	if err := h.sessionManager.Destroy(r.Context()); err != nil {
		response.InternalError(w, h.logger, "failed to destroy session token", err)
		return
	}

	response.Success(w, http.StatusOK, struct{}{})
}

func (h *AuthHandler) Me(w http.ResponseWriter, r *http.Request) {
	userID := h.sessionManager.GetInt64(r.Context(), session.DefaultUserIDKey)

	user, err := h.users.GetByID(userID)
	if errors.Is(err, database.ErrNotFound) {
		h.logger.Warn("user not found", "user_id", userID, "error", err)
		response.Error(w, http.StatusNotFound, response.APIError{
			Code:    "USER_NOT_FOUND",
			Message: "User does not exist.",
		})
		return
	}
	if err != nil {
		response.InternalError(w, h.logger, "failed to get user", err)
		return
	}

	resp := AuthResponse{ID: user.ID, Username: user.Username, CreatedAt: user.CreatedAt}
	response.Success(w, http.StatusOK, resp)
}

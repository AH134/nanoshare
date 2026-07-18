package handler

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/AH134/nanoshare/internal/database"
	"github.com/AH134/nanoshare/internal/session"
	"github.com/alexedwards/scs/v2"
	"golang.org/x/crypto/bcrypt"
)

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

const (
	errInvalidPayload = "Failed to process username or password"
	errInvalidAuth    = "Invalid username or password"
)

type AuthHandler struct {
	users          *database.UserRepository
	sessionManager *scs.SessionManager
}

func NewAuthHandler(users *database.UserRepository, sessionManager *scs.SessionManager) *AuthHandler {
	return &AuthHandler{
		users:          users,
		sessionManager: sessionManager,
	}
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, errInvalidPayload, http.StatusBadRequest)
		return
	}

	user, err := h.users.GetByUsername(req.Username)
	if err != nil {
		log.Printf("login: lookup failed for %q: %v", req.Username, err)
		http.Error(w, errInvalidAuth, http.StatusUnauthorized)
		return
	}

	reqPassword := []byte(req.Password)
	hashedPassword := []byte(user.PasswordHash)
	err = bcrypt.CompareHashAndPassword(hashedPassword, reqPassword)
	if err != nil {
		http.Error(w, errInvalidAuth, http.StatusUnauthorized)
		return
	}

	err = h.sessionManager.RenewToken(r.Context())
	if err != nil {
		log.Printf("login: renew token failed: %v", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	h.sessionManager.Put(r.Context(), session.UserIDKey, user.Id)
	w.WriteHeader(http.StatusOK)
}

func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	if err := h.sessionManager.Destroy(r.Context()); err != nil {
		log.Printf("logout: failed to destroy session data: %v", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *AuthHandler) Me(w http.ResponseWriter, r *http.Request) {
	userId := h.sessionManager.GetInt(r.Context(), session.UserIDKey)

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(map[string]int{"userId": userId}); err != nil {
		log.Printf("me: failed to encode response: %v", err)
	}

}

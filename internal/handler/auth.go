package handler

import (
	"encoding/json"
	"log"
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
	Id        int       `json:"id"`
	Username  string    `json:"username"`
	CreatedAt time.Time `json:"createdAt"`
}

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
		log.Printf("login: failed to decode json: %v", err)
		response.Error(w, http.StatusBadRequest, response.APIError{
			Code:    "BAD_REQUEST",
			Message: "Invalid JSON payload.",
		})
		return
	}

	user, err := h.users.GetByUsername(req.Username)
	if err != nil {
		log.Printf("login: failed to get user: %v", err)
		response.Error(w, http.StatusUnauthorized, response.APIError{
			Code:    "INVALID_CREDENTIALS",
			Message: "Invalid username or password.",
		})
		return
	}

	reqPassword := []byte(req.Password)
	hashedPassword := []byte(user.PasswordHash)
	err = bcrypt.CompareHashAndPassword(hashedPassword, reqPassword)
	if err != nil {
		log.Printf("login: failed to verify password: %v", err)
		response.Error(w, http.StatusUnauthorized, response.APIError{
			Code:    "INVALID_CREDENTIALS",
			Message: "Invalid username or password.",
		})
		return
	}

	err = h.sessionManager.RenewToken(r.Context())
	if err != nil {
		log.Printf("login: failed to renew session data: %v", err)
		response.Error(w, http.StatusInternalServerError, response.APIError{
			Code:    "INTERNAL_ERROR",
			Message: "An unexpected error occured. Please try again later.",
		})
		return
	}

	h.sessionManager.Put(r.Context(), session.DefaultUserIDKey, user.Id)

	resp := AuthResponse{
		Id:        user.Id,
		Username:  user.Username,
		CreatedAt: user.CreateAt,
	}
	response.Success(w, http.StatusOK, resp)
}

func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	if err := h.sessionManager.Destroy(r.Context()); err != nil {
		log.Printf("logout: failed to destroy session data: %v", err)
		response.Error(w, http.StatusInternalServerError, response.APIError{
			Code:    "INTERNAL_ERROR",
			Message: "An unexpected error occured. Please try again later.",
		})
		return
	}

	response.Success(w, http.StatusOK, struct{}{})
}

func (h *AuthHandler) Me(w http.ResponseWriter, r *http.Request) {
	userId := h.sessionManager.GetInt(r.Context(), session.DefaultUserIDKey)

	user, err := h.users.GetById(userId)
	if err != nil {
		http.Error(w, "user not found", http.StatusInternalServerError)
		response.Error(w, http.StatusNotFound, response.APIError{
			Code:    "USER_NOT_FOUND",
			Message: "User does not exist.",
		})
		return
	}

	resp := AuthResponse{Id: user.Id, Username: user.Username, CreatedAt: user.CreateAt}

	response.Success(w, http.StatusOK, resp)
}

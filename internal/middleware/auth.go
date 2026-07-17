package middleware

import (
	"log"
	"net/http"

	"github.com/AH134/nanoshare/internal/session"
	"github.com/alexedwards/scs/v2"
)

type AuthMiddleware struct {
	sessionManager *scs.SessionManager
}

func NewAuthMiddleware(sessionManager *scs.SessionManager) *AuthMiddleware {
	return &AuthMiddleware{
		sessionManager: sessionManager,
	}
}

func (m *AuthMiddleware) RequireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("AuthHandler. Request received: (%s) %s", r.Method, r.URL.Path)
		if !m.sessionManager.Exists(r.Context(), session.UserIdKey) {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
}

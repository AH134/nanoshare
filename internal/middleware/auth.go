package middleware

import (
	"net/http"

	"github.com/AH134/nanoshare/internal/session"
	"github.com/alexedwards/scs/v2"
)

func RequireAuth(sessionManager *scs.SessionManager) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if !sessionManager.Exists(r.Context(), session.UserIdKey) {
				http.Error(w, "unauthorized", http.StatusUnauthorized)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

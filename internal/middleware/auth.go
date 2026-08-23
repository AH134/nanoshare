package middleware

import (
	"net/http"

	"github.com/AH134/nanoshare/internal/response"
	"github.com/AH134/nanoshare/internal/session"
	"github.com/alexedwards/scs/v2"
)

func RequireAuth(sessionManager *scs.SessionManager) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if !sessionManager.Exists(r.Context(), session.DefaultUserIDKey) {
				response.Error(w, http.StatusUnauthorized, response.APIError{
					Code:    "UNAUTHORIZED",
					Message: "You must be logged in to access this resource.",
				})
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

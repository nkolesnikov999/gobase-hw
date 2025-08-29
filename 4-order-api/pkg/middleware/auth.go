package middleware

import (
	"api/orders/configs"
	"api/orders/pkg/jwt"
	"net/http"
	"strings"
)

// NewIsAuthed returns middleware that validates Bearer JWT and rejects unauthorized requests.
func NewIsAuthed(cfg *configs.Config) func(http.Handler) http.Handler {
	validator := jwt.NewJWT(cfg.Auth.Secret)
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
				http.Error(w, "unauthorized", http.StatusUnauthorized)
				return
			}
			token := strings.TrimPrefix(authHeader, "Bearer ")
			subject, err := validator.Parse(strings.TrimSpace(token))
			if err != nil || subject == "" {
				http.Error(w, "unauthorized", http.StatusUnauthorized)
				return
			}
			// Token is valid; proceed. We could set subject in context if needed later.
			next.ServeHTTP(w, r)
		})
	}
}

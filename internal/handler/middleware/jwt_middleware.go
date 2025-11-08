package middleware

import (
	"context"
	"net/http"
	"strings"
	"qlsvgo/internal/infrastructure/jwt"
)

type contextKey string

const userContextKey = contextKey("user")

func JWTMiddleware(secret string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			header := r.Header.Get("Authorization")
			if header == "" || !strings.HasPrefix(header, "Bearer ") {
				http.Error(w, "missing or invalid token", http.StatusUnauthorized)
				return
			}
			tokenStr := strings.TrimPrefix(header, "Bearer ")
			claims, err := jwt.ParseJWT(tokenStr, secret)
			if err != nil {
				http.Error(w, "invalid token", http.StatusUnauthorized)
				return
			}
			ctx := context.WithValue(r.Context(), userContextKey, claims)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func GetUserClaims(r *http.Request) *jwt.Claims {
	claims, _ := r.Context().Value(userContextKey).(*jwt.Claims)
	return claims
} 
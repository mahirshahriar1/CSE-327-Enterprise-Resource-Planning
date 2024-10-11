package middleware

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"erp/utils"
)

type contextKey string

const UserEmail contextKey = "email"

// JWTAuth middleware to validate JWT and extract user information
func JWTAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Authorization header missing", http.StatusUnauthorized)
			return
		}

		// Expecting "Bearer <token>"
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader {
			http.Error(w, "Bearer token missing", http.StatusUnauthorized)
			return
		}

		claims, err := utils.ValidateJWT(tokenString)
		if err != nil {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}
		// log.Println("Claims:", claims)
		// Get userID from token claims
		email, ok := claims["email"].(string) // Extract email
		if !ok {
			http.Error(w, "Invalid token claims", http.StatusUnauthorized)
			return
		}

		// Add the userID to the context
		ctx := context.WithValue(r.Context(), UserEmail, email)

		// Pass the request with updated context to the next handler
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// GetUserEmailFromContext extracts the email from the request context
func GetUserEmailFromContext(ctx context.Context) (string, error) {
	email, ok := ctx.Value(UserEmail).(string)
	if !ok {
		return "", fmt.Errorf("email not found in context")
	}
	return email, nil
}

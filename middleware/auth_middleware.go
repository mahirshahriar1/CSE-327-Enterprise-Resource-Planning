package middleware

import (
	"erp/utils"
	"log"
	"net/http"
	"strings"
)

func JWTAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get the Authorization header
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Missing authorization token", http.StatusUnauthorized)
			return
		}

		// The token is typically sent as "Bearer <token>", so we need to split it
		bearerToken := strings.Split(authHeader, " ")
		if len(bearerToken) != 2 {
			http.Error(w, "Invalid authorization token format", http.StatusUnauthorized)
			return
		}

		// Validate the JWT token
		tokenString := bearerToken[1]
		claims, err := utils.ValidateJWT(tokenString)
		if err != nil {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		// Log the token's email claim for debugging
		log.Printf("Authenticated user with email: %v", claims["email"])

		// Proceed with the request
		next.ServeHTTP(w, r)
	})
}

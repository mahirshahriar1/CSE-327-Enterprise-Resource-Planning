package utils

import (
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
)

var jwtSecret = []byte("your-secret-key")

// GenerateJWT generates a new JWT token with userID and email
func GenerateJWT(email string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email": email,
		"exp":   time.Now().Add(time.Hour * 24).Unix(), // Token expires in 24 hours
	})
	return token.SignedString(jwtSecret)
}

// ValidateJWT validates a JWT token and extracts the claims
func ValidateJWT(tokenString string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Verify token signing method etc.
		return jwtSecret, nil
	})

	if err != nil || !token.Valid {
		return nil, err
	}

	// Extract and return the claims
	if claims, ok := token.Claims.(jwt.MapClaims); ok {

		return claims, nil
	}

	return nil, fmt.Errorf("invalid token claims")
}

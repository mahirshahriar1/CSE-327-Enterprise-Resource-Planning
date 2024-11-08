// // Description: This file implements the tests for the dashboard endpoint.
package dashboard_test

// import (
// 	"erp/middleware"
// 	"erp/utils"
// 	"net/http"
// 	"net/http/httptest"
// 	"testing"

// 	"github.com/stretchr/testify/assert"
// )

// // TestDashboard_ValidToken tests the Dashboard with a valid JWT token
// func TestDashboard_ValidToken(t *testing.T) {
// 	// Generate a valid JWT token for testing
// 	tokenString, err := utils.GenerateJWT("test@example.com")
// 	assert.NoError(t, err)

// 	// Create a request to the /dashboard endpoint
// 	req := httptest.NewRequest("GET", "/dashboard", nil)
// 	req.Header.Set("Authorization", "Bearer "+tokenString) // Set the Authorization header
// 	rr := httptest.NewRecorder()

// 	// Define a simple handler to simulate the dashboard functionality
// 	handler := middleware.JWTAuth(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		// Extract email from the request context
// 		email, err := middleware.GetUserEmailFromContext(r.Context())
// 		if err != nil {
// 			http.Error(w, "Failed to get email", http.StatusInternalServerError)
// 			return
// 		}

// 		// Respond with the authenticated email
// 		w.WriteHeader(http.StatusOK)
// 		w.Write([]byte("Authenticated user: " + email))
// 	}))

// 	// Serve the request
// 	handler.ServeHTTP(rr, req)

// 	// Assert the response status code and body
// 	assert.Equal(t, http.StatusOK, rr.Code, "Expected response code to be 200 OK")
// 	assert.Equal(t, "Authenticated user: test@example.com", rr.Body.String(), "Expected authenticated email in response")
// }

// // TestDashboard_InvalidToken tests the Dashboard with an invalid JWT token
// func TestDashboard_InvalidToken(t *testing.T) {
// 	// Create a request with an invalid JWT token
// 	req := httptest.NewRequest("GET", "/dashboard", nil)
// 	req.Header.Set("Authorization", "Bearer invalid_token")
// 	rr := httptest.NewRecorder()

// 	// Define a simple handler to simulate the dashboard functionality
// 	handler := middleware.JWTAuth(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		// This handler should not be reached for invalid tokens
// 		w.WriteHeader(http.StatusOK)
// 		w.Write([]byte("This should not be called"))
// 	}))

// 	// Serve the request
// 	handler.ServeHTTP(rr, req)

// 	// Assert the response status code and error message
// 	assert.Equal(t, http.StatusUnauthorized, rr.Code, "Expected response code to be 401 Unauthorized")
// 	assert.Contains(t, rr.Body.String(), "Invalid token", "Expected an invalid token error message")
// }
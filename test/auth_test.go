package auth_handlers_test

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"erp/handlers/auth_handlers"
	"erp/middleware"
	"erp/models"
	"erp/utils"
	"net/http"
	"net/http/httptest"
	"testing"

	"golang.org/x/crypto/bcrypt"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func setupAuthHandlers(mockDB *sql.DB) *auth_handlers.AuthHandlers {
	// Initialize the UserStore with the mock DB
	userStore := &auth_handlers.DBUserStore{DB: mockDB}
	// Initialize AuthHandlers
	return &auth_handlers.AuthHandlers{UserStore: userStore}
}

func TestSignUp(t *testing.T) {
	// Create a mock database
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to open sqlmock: %s", err)
	}
	defer mockDB.Close()

	// Set up AuthHandlers
	authHandlers := setupAuthHandlers(mockDB)

	// Mock the query for checking if the user already exists (user does not exist)
	mock.ExpectQuery("SELECT email, password FROM users WHERE email=\\$1").
		WithArgs("test@example.com").
		WillReturnRows(sqlmock.NewRows([]string{})) // Simulate no user found

	// Mock the insert query
	mock.ExpectExec("INSERT INTO users").WithArgs("test@example.com", "role", "department").
		WillReturnResult(sqlmock.NewResult(1, 1)) // Simulate successful insertion

	// Prepare the request body
	reqBody := `{"email":"test@example.com", "role":"role", "department":"department"}`
	req := httptest.NewRequest("POST", "/signup", bytes.NewBufferString(reqBody))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	// Call the handler
	authHandlers.SignUp(rr, req)

	// Log the response for debugging
	t.Logf("Response: %v", rr.Body.String())

	// Assert the response status code
	assert.Equal(t, http.StatusCreated, rr.Code, "Expected response code to be 201 Created")

	// Assert the response body
	expectedBody := "User created successfully"
	assert.Equal(t, expectedBody, rr.Body.String(), "Response body mismatch")
}

func TestLogin(t *testing.T) {
	// Create a mock database
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to open sqlmock: %s", err)
	}
	defer mockDB.Close()

	// Set up AuthHandlers
	authHandlers := setupAuthHandlers(mockDB)

	// Mock the query for retrieving the hashed password
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
	mock.ExpectQuery("SELECT email, password FROM users WHERE email=\\$1").
		WithArgs("test@example.com").
		WillReturnRows(sqlmock.NewRows([]string{"email", "password"}).AddRow("test@example.com", hashedPassword))

	// Prepare the request body
	reqBody := `{"email":"test@example.com", "password":"password123"}`
	req := httptest.NewRequest("POST", "/login", bytes.NewBufferString(reqBody))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	// Call the handler
	authHandlers.Login(rr, req)

	// Assert the response status code
	assert.Equal(t, http.StatusOK, rr.Code, "Expected response code to be 200 OK")
}

func TestJWTAuth_ValidToken(t *testing.T) {
	// Generate a valid JWT token for testing
	tokenString, _ := utils.GenerateJWT("test@example.com")

	// Create a request with the token in the Authorization header
	req := httptest.NewRequest("GET", "/dashboard", nil)
	req.Header.Set("Authorization", "Bearer "+tokenString)
	rr := httptest.NewRecorder()

	// Create a handler to test the middleware
	handler := middleware.JWTAuth(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Authenticated"))
	}))

	// Serve the request
	handler.ServeHTTP(rr, req)

	// Assert the response status code
	assert.Equal(t, http.StatusOK, rr.Code, "Expected response code to be 200 OK")
	assert.Equal(t, "Authenticated", rr.Body.String())
}

func TestJWTAuth_InvalidToken(t *testing.T) {
	// Create a request with an invalid JWT token in the Authorization header
	req := httptest.NewRequest("GET", "/dashboard", nil)
	req.Header.Set("Authorization", "Bearer invalidtoken")
	rr := httptest.NewRecorder()

	// Create a handler to test the middleware
	handler := middleware.JWTAuth(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Authenticated"))
	}))

	// Serve the request
	handler.ServeHTTP(rr, req)

	// Assert the response status code
	assert.Equal(t, http.StatusUnauthorized, rr.Code, "Expected response code to be 401 Unauthorized")
}

// CheckUser verifies if a user needs to set a new password

func TestCheckUser(t *testing.T) {
	pass := "password123"
	logic := pass == "" // if password is empty, needsNewPass is true

	// Create a mock database
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to open sqlmock: %s", err)
	}
	defer mockDB.Close()

	// Set up AuthHandlers
	authHandlers := setupAuthHandlers(mockDB)

	// Mock the query for checking the user
	mock.ExpectQuery("SELECT email, password FROM users WHERE email=\\$1").
		WithArgs("test@example.com").
		WillReturnRows(sqlmock.NewRows([]string{"email", "password"}).
			AddRow("test@example.com", pass)) // We can use a valid or empty password as per our logic

	// Prepare the request body
	inputUser := models.User{
		Email: "test@example.com",
	}
	reqBody, _ := json.Marshal(inputUser)
	req := httptest.NewRequest("POST", "/check_user", bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	// Call the handler
	authHandlers.CheckUser(rr, req)

	// Assert the response status code
	assert.Equal(t, http.StatusOK, rr.Code, "Expected response code to be 200 OK")

	// Assert the response body
	var response map[string]bool
	err = json.NewDecoder(rr.Body).Decode(&response)
	assert.NoError(t, err)
	assert.Equal(t, map[string]bool{"needsNewPass": logic}, response) // Change this based on your expected logic

	// Ensure all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

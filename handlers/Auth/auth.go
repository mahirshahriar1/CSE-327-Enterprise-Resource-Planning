// Description: This file contains the handlers for user authentication, including sign-up, login, and password management.
package auth_handlers

import (
	"database/sql"
	"encoding/json"
	"erp/db"
	"erp/utils"
	"log"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

// User represents a user in the system
type User struct {
	Email        string `json:"email"`
	Password     string `json:"password,omitempty"`
	NeedsNewPass bool   `json:"needsNewPass,omitempty"`
}

// SignUpRequest represents the request structure for user sign-up
type SignUpRequest struct {
	Email      string `json:"email"`
	Role       string `json:"role"`
	Department string `json:"department"`
}

// SetNewPasswordRequest represents the request structure for setting a new password
type SetNewPasswordRequest struct {
	Email       string `json:"email"`
	NewPassword string `json:"new_password"`
}

// SignUp handles the user registration process
func SignUp(w http.ResponseWriter, r *http.Request) {
	var req SignUpRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	// Check if the user already exists
	var existingUser string
	err = db.DB.QueryRow("SELECT email FROM users WHERE email=$1", req.Email).Scan(&existingUser)
	if err != nil && err != sql.ErrNoRows {
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}
	if existingUser != "" {
		http.Error(w, "User already exists", http.StatusConflict)
		return
	}

	// Insert the new user (with a null password)
	_, err = db.DB.Exec("INSERT INTO users (email, role, department) VALUES ($1, $2, $3)", req.Email, req.Role, req.Department)
	if err != nil {
		http.Error(w, "Could not create user", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("User created successfully"))
}

// CheckUser verifies if a user needs to set a new password
func CheckUser(w http.ResponseWriter, r *http.Request) {
	var user User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	// Check if the user exists in the database
	var existingPassword sql.NullString
	err = db.DB.QueryRow("SELECT password FROM users WHERE email=$1", user.Email).Scan(&existingPassword)
	if err == sql.ErrNoRows {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	} else if err != nil {
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}

	// If no password is set (NULL in the database), ask to set a new one
	if !existingPassword.Valid || existingPassword.String == "" {
		user.NeedsNewPass = true
	} else {
		user.NeedsNewPass = false
	}

	// Respond with the user's status (needs new password or not)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

// SetNewPassword handles setting a new password for first-time login
func SetNewPassword(w http.ResponseWriter, r *http.Request) {
	var req SetNewPasswordRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	log.Printf("Setting new password for email: %s", req.Email)

	// Check if the user exists and has no password (first-time login)
	var existingUser string
	var password *string // Use a pointer to handle NULL values
	err = db.DB.QueryRow("SELECT email, password FROM users WHERE email=$1", req.Email).Scan(&existingUser, &password)
	if err == sql.ErrNoRows {
		http.Error(w, "User not found", http.StatusNotFound)
		log.Println("User not found:", req.Email)
		return
	} else if err != nil {
		http.Error(w, "Server error", http.StatusInternalServerError)
		log.Println("Error querying database:", err)
		return
	}

	// If the password pointer is nil, it means the password is NULL in the database
	if password != nil && *password != "" {
		http.Error(w, "Password already set. Use login instead.", http.StatusConflict)
		log.Println("User already has a password:", req.Email)
		return
	}

	// Hash the new password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Error setting password", http.StatusInternalServerError)
		log.Println("Error hashing password:", err)
		return
	}

	// Update the user's password in the database
	_, err = db.DB.Exec("UPDATE users SET password=$1 WHERE email=$2", string(hashedPassword), req.Email)
	if err != nil {
		http.Error(w, "Error updating password", http.StatusInternalServerError)
		log.Println("Error updating password in database:", err)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Password set successfully"))
	log.Println("Password set successfully for user:", req.Email)
}

// Login handles the authentication process for existing users
func Login(w http.ResponseWriter, r *http.Request) {
	var user User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil || user.Password == "" {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	// Retrieve the user's hashed password from the database
	var hashedPassword string
	err = db.DB.QueryRow("SELECT password FROM users WHERE email=$1", user.Email).Scan(&hashedPassword)
	if err == sql.ErrNoRows {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	} else if err != nil {
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}

	// Compare the entered password with the hashed password
	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(user.Password))
	if err != nil {
		http.Error(w, "Invalid password", http.StatusUnauthorized)
		return
	}

	// Generate JWT token
	tokenString, err := utils.GenerateJWT(user.Email)
	if err != nil {
		http.Error(w, "Could not generate token", http.StatusInternalServerError)
		return
	}

	// Return the generated token
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"token": tokenString,
	})
}


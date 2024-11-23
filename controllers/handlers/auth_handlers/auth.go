// Description: This file contains the handlers for the authentication endpoints.
package auth_handlers

import (
	"encoding/json"
	"erp/models"
	"erp/controllers/utils"
	"errors"
	"fmt"

	"log"
	"net/http"

	"github.com/gorilla/mux"
	"golang.org/x/crypto/bcrypt"
)

// AuthHandlers struct contains the user store dependency
type AuthHandlers struct {
	UserStore models.UserStore
}

// RegisterRoutes registers all the authentication routes
func (h *AuthHandlers) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/signup", h.SignUp).Methods("POST")
	router.HandleFunc("/check-user", h.CheckUser).Methods("POST")
	router.HandleFunc("/set-new-password", h.SetNewPassword).Methods("POST")
	router.HandleFunc("/login", h.Login).Methods("POST")
}

// SignUp handles the user registration process
func (h *AuthHandlers) SignUp(w http.ResponseWriter, r *http.Request) {
	var req models.SignUpRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	// Check if the user already exists
	_, err = h.UserStore.GetUserByEmail(req.Email)
	if err == nil {
		http.Error(w, "User already exists", http.StatusConflict)
		return
	} else if !errors.Is(err, ErrUserNotFound) {
		// Log the unexpected error and respond with "Server error"
		fmt.Println("Error:", err)
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}

	// Insert the new user (with name, email, role, and department)
	err = h.UserStore.CreateUser(req.Name, req.Email, req.Role, req.Department)
	if err != nil {
		http.Error(w, "Could not create user", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("User created successfully"))
}

// CheckUser verifies if a user needs to set a new password
func (h *AuthHandlers) CheckUser(w http.ResponseWriter, r *http.Request) {
	var req models.User
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	// Check if the user exists
	existingUser, err := h.UserStore.GetUserByEmail(req.Email)
	if err != nil {
		if errors.Is(err, ErrUserNotFound) {
			http.Error(w, "User not found", http.StatusNotFound)
			return
		}
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}

	// Respond with whether the user needs to set a new password
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]bool{"needsNewPass": existingUser.NeedsNewPass})
}

// SetNewPassword handles setting a new password for first-time login
func (h *AuthHandlers) SetNewPassword(w http.ResponseWriter, r *http.Request) {
	var req models.SetNewPasswordRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	// Check if the user exists and needs a new password
	existingUser, err := h.UserStore.GetUserByEmail(req.Email)
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		log.Println("User not found:", req.Email)
		return
	}

	if !existingUser.NeedsNewPass {
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
	err = h.UserStore.UpdatePassword(req.Email, string(hashedPassword))
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
func (h *AuthHandlers) Login(w http.ResponseWriter, r *http.Request) {
	var credentials models.LoginCredentials
	err := json.NewDecoder(r.Body).Decode(&credentials)
	if err != nil || credentials.Password == "" {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	// Retrieve the user's hashed password and check if the user exists
	existingUser, err := h.UserStore.GetUserByEmail(credentials.Email)
	if err != nil {
		if errors.Is(err, ErrUserNotFound) {
			http.Error(w, "User not found", http.StatusNotFound)
			return
		}
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}
	if existingUser.NeedsNewPass {
		http.Error(w, "User needs to set a new password", http.StatusUnauthorized)
		return
	}

	// Compare the provided password with the stored hashed password
	err = bcrypt.CompareHashAndPassword([]byte(existingUser.Password), []byte(credentials.Password))
	if err != nil {
		http.Error(w, "Invalid password", http.StatusUnauthorized)
		return
	}

	// Generate JWT token
	tokenString, err := utils.GenerateJWT(existingUser.Email, existingUser.Role.RoleName)
	if err != nil {
		http.Error(w, "Could not generate token", http.StatusInternalServerError)
		return
	}

	// Return the generated token along with the user's name and role
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"token": tokenString,
		"name":  existingUser.Name,
		"role":  existingUser.Role.RoleName,
	})
}

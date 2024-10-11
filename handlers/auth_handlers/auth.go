// Description: This file contains the handlers for the authentication endpoints.
package auth_handlers

import (
	"encoding/json"
	"erp/models"
	"erp/utils"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"golang.org/x/crypto/bcrypt"
)

// AuthHandlers struct contains the user store dependency
type AuthHandlers struct {
	UserStore models.UserStore
}

func (h *AuthHandlers) RegisterRoutes(router *mux.Router){
	router.HandleFunc("/signup", h.SignUp).Methods("POST")
	router.HandleFunc("/check-user" , h.CheckUser).Methods("POST")
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
	} else if err.Error() != "user not found" {
		// return the error if it's not a "user not found" error not the http error
		http.Error(w, err.Error(), http.StatusInternalServerError)
		// http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}

	// Insert the new user (with a null password)
	err = h.UserStore.CreateUser(req.Email, req.Role, req.Department)
	if err != nil {
		http.Error(w, "Could not create user", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("User created successfully"))
}

// CheckUser verifies if a user needs to set a new password
func (h *AuthHandlers) CheckUser(w http.ResponseWriter, r *http.Request) {
	var user models.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	// Check if the user exists
	existingUser, err := h.UserStore.GetUserByEmail(user.Email)
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}
	// Respond with the user's status
	w.Header().Set("Content-Type", "application/json")

	if !existingUser.NeedsNewPass {
		json.NewEncoder(w).Encode(map[string]bool{"needsNewPass": false})
		return
	}
	json.NewEncoder(w).Encode(map[string]bool{"needsNewPass": true})
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
	var user models.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil || user.Password == "" {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	// Retrieve the user's hashed password
	existingUser, err := h.UserStore.GetUserByEmail(user.Email)
	// log.Println("Existing user:", existingUser)
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}
	if existingUser.NeedsNewPass {
		http.Error(w, "User needs to set a new password", http.StatusUnauthorized)
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(existingUser.Password), []byte(user.Password))
	if err != nil {
		http.Error(w, "Invalid password", http.StatusUnauthorized)
		log.Println("Password comparison failed:", err)
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

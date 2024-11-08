package models // or package types, based on your preference

// User represents a user in the system
type User struct {
	ID           int    `json:"id,omitempty"`
	Name         string `json:"name,omitempty"`
	Email        string `json:"email"`
	Password     string `json:"password,omitempty"`
	Role         Role   `json:"role"`
	Department   string `json:"department"`
	NeedsNewPass bool   `json:"needsNewPass,omitempty"`
}

// LoginCredentials represents the structure for user login
type LoginCredentials struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// SignUpRequest represents the request structure for user sign-up
type SignUpRequest struct {
	Name       string `json:"name"`
	Email      string `json:"email"`
	Role       string `json:"role"`
	Department string `json:"department"`
}

// SetNewPasswordRequest represents the request structure for setting a new password
type SetNewPasswordRequest struct {
	Email       string `json:"email"`
	NewPassword string `json:"new_password"`
}

// UserStore defines an interface for user-related database operations
type UserStore interface {
	CreateUser(name, email, role, department string) error
	GetUserByEmail(email string) (*User, error)
	UpdatePassword(email, hashedPassword string) error
}

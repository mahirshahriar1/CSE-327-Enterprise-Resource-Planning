package types

// User represents a user in the system
type User struct {
	Email        string `json:"email"`
	Role         string `json:"role"`
	Department   string `json:"department"`
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

// UserStore defines an interface for user-related database operations
type UserStore interface {
	CreateUser(email, role, department string) error
	GetUserByEmail(email string) (*User, error)
	UpdatePassword(email, hashedPassword string) error
}

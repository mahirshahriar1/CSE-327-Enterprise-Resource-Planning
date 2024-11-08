package models

// Role represents a role in the system
type Role struct {
	ID          int    `json:"id"`
	RoleName    string `json:"role_name"`
	Permissions string `json:"permissions"`
}

// RoleStore defines an interface for role-related database operations
type RoleStore interface {
	GetRoleByID(id int) (*Role, error)
	GetRoleByName(roleName string) (*Role, error)
}

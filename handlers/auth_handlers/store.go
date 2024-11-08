package auth_handlers

import (
	"database/sql"
	"erp/models"
	"errors"
	"fmt"
)

// ErrUserNotFound is returned when a user cannot be found in the database
var ErrUserNotFound = errors.New("user not found")

// DBUserStore implements UserStore using a SQL database
type DBUserStore struct {
	DB        *sql.DB
	RoleStore models.RoleStore // RoleStore dependency to fetch roles
}

// CreateUser inserts a new user into the database with the specified name, role, and department
func (s *DBUserStore) CreateUser(name, email, roleName, department string) error {
    // Retrieve the role ID based on the role name
    role, err := s.RoleStore.GetRoleByName(roleName)
    if err != nil {
        return err // Role not found or other error
    }

    // Insert the new user with the retrieved role ID and specified name
    _, err = s.DB.Exec("INSERT INTO users (name, email, role_id, department) VALUES ($1, $2, $3, $4)", name, email, role.ID, department)
	fmt.Println("Eror in CreateUser", err)
    return err
}


// GetUserByEmail fetches a user by email along with their role information
func (s *DBUserStore) GetUserByEmail(email string) (*models.User, error) {
    var user models.User
    var roleID int
    var existingPassword sql.NullString

    // Retrieve the user's information, including the name
    err := s.DB.QueryRow("SELECT id, name, email, password, role_id, department, needs_new_pass FROM users WHERE email = $1", email).Scan(
        &user.ID, &user.Name, &user.Email, &existingPassword, &roleID, &user.Department, &user.NeedsNewPass)
    
    if err == sql.ErrNoRows {
        return nil, ErrUserNotFound // Custom error for "user not found"
    } else if err != nil {
        return nil, err // Return any other errors
    }

    user.Password = existingPassword.String
    user.NeedsNewPass = !existingPassword.Valid || existingPassword.String == ""

    // Retrieve the role by ID and assign it to the user
    role, err := s.RoleStore.GetRoleByID(roleID)
    if err != nil {
        return nil, err
    }
    user.Role = *role
    return &user, nil
}

// UpdatePassword updates the user's password in the database
func (s *DBUserStore) UpdatePassword(email, hashedPassword string) error {
	_, err := s.DB.Exec("UPDATE users SET password=$1 WHERE email=$2", hashedPassword, email)
	return err
}

// DBRoleStore implements RoleStore using a SQL database
type DBRoleStore struct {
	DB *sql.DB
}

// GetRoleByID retrieves a role by its ID
func (s *DBRoleStore) GetRoleByID(id int) (*models.Role, error) {
	var role models.Role
	err := s.DB.QueryRow("SELECT id, role_name, permissions FROM roles WHERE id=$1", id).Scan(
		&role.ID, &role.RoleName, &role.Permissions)
	if err == sql.ErrNoRows {
		return nil, errors.New("role not found")
	} else if err != nil {
		return nil, err
	}
	return &role, nil
}

// GetRoleByName retrieves a role by its name
func (s *DBRoleStore) GetRoleByName(roleName string) (*models.Role, error) {
	var role models.Role
	err := s.DB.QueryRow("SELECT id, role_name, permissions FROM roles WHERE role_name=$1", roleName).Scan(
		&role.ID, &role.RoleName, &role.Permissions)
	if err == sql.ErrNoRows {
		return nil, errors.New("role not found")
	} else if err != nil {
		return nil, err
	}
	return &role, nil
}

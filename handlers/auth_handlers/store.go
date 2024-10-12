package auth_handlers

import (
	"database/sql"
	"erp/models"
	"errors"
)

// DBUserStore implements UserStore using a SQL database
type DBUserStore struct {
	DB *sql.DB
}

// CreateUser inserts a new user into the database with a null password
func (s *DBUserStore) CreateUser(email, role, department string) error {
	_, err := s.DB.Exec("INSERT INTO users (email, role, department) VALUES ($1, $2, $3)", email, role, department)
	return err
}

// GetUserByEmail fetches a user by email
func (s *DBUserStore) GetUserByEmail(email string) (*models.User, error) {
	var user models.User
	var existingPassword sql.NullString
	err := s.DB.QueryRow("SELECT email, password FROM users WHERE email=$1", email).Scan(&user.Email, &existingPassword)
	if err == sql.ErrNoRows {
		return nil, errors.New("user not found")
	} else if err != nil {
		return nil, err
	}
	user.Password = existingPassword.String
	// If the password is NULL in the database, set NeedsNewPass to true or false
	user.NeedsNewPass = !existingPassword.Valid || existingPassword.String == ""
	return &user, nil
}

// UpdatePassword updates the user's password in the database
func (s *DBUserStore) UpdatePassword(email, hashedPassword string) error {
	_, err := s.DB.Exec("UPDATE users SET password=$1 WHERE email=$2", hashedPassword, email)
	return err
}

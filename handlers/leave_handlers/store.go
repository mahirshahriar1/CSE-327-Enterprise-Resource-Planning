package leave_handlers

import (
	"database/sql"
	"erp/models"
)

// DBLeaveStore implements LeaveStore interface
type DBLeaveStore struct {
	DB *sql.DB
}

// CreateLeave inserts a new leave request into the database
func (store *DBLeaveStore) CreateLeave(leave *models.Leave) error {
	query := "INSERT INTO leave (user_id, leave_type, start_date, end_date, status) VALUES ($1, $2, $3, $4, $5)"
	_, err := store.DB.Exec(query, leave.UserID, leave.LeaveType, leave.StartDate, leave.EndDate, leave.Status)
	return err
}

// UpdateLeaveStatus updates the status of an existing leave request
func (store *DBLeaveStore) UpdateLeaveStatus(id int, status string) error {
	query := "UPDATE leave SET status = $1 WHERE id = $2"
	_, err := store.DB.Exec(query, status, id)
	return err
}

package leave_handlers

import (
	"database/sql"
	"erp/models"
)

// DBLeaveStore provides an implementation of the LeaveStore interface using a SQL database.
type DBLeaveStore struct {
	DB *sql.DB // Database connection
}

// CreateLeave inserts a new leave request into the database.
// The `leave` parameter holds the details of the leave request, including user ID, leave type, start date, end date, and status.
// If successful, the leave request is added to the `leave` table in the database.
func (store *DBLeaveStore) CreateLeave(leave *models.Leave) error {
	query := "INSERT INTO leave (user_id, leave_type, start_date, end_date, status) VALUES ($1, $2, $3, $4, $5)"
	_, err := store.DB.Exec(query, leave.UserID, leave.LeaveType, leave.StartDate, leave.EndDate, leave.Status)
	return err
}

// UpdateLeaveStatus updates the status of an existing leave request in the database.
// The `id` parameter specifies the unique identifier of the leave to update.
// The `status` parameter indicates the new status of the leave request (e.g., "Approved", "Rejected").
// If successful, the status of the leave request is updated in the `leave` table.
func (store *DBLeaveStore) UpdateLeaveStatus(id int, status string) error {
	query := "UPDATE leave SET status = $1 WHERE id = $2"
	_, err := store.DB.Exec(query, status, id)
	return err
}

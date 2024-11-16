// Package leave_handlers provides functionality to manage leave requests,
// including creating new requests and updating their statuses using a SQL database.
package leave_handlers

import (
	"database/sql"
	"erp/models"
)

// DBLeaveStore provides an implementation of the LeaveStore interface using a SQL database.
// It allows for operations such as creating leave requests and updating their statuses.
type DBLeaveStore struct {
	DB *sql.DB // DB represents the database connection.
}

// CreateLeave inserts a new leave request into the database.
//
// Parameters:
//   - leave: A pointer to the Leave object containing the details of the leave request, including:
//   - UserID: The ID of the user requesting leave.
//   - LeaveType: The type of leave (e.g., "Vacation", "Sick Leave").
//   - StartDate: The starting date of the leave.
//   - EndDate: The ending date of the leave.
//   - Status: The status of the leave request (e.g., "Pending").
//
// Returns:
//   - error: An error if the operation fails, otherwise nil.
//
// Details:
//   - This method executes an SQL `INSERT` query to add the leave request to the `leave` table in the database.
//   - The leave request's details must be valid and conform to the database schema for successful insertion.
func (store *DBLeaveStore) CreateLeave(leave *models.Leave) error {
	query := "INSERT INTO leave (user_id, leave_type, start_date, end_date, status) VALUES ($1, $2, $3, $4, $5)"
	_, err := store.DB.Exec(query, leave.UserID, leave.LeaveType, leave.StartDate, leave.EndDate, leave.Status)
	return err
}

// UpdateLeaveStatus updates the status of an existing leave request in the database.
//
// Parameters:
//   - id: The unique identifier of the leave request to be updated.
//   - status: A string representing the new status of the leave request (e.g., "Approved", "Rejected").
//
// Returns:
//   - error: An error if the operation fails, otherwise nil.
//
// Details:
//   - This method executes an SQL `UPDATE` query to modify the `status` column of the specified leave request in the `leave` table.
//   - The leave request must exist in the database for the update to succeed.
func (store *DBLeaveStore) UpdateLeaveStatus(id int, status string) error {
	query := "UPDATE leave SET status = $1 WHERE id = $2"
	_, err := store.DB.Exec(query, status, id)
	return err
}

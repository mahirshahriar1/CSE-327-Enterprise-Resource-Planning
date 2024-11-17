// Package attendance_handlers provides database interaction for creating attendance records.
package attendance_handlers

import (
	"database/sql"
	"erp/models"
)

// DBAttendanceStore implements the AttendanceStore interface for SQL database operations.
// It handles creating attendance records in the database.
type DBAttendanceStore struct {
	DB *sql.DB // DB represents the database connection.
}

// CreateAttendance inserts a new attendance record into the database.
//
// Parameters:
//   - attendance: A pointer to the Attendance object containing the record details, including:
//     - UserID: The ID of the user marking attendance.
//     - CheckIn: The check-in time.
//     - CheckOut: The check-out time (if available; otherwise it can be zero).
//     - TotalHours: Calculated hours based on CheckIn and CheckOut.
//
// Returns:
//   - error: An error if the operation fails, otherwise nil.
//
// Details:
//   - This method executes an SQL `INSERT` query to add the attendance record to the `attendance` table.
func (store *DBAttendanceStore) CreateAttendance(attendance *models.Attendance) error {
	query := "INSERT INTO attendance (user_id, check_in, check_out, total_hours) VALUES ($1, $2, $3, $4)"
	_, err := store.DB.Exec(query, attendance.UserID, attendance.CheckIn, attendance.CheckOut, attendance.TotalHours)
	return err
}

// Package attendance_handlers provides database interaction for creating attendance records.
// using a SQL database, including fetching attendance data for specific users.
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
//   - UserID: The ID of the user marking attendance.
//   - CheckIn: The check-in time.
//   - CheckOut: The check-out time (if available; otherwise it can be zero).
//   - TotalHours: Calculated hours based on CheckIn and CheckOut.
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

// GetAttendanceByUserID retrieves all attendance records for a specific user from the database.
//
// Parameters:
//   - userID: The ID of the user whose attendance records are being retrieved.
//
// Returns:
//   - []*models.Attendance: A slice of Attendance objects for the user.
//   - error: An error if the operation fails, otherwise nil.
//
// Details:
//   - This method executes an SQL `SELECT` query to fetch attendance records from the `attendance` table
//     where the `user_id` matches the provided ID.
//   - The records are returned in the order they are found in the database.
func (store *DBAttendanceStore) GetAttendanceByUserID(userID int) ([]*models.Attendance, error) {
	// Prepare the query to fetch attendance records for the given user ID
	query := "SELECT id, user_id, check_in, check_out, total_hours FROM attendance WHERE user_id = $1"

	// Execute the query
	rows, err := store.DB.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Collect the results into a slice of Attendance objects
	var attendanceRecords []*models.Attendance
	for rows.Next() {
		var attendance models.Attendance
		if err := rows.Scan(&attendance.ID, &attendance.UserID, &attendance.CheckIn, &attendance.CheckOut, &attendance.TotalHours); err != nil {
			return nil, err
		}
		attendanceRecords = append(attendanceRecords, &attendance)
	}

	// Return the slice of attendance records
	return attendanceRecords, nil
}

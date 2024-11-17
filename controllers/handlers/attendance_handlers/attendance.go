// Package attendance_handlers provides an HTTP handler for creating attendance records.
package attendance_handlers

import (
	"encoding/json"
	"erp/models"
	"fmt"
	"net/http"
)

// CreateAttendanceRecord handles the creation of a new attendance record.
// It returns an HTTP handler function to process attendance creation requests.
//
// The handler expects a JSON payload with the following structure:
//
//	{
//	  "user_id": 1,
//	  "check_in": "2024-11-16T09:00:00Z",
//	  "check_out": "2024-11-16T17:00:00Z"
//	}
//
// Details:
//   - On success, it responds with HTTP 201 (Created) and the attendance record details in JSON format.
//   - On failure, it responds with an appropriate HTTP error status.
//
// Parameters:
//   - store: An implementation of the AttendanceStore interface to handle database operations.
//
// Returns:
//   - http.HandlerFunc: The HTTP handler function for creating attendance records.
func CreateAttendanceRecord(store models.AttendanceStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var attendance models.Attendance

		// Decode the JSON body from the request
		if err := json.NewDecoder(r.Body).Decode(&attendance); err != nil {
			http.Error(w, "Invalid request payload", http.StatusBadRequest)
			return
		}

		// Calculate total hours worked if both check-in and check-out are provided
		if !attendance.CheckIn.IsZero() && !attendance.CheckOut.IsZero() {
			duration := attendance.CheckOut.Sub(attendance.CheckIn)
			if duration < 0 {
				http.Error(w, "Check-out time cannot be before check-in time", http.StatusBadRequest)
				return
			}
			attendance.TotalHours = duration.Hours()
		}

		// Create the attendance record in the database
		if err := store.CreateAttendance(&attendance); err != nil {
			http.Error(w, fmt.Sprintf("Failed to create attendance: %v", err), http.StatusInternalServerError)
			return
		}

		// Respond with the created attendance record
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(attendance)
	}
}

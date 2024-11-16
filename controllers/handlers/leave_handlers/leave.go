// Package leave_handlers provides HTTP handlers for managing leave requests.
// It includes handlers for creating new leave requests and updating their status.
package leave_handlers

import (
	"encoding/json"
	"erp/models"
	"fmt"
	"net/http"
)

// LeaveStore defines the interface for database operations related to leave requests.
// It provides methods for creating leave requests and updating their status.
type LeaveStore interface {
	// CreateLeave inserts a new leave request into the database.
	// Parameters:
	//   - leave: A pointer to the Leave object containing the details of the leave request.
	// Returns:
	//   - error: An error if the insertion fails, otherwise nil.
	CreateLeave(leave *models.Leave) error

	// UpdateLeaveStatus updates the status of an existing leave request.
	// Parameters:
	//   - id: The unique identifier of the leave request.
	//   - status: A string representing the new status (e.g., "Approved", "Rejected").
	// Returns:
	//   - error: An error if the update fails, otherwise nil.
	UpdateLeaveStatus(id int, status string) error
}

// CreateLeaveHandler creates a new leave request in the system.
// It returns an HTTP handler function to process the creation of leave requests.
//
// The handler expects a JSON payload with the following structure:
//
//	{
//	  "user_id": 1,
//	  "leave_type": "Vacation",
//	  "start_date": "2024-11-20",
//	  "end_date": "2024-11-25"
//	}
//
// Details:
//   - The status of the new leave request is automatically set to "Pending".
//   - On success, it responds with HTTP 201 (Created) and the leave request details in JSON format.
//   - On failure, it responds with an appropriate HTTP error status.
//
// Parameters:
//   - store: An implementation of the LeaveStore interface to handle database operations.
//
// Returns:
//   - http.HandlerFunc: The HTTP handler function for creating leave requests.
func CreateLeaveHandler(store LeaveStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var leave models.Leave

		// Parse the JSON body from the request
		if err := json.NewDecoder(r.Body).Decode(&leave); err != nil {
			http.Error(w, "Invalid request payload", http.StatusBadRequest)
			return
		}

		// Default status for a new leave request is "Pending".
		leave.Status = "Pending"

		// Attempt to create the leave in the database
		if err := store.CreateLeave(&leave); err != nil {
			http.Error(w, fmt.Sprintf("Failed to create leave: %v", err), http.StatusInternalServerError)
			return
		}

		// Respond with the created leave request and a 201 status code.
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(leave)
	}
}

// UpdateLeaveStatusHandler updates the status of an existing leave request.
// It returns an HTTP handler function to process status updates for leave requests.
//
// The handler expects a JSON payload with the following structure:
//
//	{
//	  "id": 1,
//	  "status": "Approved"
//	}
//
// Details:
//   - On success, it responds with HTTP 200 (OK) and a success message.
//   - On failure, it responds with an appropriate HTTP error status.
//
// Parameters:
//   - store: An implementation of the LeaveStore interface to handle database operations.
//
// Returns:
//   - http.HandlerFunc: The HTTP handler function for updating leave request statuses.
func UpdateLeaveStatusHandler(store LeaveStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Extract the leave ID and status from the request.
		var requestData struct {
			ID     int    `json:"id"`
			Status string `json:"status"`
		}

		// Parse the JSON body from the request.
		if err := json.NewDecoder(r.Body).Decode(&requestData); err != nil {
			http.Error(w, "Invalid request payload", http.StatusBadRequest)
			return
		}

		// Attempt to update the leave status in the database.
		if err := store.UpdateLeaveStatus(requestData.ID, requestData.Status); err != nil {
			http.Error(w, fmt.Sprintf("Failed to update leave status: %v", err), http.StatusInternalServerError)
			return
		}

		// Respond with a success message.
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "Leave status updated successfully")
	}
}

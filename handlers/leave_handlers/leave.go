package leave_handlers

import (
	"encoding/json"
	"erp/models"
	"fmt"
	"net/http"
)

// LeaveStore interface for database operations.
// It provides methods for creating leave requests and updating their status.
type LeaveStore interface {
	// CreateLeave inserts a new leave request into the database.
	// The `leave` parameter holds the details of the leave request.
	CreateLeave(leave *models.Leave) error
	// UpdateLeaveStatus updates the status of an existing leave request.
	// The `id` parameter specifies the unique identifier of the leave.
	// The `status` parameter indicates the new status of the leave request (e.g., "Approved", "Rejected").
	UpdateLeaveStatus(id int, status string) error
}

// CreateLeaveHandler returns an HTTP handler function that handles the creation of a new leave request.
// It expects a JSON payload with the details of the leave request, including user ID, leave type, start date, and end date.
// The status of the new leave request is automatically set to "Pending".
//
// Example JSON request body:
//
//	{
//	  "user_id": 1,
//	  "leave_type": "Vacation",
//	  "start_date": "2024-11-20",
//	  "end_date": "2024-11-25"
//	}
//
// On success, the function returns an HTTP status code of 201 (Created) and the leave request details in JSON format.
// If the request payload is invalid or there is an error in creating the leave request, the function returns an appropriate HTTP error status.
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

// UpdateLeaveStatusHandler returns an HTTP handler function that handles updating the status of a leave request.
// It expects a JSON payload with the leave ID and the new status.
//
// Example JSON request body:
//
//	{
//	  "id": 1,
//	  "status": "Approved"
//	}
//
// On success, the function returns an HTTP status code of 200 (OK) and a success message.
// If the request payload is invalid or there is an error in updating the leave status, the function returns an appropriate HTTP error status.
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

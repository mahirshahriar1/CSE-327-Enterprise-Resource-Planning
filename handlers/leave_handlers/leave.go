package leave_handlers

import (
	"encoding/json"
	"erp/models"
	"fmt"
	"net/http"
)

// LeaveStore interface for database operations
type LeaveStore interface {
	CreateLeave(leave *models.Leave) error
	UpdateLeaveStatus(id int, status string) error
}

// CreateLeaveHandler handles the creation of a new leave request
func CreateLeaveHandler(store LeaveStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var leave models.Leave

		// Parse the JSON body
		if err := json.NewDecoder(r.Body).Decode(&leave); err != nil {
			http.Error(w, "Invalid request payload", http.StatusBadRequest)
			return
		}

		// Default status to "Pending"
		leave.Status = "Pending"

		// Attempt to create the leave in the database
		if err := store.CreateLeave(&leave); err != nil {
			http.Error(w, fmt.Sprintf("Failed to create leave: %v", err), http.StatusInternalServerError)
			return
		}

		// Respond with success
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(leave)
	}
}

// UpdateLeaveStatusHandler handles updating the status of a leave request
func UpdateLeaveStatusHandler(store LeaveStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Extract the leave ID and status from the request
		var requestData struct {
			ID     int    `json:"id"`
			Status string `json:"status"`
		}
		if err := json.NewDecoder(r.Body).Decode(&requestData); err != nil {
			http.Error(w, "Invalid request payload", http.StatusBadRequest)
			return
		}

		// Attempt to update the leave status in the database
		if err := store.UpdateLeaveStatus(requestData.ID, requestData.Status); err != nil {
			http.Error(w, fmt.Sprintf("Failed to update leave status: %v", err), http.StatusInternalServerError)
			return
		}

		// Respond with success
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "Leave status updated successfully")
	}
}

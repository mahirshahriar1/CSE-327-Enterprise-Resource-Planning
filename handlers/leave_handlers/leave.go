package leave_handlers

import (
	"encoding/json"
	"erp/models"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

// LeaveHandlers struct to hold dependencies
type LeaveHandlers struct {
	LeaveStore models.LeaveStore
}

// CreateLeave handles creating a new leave request
func (h *LeaveHandlers) CreateLeave(w http.ResponseWriter, r *http.Request) {
	var leave models.Leave
	if err := json.NewDecoder(r.Body).Decode(&leave); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}
	leave.Status = "Pending"     // Default status
	leave.StartDate = time.Now() // You can adjust how dates are parsed

	if err := h.LeaveStore.CreateLeave(&leave); err != nil {
		http.Error(w, "Failed to create leave request", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(leave)
}

// GetLeavesByUserID fetches all leave requests for a given user
func (h *LeaveHandlers) GetLeavesByUserID(w http.ResponseWriter, r *http.Request) {
	userID, err := strconv.Atoi(mux.Vars(r)["userID"])
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	leaves, err := h.LeaveStore.GetLeaveByUserID(userID)
	if err != nil {
		http.Error(w, "Failed to fetch leave requests", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(leaves)
}

// UpdateLeaveStatus handles updating the status of a leave request
func (h *LeaveHandlers) UpdateLeaveStatus(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		http.Error(w, "Invalid leave ID", http.StatusBadRequest)
		return
	}

	var req struct {
		Status string `json:"status"`
		Reason string `json:"reason"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	if err := h.LeaveStore.UpdateLeaveStatus(id, req.Status); err != nil {
		http.Error(w, "Failed to update leave status", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Leave status updated successfully"))
}

package leave_handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"erp/models"

	"github.com/stretchr/testify/assert"
)

// MockLeaveStore is a mock implementation of the LeaveStore interface.
// It simulates a database using an in-memory map for storing leave requests.
type MockLeaveStore struct {
	leaves map[int]*models.Leave // In-memory storage for leave requests.
	nextID int                   // Counter to assign unique IDs to leave requests.
}

// CreateLeave adds a new leave request to the mock store.
// Assigns a unique ID to the leave request.
//
// Parameters:
//   - leave: Pointer to the Leave object to store.
//
// Returns:
//   - error: Always nil as the operation is simulated.
func (m *MockLeaveStore) CreateLeave(leave *models.Leave) error {
	m.nextID++
	leave.ID = m.nextID
	m.leaves[leave.ID] = leave
	return nil
}

// UpdateLeaveStatus updates the status of an existing leave request in the mock store.
//
// Parameters:
//   - id: The ID of the leave request to update.
//   - status: The new status of the leave request (e.g., "Approved", "Rejected").
//
// Returns:
//   - error: "leave not found" if the leave ID does not exist.
func (m *MockLeaveStore) UpdateLeaveStatus(id int, status string) error {
	leave, exists := m.leaves[id]
	if !exists {
		return errors.New("leave not found")
	}
	leave.Status = status
	return nil
}

// TestCreateLeaveHandler verifies the CreateLeaveHandler for creating a new leave request.
// It checks whether the handler assigns an ID and default "Pending" status and responds with 201.
func TestCreateLeaveHandler(t *testing.T) {
	// Parse dates into time.Time objects.
	startDate, _ := time.Parse("2006-01-02", "2024-11-20")
	endDate, _ := time.Parse("2006-01-02", "2024-11-25")

	// Initialize the mock store and handler.
	store := &MockLeaveStore{leaves: make(map[int]*models.Leave)}
	handler := CreateLeaveHandler(store)

	// Create a sample leave request.
	leave := models.Leave{
		UserID:    1,
		LeaveType: "Vacation",
		StartDate: startDate,
		EndDate:   endDate,
	}
	body, _ := json.Marshal(leave)
	req, _ := http.NewRequest("POST", "/leaves", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	// Record the response.
	rr := httptest.NewRecorder()

	// Call the handler.
	handler(rr, req)

	// Validate the response status code.
	assert.Equal(t, http.StatusCreated, rr.Code)

	// Validate the response body.
	var createdLeave models.Leave
	json.NewDecoder(rr.Body).Decode(&createdLeave)
	assert.Equal(t, 1, createdLeave.ID)                      // Check if ID is assigned.
	assert.Equal(t, "Pending", createdLeave.Status)          // Check if default status is "Pending".
	assert.Equal(t, leave.UserID, createdLeave.UserID)       // Verify the UserID matches the input.
	assert.Equal(t, leave.LeaveType, createdLeave.LeaveType) // Verify the LeaveType matches the input.
}

// TestUpdateLeaveStatusHandler verifies the UpdateLeaveStatusHandler for updating the status of a leave request.
// It checks whether the handler correctly updates the status and responds with 200.
func TestUpdateLeaveStatusHandler(t *testing.T) {
	// Parse dates into time.Time objects.
	startDate, _ := time.Parse("2006-01-02", "2024-11-20")
	endDate, _ := time.Parse("2006-01-02", "2024-11-22")

	// Initialize the mock store and handler.
	store := &MockLeaveStore{leaves: make(map[int]*models.Leave)}
	handler := UpdateLeaveStatusHandler(store)

	// Add a leave request to the mock store.
	store.leaves[1] = &models.Leave{
		ID:        1,
		UserID:    1,
		LeaveType: "Sick Leave",
		StartDate: startDate,
		EndDate:   endDate,
		Status:    "Pending",
	}

	// Create a request to update the status of the leave request.
	updateRequest := struct {
		ID     int    `json:"id"`
		Status string `json:"status"`
	}{
		ID:     1,
		Status: "Approved",
	}
	body, _ := json.Marshal(updateRequest)
	req, _ := http.NewRequest("PUT", "/leaves/status", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	// Record the response.
	rr := httptest.NewRecorder()

	// Call the handler.
	handler(rr, req)

	// Validate the response status code.
	assert.Equal(t, http.StatusOK, rr.Code)

	// Validate the updated status in the mock store.
	assert.Equal(t, "Approved", store.leaves[1].Status) // Check if the status was updated correctly.
}

package attendance_handlers

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

// MockAttendanceStore is a mock implementation of the AttendanceStore interface.
// It simulates database operations using an in-memory map to store attendance records.
type MockAttendanceStore struct {
	attendance map[int]*models.Attendance // In-memory storage for attendance records.
	nextID     int                        // Counter to assign unique IDs to attendance records.
}

// CreateAttendance simulates adding a new attendance record to the mock store.
// It assigns a unique ID to the record and stores it in memory.
//
// Parameters:
//   - attendance: Pointer to the Attendance object to store.
//
// Returns:
//   - error: Always nil as the operation is simulated.
func (m *MockAttendanceStore) CreateAttendance(attendance *models.Attendance) error {
	m.nextID++                               // Increment the ID counter.
	attendance.ID = m.nextID                 // Assign a unique ID.
	m.attendance[attendance.ID] = attendance // Store the attendance record in memory.
	return nil
}

// GetAttendanceByUserID simulates retrieving all attendance records for a specific user.
// It iterates over the in-memory storage and collects records that match the given user ID.
//
// Parameters:
//   - userID: The ID of the user whose attendance records are to be retrieved.
//
// Returns:
//   - []*models.Attendance: A slice of Attendance records for the user.
//   - error: Always nil as the operation is simulated.
func (m *MockAttendanceStore) GetAttendanceByUserID(userID int) ([]*models.Attendance, error) {
	var records []*models.Attendance
	for _, record := range m.attendance {
		if record.UserID == userID {
			records = append(records, record)
		}
	}
	return records, nil
}

// UpdateAttendance simulates updating an existing attendance record in the mock store.
// It checks if the record exists in memory and updates it if found.
//
// Parameters:
//   - attendance: Pointer to the Attendance object with updated details.
//
// Returns:
//   - error: An error if the record is not found, otherwise nil.
func (m *MockAttendanceStore) UpdateAttendance(attendance *models.Attendance) error {
	if _, exists := m.attendance[attendance.ID]; !exists {
		return errors.New("attendance record not found")
	}
	m.attendance[attendance.ID] = attendance
	return nil
}

// DeleteAttendance is a mock implementation for compatibility with the AttendanceStore interface.
// It is a no-op implementation used to satisfy the interface requirements.
//
// Parameters:
//   - id: The ID of the attendance record to delete (not used in this mock implementation).
//
// Returns:
//   - error: Always nil as the operation is not performed.
func (m *MockAttendanceStore) DeleteAttendance(id int) error {
	return nil
}

// TestCreateAttendanceRecord verifies the CreateAttendanceRecord handler.
// It checks whether the handler creates a new attendance record with an assigned ID
// and calculates the total hours worked based on check-in and check-out times.
func TestCreateAttendanceRecord(t *testing.T) {
	// Initialize the mock store and handler.
	store := &MockAttendanceStore{attendance: make(map[int]*models.Attendance)}
	handler := CreateAttendanceRecord(store)

	// Create a sample input attendance record.
	input := models.Attendance{
		UserID:   1,
		CheckIn:  time.Now(),
		CheckOut: time.Now().Add(8 * time.Hour), // Simulate an 8-hour workday.
	}
	body, _ := json.Marshal(input)                                          // Convert the input to JSON format.
	req, _ := http.NewRequest("POST", "/attendance", bytes.NewBuffer(body)) // Create an HTTP POST request with the JSON body.
	req.Header.Set("Content-Type", "application/json")                      // Set the Content-Type header to JSON.

	// Record the HTTP response using a test recorder.
	rr := httptest.NewRecorder()

	// Call the handler with the request and response recorder.
	handler(rr, req)

	// Assert that the HTTP status code is 201 (Created).
	assert.Equal(t, http.StatusCreated, rr.Code)

	// Decode the response body into an Attendance object.
	var result models.Attendance
	json.NewDecoder(rr.Body).Decode(&result)

	// Assert that the ID is generated and not zero.
	assert.NotZero(t, result.ID)

	// Assert that the input data is correctly reflected in the response.
	assert.Equal(t, input.UserID, result.UserID) // Check the UserID matches.
	assert.Equal(t, 8.0, result.TotalHours)      // Check the total hours are calculated correctly.
}

// TestGetAttendanceByUserID verifies the GetAttendanceByUserID handler.
// It checks whether the handler retrieves attendance records for a specific user
// and returns them in the correct format.
func TestGetAttendanceByUserID(t *testing.T) {
	// Initialize the mock store with sample data.
	store := &MockAttendanceStore{
		attendance: map[int]*models.Attendance{
			1: {ID: 1, UserID: 1, CheckIn: time.Now(), CheckOut: time.Now().Add(8 * time.Hour)},
			2: {ID: 2, UserID: 1, CheckIn: time.Now(), CheckOut: time.Now().Add(4 * time.Hour)},
			3: {ID: 3, UserID: 2, CheckIn: time.Now(), CheckOut: time.Now().Add(6 * time.Hour)},
		},
	}

	// Initialize the handler.
	handler := GetAttendanceByUserID(store)

	// Create an HTTP GET request for user ID 1.
	req, _ := http.NewRequest("GET", "/attendance?user_id=1", nil)

	// Record the HTTP response using a test recorder.
	rr := httptest.NewRecorder()

	// Call the handler with the request and response recorder.
	handler(rr, req)

	// Assert that the HTTP status code is 200 (OK).
	assert.Equal(t, http.StatusOK, rr.Code)

	// Decode the response body into a slice of Attendance objects.
	var results []*models.Attendance
	json.NewDecoder(rr.Body).Decode(&results)

	// Assert that the correct number of records are returned.
	assert.Len(t, results, 2) // Two records exist for user ID 1.

	// Assert that all returned records have the correct UserID.
	for _, record := range results {
		assert.Equal(t, 1, record.UserID)
	}
}

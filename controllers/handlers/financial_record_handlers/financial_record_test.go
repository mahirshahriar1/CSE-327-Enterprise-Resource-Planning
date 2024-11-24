// Package financial_record_handlers provides HTTP handlers for managing financial records.
// It includes functions for creating, retrieving, updating, and deleting financial records.
// The handlers interact with a store for data persistence, and HTTP routes are registered using the Gorilla mux router.
//
// This package is intended to support CRUD (Create, Read, Update, Delete) operations for financial records
// in an ERP system.

package financial_record_handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"erp/models"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

// TestCreateRecord tests the creation of a new financial record.
// It sends a POST request with a sample financial record and asserts that the response
// status is 201 (Created) and the returned record has the correct ID and details.
//
// Params:
// - t: *testing.T - The testing object used to report the test results.
//
// It mocks the store method CreateFinancialRecord to return a successful creation.
// The request is made with a sample FinancialRecord and verified against the response body.
func TestCreateRecord(t *testing.T) {
	// Mock data for financial record creation
	record := &models.FinancialRecord{
		TransactionID: 123,
		AccountID:     456,
		Amount:        1000.00,
		TransactionDate: time.Date(2024, time.November, 17, 0, 0, 0, 0, time.UTC),
		TransactionType: "Credit",
		Description: "Payment received",
	}

	// Mock store that returns no error for CreateFinancialRecord
	mockStore := &MockFinancialRecordStore{
		CreateFinancialRecordFn: func(record *models.FinancialRecord) error {
			record.ID = 1
			return nil
		},
	}

	// Create mock HTTP request
	reqBody, _ := json.Marshal(record)
	req, err := http.NewRequest("POST", "/records", bytes.NewBuffer(reqBody))
	if err != nil {
		t.Fatalf("Could not create request: %v", err)
	}

	// Create mock HTTP recorder
	rr := httptest.NewRecorder()

	// Register the handler with routes
	r := mux.NewRouter()
	RegisterRoutes(r, mockStore)

	// Serve the request
	r.ServeHTTP(rr, req)

	// Assert the response status
	assert.Equal(t, http.StatusCreated, rr.Code)

	// Assert the response body contains the created record
	var createdRecord models.FinancialRecord
	err = json.NewDecoder(rr.Body).Decode(&createdRecord)
	if err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}
	assert.Equal(t, 1, createdRecord.ID)
	assert.Equal(t, record.TransactionID, createdRecord.TransactionID)
	assert.Equal(t, record.AccountID, createdRecord.AccountID)
	assert.Equal(t, record.Amount, createdRecord.Amount)
}

// TestGetRecord tests the retrieval of a financial record by its ID.
// It sends a GET request and asserts that the response status is 200 (OK)
// and that the returned record matches the expected values.
//
// Params:
// - t: *testing.T - The testing object used to report the test results.
//
// This test simulates retrieving a financial record from the mock store by its ID.
// The response is verified to ensure the correct financial record is returned.
func TestGetRecord(t *testing.T) {
	// Mock data for retrieving a financial record
	record := &models.FinancialRecord{
		ID:             1,
		TransactionID:  123,
		AccountID:      456,
		Amount:         1000.00,
		TransactionDate: time.Date(2024, time.November, 17, 0, 0, 0, 0, time.UTC),
		TransactionType: "Credit",
		Description:    "Payment received",
	}

	// Mock store with data retrieval method
	mockStore := &MockFinancialRecordStore{
		GetFinancialRecordByIDFn: func(id int) (*models.FinancialRecord, error) {
			if id == 1 {
				return record, nil
			}
			return nil, fmt.Errorf("record not found")
		},
	}

	// Create mock HTTP request to get the financial record with ID 1
	req, err := http.NewRequest("GET", "/records/1", nil)
	if err != nil {
		t.Fatalf("Could not create request: %v", err)
	}

	// Create mock HTTP recorder
	rr := httptest.NewRecorder()

	// Register the handler with routes
	r := mux.NewRouter()
	RegisterRoutes(r, mockStore)

	// Serve the request
	r.ServeHTTP(rr, req)

	// Assert the response status is OK (200)
	assert.Equal(t, http.StatusOK, rr.Code)

	// Assert the response body matches the expected financial record
	var fetchedRecord models.FinancialRecord
	err = json.NewDecoder(rr.Body).Decode(&fetchedRecord)
	if err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}
	assert.Equal(t, record.ID, fetchedRecord.ID)
}

// TestUpdateRecord tests the updating of an existing financial record.
// It sends a PUT request with updated data and asserts that the response
// status is 200 (OK) to indicate a successful update.
//
// Params:
// - t: *testing.T - The testing object used to report the test results.
//
// This test simulates updating an existing financial record and verifies that
// the update was successful by checking the HTTP response status.
func TestUpdateRecord(t *testing.T) {
	// Mock data for updating the financial record
	record := &models.FinancialRecord{
		ID:             1,
		TransactionID:  123,
		AccountID:      456,
		Amount:         1000.00,
		TransactionDate: time.Date(2024, time.November, 17, 0, 0, 0, 0, time.UTC),
		TransactionType: "Credit",
		Description:    "Payment received",
	}

	// Mock store with update method
	mockStore := &MockFinancialRecordStore{
		UpdateFinancialRecordFn: func(record *models.FinancialRecord) error {
			// Assume the update is successful
			return nil
		},
	}

	// Create mock HTTP request with updated data
	reqBody, _ := json.Marshal(record)
	req, err := http.NewRequest("PUT", "/records/1", bytes.NewBuffer(reqBody))
	if err != nil {
		t.Fatalf("Could not create request: %v", err)
	}

	// Create mock HTTP recorder
	rr := httptest.NewRecorder()

	// Register the handler with routes
	r := mux.NewRouter()
	RegisterRoutes(r, mockStore)

	// Serve the request
	r.ServeHTTP(rr, req)

	// Assert the response status is OK (200)
	assert.Equal(t, http.StatusOK, rr.Code)
}

// TestDeleteRecord tests the deletion of a financial record.
// It sends a DELETE request and asserts that the response status is 204 (No Content)
// to indicate a successful deletion.
//
// Params:
// - t: *testing.T - The testing object used to report the test results.
//
// This test simulates deleting an existing financial record and verifies that
// the deletion was successful by checking the HTTP response status.
func TestDeleteRecord(t *testing.T) {
	// Mock store with delete method
	mockStore := &MockFinancialRecordStore{
		DeleteFinancialRecordFn: func(id int) error {
			// Assume the delete is successful
			return nil
		},
	}

	// Create mock HTTP request to delete the financial record with ID 1
	req, err := http.NewRequest("DELETE", "/records/1", nil)
	if err != nil {
		t.Fatalf("Could not create request: %v", err)
	}

	// Create mock HTTP recorder
	rr := httptest.NewRecorder()

	// Register the handler with routes
	r := mux.NewRouter()
	RegisterRoutes(r, mockStore)

	// Serve the request
	r.ServeHTTP(rr, req)

	// Assert the response status is No Content (204)
	assert.Equal(t, http.StatusNoContent, rr.Code)
}

// MockFinancialRecordStore is a mock implementation of the FinancialRecordStore interface.
// It is used to simulate interactions with a data store during testing.
type MockFinancialRecordStore struct {
	CreateFinancialRecordFn    func(record *models.FinancialRecord) error
	GetFinancialRecordByIDFn   func(id int) (*models.FinancialRecord, error)
	UpdateFinancialRecordFn    func(record *models.FinancialRecord) error
	DeleteFinancialRecordFn    func(id int) error
	GetAllFinancialRecordsFn   func() ([]models.FinancialRecord, error) // Added this function
}

// CreateFinancialRecord simulates the creation of a financial record in the store.
// It invokes the mock function CreateFinancialRecordFn.
func (m *MockFinancialRecordStore) CreateFinancialRecord(record *models.FinancialRecord) error {
	return m.CreateFinancialRecordFn(record)
}

// GetFinancialRecordByID retrieves a financial record by its ID from the mock store.
// It invokes the mock function GetFinancialRecordByIDFn.
func (m *MockFinancialRecordStore) GetFinancialRecordByID(id int) (*models.FinancialRecord, error) {
	return m.GetFinancialRecordByIDFn(id)
}

// UpdateFinancialRecord simulates the updating of a financial record in the store.
// It invokes the mock function UpdateFinancialRecordFn.
func (m *MockFinancialRecordStore) UpdateFinancialRecord(record *models.FinancialRecord) error {
	return m.UpdateFinancialRecordFn(record)
}

// DeleteFinancialRecord simulates the deletion of a financial record in the store.
// It invokes the mock function DeleteFinancialRecordFn.
func (m *MockFinancialRecordStore) DeleteFinancialRecord(id int) error {
	return m.DeleteFinancialRecordFn(id)
}

// GetAllFinancialRecords retrieves all financial records from the mock store.
// It invokes the mock function GetAllFinancialRecordsFn.
func (m *MockFinancialRecordStore) GetAllFinancialRecords() ([]models.FinancialRecord, error) {
	return m.GetAllFinancialRecordsFn()
}

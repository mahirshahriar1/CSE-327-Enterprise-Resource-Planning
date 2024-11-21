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

func TestCreateRecord(t *testing.T) {
	// Mock data
	record := &models.FinancialRecord{
		TransactionID: 123,
		AccountID:     456,
		Amount:        1000.00,
		TransactionDate: time.Date(2024, time.November, 17, 0, 0, 0, 0, time.UTC),
		TransactionType: "Credit",
		Description: "Payment received",
	}

	// Create a mock store that returns no error for CreateFinancialRecord
	mockStore := &MockFinancialRecordStore{
		CreateFinancialRecordFn: func(record *models.FinancialRecord) error {
			record.ID = 1
			return nil
		},
	}

	// Create a mock HTTP request
	reqBody, _ := json.Marshal(record)
	req, err := http.NewRequest("POST", "/records", bytes.NewBuffer(reqBody))
	if err != nil {
		t.Fatalf("Could not create request: %v", err)
	}

	// Create a mock HTTP recorder
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

func TestGetRecord(t *testing.T) {
	// Mock data
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

	// Create mock HTTP request
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

	// Assert the response status
	assert.Equal(t, http.StatusOK, rr.Code)

	// Assert the response body
	var fetchedRecord models.FinancialRecord
	err = json.NewDecoder(rr.Body).Decode(&fetchedRecord)
	if err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}
	assert.Equal(t, record.ID, fetchedRecord.ID)
}

func TestUpdateRecord(t *testing.T) {
	// Mock data
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

	// Create mock HTTP request
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

	// Assert the response status
	assert.Equal(t, http.StatusOK, rr.Code)
}

func TestDeleteRecord(t *testing.T) {
	// Mock store with delete method
	mockStore := &MockFinancialRecordStore{
		DeleteFinancialRecordFn: func(id int) error {
			// Assume the delete is successful
			return nil
		},
	}

	// Create mock HTTP request
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

	// Assert the response status
	assert.Equal(t, http.StatusNoContent, rr.Code)
}

// Mock Financial Record Store
type MockFinancialRecordStore struct {
	CreateFinancialRecordFn    func(record *models.FinancialRecord) error
	GetFinancialRecordByIDFn   func(id int) (*models.FinancialRecord, error)
	UpdateFinancialRecordFn    func(record *models.FinancialRecord) error
	DeleteFinancialRecordFn    func(id int) error
	GetAllFinancialRecordsFn   func() ([]models.FinancialRecord, error) // Added this function
}

func (m *MockFinancialRecordStore) CreateFinancialRecord(record *models.FinancialRecord) error {
	return m.CreateFinancialRecordFn(record)
}

func (m *MockFinancialRecordStore) GetFinancialRecordByID(id int) (*models.FinancialRecord, error) {
	return m.GetFinancialRecordByIDFn(id)
}

func (m *MockFinancialRecordStore) UpdateFinancialRecord(record *models.FinancialRecord) error {
	return m.UpdateFinancialRecordFn(record)
}

func (m *MockFinancialRecordStore) DeleteFinancialRecord(id int) error {
	return m.DeleteFinancialRecordFn(id)
}


func (m *MockFinancialRecordStore) GetAllFinancialRecords() ([]models.FinancialRecord, error) {
	return m.GetAllFinancialRecordsFn()
}

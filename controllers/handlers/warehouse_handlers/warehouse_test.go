// Package warehouse_handlers_test contains unit tests for warehouse-related HTTP handlers.
package warehouse_handlers

import (
	"bytes"
	"encoding/json"
	"erp/models"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

// TestCreateWarehouse tests the CreateWarehouse handler.
func TestCreateWarehouse(t *testing.T) {
	// Set up mock database
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("could not create mock db: %v", err)
	}
	defer db.Close()

	store := &DBWarehouseStore{DB: db}
	handler := &WarehouseHandlers{WarehouseStore: store}

	// Sample warehouse data
	warehouse := &models.Warehouse{
		Name:     "Test Warehouse",
		Capacity: 500,
		Location: "Test Location",
	}

	// Mock database behavior
	mock.ExpectExec("INSERT INTO warehouses").
		WithArgs(warehouse.Name, warehouse.Capacity, warehouse.Location).
		WillReturnResult(sqlmock.NewResult(1, 1))

	// Create HTTP request and recorder
	body, _ := json.Marshal(warehouse)
	req, _ := http.NewRequest("POST", "/warehouses", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	// Call the handler
	handler.CreateWarehouse(rec, req)

	// Assert that no error occurred and response is correct
	assert.Equal(t, http.StatusCreated, rec.Code)
	assert.Equal(t, "Warehouse created successfully", rec.Body.String())

	// Assert that the expected query was executed
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unmet expectations: %v", err)
	}
}

// TestGetWarehouseByID tests the GetWarehouseByID handler.
func TestGetWarehouseByID(t *testing.T) {
	// Set up mock database
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("could not create mock db: %v", err)
	}
	defer db.Close()

	store := &DBWarehouseStore{DB: db}
	handler := &WarehouseHandlers{WarehouseStore: store}

	// Sample warehouse data
	warehouse := &models.Warehouse{
		ID:       1,
		Name:     "Test Warehouse",
		Capacity: 500,
		Location: "Test Location",
	}

	// Mock database behavior
	mock.ExpectQuery("SELECT id, name, capacity, location FROM warehouses WHERE id = \\$1").
		WithArgs(warehouse.ID).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name", "capacity", "location"}).
			AddRow(warehouse.ID, warehouse.Name, warehouse.Capacity, warehouse.Location))

	// Create HTTP request and recorder
	req, _ := http.NewRequest("GET", "/warehouses/1", nil)
	rec := httptest.NewRecorder()

	// Add mux variables to the request
	req = mux.SetURLVars(req, map[string]string{"id": "1"})

	// Call the handler
	handler.GetWarehouseByID(rec, req)

	// Assert that no error occurred and response is correct
	assert.Equal(t, http.StatusOK, rec.Code)
	expectedBody, _ := json.Marshal(warehouse)
	assert.JSONEq(t, string(expectedBody), rec.Body.String())

	// Assert that the expected query was executed
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unmet expectations: %v", err)
	}
}

// TestUpdateWarehouse tests the UpdateWarehouse handler.
func TestUpdateWarehouse(t *testing.T) {
	// Set up mock database
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("could not create mock db: %v", err)
	}
	defer db.Close()

	store := &DBWarehouseStore{DB: db}
	handler := &WarehouseHandlers{WarehouseStore: store}

	// Sample warehouse data
	warehouse := &models.Warehouse{
		ID:       1,
		Name:     "Updated Warehouse",
		Capacity: 600,
		Location: "Updated Location",
	}

	// Mock database behavior
	mock.ExpectExec("UPDATE warehouses SET name = \\$1, capacity = \\$2, location = \\$3 WHERE id = \\$4").
		WithArgs(warehouse.Name, warehouse.Capacity, warehouse.Location, warehouse.ID).
		WillReturnResult(sqlmock.NewResult(1, 1))

	// Create HTTP request and recorder
	body, _ := json.Marshal(warehouse)
	req, _ := http.NewRequest("PUT", "/warehouses/1", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	// Add mux variables to the request
	req = mux.SetURLVars(req, map[string]string{"id": "1"})

	// Call the handler
	handler.UpdateWarehouse(rec, req)

	// Assert that no error occurred and response is correct
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, "Warehouse updated successfully", rec.Body.String())

	// Assert that the expected query was executed
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unmet expectations: %v", err)
	}
}

// TestDeleteWarehouse tests the DeleteWarehouse handler.
func TestDeleteWarehouse(t *testing.T) {
	// Set up mock database
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("could not create mock db: %v", err)
	}
	defer db.Close()

	store := &DBWarehouseStore{DB: db}
	handler := &WarehouseHandlers{WarehouseStore: store}

	// Mock database behavior
	mock.ExpectExec("DELETE FROM warehouses WHERE id = \\$1").
		WithArgs(1).
		WillReturnResult(sqlmock.NewResult(1, 1))

	// Create HTTP request and recorder
	req, _ := http.NewRequest("DELETE", "/warehouses/1", nil)
	rec := httptest.NewRecorder()

	// Add mux variables to the request
	req = mux.SetURLVars(req, map[string]string{"id": "1"})

	// Call the handler
	handler.DeleteWarehouse(rec, req)

	// Assert that no error occurred and response is correct
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, "Warehouse deleted successfully", rec.Body.String())

	// Assert that the expected query was executed
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unmet expectations: %v", err)
	}
}

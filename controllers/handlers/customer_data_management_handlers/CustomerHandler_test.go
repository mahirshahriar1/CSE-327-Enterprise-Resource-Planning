// Package customer_data_management_handlers_test provides comprehensive tests for the customer data management HTTP handlers.
// These tests ensure correctness by validating customer creation, retrieval, updating, and deletion.
// A mock data store is used to simulate a database in an isolated environment.

package customer_data_management_handlers_test

import (
	"bytes"
	"encoding/json"
	"erp/models"
	"erp/controllers/handlers/customer_data_management_handlers"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

// MockCustomerStore provides a mock implementation of the CustomerStore interface.
//
// This mock simulates database operations, allowing unit tests to focus
// on the handler logic without depending on a real database connection.
type MockCustomerStore struct {
	customers map[int]*models.Customer // Stores mock customer data
	nextID    int                      // Tracks the next available customer ID
}

// NewMockCustomerStore initializes a new instance of the MockCustomerStore.
//
// Returns:
//   - A pointer to a MockCustomerStore with an empty customer map and ID counter set to 1.
func NewMockCustomerStore() *MockCustomerStore {
	return &MockCustomerStore{
		customers: make(map[int]*models.Customer),
		nextID:    1,
	}
}

// CreateCustomer simulates adding a new customer to the mock store.
//
// Parameters:
//   - customer: Pointer to the Customer object to be added.
//
// Returns:
//   - Always returns nil as it assumes no errors in a mock setup.
func (m *MockCustomerStore) CreateCustomer(customer *models.Customer) error {
	customer.ID = m.nextID
	m.customers[m.nextID] = customer
	m.nextID++
	return nil
}

// GetCustomerByID simulates fetching a customer by their ID.
//
// Parameters:
//   - id: The unique identifier of the customer.
//
// Returns:
//   - The customer object if found.
//   - models.ErrNotFound if no customer exists with the given ID.
func (m *MockCustomerStore) GetCustomerByID(id int) (*models.Customer, error) {
	customer, exists := m.customers[id]
	if !exists {
		return nil, models.ErrNotFound
	}
	return customer, nil
}

// UpdateCustomer simulates updating an existing customer's data.
//
// Parameters:
//   - customer: Pointer to the updated Customer object.
//
// Returns:
//   - nil if the update is successful.
//   - models.ErrNotFound if no customer exists with the given ID.
func (m *MockCustomerStore) UpdateCustomer(customer *models.Customer) error {
	_, exists := m.customers[customer.ID]
	if !exists {
		return models.ErrNotFound
	}
	m.customers[customer.ID] = customer
	return nil
}

// DeleteCustomer simulates deleting a customer by their ID.
//
// Parameters:
//   - id: The unique identifier of the customer to be deleted.
//
// Returns:
//   - nil if the deletion is successful.
//   - models.ErrNotFound if no customer exists with the given ID.
func (m *MockCustomerStore) DeleteCustomer(id int) error {
	_, exists := m.customers[id]
	if !exists {
		return models.ErrNotFound
	}
	delete(m.customers, id)
	return nil
}

// TestCreateCustomerHandler validates the CreateCustomerHandler functionality.
//
// Steps:
//   - Simulate an HTTP POST request to create a new customer.
//   - Verify the response status and ensure the customer data is correctly returned.
func TestCreateCustomerHandler(t *testing.T) {
	store := NewMockCustomerStore()
	handler := customer_data_management_handlers.CustomerHandlers{Store: store}

	// Input data for a new customer
	newCustomer := &models.Customer{Name: "Test Customer", Contact: "1234567890", OrderHistory: "Order 1, Order 2"}
	payload, _ := json.Marshal(newCustomer)

	// Simulate the HTTP POST request
	req, _ := http.NewRequest(http.MethodPost, "/customers", bytes.NewBuffer(payload))
	rec := httptest.NewRecorder()

	// Invoke the handler
	handler.CreateCustomerHandler(rec, req)

	// Assertions
	assert.Equal(t, http.StatusCreated, rec.Code, "Expected status code 201 Created")
	var createdCustomer models.Customer
	json.NewDecoder(rec.Body).Decode(&createdCustomer)
	assert.Equal(t, newCustomer.Name, createdCustomer.Name, "Customer name mismatch")
	assert.Equal(t, newCustomer.Contact, createdCustomer.Contact, "Customer contact mismatch")
	assert.Equal(t, newCustomer.OrderHistory, createdCustomer.OrderHistory, "Customer order history mismatch")
}

// TestGetCustomerByIDHandler validates the GetCustomerByIDHandler functionality.
//
// Steps:
//   - Add a customer to the mock store.
//   - Simulate an HTTP GET request to fetch the customer by ID.
//   - Verify the response status and customer data.
func TestGetCustomerByIDHandler(t *testing.T) {
	store := NewMockCustomerStore()
	handler := customer_data_management_handlers.CustomerHandlers{Store: store}

	// Add a customer to the mock store
	store.CreateCustomer(&models.Customer{Name: "Existing Customer", Contact: "9876543210", OrderHistory: "Order 3"})

	// Simulate the HTTP GET request
	req, _ := http.NewRequest(http.MethodGet, "/customers/1", nil)
	req = mux.SetURLVars(req, map[string]string{"id": "1"})
	rec := httptest.NewRecorder()

	// Invoke the handler
	handler.GetCustomerByIDHandler(rec, req)

	// Assertions
	assert.Equal(t, http.StatusOK, rec.Code, "Expected status code 200 OK")
	var retrievedCustomer models.Customer
	json.NewDecoder(rec.Body).Decode(&retrievedCustomer)
	assert.Equal(t, "Existing Customer", retrievedCustomer.Name, "Customer name mismatch")
	assert.Equal(t, "9876543210", retrievedCustomer.Contact, "Customer contact mismatch")
	assert.Equal(t, "Order 3", retrievedCustomer.OrderHistory, "Customer order history mismatch")
}

// TestUpdateCustomerHandler validates the UpdateCustomerHandler functionality.
//
// Steps:
//   - Add a customer to the mock store.
//   - Simulate an HTTP PUT request to update the customer data.
//   - Verify the response status and updated data.
func TestUpdateCustomerHandler(t *testing.T) {
	store := NewMockCustomerStore()
	handler := customer_data_management_handlers.CustomerHandlers{Store: store}

	// Add a customer to the mock store
	store.CreateCustomer(&models.Customer{Name: "Old Name", Contact: "0000000000", OrderHistory: "Order A"})

	// Updated customer data
	updatedCustomer := &models.Customer{ID: 1, Name: "Updated Name", Contact: "9999999999", OrderHistory: "Order B"}
	payload, _ := json.Marshal(updatedCustomer)

	// Simulate the HTTP PUT request
	req, _ := http.NewRequest(http.MethodPut, "/customers/1", bytes.NewBuffer(payload))
	req = mux.SetURLVars(req, map[string]string{"id": "1"})
	rec := httptest.NewRecorder()

	// Invoke the handler
	handler.UpdateCustomerHandler(rec, req)

	// Assertions
	assert.Equal(t, http.StatusOK, rec.Code, "Expected status code 200 OK")
	var updatedResult models.Customer
	json.NewDecoder(rec.Body).Decode(&updatedResult)
	assert.Equal(t, updatedCustomer.Name, updatedResult.Name, "Customer name mismatch")
	assert.Equal(t, updatedCustomer.Contact, updatedResult.Contact, "Customer contact mismatch")
	assert.Equal(t, updatedCustomer.OrderHistory, updatedResult.OrderHistory, "Customer order history mismatch")
}

// TestDeleteCustomerHandler validates the DeleteCustomerHandler functionality.
//
// Steps:
//   - Add a customer to the mock store.
//   - Simulate an HTTP DELETE request to remove the customer by ID.
//   - Verify the response status and ensure the customer is deleted.
func TestDeleteCustomerHandler(t *testing.T) {
	store := NewMockCustomerStore()
	handler := customer_data_management_handlers.CustomerHandlers{Store: store}

	// Add a customer to the mock store
	store.CreateCustomer(&models.Customer{Name: "To Be Deleted", Contact: "1111111111", OrderHistory: "Order X"})

	// Simulate the HTTP DELETE request
	req, _ := http.NewRequest(http.MethodDelete, "/customers/1", nil)
	req = mux.SetURLVars(req, map[string]string{"id": "1"})
	rec := httptest.NewRecorder()

	// Invoke the handler
	handler.DeleteCustomerHandler(rec, req)

	// Assertions
	assert.Equal(t, http.StatusNoContent, rec.Code, "Expected status code 204 No Content")
	_, err := store.GetCustomerByID(1)
	assert.Equal(t, models.ErrNotFound, err, "Expected the customer to be deleted")
}

// Package invoice_handlers_test provides comprehensive tests for the invoice data management HTTP handlers.
// These tests ensure correctness by validating invoice creation, retrieval, updating, and deletion.
// A mock data store is used to simulate a database in an isolated environment.

package invoice_handlers

import (
	"bytes"
	"encoding/json"

	"erp/models"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

// MockInvoiceStore provides a mock implementation of the InvoiceStore interface.
//
// This mock simulates database operations, allowing unit tests to focus
// on the handler logic without depending on a real database connection.
type MockInvoiceStore struct {
	invoices map[int]*models.Invoice // Stores mock invoice data
	nextID   int                     // Tracks the next available invoice ID
}

// NewMockInvoiceStore initializes a new instance of the MockInvoiceStore.
//
// Returns:
//   - A pointer to a MockInvoiceStore with an empty invoice map and ID counter set to 1.
func NewMockInvoiceStore() *MockInvoiceStore {
	return &MockInvoiceStore{
		invoices: make(map[int]*models.Invoice),
		nextID:   1,
	}
}

// CreateInvoice simulates adding a new invoice to the mock store.
//
// Parameters:
//   - invoice: Pointer to the Invoice object to be added.
//
// Returns:
//   - Always returns nil as it assumes no errors in a mock setup.
func (m *MockInvoiceStore) CreateInvoice(invoice *models.Invoice) error {
	invoice.ID = m.nextID
	m.invoices[m.nextID] = invoice
	m.nextID++
	return nil
}

// GetInvoiceByID simulates fetching an invoice by its ID.
//
// Parameters:
//   - id: The unique identifier of the invoice.
//
// Returns:
//   - The invoice object if found.
//   - models.ErrNotFound if no invoice exists with the given ID.
func (m *MockInvoiceStore) GetInvoiceByID(id int) (*models.Invoice, error) {
	invoice, exists := m.invoices[id]
	if !exists {
		return nil, models.ErrNotFound
	}
	return invoice, nil
}

// UpdateInvoice simulates updating an existing invoice's data.
//
// Parameters:
//   - invoice: Pointer to the updated Invoice object.
//
// Returns:
//   - nil if the update is successful.
//   - models.ErrNotFound if no invoice exists with the given ID.
func (m *MockInvoiceStore) UpdateInvoice(invoice *models.Invoice) error {
	_, exists := m.invoices[invoice.ID]
	if !exists {
		return models.ErrNotFound
	}
	m.invoices[invoice.ID] = invoice
	return nil
}

// DeleteInvoice simulates deleting an invoice by its ID.
//
// Parameters:
//   - id: The unique identifier of the invoice to be deleted.
//
// Returns:
//   - nil if the deletion is successful.
//   - models.ErrNotFound if no invoice exists with the given ID.
func (m *MockInvoiceStore) DeleteInvoice(id int) error {
	_, exists := m.invoices[id]
	if !exists {
		return models.ErrNotFound
	}
	delete(m.invoices, id)
	return nil
}

// TestCreateInvoiceHandler validates the CreateInvoiceHandler functionality.
//
// Steps:
//   - Simulate an HTTP POST request to create a new invoice.
//   - Verify the response status and ensure the invoice data is correctly returned.
func TestCreateInvoiceHandler(t *testing.T) {
	store := NewMockInvoiceStore()
	handler := InvoiceHandlers{Store: store}

	// Input data for a new invoice
	newInvoice := &models.Invoice{SalesOrderID: 1, CustomerID: 123, Amount: 250.75, Status: "Pending"}
	payload, _ := json.Marshal(newInvoice)

	// Simulate the HTTP POST request
	req, _ := http.NewRequest(http.MethodPost, "/invoices", bytes.NewBuffer(payload))
	rec := httptest.NewRecorder()

	// Invoke the handler
	handler.CreateInvoiceHandler(rec, req)

	// Assertions
	assert.Equal(t, http.StatusCreated, rec.Code, "Expected status code 201 Created")
	var createdInvoice models.Invoice
	json.NewDecoder(rec.Body).Decode(&createdInvoice)
	assert.Equal(t, newInvoice.SalesOrderID, createdInvoice.SalesOrderID, "SalesOrderID mismatch")
	assert.Equal(t, newInvoice.CustomerID, createdInvoice.CustomerID, "CustomerID mismatch")
	assert.Equal(t, newInvoice.Amount, createdInvoice.Amount, "Amount mismatch")
	assert.Equal(t, newInvoice.Status, createdInvoice.Status, "Status mismatch")
}

// TestGetInvoiceByIDHandler validates the GetInvoiceByIDHandler functionality.
//
// Steps:
//   - Add an invoice to the mock store.
//   - Simulate an HTTP GET request to fetch the invoice by ID.
//   - Verify the response status and invoice data.
func TestGetInvoiceByIDHandler(t *testing.T) {
	store := NewMockInvoiceStore()
	handler := InvoiceHandlers{Store: store}

	// Add an invoice to the mock store
	store.CreateInvoice(&models.Invoice{SalesOrderID: 2, CustomerID: 456, Amount: 500.00, Status: "Paid"})

	// Simulate the HTTP GET request
	req, _ := http.NewRequest(http.MethodGet, "/invoices/1", nil)
	req = mux.SetURLVars(req, map[string]string{"id": "1"})
	rec := httptest.NewRecorder()

	// Invoke the handler
	handler.GetInvoiceByIDHandler(rec, req)

	// Assertions
	assert.Equal(t, http.StatusOK, rec.Code, "Expected status code 200 OK")
	var retrievedInvoice models.Invoice
	json.NewDecoder(rec.Body).Decode(&retrievedInvoice)
	assert.Equal(t, 2, retrievedInvoice.SalesOrderID, "SalesOrderID mismatch")
	assert.Equal(t, 456, retrievedInvoice.CustomerID, "CustomerID mismatch")
	assert.Equal(t, 500.00, retrievedInvoice.Amount, "Amount mismatch")
	assert.Equal(t, "Paid", retrievedInvoice.Status, "Status mismatch")
}

// TestUpdateInvoiceHandler validates the UpdateInvoiceHandler functionality.
//
// Steps:
//   - Add an invoice to the mock store.
//   - Simulate an HTTP PUT request to update the invoice data.
//   - Verify the response status and updated data.
func TestUpdateInvoiceHandler(t *testing.T) {
	store := NewMockInvoiceStore()
	handler := InvoiceHandlers{Store: store}

	// Add an invoice to the mock store
	store.CreateInvoice(&models.Invoice{SalesOrderID: 3, CustomerID: 789, Amount: 150.00, Status: "Pending"})

	// Updated invoice data
	updatedInvoice := &models.Invoice{ID: 1, SalesOrderID: 4, CustomerID: 890, Amount: 300.00, Status: "Paid"}
	payload, _ := json.Marshal(updatedInvoice)

	// Simulate the HTTP PUT request
	req, _ := http.NewRequest(http.MethodPut, "/invoices/1", bytes.NewBuffer(payload))
	req = mux.SetURLVars(req, map[string]string{"id": "1"})
	rec := httptest.NewRecorder()

	// Invoke the handler
	handler.UpdateInvoiceHandler(rec, req)

	// Assertions
	assert.Equal(t, http.StatusOK, rec.Code, "Expected status code 200 OK")
	var updatedResult models.Invoice
	json.NewDecoder(rec.Body).Decode(&updatedResult)
	assert.Equal(t, updatedInvoice.SalesOrderID, updatedResult.SalesOrderID, "SalesOrderID mismatch")
	assert.Equal(t, updatedInvoice.CustomerID, updatedResult.CustomerID, "CustomerID mismatch")
	assert.Equal(t, updatedInvoice.Amount, updatedResult.Amount, "Amount mismatch")
	assert.Equal(t, updatedInvoice.Status, updatedResult.Status, "Status mismatch")
}

// TestDeleteInvoiceHandler validates the DeleteInvoiceHandler functionality.
//
// Steps:
//   - Add an invoice to the mock store.
//   - Simulate an HTTP DELETE request to remove the invoice by ID.
//   - Verify the response status and ensure the invoice is deleted.
func TestDeleteInvoiceHandler(t *testing.T) {
	store := NewMockInvoiceStore()
	handler := InvoiceHandlers{Store: store}

	// Add an invoice to the mock store
	store.CreateInvoice(&models.Invoice{SalesOrderID: 5, CustomerID: 123, Amount: 700.00, Status: "Unpaid"})

	// Simulate the HTTP DELETE request
	req, _ := http.NewRequest(http.MethodDelete, "/invoices/1", nil)
	req = mux.SetURLVars(req, map[string]string{"id": "1"})
	rec := httptest.NewRecorder()

	// Invoke the handler
	handler.DeleteInvoiceHandler(rec, req)

	// Assertions
	assert.Equal(t, http.StatusNoContent, rec.Code, "Expected status code 204 No Content")
	_, err := store.GetInvoiceByID(1)
	assert.Equal(t, models.ErrNotFound, err, "Expected the invoice to be deleted")
}

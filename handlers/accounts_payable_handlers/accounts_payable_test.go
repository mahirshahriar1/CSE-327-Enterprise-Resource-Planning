// Package accounts_payable_handlers contains tests for CRUD operations on accounts payable.
// This package uses a mock implementation of the PaymentStore interface for testing purposes
// to ensure the handlers behave as expected.
package accounts_payable_handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"erp/models"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

// MockPaymentStore is a mock implementation of the PaymentStore interface.
// It uses an in-memory map to simulate payment storage for testing purposes.
type MockPaymentStore struct {
	payments map[int]*models.Payment // In-memory storage for payments.
	nextID   int                     // Counter to assign unique IDs to payments.
}

// CreatePayment adds a new payment to the mock store.
// Assigns a unique ID to the payment.
//
// Parameters:
//   - payment: Pointer to the Payment object to store.
//
// Returns:
//   - error: Always nil as the operation is simulated.
func (m *MockPaymentStore) CreatePayment(payment *models.Payment) error {
	m.nextID++
	payment.ID = m.nextID
	m.payments[payment.ID] = payment
	return nil
}

// GetPaymentByID retrieves a payment by its ID from the mock store.
//
// Parameters:
//   - id: The ID of the payment to retrieve.
//
// Returns:
//   - *Payment: Pointer to the payment object if found.
//   - error: "payment not found" if the ID does not exist in the store.
func (m *MockPaymentStore) GetPaymentByID(id int) (*models.Payment, error) {
	payment, exists := m.payments[id]
	if !exists {
		return nil, errors.New("payment not found")
	}
	return payment, nil
}

// UpdatePayment updates an existing payment in the mock store.
//
// Parameters:
//   - payment: Pointer to the Payment object with updated details.
//
// Returns:
//   - error: "payment not found" if the payment ID does not exist.
func (m *MockPaymentStore) UpdatePayment(payment *models.Payment) error {
	_, exists := m.payments[payment.ID]
	if !exists {
		return errors.New("payment not found")
	}
	m.payments[payment.ID] = payment
	return nil
}

// DeletePayment removes a payment by ID from the mock store.
//
// Parameters:
//   - id: The ID of the payment to delete.
//
// Returns:
//   - error: "payment not found" if the ID does not exist in the store.
func (m *MockPaymentStore) DeletePayment(id int) error {
	_, exists := m.payments[id]
	if !exists {
		return errors.New("payment not found")
	}
	delete(m.payments, id)
	return nil
}

// TestCreateBill verifies the CreateBill handler for creating a new bill.
// It tests whether the handler correctly assigns an ID and returns a 201 status code.
func TestCreateBill(t *testing.T) {
	store := &MockPaymentStore{payments: make(map[int]*models.Payment)}
	handler := &AccountsPayableHandler{PaymentStore: store}

	// Create a sample payment request
	payment := models.Payment{
		InvoiceID:     123,
		Amount:        100.50,
		PaymentDate:   time.Now(),
		PaymentMethod: "credit_card",
	}
	body, _ := json.Marshal(payment)
	req, err := http.NewRequest("POST", "/accounts_payable", bytes.NewBuffer(body))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	// Record the response
	rr := httptest.NewRecorder()
	handler.CreateBill(rr, req)

	// Validate the response status code
	assert.Equal(t, http.StatusCreated, rr.Code)

	// Validate the response body
	var createdPayment models.Payment
	json.NewDecoder(rr.Body).Decode(&createdPayment)
	assert.Equal(t, 1, createdPayment.ID)
	assert.Equal(t, payment.InvoiceID, createdPayment.InvoiceID)
	assert.Equal(t, payment.Amount, createdPayment.Amount)
}

// TestGetBill verifies the GetBill handler for retrieving a bill by ID.
// It tests whether the handler correctly fetches the payment details and returns a 200 status code.
func TestGetBill(t *testing.T) {
	store := &MockPaymentStore{payments: make(map[int]*models.Payment)}
	handler := &AccountsPayableHandler{PaymentStore: store}

	// Add a payment to the mock store
	store.nextID = 1
	payment := &models.Payment{
		ID:            1,
		InvoiceID:     123,
		Amount:        100.50,
		PaymentDate:   time.Now(),
		PaymentMethod: "credit_card",
	}
	store.payments[1] = payment

	// Create a GET request
	req, err := http.NewRequest("GET", "/accounts_payable/1", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()

	// Attach the handler to a router and serve the request
	router := mux.NewRouter()
	router.HandleFunc("/accounts_payable/{id}", handler.GetBill).Methods("GET")
	router.ServeHTTP(rr, req)

	// Validate the response status code
	assert.Equal(t, http.StatusOK, rr.Code)

	// Validate the response body
	var gotPayment models.Payment
	json.NewDecoder(rr.Body).Decode(&gotPayment)
	assert.Equal(t, payment.ID, gotPayment.ID)
	assert.Equal(t, payment.InvoiceID, gotPayment.InvoiceID)
}

// TestUpdateBill verifies the UpdateBill handler for modifying an existing bill.
// It tests whether the handler updates the payment details and returns a 200 status code.
func TestUpdateBill(t *testing.T) {
	store := &MockPaymentStore{payments: make(map[int]*models.Payment)}
	handler := &AccountsPayableHandler{PaymentStore: store}

	// Add a payment to the mock store
	store.nextID = 1
	store.payments[1] = &models.Payment{
		ID:            1,
		InvoiceID:     123,
		Amount:        100.50,
		PaymentDate:   time.Now(),
		PaymentMethod: "credit_card",
	}

	// Create a request with updated payment details
	updatedPayment := models.Payment{
		InvoiceID:     123,
		Amount:        200.00,
		PaymentDate:   time.Now(),
		PaymentMethod: "bank_transfer",
	}
	body, _ := json.Marshal(updatedPayment)
	req, err := http.NewRequest("PUT", "/accounts_payable/1", bytes.NewBuffer(body))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	// Record the response
	rr := httptest.NewRecorder()

	// Attach the handler to a router and serve the request
	router := mux.NewRouter()
	router.HandleFunc("/accounts_payable/{id}", handler.UpdateBill).Methods("PUT")
	router.ServeHTTP(rr, req)

	// Validate the response status code
	assert.Equal(t, http.StatusOK, rr.Code)

	// Validate the updated payment details
	var gotPayment models.Payment
	json.NewDecoder(rr.Body).Decode(&gotPayment)
	assert.Equal(t, updatedPayment.Amount, gotPayment.Amount)
}

// TestDeleteBill verifies the DeleteBill handler for removing a bill by ID.
// It tests whether the handler deletes the payment and returns a 204 status code.
func TestDeleteBill(t *testing.T) {
	store := &MockPaymentStore{payments: make(map[int]*models.Payment)}
	handler := &AccountsPayableHandler{PaymentStore: store}

	// Add a payment to the mock store
	store.nextID = 1
	store.payments[1] = &models.Payment{
		ID:            1,
		InvoiceID:     123,
		Amount:        100.50,
		PaymentDate:   time.Now(),
		PaymentMethod: "credit_card",
	}

	// Create a DELETE request
	req, err := http.NewRequest("DELETE", "/accounts_payable/1", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()

	// Attach the handler to a router and serve the request
	router := mux.NewRouter()
	router.HandleFunc("/accounts_payable/{id}", handler.DeleteBill).Methods("DELETE")
	router.ServeHTTP(rr, req)

	// Validate the response status code
	assert.Equal(t, http.StatusNoContent, rr.Code)

	// Validate that the payment was deleted
	_, exists := store.payments[1]
	assert.False(t, exists)
}

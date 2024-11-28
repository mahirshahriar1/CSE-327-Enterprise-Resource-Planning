// Package accounts_payable_handlers provides handlers for CRUD operations
// on accounts payable resources. It includes test cases for validating
// the functionality of these handlers using a mock PaymentStore implementation.
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
// It simulates a storage layer using an in-memory map for testing purposes.
type MockPaymentStore struct {
	payments map[int]*models.Payment // In-memory storage for payments
	nextID   int                     // Counter for unique payment IDs
}

// CreatePayment adds a new payment to the mock store.
//
// Parameters:
//   - payment: A pointer to the Payment object to be stored.
//
// Returns:
//   - error: Always nil, as this is a simulated operation.
func (m *MockPaymentStore) CreatePayment(payment *models.Payment) error {
	m.nextID++
	payment.ID = m.nextID
	m.payments[payment.ID] = payment
	return nil
}

// GetPaymentByID retrieves a payment from the mock store by its ID.
//
// Parameters:
//   - id: The unique ID of the payment to retrieve.
//
// Returns:
//   - *Payment: Pointer to the retrieved payment, if found.
//   - error: "payment not found" if no payment exists with the given ID.
func (m *MockPaymentStore) GetPaymentByID(id int) (*models.Payment, error) {
	payment, exists := m.payments[id]
	if !exists {
		return nil, errors.New("payment not found")
	}
	return payment, nil
}

// UpdatePayment modifies an existing payment in the mock store.
//
// Parameters:
//   - payment: Pointer to the Payment object with updated details.
//
// Returns:
//   - error: "payment not found" if the payment ID does not exist in the store.
func (m *MockPaymentStore) UpdatePayment(payment *models.Payment) error {
	_, exists := m.payments[payment.ID]
	if !exists {
		return errors.New("payment not found")
	}
	m.payments[payment.ID] = payment
	return nil
}

// DeletePayment removes a payment from the mock store by its ID.
//
// Parameters:
//   - id: The unique ID of the payment to delete.
//
// Returns:
//   - error: "payment not found" if no payment exists with the given ID.
func (m *MockPaymentStore) DeletePayment(id int) error {
	_, exists := m.payments[id]
	if !exists {
		return errors.New("payment not found")
	}
	delete(m.payments, id)
	return nil
}

// TestCreateBill tests the CreateBill handler for adding a new payment.
//
// Steps:
//   - Creates a sample payment request.
//   - Simulates an HTTP POST request to the CreateBill endpoint.
//   - Validates the response status code and the returned payment details.
func TestCreateBill(t *testing.T) {
	store := &MockPaymentStore{payments: make(map[int]*models.Payment)}
	handler := &AccountsPayableHandler{PaymentStore: store}

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

	rr := httptest.NewRecorder()
	handler.CreateBill(rr, req)

	assert.Equal(t, http.StatusCreated, rr.Code)

	var createdPayment models.Payment
	json.NewDecoder(rr.Body).Decode(&createdPayment)
	assert.Equal(t, 1, createdPayment.ID)
	assert.Equal(t, payment.InvoiceID, createdPayment.InvoiceID)
	assert.Equal(t, payment.Amount, createdPayment.Amount)
}

// TestGetBill tests the GetBill handler for fetching a payment by ID.
//
// Steps:
//   - Adds a sample payment to the mock store.
//   - Simulates an HTTP GET request to the GetBill endpoint.
//   - Validates the response status code and the returned payment details.
func TestGetBill(t *testing.T) {
	store := &MockPaymentStore{payments: make(map[int]*models.Payment)}
	handler := &AccountsPayableHandler{PaymentStore: store}

	store.nextID = 1
	payment := &models.Payment{
		ID:            1,
		InvoiceID:     123,
		Amount:        100.50,
		PaymentDate:   time.Now(),
		PaymentMethod: "credit_card",
	}
	store.payments[1] = payment

	req, err := http.NewRequest("GET", "/accounts_payable/1", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()

	router := mux.NewRouter()
	router.HandleFunc("/accounts_payable/{id}", handler.GetBill).Methods("GET")
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	var gotPayment models.Payment
	json.NewDecoder(rr.Body).Decode(&gotPayment)
	assert.Equal(t, payment.ID, gotPayment.ID)
	assert.Equal(t, payment.InvoiceID, gotPayment.InvoiceID)
}

// TestUpdateBill tests the UpdateBill handler for modifying an existing payment.
//
// Steps:
//   - Adds a sample payment to the mock store.
//   - Simulates an HTTP PUT request to the UpdateBill endpoint with updated payment details.
//   - Validates the response status code and the updated payment details.
func TestUpdateBill(t *testing.T) {
	store := &MockPaymentStore{payments: make(map[int]*models.Payment)}
	handler := &AccountsPayableHandler{PaymentStore: store}

	store.nextID = 1
	store.payments[1] = &models.Payment{
		ID:            1,
		InvoiceID:     123,
		Amount:        100.50,
		PaymentDate:   time.Now(),
		PaymentMethod: "credit_card",
	}

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

	rr := httptest.NewRecorder()

	router := mux.NewRouter()
	router.HandleFunc("/accounts_payable/{id}", handler.UpdateBill).Methods("PUT")
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	var gotPayment models.Payment
	json.NewDecoder(rr.Body).Decode(&gotPayment)
	assert.Equal(t, updatedPayment.Amount, gotPayment.Amount)
}

// TestDeleteBill tests the DeleteBill handler for removing a payment by ID.
//
// Steps:
//   - Adds a sample payment to the mock store.
//   - Simulates an HTTP DELETE request to the DeleteBill endpoint.
//   - Validates the response status code and ensures the payment is deleted.
func TestDeleteBill(t *testing.T) {
	store := &MockPaymentStore{payments: make(map[int]*models.Payment)}
	handler := &AccountsPayableHandler{PaymentStore: store}

	store.nextID = 1
	store.payments[1] = &models.Payment{
		ID:            1,
		InvoiceID:     123,
		Amount:        100.50,
		PaymentDate:   time.Now(),
		PaymentMethod: "credit_card",
	}

	req, err := http.NewRequest("DELETE", "/accounts_payable/1", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()

	router := mux.NewRouter()
	router.HandleFunc("/accounts_payable/{id}", handler.DeleteBill).Methods("DELETE")
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusNoContent, rr.Code)

	_, exists := store.payments[1]
	assert.False(t, exists)
}

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

type MockPaymentStore struct {
	payments map[int]*models.Payment
	nextID   int
}

func (m *MockPaymentStore) CreatePayment(payment *models.Payment) error {
	m.nextID++
	payment.ID = m.nextID
	m.payments[payment.ID] = payment
	return nil
}

func (m *MockPaymentStore) GetPaymentByID(id int) (*models.Payment, error) {
	payment, exists := m.payments[id]
	if !exists {
		return nil, errors.New("payment not found")
	}
	return payment, nil
}

func (m *MockPaymentStore) UpdatePayment(payment *models.Payment) error {
	_, exists := m.payments[payment.ID]
	if !exists {
		return errors.New("payment not found")
	}
	m.payments[payment.ID] = payment
	return nil
}

func (m *MockPaymentStore) DeletePayment(id int) error {
	_, exists := m.payments[id]
	if !exists {
		return errors.New("payment not found")
	}
	delete(m.payments, id)
	return nil
}

func TestCreateBill(t *testing.T) {
	store := &MockPaymentStore{payments: make(map[int]*models.Payment)}
	handler := &AccountsPayableHandler{PaymentStore: store}

	// Create a request to pass to the handler
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

	// Check the response status code
	assert.Equal(t, http.StatusCreated, rr.Code)

	// Check the response body for ID assignment
	var createdPayment models.Payment
	json.NewDecoder(rr.Body).Decode(&createdPayment)
	assert.Equal(t, 1, createdPayment.ID)
	assert.Equal(t, payment.InvoiceID, createdPayment.InvoiceID)
	assert.Equal(t, payment.Amount, createdPayment.Amount)
}

func TestGetBill(t *testing.T) {
	store := &MockPaymentStore{payments: make(map[int]*models.Payment)}
	handler := &AccountsPayableHandler{PaymentStore: store}

	// Add a payment to the store for retrieval
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

	// Set path parameters
	router := mux.NewRouter()
	router.HandleFunc("/accounts_payable/{id}", handler.GetBill).Methods("GET")
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	var gotPayment models.Payment
	json.NewDecoder(rr.Body).Decode(&gotPayment)
	assert.Equal(t, payment.ID, gotPayment.ID)
	assert.Equal(t, payment.InvoiceID, gotPayment.InvoiceID)
}

func TestUpdateBill(t *testing.T) {
	store := &MockPaymentStore{payments: make(map[int]*models.Payment)}
	handler := &AccountsPayableHandler{PaymentStore: store}

	// Add a payment to the store to update
	store.nextID = 1
	store.payments[1] = &models.Payment{
		ID:            1,
		InvoiceID:     123,
		Amount:        100.50,
		PaymentDate:   time.Now(),
		PaymentMethod: "credit_card",
	}

	// Create an updated payment
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

	// Set path parameters
	router := mux.NewRouter()
	router.HandleFunc("/accounts_payable/{id}", handler.UpdateBill).Methods("PUT")
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	var gotPayment models.Payment
	json.NewDecoder(rr.Body).Decode(&gotPayment)
	assert.Equal(t, updatedPayment.Amount, gotPayment.Amount)
}

func TestDeleteBill(t *testing.T) {
	store := &MockPaymentStore{payments: make(map[int]*models.Payment)}
	handler := &AccountsPayableHandler{PaymentStore: store}

	// Add a payment to the store for deletion
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

	// Set path parameters
	router := mux.NewRouter()
	router.HandleFunc("/accounts_payable/{id}", handler.DeleteBill).Methods("DELETE")
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusNoContent, rr.Code)
	_, exists := store.payments[1]
	assert.False(t, exists)
}

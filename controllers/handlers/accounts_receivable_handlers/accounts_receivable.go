// Package accounts_receivable_handlers provides HTTP handlers for managing accounts receivable,
// including the creation, retrieval, updating, and deletion of payment records.
package accounts_receivable_handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"erp/models"

	"github.com/gorilla/mux"
)

// AccountsReceivableHandler handles HTTP requests for managing accounts receivable.
// It interacts with the ReceivableStore for receivable records and the TransactionStore
// for associated financial transactions.
type AccountsReceivableHandler struct {
	ReceivableStore  models.ReceivableStore          // Store for managing receivable records.
	TransactionStore models.FinancialTransactionStore // Store for managing related financial transactions.
}

// RegisterRoutes registers HTTP routes for accounts receivable handlers.
// It maps CRUD operations for payment records to the appropriate HTTP methods.
//
// Parameters:
//   - router: The Gorilla Mux router to which the routes are registered.
//   - receivableStore: The store interface for managing receivable records.
//   - transactionStore: The store interface for managing financial transactions.
func RegisterRoutes(router *mux.Router, receivableStore models.ReceivableStore, transactionStore models.FinancialTransactionStore) {
	handler := &AccountsReceivableHandler{ReceivableStore: receivableStore, TransactionStore: transactionStore}

	router.HandleFunc("", handler.CreatePayment).Methods("POST")
	router.HandleFunc("/{id}", handler.GetPayment).Methods("GET")
	router.HandleFunc("/{id}", handler.UpdatePayment).Methods("PUT")
	router.HandleFunc("/{id}", handler.DeletePayment).Methods("DELETE")
}

// CreatePayment creates a new payment record and stores it in the accounts receivable system.
//
// HTTP Method: POST
// URL Path: / (root path of accounts receivable routes)
//
// Request Body:
//   - JSON representation of a Receivable object (excluding the payment date).
//
// Response:
//   - Status Code: 201 (Created) if the payment is successfully created.
//   - JSON representation of the created payment on success.
//   - Status Code: 400 (Bad Request) if the input data is invalid.
//   - Status Code: 500 (Internal Server Error) if the payment could not be saved.
func (h *AccountsReceivableHandler) CreatePayment(w http.ResponseWriter, r *http.Request) {
	var receivable models.Receivable
	if err := json.NewDecoder(r.Body).Decode(&receivable); err != nil {
		http.Error(w, "Invalid input data", http.StatusBadRequest)
		return
	}

	receivable.PaymentDate = time.Now()
	if err := h.ReceivableStore.CreateReceivable(&receivable); err != nil {
		http.Error(w, fmt.Sprintf("Failed to create payment: %v", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(receivable); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}

// GetPayment retrieves a payment record by its ID.
//
// HTTP Method: GET
// URL Path: /{id} (ID of the payment in the path)
//
// Response:
//   - Status Code: 200 (OK) with the payment data in JSON format if found.
//   - Status Code: 400 (Bad Request) if the ID is invalid.
//   - Status Code: 404 (Not Found) if the payment with the specified ID does not exist.
func (h *AccountsReceivableHandler) GetPayment(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		http.Error(w, "Invalid payment ID", http.StatusBadRequest)
		return
	}

	payment, err := h.ReceivableStore.GetReceivableByID(id)
	if err != nil {
		http.Error(w, fmt.Sprintf("Payment not found: %v", err), http.StatusNotFound)
		return
	}

	if err := json.NewEncoder(w).Encode(payment); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}

// UpdatePayment updates an existing payment record with new data.
//
// HTTP Method: PUT
// URL Path: /{id} (ID of the payment in the path)
//
// Request Body:
//   - JSON representation of a Receivable object (excluding the ID, which is taken from the URL).
//
// Response:
//   - Status Code: 200 (OK) with the updated payment data in JSON format if successful.
//   - Status Code: 400 (Bad Request) if the ID or input data is invalid.
//   - Status Code: 500 (Internal Server Error) if the update operation fails.
func (h *AccountsReceivableHandler) UpdatePayment(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		http.Error(w, "Invalid payment ID", http.StatusBadRequest)
		return
	}

	var receivable models.Receivable
	if err := json.NewDecoder(r.Body).Decode(&receivable); err != nil {
		http.Error(w, "Invalid input data", http.StatusBadRequest)
		return
	}

	receivable.ID = id
	if err := h.ReceivableStore.UpdateReceivable(&receivable); err != nil {
		http.Error(w, fmt.Sprintf("Failed to update payment: %v", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(receivable); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}

// DeletePayment deletes a payment record identified by its ID.
//
// HTTP Method: DELETE
// URL Path: /{id} (ID of the payment in the path)
//
// Response:
//   - Status Code: 204 (No Content) if the payment is successfully deleted.
//   - Status Code: 400 (Bad Request) if the ID is invalid.
//   - Status Code: 500 (Internal Server Error) if the deletion operation fails.
func (h *AccountsReceivableHandler) DeletePayment(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		http.Error(w, "Invalid payment ID", http.StatusBadRequest)
		return
	}

	if err := h.ReceivableStore.DeleteReceivable(id); err != nil {
		http.Error(w, fmt.Sprintf("Failed to delete payment: %v", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// Package accounts_receivable_handlers provides HTTP handlers for managing accounts receivable,
// including managing payments received from customers.
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

// AccountsReceivableHandler struct handles HTTP requests for managing receivables.
type AccountsReceivableHandler struct {
	ReceivableStore     models.ReceivableStore                // Store for managing receivable records.
	TransactionStore    models.FinancialTransactionStore     // Store for associated financial transactions.
}

// RegisterRoutes registers the routes for accounts receivable handlers.
func RegisterRoutes(router *mux.Router, receivableStore models.ReceivableStore, transactionStore models.FinancialTransactionStore) {
	handler := &AccountsReceivableHandler{ReceivableStore: receivableStore, TransactionStore: transactionStore}

	router.HandleFunc("", handler.CreatePayment).Methods("POST")
	router.HandleFunc("/{id}", handler.GetPayment).Methods("GET")
	router.HandleFunc("/{id}", handler.UpdatePayment).Methods("PUT")
	router.HandleFunc("/{id}", handler.DeletePayment).Methods("DELETE")
}

// CreatePayment creates a new payment record.
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

// GetPayment retrieves a payment by ID.
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

// UpdatePayment updates an existing payment record.
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

// DeletePayment deletes a payment record by ID.
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

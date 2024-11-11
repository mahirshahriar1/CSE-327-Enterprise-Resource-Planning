// Package general_ledger_handlers provides HTTP handlers to perform CRUD operations
// on the general ledger within an ERP system.
package general_ledger_handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"erp/models"

	"github.com/gorilla/mux"
)

// GeneralLedgerHandler struct provides HTTP handlers for interacting with financial
// transactions stored in the general ledger. It uses a FinancialTransactionStore
// interface to perform data storage operations.
type GeneralLedgerHandler struct {
	Store models.FinancialTransactionStore
}

// RegisterRoutes maps general ledger routes to their respective handler functions.
// It keeps route definitions modular, enabling easy modifications and scalability.
// The routes are registered with the provided router and are bound to the specified store.
func RegisterRoutes(router *mux.Router, store models.FinancialTransactionStore) {
	handler := &GeneralLedgerHandler{Store: store}

	router.HandleFunc("", handler.CreateTransaction).Methods("POST")
	router.HandleFunc("/{id}", handler.GetTransaction).Methods("GET")
	router.HandleFunc("/{id}", handler.UpdateTransaction).Methods("PUT")
	router.HandleFunc("/{id}", handler.DeleteTransaction).Methods("DELETE")
}

// CreateTransaction is an HTTP handler that creates a new financial transaction
// in the general ledger. It reads transaction data from the request body, assigns
// the current time as the transaction date, and saves it to the database.
// On success, it returns the created transaction in JSON format with a 201 status code.
func (h *GeneralLedgerHandler) CreateTransaction(w http.ResponseWriter, r *http.Request) {
	var transaction models.FinancialTransaction
	if err := json.NewDecoder(r.Body).Decode(&transaction); err != nil {
		http.Error(w, "Invalid input data", http.StatusBadRequest)
		return
	}

	transaction.TransactionDate = time.Now()
	if err := h.Store.CreateTransaction(&transaction); err != nil {
		http.Error(w, fmt.Sprintf("Failed to create transaction: %v", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(transaction); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}

// GetTransaction retrieves and returns a financial transaction by its ID.
// The ID is parsed from the URL path, and the transaction data is fetched from the
// database. If found, it returns the transaction in JSON format with a 200 status.
// If not found or if the ID is invalid, it returns an appropriate error message.
func (h *GeneralLedgerHandler) GetTransaction(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		http.Error(w, "Invalid transaction ID", http.StatusBadRequest)
		return
	}

	transaction, err := h.Store.GetTransactionByID(id)
	if err != nil {
		http.Error(w, fmt.Sprintf("Transaction not found: %v", err), http.StatusNotFound)
		return
	}

	if err := json.NewEncoder(w).Encode(transaction); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}

// UpdateTransaction modifies an existing financial transaction based on the given ID.
// It reads the updated transaction data from the request body, updates the record in the
// database, and returns the updated transaction in JSON format with a 200 status.
// If the ID is invalid or the update fails, it returns an error message.
func (h *GeneralLedgerHandler) UpdateTransaction(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		http.Error(w, "Invalid transaction ID", http.StatusBadRequest)
		return
	}

	var transaction models.FinancialTransaction
	if err := json.NewDecoder(r.Body).Decode(&transaction); err != nil {
		http.Error(w, "Invalid input data", http.StatusBadRequest)
		return
	}

	transaction.ID = id
	if err := h.Store.UpdateTransaction(&transaction); err != nil {
		http.Error(w, fmt.Sprintf("Failed to update transaction: %v", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(transaction); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}

// DeleteTransaction removes a financial transaction identified by its ID.
// The ID is extracted from the URL path, and the transaction is deleted from the database.
// On successful deletion, it returns a 204 status with no content. If the ID is invalid
// or deletion fails, it returns an appropriate error message.
func (h *GeneralLedgerHandler) DeleteTransaction(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		http.Error(w, "Invalid transaction ID", http.StatusBadRequest)
		return
	}

	if err := h.Store.DeleteTransaction(id); err != nil {
		http.Error(w, fmt.Sprintf("Failed to delete transaction: %v", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

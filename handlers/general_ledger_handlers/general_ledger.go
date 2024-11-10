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

// GeneralLedgerHandler provides HTTP handlers for general ledger operations
type GeneralLedgerHandler struct {
	Store models.FinancialTransactionStore
}

// RegisterRoutes maps general ledger routes to their respective handlers
// This function is called during route initialization to keep route definitions modular and scalable.
func RegisterRoutes(router *mux.Router, store models.FinancialTransactionStore) {
	handler := &GeneralLedgerHandler{Store: store}

	router.HandleFunc("", handler.createTransaction).Methods("POST")
	router.HandleFunc("/{id}", handler.getTransaction).Methods("GET")
	router.HandleFunc("/{id}", handler.updateTransaction).Methods("PUT")
	router.HandleFunc("/{id}", handler.deleteTransaction).Methods("DELETE")
}

// createTransaction handles the creation of a new financial transaction
func (h *GeneralLedgerHandler) createTransaction(w http.ResponseWriter, r *http.Request) {
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

// getTransaction retrieves a financial transaction by ID
func (h *GeneralLedgerHandler) getTransaction(w http.ResponseWriter, r *http.Request) {
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

// updateTransaction modifies an existing financial transaction
func (h *GeneralLedgerHandler) updateTransaction(w http.ResponseWriter, r *http.Request) {
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

// deleteTransaction removes a financial transaction by ID
func (h *GeneralLedgerHandler) deleteTransaction(w http.ResponseWriter, r *http.Request) {
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

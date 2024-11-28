// Package general_ledger_handlers provides HTTP handlers to perform CRUD operations
// on the general ledger within an ERP system. It defines routes for managing financial transactions
// stored in the general ledger and interacts with a database through the FinancialTransactionStore interface.
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
	Store models.FinancialTransactionStore // Store defines the interface for managing transactions in the database.
}

// RegisterRoutes maps general ledger routes to their respective handler functions.
// It keeps route definitions modular, enabling easy modifications and scalability.
//
// Parameters:
//   - router: The HTTP router (from the Gorilla Mux library) where the routes are registered.
//   - store: An implementation of the FinancialTransactionStore interface for managing transaction data.
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
//
// HTTP Method: POST
// URL Path: / (root path of general ledger routes)
//
// Request Body:
//   - JSON representation of a FinancialTransaction object (excluding the transaction date).
//
// Response:
//   - Status Code: 201 (Created) if the transaction is successfully created.
//   - JSON representation of the created transaction on success.
//   - Status Code: 400 (Bad Request) if the input data is invalid.
//   - Status Code: 500 (Internal Server Error) if the transaction could not be saved.
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
// database.
//
// HTTP Method: GET
// URL Path: /{id} (ID of the transaction in the path)
//
// Response:
//   - Status Code: 200 (OK) with the transaction data in JSON format if found.
//   - Status Code: 400 (Bad Request) if the ID is invalid.
//   - Status Code: 404 (Not Found) if the transaction with the specified ID does not exist.
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
// database, and returns the updated transaction in JSON format.
//
// HTTP Method: PUT
// URL Path: /{id} (ID of the transaction in the path)
//
// Request Body:
//   - JSON representation of a FinancialTransaction object (excluding the ID, which is taken from the URL).
//
// Response:
//   - Status Code: 200 (OK) with the updated transaction data in JSON format if successful.
//   - Status Code: 400 (Bad Request) if the ID or input data is invalid.
//   - Status Code: 500 (Internal Server Error) if the update operation fails.
func (h *GeneralLedgerHandler) UpdateTransaction(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		http.Error(w, "Invalid transaction ID", http.StatusBadRequest)
		return
	}
	fmt.Println("ID: ", id)

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
//
// HTTP Method: DELETE
// URL Path: /{id} (ID of the transaction in the path)
//
// Response:
//   - Status Code: 204 (No Content) if the transaction is successfully deleted.
//   - Status Code: 400 (Bad Request) if the ID is invalid.
//   - Status Code: 500 (Internal Server Error) if the deletion operation fails.
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

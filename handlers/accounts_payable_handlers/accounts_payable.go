// Package accounts_payable_handlers provides HTTP handlers to perform CRUD operations
// on accounts payable within an ERP system. These handlers allow management of payable bills and
// transactions in the system.
package accounts_payable_handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"erp/models"

	"github.com/gorilla/mux"
)

// AccountsPayableHandler struct provides HTTP handlers for managing accounts payable.
// It interacts with the PaymentStore to manage bills and the FinancialTransactionStore
// for related financial transactions.
type AccountsPayableHandler struct {
	PaymentStore     models.PaymentStore                // PaymentStore manages payable bill records.
	TransactionStore models.FinancialTransactionStore // TransactionStore manages associated financial transactions.
}

// RegisterRoutes maps accounts payable routes to their respective handler functions.
// This function registers the routes to the provided router, associating them with
// the appropriate handler methods.
//
// Parameters:
//   - router: The HTTP router (from the Gorilla Mux library) to which the routes are registered.
//   - paymentStore: An implementation of the PaymentStore interface for managing payments.
//   - transactionStore: An implementation of the FinancialTransactionStore interface for managing transactions.
func RegisterRoutes(router *mux.Router, paymentStore models.PaymentStore, transactionStore models.FinancialTransactionStore) {
	handler := &AccountsPayableHandler{PaymentStore: paymentStore, TransactionStore: transactionStore}

	router.HandleFunc("", handler.CreateBill).Methods("POST")
	router.HandleFunc("/{id}", handler.GetBill).Methods("GET")
	router.HandleFunc("/{id}", handler.UpdateBill).Methods("PUT")
	router.HandleFunc("/{id}", handler.DeleteBill).Methods("DELETE")
}

// CreateBill creates a new payable bill entry in the system. The bill data is extracted
// from the request body, and the current time is assigned as the payment date before
// saving it to the database.
//
// HTTP Method: POST
// URL Path: / (root path of accounts payable routes)
//
// Request Body:
//   - JSON representation of a Payment object.
//
// Response:
//   - Status Code: 201 (Created) with the created bill in JSON format.
//   - Status Code: 400 (Bad Request) if the input data is invalid.
//   - Status Code: 500 (Internal Server Error) if the bill creation fails.
func (h *AccountsPayableHandler) CreateBill(w http.ResponseWriter, r *http.Request) {
	var payment models.Payment
	if err := json.NewDecoder(r.Body).Decode(&payment); err != nil {
		http.Error(w, "Invalid input data", http.StatusBadRequest)
		return
	}

	payment.PaymentDate = time.Now() // Set the payment date to the current time.
	if err := h.PaymentStore.CreatePayment(&payment); err != nil {
		http.Error(w, fmt.Sprintf("Failed to create payment: %v", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(payment); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}

// GetBill retrieves and returns a bill by its ID. The ID is parsed from the URL path,
// and the bill is fetched from the database.
//
// HTTP Method: GET
// URL Path: /{id} (ID of the bill in the path)
//
// Response:
//   - Status Code: 200 (OK) with the bill in JSON format if found.
//   - Status Code: 400 (Bad Request) if the ID is invalid.
//   - Status Code: 404 (Not Found) if the bill does not exist.
func (h *AccountsPayableHandler) GetBill(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		http.Error(w, "Invalid bill ID", http.StatusBadRequest)
		return
	}

	bill, err := h.PaymentStore.GetPaymentByID(id)
	if err != nil {
		http.Error(w, fmt.Sprintf("Bill not found: %v", err), http.StatusNotFound)
		return
	}

	if err := json.NewEncoder(w).Encode(bill); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}

// UpdateBill modifies an existing bill's details. The updated bill data is read from the
// request body and the corresponding record is updated in the database.
//
// HTTP Method: PUT
// URL Path: /{id} (ID of the bill in the path)
//
// Request Body:
//   - JSON representation of a Payment object (excluding the ID, which is taken from the URL).
//
// Response:
//   - Status Code: 200 (OK) with the updated bill in JSON format.
//   - Status Code: 400 (Bad Request) if the ID or input data is invalid.
//   - Status Code: 500 (Internal Server Error) if the update operation fails.
func (h *AccountsPayableHandler) UpdateBill(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		http.Error(w, "Invalid bill ID", http.StatusBadRequest)
		return
	}

	var payment models.Payment
	if err := json.NewDecoder(r.Body).Decode(&payment); err != nil {
		http.Error(w, "Invalid input data", http.StatusBadRequest)
		return
	}

	payment.ID = id
	if err := h.PaymentStore.UpdatePayment(&payment); err != nil {
		http.Error(w, fmt.Sprintf("Failed to update bill: %v", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(payment); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}

// DeleteBill deletes a bill by its ID. The ID is extracted from the URL path, and the
// corresponding bill record is deleted from the database.
//
// HTTP Method: DELETE
// URL Path: /{id} (ID of the bill in the path)
//
// Response:
//   - Status Code: 204 (No Content) if the bill is successfully deleted.
//   - Status Code: 400 (Bad Request) if the ID is invalid.
//   - Status Code: 500 (Internal Server Error) if the deletion operation fails.
func (h *AccountsPayableHandler) DeleteBill(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		http.Error(w, "Invalid bill ID", http.StatusBadRequest)
		return
	}

	if err := h.PaymentStore.DeletePayment(id); err != nil {
		http.Error(w, fmt.Sprintf("Failed to delete bill: %v", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

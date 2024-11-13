// Package accounts_payable_handlers provides HTTP handlers to perform CRUD operations
// on accounts payable within an ERP system.
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
type AccountsPayableHandler struct {
    PaymentStore       models.PaymentStore
    TransactionStore   models.FinancialTransactionStore
}

// RegisterRoutes maps accounts payable routes to their respective handler functions.
func RegisterRoutes(router *mux.Router, paymentStore models.PaymentStore, transactionStore models.FinancialTransactionStore) {
    handler := &AccountsPayableHandler{PaymentStore: paymentStore, TransactionStore: transactionStore}

    router.HandleFunc("", handler.CreateBill).Methods("POST")
    router.HandleFunc("/{id}", handler.GetBill).Methods("GET")
    router.HandleFunc("/{id}", handler.UpdateBill).Methods("PUT")
    router.HandleFunc("/{id}", handler.DeleteBill).Methods("DELETE")
}

// CreateBill creates a new payable bill entry.
func (h *AccountsPayableHandler) CreateBill(w http.ResponseWriter, r *http.Request) {
    var payment models.Payment
    if err := json.NewDecoder(r.Body).Decode(&payment); err != nil {
        http.Error(w, "Invalid input data", http.StatusBadRequest)
        return
    }

    payment.PaymentDate = time.Now() // Set payment date
    if err := h.PaymentStore.CreatePayment(&payment); err != nil {
        http.Error(w, fmt.Sprintf("Failed to create payment: %v", err), http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusCreated)
    if err := json.NewEncoder(w).Encode(payment); err != nil {
        http.Error(w, "Failed to encode response", http.StatusInternalServerError)
    }
}

// GetBill retrieves a bill by its ID.
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

// UpdateBill modifies an existing bill's details.
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

// DeleteBill deletes a bill by its ID.
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

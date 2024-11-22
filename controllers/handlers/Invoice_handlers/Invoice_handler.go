package invoice_handlers

import (
    "encoding/json"
    "erp/models"
    "net/http"
    "strconv"
    "github.com/gorilla/mux"
)

// InvoiceHandlers struct to hold the invoice store dependency
type InvoiceHandlers struct {
    Store models.InvoiceStore // Interface for interacting with the invoice store
}

// CreateInvoiceHandler handles the creation of a new invoice
func (h *InvoiceHandlers) CreateInvoiceHandler(w http.ResponseWriter, r *http.Request) {
    var invoice models.Invoice

    // Decode the JSON body into the invoice struct
    err := json.NewDecoder(r.Body).Decode(&invoice)
    if err != nil {
        http.Error(w, "Invalid request payload", http.StatusBadRequest)
        return
    }

    // Call the store to create the invoice
    err = h.Store.CreateInvoice(&invoice)
    if err != nil {
        http.Error(w, "Failed to create invoice", http.StatusInternalServerError)
        return
    }

    // Respond with the created invoice
    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(invoice)
}

// GetInvoiceByIDHandler handles the retrieval of an invoice by ID
func (h *InvoiceHandlers) GetInvoiceByIDHandler(w http.ResponseWriter, r *http.Request) {
    // Extract the "id" from the URL
    vars := mux.Vars(r)
    id, err := strconv.Atoi(vars["id"])
    if err != nil {
        http.Error(w, "Invalid invoice ID", http.StatusBadRequest)
        return
    }

    // Call the store to fetch the invoice by ID
    invoice, err := h.Store.GetInvoiceByID(id)
    if err != nil {
        http.Error(w, "Invoice not found", http.StatusNotFound)
        return
    }

    // Respond with the fetched invoice
    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(invoice)
}





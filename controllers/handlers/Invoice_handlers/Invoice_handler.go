// Package invoice_handlers provides HTTP handlers for managing invoice data.
// It includes functionality for creating, retrieving, updating, and deleting invoices.
package invoice_handlers

import (
	"encoding/json"
	"erp/models"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

// InvoiceHandlers is a struct that provides methods to handle invoice-related HTTP requests.
// It interacts with a data store through the InvoiceStore interface.
type InvoiceHandlers struct {
	Store models.InvoiceStore // Interface for interacting with the invoice data store
}

// CreateInvoiceHandler handles HTTP POST requests for creating a new invoice.
//
// Request Body:
//   - JSON object representing an invoice.
//
// Response:
//   - 201 Created: If the invoice is successfully created, returns the invoice object as JSON.
//   - 400 Bad Request: If the request payload is invalid.
//   - 500 Internal Server Error: If an error occurs while creating the invoice.
func (h *InvoiceHandlers) CreateInvoiceHandler(w http.ResponseWriter, r *http.Request) {
	var invoice models.Invoice

	// Decode JSON body into the invoice struct
	err := json.NewDecoder(r.Body).Decode(&invoice)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Create the invoice in the database
	err = h.Store.CreateInvoice(&invoice)
	if err != nil {
		http.Error(w, "Failed to create invoice", http.StatusInternalServerError)
		return
	}

	// Respond with the created invoice object
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(invoice)
}

// GetInvoiceByIDHandler handles HTTP GET requests to fetch an invoice by its ID.
//
// URL Parameters:
//   - id: Invoice ID (integer).
//
// Response:
//   - 200 OK: Returns the invoice object as JSON.
//   - 400 Bad Request: If the provided ID is invalid.
//   - 404 Not Found: If no invoice with the given ID exists.
func (h *InvoiceHandlers) GetInvoiceByIDHandler(w http.ResponseWriter, r *http.Request) {
	// Extract the "id" variable from the URL
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid invoice ID", http.StatusBadRequest)
		return
	}

	// Fetch the invoice by ID
	invoice, err := h.Store.GetInvoiceByID(id)
	if err != nil {
		http.Error(w, "Invoice not found", http.StatusNotFound)
		return
	}

	// Respond with the invoice object
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(invoice)
}

// UpdateInvoiceHandler handles HTTP PUT requests to update an existing invoice.
//
// URL Parameters:
//   - id: Invoice ID (integer).
//
// Request Body:
//   - JSON object representing the updated invoice data.
//
// Response:
//   - 200 OK: If the update is successful, returns the updated invoice object as JSON.
//   - 400 Bad Request: If the ID is invalid or the request payload is malformed.
//   - 500 Internal Server Error: If an error occurs while updating the invoice.
func (h *InvoiceHandlers) UpdateInvoiceHandler(w http.ResponseWriter, r *http.Request) {
	// Extract the "id" variable from the URL
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid invoice ID", http.StatusBadRequest)
		return
	}

	var invoice models.Invoice
	// Decode JSON body into the invoice struct
	err = json.NewDecoder(r.Body).Decode(&invoice)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Ensure the invoice ID matches the URL parameter
	invoice.ID = id

	// Update the invoice data in the store
	err = h.Store.UpdateInvoice(&invoice)
	if err != nil {
		http.Error(w, "Failed to update invoice", http.StatusInternalServerError)
		return
	}

	// Respond with the updated invoice object
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(invoice)
}

// DeleteInvoiceHandler handles HTTP DELETE requests to remove an invoice by its ID.
//
// URL Parameters:
//   - id: Invoice ID (integer).
//
// Response:
//   - 204 No Content: If the deletion is successful.
//   - 400 Bad Request: If the provided ID is invalid.
//   - 500 Internal Server Error: If an error occurs while deleting the invoice.
func (h *InvoiceHandlers) DeleteInvoiceHandler(w http.ResponseWriter, r *http.Request) {
	// Extract the "id" variable from the URL
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid invoice ID", http.StatusBadRequest)
		return
	}

	// Delete the invoice by ID
	err = h.Store.DeleteInvoice(id)
	if err != nil {
		http.Error(w, "Failed to delete invoice", http.StatusInternalServerError)
		return
	}

	// Respond with no content
	w.WriteHeader(http.StatusNoContent)
}

// Package customer_data_management_handlers provides HTTP handlers for managing customer data.
// It includes functionality for creating, retrieving, updating, and deleting customer records.
package customer_data_management_handlers

import (
	"encoding/json"
	"erp/models"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

// CustomerHandlers is a struct that provides methods to handle customer-related HTTP requests.
// It interacts with a data store through the CustomerStore interface.
type CustomerHandlers struct {
	Store models.CustomerStore // Interface for interacting with the customer data store
}

// CreateCustomerHandler handles HTTP POST requests for creating a new customer.
//
// Request Body:
//   - JSON object representing a customer.
//
// Response:
//   - 201 Created: If the customer is successfully created, returns the customer object as JSON.
//   - 400 Bad Request: If the request payload is invalid.
//   - 500 Internal Server Error: If an error occurs while creating the customer.
func (h *CustomerHandlers) CreateCustomerHandler(w http.ResponseWriter, r *http.Request) {
	var customer models.Customer

	// Decode JSON body into the customer struct
	err := json.NewDecoder(r.Body).Decode(&customer)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Create the customer in the database
	err = h.Store.CreateCustomer(&customer)
	if err != nil {
		http.Error(w, "Failed to create customer", http.StatusInternalServerError)
		return
	}

	// Respond with the created customer object
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(customer)
}

// GetCustomerByIDHandler handles HTTP GET requests to fetch a customer by their ID.
//
// URL Parameters:
//   - id: Customer ID (integer).
//
// Response:
//   - 200 OK: Returns the customer object as JSON.
//   - 400 Bad Request: If the provided ID is invalid.
//   - 404 Not Found: If no customer with the given ID exists.
func (h *CustomerHandlers) GetCustomerByIDHandler(w http.ResponseWriter, r *http.Request) {
	// Extract the "id" variable from the URL
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid customer ID", http.StatusBadRequest)
		return
	}

	// Fetch the customer by ID
	customer, err := h.Store.GetCustomerByID(id)
	if err != nil {
		http.Error(w, "Customer not found", http.StatusNotFound)
		return
	}

	// Respond with the customer object
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(customer)
}

// UpdateCustomerHandler handles HTTP PUT requests to update an existing customer's data.
//
// URL Parameters:
//   - id: Customer ID (integer).
//
// Request Body:
//   - JSON object representing the updated customer data.
//
// Response:
//   - 200 OK: If the update is successful, returns the updated customer object as JSON.
//   - 400 Bad Request: If the ID is invalid or the request payload is malformed.
//   - 500 Internal Server Error: If an error occurs while updating the customer.
func (h *CustomerHandlers) UpdateCustomerHandler(w http.ResponseWriter, r *http.Request) {
	// Extract the "id" variable from the URL
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid customer ID", http.StatusBadRequest)
		return
	}

	var customer models.Customer
	// Decode JSON body into the customer struct
	err = json.NewDecoder(r.Body).Decode(&customer)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Ensure the customer ID matches the URL parameter
	customer.ID = id

	// Update the customer data in the store
	err = h.Store.UpdateCustomer(&customer)
	if err != nil {
		http.Error(w, "Failed to update customer", http.StatusInternalServerError)
		return
	}

	// Respond with the updated customer object
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(customer)
}

// DeleteCustomerHandler handles HTTP DELETE requests to remove a customer by their ID.
//
// URL Parameters:
//   - id: Customer ID (integer).
//
// Response:
//   - 204 No Content: If the deletion is successful.
//   - 400 Bad Request: If the provided ID is invalid.
//   - 500 Internal Server Error: If an error occurs while deleting the customer.
func (h *CustomerHandlers) DeleteCustomerHandler(w http.ResponseWriter, r *http.Request) {
	// Extract the "id" variable from the URL
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid customer ID", http.StatusBadRequest)
		return
	}

	// Delete the customer by ID
	err = h.Store.DeleteCustomer(id)
	if err != nil {
		http.Error(w, "Failed to delete customer", http.StatusInternalServerError)
		return
	}

	// Respond with no content
	w.WriteHeader(http.StatusNoContent)
}

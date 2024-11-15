package customer_data_management_handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"erp/models"
	"github.com/gorilla/mux"
)

// CustomerHandlers struct to handle customer requests
type CustomerHandlers struct {
	Store models.CustomerStore
}

// CreateCustomerHandler handles customer creation requests.
func (h *CustomerHandlers) CreateCustomerHandler(w http.ResponseWriter, r *http.Request) {
	var customer models.Customer
	err := json.NewDecoder(r.Body).Decode(&customer)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	err = h.Store.CreateCustomer(&customer)
	if err != nil {
		http.Error(w, "Failed to create customer", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(customer)
}

// GetCustomerByIDHandler handles requests to fetch a customer by ID.
func (h *CustomerHandlers) GetCustomerByIDHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid customer ID", http.StatusBadRequest)
		return
	}

	customer, err := h.Store.GetCustomerByID(id)
	if err != nil {
		http.Error(w, "Customer not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(customer)
}

// UpdateCustomerHandler handles requests to update an existing customer.
func (h *CustomerHandlers) UpdateCustomerHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid customer ID", http.StatusBadRequest)
		return
	}

	var customer models.Customer
	err = json.NewDecoder(r.Body).Decode(&customer)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}
	customer.ID = id // Ensure the customer ID matches the URL parameter

	err = h.Store.UpdateCustomer(&customer)
	if err != nil {
		http.Error(w, "Failed to update customer", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(customer)
}

// DeleteCustomerHandler handles requests to delete a customer by ID.
func (h *CustomerHandlers) DeleteCustomerHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid customer ID", http.StatusBadRequest)
		return
	}

	err = h.Store.DeleteCustomer(id)
	if err != nil {
		http.Error(w, "Failed to delete customer", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}


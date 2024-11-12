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

// Package warehouse_handlers contains the handlers for warehouse-related HTTP endpoints.
package warehouse_handlers

import (
	"encoding/json"
	"erp/models"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

// WarehouseHandlers contains dependencies for handling warehouse-related requests.
type WarehouseHandlers struct {
	WarehouseStore models.WarehouseStore
}

// RegisterRoutes registers all the warehouse routes for the HTTP server.
func (h *WarehouseHandlers) RegisterRoutes(router *mux.Router) {
	// Register route for creating a new warehouse
	router.HandleFunc("/warehouses", h.CreateWarehouse).Methods("POST")

	// Register route for fetching a warehouse by ID
	router.HandleFunc("/warehouses/{id:[0-9]+}", h.GetWarehouseByID).Methods("GET")

	// Register route for updating an existing warehouse
	router.HandleFunc("/warehouses/{id:[0-9]+}", h.UpdateWarehouse).Methods("PUT")

	// Register route for deleting a warehouse by ID
	router.HandleFunc("/warehouses/{id:[0-9]+}", h.DeleteWarehouse).Methods("DELETE")
}

// CreateWarehouse handles the creation of a new warehouse.
// It decodes the request body into a Warehouse struct and stores it in the database.
func (h *WarehouseHandlers) CreateWarehouse(w http.ResponseWriter, r *http.Request) {
	var req models.Warehouse
	// Decode the JSON body into the Warehouse struct
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	// Attempt to create the warehouse in the database
	err = h.WarehouseStore.CreateWarehouse(&req)
	if err != nil {
		http.Error(w, "Could not create warehouse", http.StatusInternalServerError)
		return
	}

	// Return success response
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("Warehouse created successfully"))
}

// GetWarehouseByID handles retrieving a warehouse by its ID.
// It returns the warehouse details in JSON format if found, or an error if not.
func (h *WarehouseHandlers) GetWarehouseByID(w http.ResponseWriter, r *http.Request) {
	// Extract warehouse ID from the URL parameters
	params := mux.Vars(r)
	warehouseID, err := strconv.Atoi(params["id"])
	if err != nil {
		http.Error(w, "Invalid warehouse ID", http.StatusBadRequest)
		return
	}

	// Retrieve the warehouse from the database
	warehouse, err := h.WarehouseStore.GetWarehouseByID(warehouseID)
	if err != nil {
		http.Error(w, "Warehouse not found", http.StatusNotFound)
		return
	}

	// Respond with the warehouse in JSON format
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(warehouse)
}

// UpdateWarehouse handles updating an existing warehouse by ID.
// It decodes the request body, updates the warehouse in the database, and returns a success response.
func (h *WarehouseHandlers) UpdateWarehouse(w http.ResponseWriter, r *http.Request) {
	// Extract warehouse ID from the URL parameters
	params := mux.Vars(r)
	warehouseID, err := strconv.Atoi(params["id"])
	if err != nil {
		http.Error(w, "Invalid warehouse ID", http.StatusBadRequest)
		return
	}

	var req models.Warehouse
	// Decode the request body into the Warehouse struct
	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	// Set the warehouse ID for the update operation
	req.ID = warehouseID

	// Attempt to update the warehouse in the database
	err = h.WarehouseStore.UpdateWarehouse(&req)
	if err != nil {
		http.Error(w, "Could not update warehouse", http.StatusInternalServerError)
		return
	}

	// Return success response
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Warehouse updated successfully"))
}

// DeleteWarehouse handles deleting a warehouse by ID.
// It deletes the warehouse from the database and returns a success message.
func (h *WarehouseHandlers) DeleteWarehouse(w http.ResponseWriter, r *http.Request) {
	// Extract warehouse ID from the URL parameters
	params := mux.Vars(r)
	warehouseID, err := strconv.Atoi(params["id"])
	if err != nil {
		http.Error(w, "Invalid warehouse ID", http.StatusBadRequest)
		return
	}

	// Attempt to delete the warehouse from the database
	err = h.WarehouseStore.DeleteWarehouse(warehouseID)
	if err != nil {
		http.Error(w, "Could not delete warehouse", http.StatusInternalServerError)
		return
	}

	// Return success response
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Warehouse deleted successfully"))
}

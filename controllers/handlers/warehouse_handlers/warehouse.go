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

// RegisterRoutes registers all the warehouse-related routes for the HTTP server.
//
// This method sets up routes for creating, retrieving, updating, and deleting warehouses.
// It takes a router object from the mux package for route registration.
//
// URL Paths:
// - POST /warehouses: Create a new warehouse
// - GET /warehouses/{id}: Retrieve a warehouse by ID
// - PUT /warehouses/{id}: Update an existing warehouse by ID
// - DELETE /warehouses/{id}: Delete a warehouse by ID
func (h *WarehouseHandlers) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/warehouses", h.CreateWarehouse).Methods("POST")
	router.HandleFunc("/warehouses/{id:[0-9]+}", h.GetWarehouseByID).Methods("GET")
	router.HandleFunc("/warehouses/{id:[0-9]+}", h.UpdateWarehouse).Methods("PUT")
	router.HandleFunc("/warehouses/{id:[0-9]+}", h.DeleteWarehouse).Methods("DELETE")
}

// CreateWarehouse handles the creation of a new warehouse.
//
// This handler reads the incoming request body, decodes it into a Warehouse struct,
// and attempts to store it in the database. On successful creation, it returns
// a status code 201 Created. If an error occurs, it responds with an appropriate
// status code and error message.
//
// HTTP Method: POST
// URL Path: /warehouses
//
// Request Body:
// - JSON representation of a Warehouse object.
//
// Response:
// - Status Code: 201 (Created) if the warehouse is successfully created.
// - Status Code: 400 (Bad Request) if the request body is invalid.
// - Status Code: 500 (Internal Server Error) if the creation fails.
func (h *WarehouseHandlers) CreateWarehouse(w http.ResponseWriter, r *http.Request) {
	var req models.Warehouse
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	err = h.WarehouseStore.CreateWarehouse(&req)
	if err != nil {
		http.Error(w, "Could not create warehouse", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("Warehouse created successfully"))
}

// GetWarehouseByID handles retrieving a warehouse by its ID.
//
// This handler extracts the warehouse ID from the URL path, retrieves the warehouse
// from the database, and responds with the warehouse details in JSON format if found.
// If the warehouse is not found or an error occurs, it responds with an appropriate
// status code and error message.
//
// HTTP Method: GET
// URL Path: /warehouses/{id}
//
// Response:
// - Status Code: 200 (OK) and the warehouse details in JSON if the warehouse is found.
// - Status Code: 400 (Bad Request) if the ID is invalid.
// - Status Code: 404 (Not Found) if the warehouse is not found.
func (h *WarehouseHandlers) GetWarehouseByID(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	warehouseID, err := strconv.Atoi(params["id"])
	if err != nil {
		http.Error(w, "Invalid warehouse ID", http.StatusBadRequest)
		return
	}

	warehouse, err := h.WarehouseStore.GetWarehouseByID(warehouseID)
	if err != nil {
		http.Error(w, "Warehouse not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(warehouse)
}

// UpdateWarehouse handles updating an existing warehouse by ID.
//
// This handler extracts the warehouse ID from the URL path, decodes the request body
// into a Warehouse struct, updates the warehouse in the database, and returns a success
// response. If an error occurs, it responds with an appropriate status code and error
// message.
//
// HTTP Method: PUT
// URL Path: /warehouses/{id}
//
// Request Body:
// - JSON representation of a Warehouse object to update.
//
// Response:
// - Status Code: 200 (OK) if the warehouse is successfully updated.
// - Status Code: 400 (Bad Request) if the request body or ID is invalid.
// - Status Code: 500 (Internal Server Error) if the update fails.
func (h *WarehouseHandlers) UpdateWarehouse(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	warehouseID, err := strconv.Atoi(params["id"])
	if err != nil {
		http.Error(w, "Invalid warehouse ID", http.StatusBadRequest)
		return
	}

	var req models.Warehouse
	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	req.ID = warehouseID
	err = h.WarehouseStore.UpdateWarehouse(&req)
	if err != nil {
		http.Error(w, "Could not update warehouse", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Warehouse updated successfully"))
}

// DeleteWarehouse handles deleting a warehouse by its ID.
//
// This handler extracts the warehouse ID from the URL path, deletes the warehouse
// from the database, and returns a success message. If any error occurs, it
// responds with an appropriate status code and error message.
//
// HTTP Method: DELETE
// URL Path: /warehouses/{id}
//
// Response:
// - Status Code: 200 (OK) if the warehouse is successfully deleted.
// - Status Code: 400 (Bad Request) if the ID is invalid.
// - Status Code: 500 (Internal Server Error) if the deletion fails.
func (h *WarehouseHandlers) DeleteWarehouse(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	warehouseID, err := strconv.Atoi(params["id"])
	if err != nil {
		http.Error(w, "Invalid warehouse ID", http.StatusBadRequest)
		return
	}

	err = h.WarehouseStore.DeleteWarehouse(warehouseID)
	if err != nil {
		http.Error(w, "Could not delete warehouse", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Warehouse deleted successfully"))
}

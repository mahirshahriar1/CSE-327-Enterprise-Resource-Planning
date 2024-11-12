// Package stock_handlers contains the handlers for stock-related HTTP endpoints.
package stock_handlers

import (
	"encoding/json"
	"erp/models"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

// StockHandlers contains dependencies for handling stock-related requests.
type StockHandlers struct {
	StockStore models.StockStore
}

// RegisterRoutes registers all the stock routes for the HTTP server.
func (h *StockHandlers) RegisterRoutes(router *mux.Router) {
	// Register route for creating a new stock entry
	router.HandleFunc("/stock", h.CreateStock).Methods("POST")

	// Register route for fetching stock by product ID
	router.HandleFunc("/stock/product/{product_id:[0-9]+}", h.GetStockByProductID).Methods("GET")

	// Register route for updating an existing stock entry
	router.HandleFunc("/stock/{id:[0-9]+}", h.UpdateStock).Methods("PUT")

	// Register route for deleting a stock entry by ID
	router.HandleFunc("/stock/{id:[0-9]+}", h.DeleteStock).Methods("DELETE")
}

// CreateStock handles the creation of a new stock entry.
// It decodes the request body into a Stock struct and stores it in the database.
func (h *StockHandlers) CreateStock(w http.ResponseWriter, r *http.Request) {
	var req models.Stock
	// Decode the JSON body into the Stock struct
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	// Attempt to create the stock entry in the database
	err = h.StockStore.CreateStock(&req)
	if err != nil {
		http.Error(w, "Could not create stock", http.StatusInternalServerError)
		return
	}

	// Return success response
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("Stock created successfully"))
}

// GetStockByProductID handles retrieving stock information by product ID.
// It returns the stock details in JSON format if found, or an error if not.
func (h *StockHandlers) GetStockByProductID(w http.ResponseWriter, r *http.Request) {
	// Extract product ID from the URL parameters
	params := mux.Vars(r)
	productID, err := strconv.Atoi(params["product_id"])
	if err != nil {
		http.Error(w, "Invalid product ID", http.StatusBadRequest)
		return
	}

	// Retrieve the stock entry by product ID from the database
	stock, err := h.StockStore.GetStockByProductID(productID)
	if err != nil {
		http.Error(w, "Stock not found for the given product ID", http.StatusNotFound)
		return
	}

	// Respond with the stock entry in JSON format
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(stock)
}

// UpdateStock handles updating an existing stock entry by ID.
// It decodes the request body, updates the stock entry in the database, and returns a success response.
func (h *StockHandlers) UpdateStock(w http.ResponseWriter, r *http.Request) {
	// Extract stock ID from the URL parameters
	params := mux.Vars(r)
	stockID, err := strconv.Atoi(params["id"])
	if err != nil {
		http.Error(w, "Invalid stock ID", http.StatusBadRequest)
		return
	}

	var req models.Stock
	// Decode the request body into the Stock struct
	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	// Set the stock ID for the update operation
	req.ID = stockID

	// Attempt to update the stock entry in the database
	err = h.StockStore.UpdateStock(&req)
	if err != nil {
		http.Error(w, "Could not update stock", http.StatusInternalServerError)
		return
	}

	// Return success response
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Stock updated successfully"))
}

// DeleteStock handles deleting a stock entry by ID.
// It deletes the stock entry from the database and returns a success message.
func (h *StockHandlers) DeleteStock(w http.ResponseWriter, r *http.Request) {
	// Extract stock ID from the URL parameters
	params := mux.Vars(r)
	stockID, err := strconv.Atoi(params["id"])
	if err != nil {
		http.Error(w, "Invalid stock ID", http.StatusBadRequest)
		return
	}

	// Attempt to delete the stock entry from the database
	err = h.StockStore.DeleteStock(stockID)
	if err != nil {
		http.Error(w, "Could not delete stock", http.StatusInternalServerError)
		return
	}

	// Return success response
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Stock deleted successfully"))
}

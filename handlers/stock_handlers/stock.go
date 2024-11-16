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

// RegisterRoutes registers all the stock-related routes for the HTTP server.
//
// This method sets up routes for creating, retrieving, updating, and deleting stock entries.
// It uses the router object from the mux package for route registration.
//
// URL Paths:
// - POST /stock: Create a new stock entry
// - GET /stock/product/{product_id}: Retrieve stock by product ID
// - PUT /stock/{id}: Update an existing stock entry by ID
// - DELETE /stock/{id}: Delete a stock entry by ID
func (h *StockHandlers) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/stock", h.CreateStock).Methods("POST")
	router.HandleFunc("/stock/product/{product_id:[0-9]+}", h.GetStockByProductID).Methods("GET")
	router.HandleFunc("/stock/{id:[0-9]+}", h.UpdateStock).Methods("PUT")
	router.HandleFunc("/stock/{id:[0-9]+}", h.DeleteStock).Methods("DELETE")
}

// CreateStock handles the creation of a new stock entry.
//
// This handler reads the incoming request body, decodes it into a Stock struct,
// and attempts to store it in the database. On successful creation, it returns
// a status code 201 Created. If an error occurs, it responds with an appropriate
// status code and error message.
//
// HTTP Method: POST
// URL Path: /stock
//
// Request Body:
// - JSON representation of a Stock object.
//
// Response:
// - Status Code: 201 (Created) if the stock is successfully created.
// - Status Code: 400 (Bad Request) if the request body is invalid.
// - Status Code: 500 (Internal Server Error) if the creation fails.
func (h *StockHandlers) CreateStock(w http.ResponseWriter, r *http.Request) {
	var req models.Stock
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	err = h.StockStore.CreateStock(&req)
	if err != nil {
		http.Error(w, "Could not create stock", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("Stock created successfully"))
}

// GetStockByProductID handles retrieving stock information by product ID.
//
// This handler extracts the product ID from the URL path, retrieves the stock
// from the database, and responds with the stock details in JSON format if found.
// If the stock is not found or an error occurs, it responds with an appropriate
// status code and error message.
//
// HTTP Method: GET
// URL Path: /stock/product/{product_id}
//
// Response:
// - Status Code: 200 (OK) and the stock details in JSON if found.
// - Status Code: 400 (Bad Request) if the product ID is invalid.
// - Status Code: 404 (Not Found) if the stock is not found.
func (h *StockHandlers) GetStockByProductID(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	productID, err := strconv.Atoi(params["product_id"])
	if err != nil {
		http.Error(w, "Invalid product ID", http.StatusBadRequest)
		return
	}

	stock, err := h.StockStore.GetStockByProductID(productID)
	if err != nil {
		http.Error(w, "Stock not found for the given product ID", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(stock)
}

// UpdateStock handles updating an existing stock entry by ID.
//
// This handler extracts the stock ID from the URL path, decodes the request body
// into a Stock struct, updates the stock in the database, and returns a success
// response. If an error occurs, it responds with an appropriate status code and
// error message.
//
// HTTP Method: PUT
// URL Path: /stock/{id}
//
// Request Body:
// - JSON representation of a Stock object to update.
//
// Response:
// - Status Code: 200 (OK) if the stock is successfully updated.
// - Status Code: 400 (Bad Request) if the request body or stock ID is invalid.
// - Status Code: 500 (Internal Server Error) if the update fails.
func (h *StockHandlers) UpdateStock(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	stockID, err := strconv.Atoi(params["id"])
	if err != nil {
		http.Error(w, "Invalid stock ID", http.StatusBadRequest)
		return
	}

	var req models.Stock
	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	req.ID = stockID
	err = h.StockStore.UpdateStock(&req)
	if err != nil {
		http.Error(w, "Could not update stock", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Stock updated successfully"))
}

// DeleteStock handles deleting a stock entry by ID.
//
// This handler extracts the stock ID from the URL path, deletes the stock
// from the database, and returns a success message. If an error occurs, it
// responds with an appropriate status code and error message.
//
// HTTP Method: DELETE
// URL Path: /stock/{id}
//
// Response:
// - Status Code: 200 (OK) if the stock is successfully deleted.
// - Status Code: 400 (Bad Request) if the stock ID is invalid.
// - Status Code: 500 (Internal Server Error) if the deletion fails.
func (h *StockHandlers) DeleteStock(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	stockID, err := strconv.Atoi(params["id"])
	if err != nil {
		http.Error(w, "Invalid stock ID", http.StatusBadRequest)
		return
	}

	err = h.StockStore.DeleteStock(stockID)
	if err != nil {
		http.Error(w, "Could not delete stock", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Stock deleted successfully"))
}

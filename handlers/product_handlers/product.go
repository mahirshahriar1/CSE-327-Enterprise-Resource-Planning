// product_handlers.go
// Description: This file contains the handlers for the product-related endpoints.

package product_handlers

import (
	"encoding/json"
	"erp/models"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

// ProductHandlers struct contains the product store dependency
type ProductHandlers struct {
	ProductStore models.ProductStore
}

// RegisterRoutes registers all the product routes
func (h *ProductHandlers) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/products", h.CreateProduct).Methods("POST")               // Create a new product
	router.HandleFunc("/products/{id:[0-9]+}", h.GetProductByID).Methods("GET")   // Get product by ID
	router.HandleFunc("/products/{id:[0-9]+}", h.UpdateProduct).Methods("PUT")    // Update product
	router.HandleFunc("/products/{id:[0-9]+}", h.DeleteProduct).Methods("DELETE") // Delete product
}

// CreateProduct handles product creation
func (h *ProductHandlers) CreateProduct(w http.ResponseWriter, r *http.Request) {
	var req models.Product
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	// Create the product in the database
	err = h.ProductStore.CreateProduct(&req)
	if err != nil {
		http.Error(w, "Could not create product", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("Product created successfully"))
}

// GetProductByID handles retrieving a product by ID
func (h *ProductHandlers) GetProductByID(w http.ResponseWriter, r *http.Request) {
	// Extract the product ID from the URL
	params := mux.Vars(r)
	productID, err := strconv.Atoi(params["id"])
	if err != nil {
		http.Error(w, "Invalid product ID", http.StatusBadRequest)
		return
	}

	// Retrieve the product from the database
	product, err := h.ProductStore.GetProductByID(productID)
	if err != nil {
		http.Error(w, "Product not found", http.StatusNotFound)
		return
	}

	// Respond with the product details in JSON format
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(product)
}

// UpdateProduct handles updating a product
func (h *ProductHandlers) UpdateProduct(w http.ResponseWriter, r *http.Request) {
	// Extract the product ID from the URL
	params := mux.Vars(r)
	productID, err := strconv.Atoi(params["id"])
	if err != nil {
		http.Error(w, "Invalid product ID", http.StatusBadRequest)
		return
	}

	var req models.Product
	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	// Set the ID of the product to update
	req.ID = productID

	// Update the product in the database
	err = h.ProductStore.UpdateProduct(&req)
	if err != nil {
		http.Error(w, "Could not update product", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Product updated successfully"))
}

// DeleteProduct handles deleting a product by ID
func (h *ProductHandlers) DeleteProduct(w http.ResponseWriter, r *http.Request) {
	// Extract the product ID from the URL
	params := mux.Vars(r)
	productID, err := strconv.Atoi(params["id"])
	if err != nil {
		http.Error(w, "Invalid product ID", http.StatusBadRequest)
		return
	}

	// Delete the product from the database
	err = h.ProductStore.DeleteProduct(productID)
	if err != nil {
		http.Error(w, "Could not delete product", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Product deleted successfully"))
}

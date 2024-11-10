// Package product_handlers contains the handlers for product-related HTTP endpoints.
package product_handlers

import (
	"encoding/json"
	"erp/models"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

// ProductHandlers contains dependencies for handling product-related requests.
type ProductHandlers struct {
	ProductStore models.ProductStore
}

// RegisterRoutes registers all the product routes for the HTTP server.
func (h *ProductHandlers) RegisterRoutes(router *mux.Router) {
	// Register route for creating a new product
	router.HandleFunc("/products", h.CreateProduct).Methods("POST")

	// Register route for fetching a product by ID
	router.HandleFunc("/products/{id:[0-9]+}", h.GetProductByID).Methods("GET")

	// Register route for updating an existing product
	router.HandleFunc("/products/{id:[0-9]+}", h.UpdateProduct).Methods("PUT")

	// Register route for deleting a product by ID
	router.HandleFunc("/products/{id:[0-9]+}", h.DeleteProduct).Methods("DELETE")
}

// CreateProduct handles the creation of a new product.
// It decodes the request body into a Product struct and stores it in the database.
func (h *ProductHandlers) CreateProduct(w http.ResponseWriter, r *http.Request) {
	var req models.Product
	// Decode the JSON body into the Product struct
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	// Attempt to create the product in the database
	err = h.ProductStore.CreateProduct(&req)
	if err != nil {
		http.Error(w, "Could not create product", http.StatusInternalServerError)
		return
	}

	// Return success response
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("Product created successfully"))
}

// GetProductByID handles retrieving a product by its ID.
// It returns the product details in JSON format if found, or an error if not.
func (h *ProductHandlers) GetProductByID(w http.ResponseWriter, r *http.Request) {
	// Extract product ID from the URL parameters
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

	// Respond with the product in JSON format
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(product)
}

// UpdateProduct handles updating an existing product by ID.
// It decodes the request body, updates the product in the database, and returns a success response.
func (h *ProductHandlers) UpdateProduct(w http.ResponseWriter, r *http.Request) {
	// Extract product ID from the URL parameters
	params := mux.Vars(r)
	productID, err := strconv.Atoi(params["id"])
	if err != nil {
		http.Error(w, "Invalid product ID", http.StatusBadRequest)
		return
	}

	var req models.Product
	// Decode the request body into the Product struct
	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	// Set the product ID for the update operation
	req.ID = productID

	// Attempt to update the product in the database
	err = h.ProductStore.UpdateProduct(&req)
	if err != nil {
		http.Error(w, "Could not update product", http.StatusInternalServerError)
		return
	}

	// Return success response
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Product updated successfully"))
}

// DeleteProduct handles deleting a product by ID.
// It deletes the product from the database and returns a success message.
func (h *ProductHandlers) DeleteProduct(w http.ResponseWriter, r *http.Request) {
	// Extract product ID from the URL parameters
	params := mux.Vars(r)
	productID, err := strconv.Atoi(params["id"])
	if err != nil {
		http.Error(w, "Invalid product ID", http.StatusBadRequest)
		return
	}

	// Attempt to delete the product from the database
	err = h.ProductStore.DeleteProduct(productID)
	if err != nil {
		http.Error(w, "Could not delete product", http.StatusInternalServerError)
		return
	}

	// Return success response
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Product deleted successfully"))
}

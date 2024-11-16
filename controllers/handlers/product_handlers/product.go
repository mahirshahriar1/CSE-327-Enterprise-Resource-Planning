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

// RegisterRoutes registers all the product-related routes for the HTTP server.
//
// This method sets up routes for creating, retrieving, updating, and deleting products.
// It takes a router object from the mux package for route registration.
//
// URL Paths:
// - POST /products: Create a new product
// - GET /products/{id}: Retrieve a product by ID
// - PUT /products/{id}: Update an existing product by ID
// - DELETE /products/{id}: Delete a product by ID
func (h *ProductHandlers) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/products", h.CreateProduct).Methods("POST")
	router.HandleFunc("/products/{id:[0-9]+}", h.GetProductByID).Methods("GET")
	router.HandleFunc("/products/{id:[0-9]+}", h.UpdateProduct).Methods("PUT")
	router.HandleFunc("/products/{id:[0-9]+}", h.DeleteProduct).Methods("DELETE")
}

// CreateProduct handles the creation of a new product.
//
// This handler reads the incoming request body, decodes it into a Product struct,
// and attempts to store it in the database. On successful creation, it returns
// a status code 201 Created. If an error occurs, it responds with an appropriate
// status code and error message.
//
// HTTP Method: POST
// URL Path: /products
//
// Request Body:
// - JSON representation of a Product object.
//
// Response:
// - Status Code: 201 (Created) if the product is successfully created.
// - Status Code: 400 (Bad Request) if the request body is invalid.
// - Status Code: 500 (Internal Server Error) if the creation fails.
func (h *ProductHandlers) CreateProduct(w http.ResponseWriter, r *http.Request) {
	var req models.Product
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	err = h.ProductStore.CreateProduct(&req)
	if err != nil {
		http.Error(w, "Could not create product", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("Product created successfully"))
}

// GetProductByID handles retrieving a product by its ID.
//
// This handler extracts the product ID from the URL path, retrieves the product
// from the database, and responds with the product details in JSON format if found.
// If the product is not found or an error occurs, it responds with an appropriate
// status code and error message.
//
// HTTP Method: GET
// URL Path: /products/{id}
//
// Response:
// - Status Code: 200 (OK) and the product details in JSON if the product is found.
// - Status Code: 400 (Bad Request) if the ID is invalid.
// - Status Code: 404 (Not Found) if the product is not found.
func (h *ProductHandlers) GetProductByID(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	productID, err := strconv.Atoi(params["id"])
	if err != nil {
		http.Error(w, "Invalid product ID", http.StatusBadRequest)
		return
	}

	product, err := h.ProductStore.GetProductByID(productID)
	if err != nil {
		http.Error(w, "Product not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(product)
}

// UpdateProduct handles updating an existing product by ID.
//
// This handler extracts the product ID from the URL path, decodes the request body
// into a Product struct, updates the product in the database, and returns a success
// response. If an error occurs, it responds with an appropriate status code and error
// message.
//
// HTTP Method: PUT
// URL Path: /products/{id}
//
// Request Body:
// - JSON representation of a Product object to update.
//
// Response:
// - Status Code: 200 (OK) if the product is successfully updated.
// - Status Code: 400 (Bad Request) if the request body or ID is invalid.
// - Status Code: 500 (Internal Server Error) if the update fails.
func (h *ProductHandlers) UpdateProduct(w http.ResponseWriter, r *http.Request) {
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

	req.ID = productID
	err = h.ProductStore.UpdateProduct(&req)
	if err != nil {
		http.Error(w, "Could not update product", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Product updated successfully"))
}

// DeleteProduct handles deleting a product by its ID.
//
// This handler extracts the product ID from the URL path, deletes the product
// from the database, and returns a success message. If any error occurs, it
// responds with an appropriate status code and error message.
//
// HTTP Method: DELETE
// URL Path: /products/{id}
//
// Response:
// - Status Code: 200 (OK) if the product is successfully deleted.
// - Status Code: 400 (Bad Request) if the ID is invalid.
// - Status Code: 500 (Internal Server Error) if the deletion fails.
func (h *ProductHandlers) DeleteProduct(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	productID, err := strconv.Atoi(params["id"])
	if err != nil {
		http.Error(w, "Invalid product ID", http.StatusBadRequest)
		return
	}

	err = h.ProductStore.DeleteProduct(productID)
	if err != nil {
		http.Error(w, "Could not delete product", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Product deleted successfully"))
}

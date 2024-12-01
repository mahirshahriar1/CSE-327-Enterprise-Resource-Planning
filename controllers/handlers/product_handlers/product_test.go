// Package product_handlers_test contains unit tests for product-related HTTP handlers.
package product_handlers_test

import (
	"bytes"
	"encoding/json"
	"erp/controllers/handlers/product_handlers"
	"erp/models"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

// TestCreateProduct verifies the behavior of the CreateProduct handler.
//
// This test sets up a mock database, simulates a POST request to create a new product,
// and ensures the handler returns the expected response and interacts with the database correctly.
func TestCreateProduct(t *testing.T) {
	// Set up mock database
	db, mock, err := sqlmock.New()
	assert.NoError(t, err, "failed to create mock database")
	defer db.Close()

	store := product_handlers.NewDBProductStore(db)
	handler := &product_handlers.ProductHandlers{ProductStore: store}

	// Sample product data
	product := &models.Product{
		Name:   "Test Product",
		Brand:  "Test Brand",
		Season: "Summer",
		Price:  100.50,
	}

	// Mock database behavior
	mock.ExpectQuery(`INSERT INTO products \(name, brand, season, price\)`).
		WithArgs(product.Name, product.Brand, product.Season, product.Price).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

	// Create HTTP request and recorder
	body, _ := json.Marshal(product)
	req := httptest.NewRequest(http.MethodPost, "/products", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	// Call the handler
	handler.CreateProduct(rec, req)

	// Verify response
	assert.Equal(t, http.StatusCreated, rec.Code)
	assert.Equal(t, "Product created successfully", rec.Body.String())

	// Verify expectations
	assert.NoError(t, mock.ExpectationsWereMet(), "unmet mock database expectations")
}

// TestGetProductByID verifies the behavior of the GetProductByID handler.
//
// This test sets up a mock database, simulates a GET request for a product by ID,
// and ensures the handler returns the expected product details or error response.
func TestGetProductByID(t *testing.T) {
	// Set up mock database
	db, mock, err := sqlmock.New()
	assert.NoError(t, err, "failed to create mock database")
	defer db.Close()

	store := product_handlers.NewDBProductStore(db)
	handler := &product_handlers.ProductHandlers{ProductStore: store}

	// Sample product data
	product := &models.Product{
		ID:     1,
		Name:   "Test Product",
		Brand:  "Test Brand",
		Season: "Summer",
		Price:  100.50,
	}

	// Mock database behavior
	mock.ExpectQuery(`SELECT id, name, brand, season, price FROM products WHERE id = \$1`).
		WithArgs(product.ID).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name", "brand", "season", "price"}).
			AddRow(product.ID, product.Name, product.Brand, product.Season, product.Price))

	// Create HTTP request and recorder
	req := httptest.NewRequest(http.MethodGet, "/products/1", nil)
	rec := httptest.NewRecorder()
	req = mux.SetURLVars(req, map[string]string{"id": "1"})

	// Call the handler
	handler.GetProductByID(rec, req)

	// Verify response
	assert.Equal(t, http.StatusOK, rec.Code)
	expectedBody, _ := json.Marshal(product)
	assert.JSONEq(t, string(expectedBody), rec.Body.String())

	// Verify mock expectations
	assert.NoError(t, mock.ExpectationsWereMet(), "unmet mock database expectations")
}

// TestUpdateProduct verifies the behavior of the UpdateProduct handler.
//
// This test sets up a mock database, simulates a PUT request to update a product,
// and ensures the handler processes the request correctly and updates the database.
func TestUpdateProduct(t *testing.T) {
	// Set up mock database
	db, mock, err := sqlmock.New()
	assert.NoError(t, err, "failed to create mock database")
	defer db.Close()

	store := product_handlers.NewDBProductStore(db)
	handler := &product_handlers.ProductHandlers{ProductStore: store}

	// Sample product data
	product := &models.Product{
		ID:     1,
		Name:   "Updated Product",
		Brand:  "Updated Brand",
		Season: "Winter",
		Price:  120.75,
	}

	// Mock database behavior
	mock.ExpectExec(`UPDATE products SET name = \$1, brand = \$2, season = \$3, price = \$4 WHERE id = \$5`).
		WithArgs(product.Name, product.Brand, product.Season, product.Price, product.ID).
		WillReturnResult(sqlmock.NewResult(0, 1)) // Simulate one row affected

	// Create HTTP request and recorder
	body, _ := json.Marshal(models.Product{
		Name:   product.Name,
		Brand:  product.Brand,
		Season: product.Season,
		Price:  product.Price,
	})
	req := httptest.NewRequest(http.MethodPut, "/products/1", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	req = mux.SetURLVars(req, map[string]string{"id": "1"})

	// Call the handler
	handler.UpdateProduct(rec, req)

	// Verify response
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, "Product updated successfully", rec.Body.String())

	// Verify mock expectations
	assert.NoError(t, mock.ExpectationsWereMet(), "unmet mock database expectations")
}

// TestDeleteProduct verifies the behavior of the DeleteProduct handler.
//
// This test sets up a mock database, simulates a DELETE request for a product by ID,
// and ensures the handler interacts correctly with the database and responds appropriately.
func TestDeleteProduct(t *testing.T) {
	// Set up mock database
	db, mock, err := sqlmock.New()
	assert.NoError(t, err, "failed to create mock database")
	defer db.Close()

	store := product_handlers.NewDBProductStore(db)
	handler := &product_handlers.ProductHandlers{ProductStore: store}

	// Mock database behavior
	mock.ExpectExec(`DELETE FROM products WHERE id = \$1`).
		WithArgs(1).
		WillReturnResult(sqlmock.NewResult(0, 1)) // Simulate one row affected

	// Create HTTP request and recorder
	req := httptest.NewRequest(http.MethodDelete, "/products/1", nil)
	rec := httptest.NewRecorder()
	req = mux.SetURLVars(req, map[string]string{"id": "1"})

	// Call the handler
	handler.DeleteProduct(rec, req)

	// Verify response
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, "Product deleted successfully", rec.Body.String())

	// Verify mock expectations
	assert.NoError(t, mock.ExpectationsWereMet(), "unmet mock database expectations")
}

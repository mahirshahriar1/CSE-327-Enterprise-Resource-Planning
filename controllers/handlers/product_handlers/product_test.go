// Package product_handlers_test contains unit tests for product-related HTTP handlers.
package product_handlers

import (
	"bytes"
	"encoding/json"
	"erp/models"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

// TestCreateProduct tests the CreateProduct handler.
func TestCreateProduct(t *testing.T) {
	// Set up mock database
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("could not create mock db: %v", err)
	}
	defer db.Close()

	store := &DBProductStore{DB: db}
	handler := &ProductHandlers{ProductStore: store}

	// Sample product data
	product := &models.Product{
		Name:   "Test Product",
		Brand:  "Test Brand",
		Season: "Summer",
		Price:  100.50,
	}

	// Mock database behavior
	mock.ExpectExec("INSERT INTO products").
		WithArgs(product.Name, product.Brand, product.Season, product.Price).
		WillReturnResult(sqlmock.NewResult(1, 1))

	// Create HTTP request and recorder
	body, _ := json.Marshal(product)
	req, _ := http.NewRequest("POST", "/products", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	// Call the handler
	handler.CreateProduct(rec, req)

	// Assert that no error occurred and response is correct
	assert.Equal(t, http.StatusCreated, rec.Code)
	assert.Equal(t, "Product created successfully", rec.Body.String())

	// Assert that the expected query was executed
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unmet expectations: %v", err)
	}
}

// TestGetProductByID tests the GetProductByID handler.
func TestGetProductByID(t *testing.T) {
	// Set up mock database
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("could not create mock db: %v", err)
	}
	defer db.Close()

	store := &DBProductStore{DB: db}
	handler := &ProductHandlers{ProductStore: store}

	// Sample product data
	product := &models.Product{
		ID:     1,
		Name:   "Test Product",
		Brand:  "Test Brand",
		Season: "Summer",
		Price:  100.50,
	}

	// Mock database behavior
	mock.ExpectQuery("SELECT id, name, brand, season, price FROM products WHERE id = \\$1").
		WithArgs(product.ID).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name", "brand", "season", "price"}).
			AddRow(product.ID, product.Name, product.Brand, product.Season, product.Price))

	// Create HTTP request and recorder
	req, _ := http.NewRequest("GET", "/products/1", nil)
	rec := httptest.NewRecorder()

	// Add mux variables to the request
	req = mux.SetURLVars(req, map[string]string{"id": "1"})

	// Call the handler
	handler.GetProductByID(rec, req)

	// Assert that no error occurred and response is correct
	assert.Equal(t, http.StatusOK, rec.Code)
	expectedBody, _ := json.Marshal(product)
	assert.JSONEq(t, string(expectedBody), rec.Body.String())

	// Assert that the expected query was executed
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unmet expectations: %v", err)
	}
}

// TestUpdateProduct tests the UpdateProduct handler.
func TestUpdateProduct(t *testing.T) {
	// Set up mock database
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("could not create mock db: %v", err)
	}
	defer db.Close()

	store := &DBProductStore{DB: db}
	handler := &ProductHandlers{ProductStore: store}

	// Sample product data
	product := &models.Product{
		ID:     1,
		Name:   "Updated Product",
		Brand:  "Updated Brand",
		Season: "Winter",
		Price:  120.75,
	}

	// Mock database behavior
	mock.ExpectExec("UPDATE products SET name = \\$1, brand = \\$2, season = \\$3, price = \\$4 WHERE id = \\$5").
		WithArgs(product.Name, product.Brand, product.Season, product.Price, product.ID).
		WillReturnResult(sqlmock.NewResult(1, 1))

	// Create HTTP request and recorder
	body, _ := json.Marshal(product)
	req, _ := http.NewRequest("PUT", "/products/1", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	// Add mux variables to the request
	req = mux.SetURLVars(req, map[string]string{"id": "1"})

	// Call the handler
	handler.UpdateProduct(rec, req)

	// Assert that no error occurred and response is correct
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, "Product updated successfully", rec.Body.String())

	// Assert that the expected query was executed
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unmet expectations: %v", err)
	}
}

// TestDeleteProduct tests the DeleteProduct handler.
func TestDeleteProduct(t *testing.T) {
	// Set up mock database
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("could not create mock db: %v", err)
	}
	defer db.Close()

	store := &DBProductStore{DB: db}
	handler := &ProductHandlers{ProductStore: store}

	// Mock database behavior
	mock.ExpectExec("DELETE FROM products WHERE id = \\$1").
		WithArgs(1).
		WillReturnResult(sqlmock.NewResult(1, 1))

	// Create HTTP request and recorder
	req, _ := http.NewRequest("DELETE", "/products/1", nil)
	rec := httptest.NewRecorder()

	// Add mux variables to the request
	req = mux.SetURLVars(req, map[string]string{"id": "1"})

	// Call the handler
	handler.DeleteProduct(rec, req)

	// Assert that no error occurred and response is correct
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, "Product deleted successfully", rec.Body.String())

	// Assert that the expected query was executed
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unmet expectations: %v", err)
	}
}

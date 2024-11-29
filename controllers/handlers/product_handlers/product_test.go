package product_handlers_test

import (
	"bytes"
	"encoding/json"
	"erp/controllers/handlers/product_handlers"
	"erp/models"

	// "erp/product_handlers"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

// MockProductStore is a mock implementation of the ProductStore interface for testing purposes.
type MockProductStore struct {
	CreateProductFn  func(product *models.Product) error
	GetProductByIDFn func(id int) (*models.Product, error)
	UpdateProductFn  func(product *models.Product) error
	DeleteProductFn  func(id int) error
}

func (m *MockProductStore) CreateProduct(product *models.Product) error {
	return m.CreateProductFn(product)
}

func (m *MockProductStore) GetProductByID(id int) (*models.Product, error) {
	return m.GetProductByIDFn(id)
}

func (m *MockProductStore) UpdateProduct(product *models.Product) error {
	return m.UpdateProductFn(product)
}

func (m *MockProductStore) DeleteProduct(id int) error {
	return m.DeleteProductFn(id)
}

// TestCreateProduct tests the CreateProduct handler.
func TestCreateProduct(t *testing.T) {
	mockStore := &MockProductStore{
		CreateProductFn: func(product *models.Product) error {
			product.ID = 1
			return nil
		},
	}

	handler := product_handlers.ProductHandlers{ProductStore: mockStore}

	t.Run("Successful creation", func(t *testing.T) {
		product := models.Product{Name: "Shirt", Brand: "BrandX", Season: "Summer", Price: 49.99}
		body, _ := json.Marshal(product)

		req := httptest.NewRequest(http.MethodPost, "/products", bytes.NewReader(body))
		w := httptest.NewRecorder()

		handler.CreateProduct(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)
		assert.Equal(t, "Product created successfully", w.Body.String())
	})

	t.Run("Invalid input", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/products", bytes.NewReader([]byte("invalid")))
		w := httptest.NewRecorder()

		handler.CreateProduct(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Equal(t, "Invalid input\n", w.Body.String())
	})

	t.Run("Store error", func(t *testing.T) {
		mockStore.CreateProductFn = func(product *models.Product) error {
			return errors.New("database error")
		}

		product := models.Product{Name: "Shirt", Brand: "BrandX", Season: "Summer", Price: 49.99}
		body, _ := json.Marshal(product)

		req := httptest.NewRequest(http.MethodPost, "/products", bytes.NewReader(body))
		w := httptest.NewRecorder()

		handler.CreateProduct(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.Equal(t, "Could not create product\n", w.Body.String())
	})
}

// TestGetProductByID tests the GetProductByID handler.
func TestGetProductByID(t *testing.T) {
	mockStore := &MockProductStore{
		GetProductByIDFn: func(id int) (*models.Product, error) {
			if id == 1 {
				return &models.Product{ID: 1, Name: "Shirt", Brand: "BrandX", Season: "Summer", Price: 49.99}, nil
			}
			return nil, errors.New("not found")
		},
	}

	handler := product_handlers.ProductHandlers{ProductStore: mockStore}

	t.Run("Successful retrieval", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/products/1", nil)
		req = mux.SetURLVars(req, map[string]string{"id": "1"})
		w := httptest.NewRecorder()

		handler.GetProductByID(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var product models.Product
		err := json.NewDecoder(w.Body).Decode(&product)
		assert.NoError(t, err)
		assert.Equal(t, "Shirt", product.Name)
	})

	t.Run("Invalid ID", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/products/abc", nil)
		w := httptest.NewRecorder()

		handler.GetProductByID(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Equal(t, "Invalid product ID\n", w.Body.String())
	})

	t.Run("Product not found", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/products/2", nil)
		req = mux.SetURLVars(req, map[string]string{"id": "2"})
		w := httptest.NewRecorder()

		handler.GetProductByID(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
		assert.Equal(t, "Product not found\n", w.Body.String())
	})
}

// TestUpdateProduct tests the UpdateProduct handler.
func TestUpdateProduct(t *testing.T) {
	mockStore := &MockProductStore{
		UpdateProductFn: func(product *models.Product) error {
			if product.ID == 1 {
				return nil
			}
			return errors.New("not found")
		},
	}

	handler := product_handlers.ProductHandlers{ProductStore: mockStore}

	t.Run("Successful update", func(t *testing.T) {
		product := models.Product{Name: "Shirt Updated", Brand: "BrandX", Season: "Summer", Price: 59.99}
		body, _ := json.Marshal(product)

		req := httptest.NewRequest(http.MethodPut, "/products/1", bytes.NewReader(body))
		req = mux.SetURLVars(req, map[string]string{"id": "1"})
		w := httptest.NewRecorder()

		handler.UpdateProduct(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, "Product updated successfully", w.Body.String())
	})

	t.Run("Invalid ID", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPut, "/products/abc", nil)
		w := httptest.NewRecorder()

		handler.UpdateProduct(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Equal(t, "Invalid product ID\n", w.Body.String())
	})

	t.Run("Update error", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPut, "/products/2", bytes.NewReader([]byte(`{"name":"Invalid Product"}`)))
		req = mux.SetURLVars(req, map[string]string{"id": "2"})
		w := httptest.NewRecorder()

		handler.UpdateProduct(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.Equal(t, "Could not update product\n", w.Body.String())
	})
}

// TestDeleteProduct tests the DeleteProduct handler.
func TestDeleteProduct(t *testing.T) {
	mockStore := &MockProductStore{
		DeleteProductFn: func(id int) error {
			if id == 1 {
				return nil
			}
			return errors.New("not found")
		},
	}

	handler := product_handlers.ProductHandlers{ProductStore: mockStore}

	t.Run("Successful deletion", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodDelete, "/products/1", nil)
		req = mux.SetURLVars(req, map[string]string{"id": "1"})
		w := httptest.NewRecorder()

		handler.DeleteProduct(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, "Product deleted successfully", w.Body.String())
	})

	t.Run("Invalid ID", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodDelete, "/products/abc", nil)
		w := httptest.NewRecorder()

		handler.DeleteProduct(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Equal(t, "Invalid product ID\n", w.Body.String())
	})

	t.Run("Deletion error", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodDelete, "/products/2", nil)
		req = mux.SetURLVars(req, map[string]string{"id": "2"})
		w := httptest.NewRecorder()

		handler.DeleteProduct(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.Equal(t, "Could not delete product\n", w.Body.String())
	})
}

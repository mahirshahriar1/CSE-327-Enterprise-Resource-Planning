// Package stock_handlers_test contains unit tests for stock-related HTTP handlers.
package stock_handlers_test

import (
	"bytes"
	"encoding/json"
	"erp/controllers/handlers/stock_handlers"
	"erp/models"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockStockStore is a mock implementation of the models.StockStore interface for testing.
type MockStockStore struct {
	mock.Mock
}

func (m *MockStockStore) CreateStock(stock *models.Stock) error {
	args := m.Called(stock)
	return args.Error(0)
}

func (m *MockStockStore) GetStockByProductID(productID int) (*models.Stock, error) {
	args := m.Called(productID)
	return args.Get(0).(*models.Stock), args.Error(1)
}

func (m *MockStockStore) UpdateStock(stock *models.Stock) error {
	args := m.Called(stock)
	return args.Error(0)
}

func (m *MockStockStore) DeleteStock(id int) error {
	args := m.Called(id)
	return args.Error(0)
}

// TestStockHandlers tests the stock-related HTTP handlers.
func TestStockHandlers(t *testing.T) {
	mockStore := new(MockStockStore)
	handler := &stock_handlers.StockHandlers{StockStore: mockStore}
	router := mux.NewRouter()
	handler.RegisterRoutes(router)

	t.Run("CreateStock", func(t *testing.T) {
		stock := &models.Stock{ProductID: 1, Quantity: 10, WarehouseID: 1, Location: "A1"}
		mockStore.On("CreateStock", stock).Return(nil)

		body, _ := json.Marshal(stock)
		req := httptest.NewRequest(http.MethodPost, "/stock", bytes.NewReader(body))
		rec := httptest.NewRecorder()

		router.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusCreated, rec.Code)
		assert.Equal(t, "Stock created successfully", rec.Body.String())
		mockStore.AssertCalled(t, "CreateStock", stock)
	})

	t.Run("GetStockByProductID", func(t *testing.T) {
		productID := 1
		expectedStock := &models.Stock{ID: 1, ProductID: productID, Quantity: 20, WarehouseID: 2, Location: "B2"}
		mockStore.On("GetStockByProductID", productID).Return(expectedStock, nil)

		req := httptest.NewRequest(http.MethodGet, "/stock/product/"+strconv.Itoa(productID), nil)
		rec := httptest.NewRecorder()

		router.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)

		var result models.Stock
		json.Unmarshal(rec.Body.Bytes(), &result)
		assert.Equal(t, *expectedStock, result)
		mockStore.AssertCalled(t, "GetStockByProductID", productID)
	})

	t.Run("UpdateStock", func(t *testing.T) {
		stockID := 1
		stock := &models.Stock{ID: stockID, ProductID: 2, Quantity: 15, WarehouseID: 3, Location: "C3"}
		mockStore.On("UpdateStock", stock).Return(nil)

		body, _ := json.Marshal(stock)
		req := httptest.NewRequest(http.MethodPut, "/stock/"+strconv.Itoa(stockID), bytes.NewReader(body))
		rec := httptest.NewRecorder()

		router.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, "Stock updated successfully", rec.Body.String())
		mockStore.AssertCalled(t, "UpdateStock", stock)
	})

	t.Run("DeleteStock", func(t *testing.T) {
		stockID := 1
		mockStore.On("DeleteStock", stockID).Return(nil)

		req := httptest.NewRequest(http.MethodDelete, "/stock/"+strconv.Itoa(stockID), nil)
		rec := httptest.NewRecorder()

		router.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, "Stock deleted successfully", rec.Body.String())
		mockStore.AssertCalled(t, "DeleteStock", stockID)
	})
}

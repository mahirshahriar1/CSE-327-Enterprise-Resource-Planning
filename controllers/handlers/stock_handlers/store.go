// Package stock_handlers contains the database store implementations for managing stock records.
package stock_handlers

import (
	"database/sql"
	"erp/models"
	"fmt"
)

// DBStockStore implements the StockStore interface for database operations.
type DBStockStore struct {
	DB *sql.DB
}

// NewDBStockStore initializes a new DBStockStore instance.
//
// Parameters:
// - db: A *sql.DB instance representing the database connection.
//
// Returns:
// - A pointer to an instance of DBStockStore.
func NewDBStockStore(db *sql.DB) *DBStockStore {
	return &DBStockStore{DB: db}
}

// CreateStock inserts a new stock record into the database.
//
// Parameters:
// - stock: A pointer to the Stock struct containing the stock details to insert.
//
// Returns:
// - An error if the insertion fails, otherwise nil.
func (s *DBStockStore) CreateStock(stock *models.Stock) error {
	query := `
		INSERT INTO stock (product_id, quantity, warehouse_id, location)
		VALUES ($1, $2, $3, $4)
	`
	_, err := s.DB.Exec(query, stock.ProductID, stock.Quantity, stock.WarehouseID, stock.Location)
	if err != nil {
		return fmt.Errorf("failed to insert stock: %w", err)
	}
	return nil
}

// GetStockByProductID retrieves a stock record from the database by product ID.
//
// Parameters:
// - productID: An integer representing the product ID.
//
// Returns:
// - A pointer to the Stock struct if found.
// - An error if no record is found or if the query fails.
func (s *DBStockStore) GetStockByProductID(productID int) (*models.Stock, error) {
	query := `
		SELECT id, product_id, quantity, warehouse_id, location
		FROM stock
		WHERE product_id = $1
	`
	row := s.DB.QueryRow(query, productID)

	var stock models.Stock
	err := row.Scan(&stock.ID, &stock.ProductID, &stock.Quantity, &stock.WarehouseID, &stock.Location)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("no stock found for product ID %d", productID)
		}
		return nil, fmt.Errorf("failed to retrieve stock: %w", err)
	}

	return &stock, nil
}

// UpdateStock updates an existing stock record in the database.
//
// Parameters:
// - stock: A pointer to the Stock struct containing the updated stock details.
//
// Returns:
// - An error if the update fails, otherwise nil.
func (s *DBStockStore) UpdateStock(stock *models.Stock) error {
	query := `
		UPDATE stock
		SET product_id = $1, quantity = $2, warehouse_id = $3, location = $4
		WHERE id = $5
	`
	_, err := s.DB.Exec(query, stock.ProductID, stock.Quantity, stock.WarehouseID, stock.Location, stock.ID)
	if err != nil {
		return fmt.Errorf("failed to update stock with ID %d: %w", stock.ID, err)
	}
	return nil
}

// DeleteStock removes a stock record from the database by ID.
//
// Parameters:
// - id: An integer representing the stock ID to delete.
//
// Returns:
// - An error if the deletion fails, otherwise nil.
func (s *DBStockStore) DeleteStock(id int) error {
	query := `
		DELETE FROM stock
		WHERE id = $1
	`
	_, err := s.DB.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete stock with ID %d: %w", id, err)
	}
	return nil
}

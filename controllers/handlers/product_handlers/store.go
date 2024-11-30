// Package product_handlers contains the database store implementations for managing product records.
package product_handlers

import (
	"database/sql"
	"erp/models"
	"fmt"
)

// DBProductStore implements the ProductStore interface for database operations.
type DBProductStore struct {
	DB *sql.DB
}

// NewDBProductStore initializes a new DBProductStore instance.
//
// Parameters:
// - db: A *sql.DB instance representing the database connection.
//
// Returns:
// - A pointer to an instance of DBProductStore.
func NewDBProductStore(db *sql.DB) *DBProductStore {
	return &DBProductStore{DB: db}
}

// CreateProduct inserts a new product record into the database.
//
// Parameters:
// - product: A pointer to the Product struct containing the product details to insert.
//
// Returns:
// - An error if the insertion fails, otherwise nil.
func (s *DBProductStore) CreateProduct(product *models.Product) error {
	query := `
		INSERT INTO products (name, brand, season, price)
		VALUES ($1, $2, $3, $4)
		RETURNING id
	`
	err := s.DB.QueryRow(query, product.Name, product.Brand, product.Season, product.Price).Scan(&product.ID)
	if err != nil {
		return fmt.Errorf("failed to insert product: %w", err)
	}
	return nil
}

// GetProductByID retrieves a product record from the database by ID.
//
// Parameters:
// - id: An integer representing the product ID.
//
// Returns:
// - A pointer to the Product struct if found.
// - An error if no record is found or if the query fails.
func (s *DBProductStore) GetProductByID(id int) (*models.Product, error) {
	query := `
		SELECT id, name, brand, season, price
		FROM products
		WHERE id = $1
	`
	row := s.DB.QueryRow(query, id)

	var product models.Product
	err := row.Scan(&product.ID, &product.Name, &product.Brand, &product.Season, &product.Price)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("no product found with ID %d", id)
		}
		return nil, fmt.Errorf("failed to retrieve product: %w", err)
	}

	return &product, nil
}

// UpdateProduct updates an existing product record in the database.
//
// Parameters:
// - product: A pointer to the Product struct containing the updated product details.
//
// Returns:
// - An error if the update fails, otherwise nil.
func (s *DBProductStore) UpdateProduct(product *models.Product) error {
	query := `
		UPDATE products
		SET name = $1, brand = $2, season = $3, price = $4
		WHERE id = $5
	`
	result, err := s.DB.Exec(query, product.Name, product.Brand, product.Season, product.Price, product.ID)
	if err != nil {
		return fmt.Errorf("failed to update product: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to check rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("no product found with ID %d", product.ID)
	}

	return nil
}

// DeleteProduct removes a product record from the database by ID.
//
// Parameters:
// - id: An integer representing the product ID to delete.
//
// Returns:
// - An error if the deletion fails, otherwise nil.
func (s *DBProductStore) DeleteProduct(id int) error {
	query := `
		DELETE FROM products
		WHERE id = $1
	`
	result, err := s.DB.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete product with ID %d: %w", id, err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to check rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("no product found with ID %d", id)
	}

	return nil
}

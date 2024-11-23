package product_handlers

import (
	"database/sql"
	"erp/models"
	"fmt"
)

// ProductStoreImpl implements the ProductStore interface using a SQL database.
//
// It provides methods for creating, retrieving, updating, and deleting products
// in the database using SQL queries. This implementation assumes that a valid
// *sql.DB connection is provided.
type ProductStoreImpl struct {
	DB *sql.DB
}

// NewProductStoreImpl creates a new instance of ProductStoreImpl.
//
// It accepts a database connection and returns a ProductStore instance that
// can be used for managing products in the inventory.
//
// Parameters:
// - db: A pointer to an open SQL database connection.
//
// Returns:
// - A ProductStore implementation.
func NewProductStoreImpl(db *sql.DB) *ProductStoreImpl {
	return &ProductStoreImpl{
		DB: db,
	}
}

// CreateProduct adds a new product to the database.
//
// It inserts the provided product details into the products table. If an error
// occurs during insertion, the error is returned.
//
// Parameters:
// - product: A pointer to the Product object containing product details.
//
// Returns:
// - An error if the insertion fails; nil otherwise.
func (store *ProductStoreImpl) CreateProduct(product *models.Product) error {
	query := `INSERT INTO products (name, brand, season, price) VALUES ($1, $2, $3, $4) RETURNING id`
	err := store.DB.QueryRow(query, product.Name, product.Brand, product.Season, product.Price).Scan(&product.ID)
	if err != nil {
		return fmt.Errorf("failed to create product: %w", err)
	}
	return nil
}

// GetProductByID retrieves a product by its ID from the database.
//
// It queries the products table for a product with the specified ID and returns
// a pointer to the Product object if found. If no product is found or an error
// occurs, an error is returned.
//
// Parameters:
// - id: The ID of the product to retrieve.
//
// Returns:
// - A pointer to the Product object if found; nil and an error otherwise.
func (store *ProductStoreImpl) GetProductByID(id int) (*models.Product, error) {
	query := `SELECT id, name, brand, season, price FROM products WHERE id = $1`
	row := store.DB.QueryRow(query, id)

	product := &models.Product{}
	err := row.Scan(&product.ID, &product.Name, &product.Brand, &product.Season, &product.Price)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("product with ID %d not found", id)
		}
		return nil, fmt.Errorf("failed to retrieve product: %w", err)
	}
	return product, nil
}

// UpdateProduct updates an existing product in the database.
//
// It modifies the product details with the specified ID based on the provided
// Product object. If an error occurs during the update, it is returned.
//
// Parameters:
// - product: A pointer to the Product object containing updated product details.
//
// Returns:
// - An error if the update fails; nil otherwise.
func (store *ProductStoreImpl) UpdateProduct(product *models.Product) error {
	query := `UPDATE products SET name = $1, brand = $2, season = $3, price = $4 WHERE id = $5`
	result, err := store.DB.Exec(query, product.Name, product.Brand, product.Season, product.Price, product.ID)
	if err != nil {
		return fmt.Errorf("failed to update product: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to check rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("no product with ID %d found", product.ID)
	}
	return nil
}

// DeleteProduct removes a product by its ID from the database.
//
// It deletes the product record with the specified ID. If an error occurs during
// deletion or no record is found, an error is returned.
//
// Parameters:
// - id: The ID of the product to delete.
//
// Returns:
// - An error if the deletion fails or if no matching record is found.
func (store *ProductStoreImpl) DeleteProduct(id int) error {
	query := `DELETE FROM products WHERE id = $1`
	result, err := store.DB.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete product: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to check rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("no product with ID %d found", id)
	}
	return nil
}

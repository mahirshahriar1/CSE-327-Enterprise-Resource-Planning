// Package warehouse_store contains the database operations for warehouse management.
package warehouse_handlers

import (
	"database/sql"
	"erp/models"
	"errors"
)

// DBWarehouseStore implements WarehouseStore using a SQL database.
type DBWarehouseStore struct {
	DB *sql.DB
}

// CreateWarehouse inserts a new warehouse into the database.
//
// This method stores the provided Warehouse object in the database.
// It returns an error if the operation fails.
//
// Parameters:
// - warehouse: A pointer to the Warehouse object to be created.
//
// Returns:
// - nil if the warehouse is created successfully.
// - An error if the creation fails.
func (s *DBWarehouseStore) CreateWarehouse(warehouse *models.Warehouse) error {
	_, err := s.DB.Exec(
		"INSERT INTO warehouses (name, capacity, location) VALUES ($1, $2, $3)",
		warehouse.Name, warehouse.Capacity, warehouse.Location,
	)
	if err != nil {
		return errors.New("failed to create warehouse: " + err.Error())
	}
	return nil
}

// GetWarehouseByID retrieves a warehouse from the database by its ID.
//
// Parameters:
// - id: The ID of the warehouse to retrieve.
//
// Returns:
// - A pointer to the Warehouse object if found.
// - An error if the warehouse is not found or the operation fails.
func (s *DBWarehouseStore) GetWarehouseByID(id int) (*models.Warehouse, error) {
	var warehouse models.Warehouse
	err := s.DB.QueryRow(
		"SELECT id, name, capacity, location FROM warehouses WHERE id = $1",
		id,
	).Scan(&warehouse.ID, &warehouse.Name, &warehouse.Capacity, &warehouse.Location)

	if err == sql.ErrNoRows {
		return nil, errors.New("warehouse not found")
	} else if err != nil {
		return nil, errors.New("failed to retrieve warehouse: " + err.Error())
	}
	return &warehouse, nil
}

// UpdateWarehouse updates an existing warehouse in the database.
//
// This method modifies the details of an existing warehouse based on the provided Warehouse object.
// It returns an error if the update fails.
//
// Parameters:
// - warehouse: A pointer to the Warehouse object containing updated details.
//
// Returns:
// - nil if the warehouse is updated successfully.
// - An error if the update fails.
func (s *DBWarehouseStore) UpdateWarehouse(warehouse *models.Warehouse) error {
	_, err := s.DB.Exec(
		"UPDATE warehouses SET name = $1, capacity = $2, location = $3 WHERE id = $4",
		warehouse.Name, warehouse.Capacity, warehouse.Location, warehouse.ID,
	)
	if err != nil {
		return errors.New("failed to update warehouse: " + err.Error())
	}
	return nil
}

// DeleteWarehouse removes a warehouse from the database by its ID.
//
// Parameters:
// - id: The ID of the warehouse to delete.
//
// Returns:
// - nil if the warehouse is deleted successfully.
// - An error if the deletion fails.
func (s *DBWarehouseStore) DeleteWarehouse(id int) error {
	_, err := s.DB.Exec("DELETE FROM warehouses WHERE id = $1", id)
	if err != nil {
		return errors.New("failed to delete warehouse: " + err.Error())
	}
	return nil
}

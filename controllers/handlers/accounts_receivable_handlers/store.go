// Package accounts_receivable_handlers provides SQL-backed methods to manage accounts receivable.
// It defines a DBReceivableStore struct for database operations on receivable records.
package accounts_receivable_handlers

import (
	"database/sql"
	"erp/models"
	"fmt"
)

// DBReceivableStore provides SQL-backed methods for managing receivables.
// It uses a database connection to perform CRUD operations.
type DBReceivableStore struct {
	// DB is the database connection used for executing SQL queries.
	DB *sql.DB
}

// CreateReceivable inserts a new receivable into the database and assigns the generated ID to the input receivable.
//
// Parameters:
//   - receivable: A pointer to the Receivable object containing the details of the receivable to be created.
//
// Returns:
//   - An error if the operation fails, or nil if the receivable is successfully created.
func (store *DBReceivableStore) CreateReceivable(receivable *models.Receivable) error {
	return store.DB.QueryRow(
		"INSERT INTO receivables (customer_name, amount, due_date, invoice_number) VALUES ($1, $2, $3, $4) RETURNING id",
		receivable.CustomerName, receivable.Amount, receivable.DueDate, receivable.InvoiceNumber,
	).Scan(&receivable.ID)
}

// GetReceivableByID retrieves a receivable record from the database using its ID.
//
// Parameters:
//   - id: The unique identifier of the receivable to be retrieved.
//
// Returns:
//   - A pointer to the Receivable object if found.
//   - An error if the receivable does not exist or if the operation fails.
func (store *DBReceivableStore) GetReceivableByID(id int) (*models.Receivable, error) {
	row := store.DB.QueryRow("SELECT id, customer_name, amount, due_date, invoice_number FROM receivables WHERE id = $1", id)

	var receivable models.Receivable
	err := row.Scan(&receivable.ID, &receivable.CustomerName, &receivable.Amount, &receivable.DueDate, &receivable.InvoiceNumber)
	if err != nil {
		return nil, err
	}
	return &receivable, nil
}

// UpdateReceivable updates an existing receivable record in the database.
//
// Parameters:
//   - receivable: A pointer to the Receivable object containing updated details. The ID field must be set.
//
// Returns:
//   - An error if the operation fails, or if no rows are affected (indicating the receivable does not exist).
func (store *DBReceivableStore) UpdateReceivable(receivable *models.Receivable) error {
	result, err := store.DB.Exec(
		"UPDATE receivables SET customer_name = $1, amount = $2, due_date = $3, invoice_number = $4 WHERE id = $5",
		receivable.CustomerName, receivable.Amount, receivable.DueDate, receivable.InvoiceNumber, receivable.ID,
	)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return fmt.Errorf("receivable with ID %d does not exist", receivable.ID)
	}

	return nil
}

// DeleteReceivable deletes a receivable record from the database using its ID.
//
// Parameters:
//   - id: The unique identifier of the receivable to be deleted.
//
// Returns:
//   - An error if the operation fails, or if no rows are affected (indicating the receivable does not exist).
func (store *DBReceivableStore) DeleteReceivable(id int) error {
	result, err := store.DB.Exec("DELETE FROM receivables WHERE id = $1", id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return fmt.Errorf("receivable with ID %d does not exist", id)
	}

	return nil
}

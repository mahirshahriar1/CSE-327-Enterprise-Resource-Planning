// Package accounts_receivable_handlers provides SQL-backed methods to manage receivables.
package accounts_receivable_handlers

import (
	"database/sql"
	"erp/models"
	"fmt"
)

// DBReceivableStore provides SQL-backed methods for managing receivables.
type DBReceivableStore struct {
	DB *sql.DB // DB represents the database connection.
}

// CreateReceivable inserts a new receivable into the database.
func (store *DBReceivableStore) CreateReceivable(receivable *models.Receivable) error {
	return store.DB.QueryRow(
		"INSERT INTO receivables (customer_name, amount, due_date, invoice_number) VALUES ($1, $2, $3, $4) RETURNING id",
		receivable.CustomerName, receivable.Amount, receivable.DueDate, receivable.InvoiceNumber,
	).Scan(&receivable.ID)
}

// GetReceivableByID retrieves a receivable by its ID.
func (store *DBReceivableStore) GetReceivableByID(id int) (*models.Receivable, error) {
	row := store.DB.QueryRow("SELECT id, customer_name, amount, due_date, invoice_number FROM receivables WHERE id = $1", id)

	var receivable models.Receivable
	err := row.Scan(&receivable.ID, &receivable.CustomerName, &receivable.Amount, &receivable.DueDate, &receivable.InvoiceNumber)
	if err != nil {
		return nil, err
	}
	return &receivable, nil
}

// UpdateReceivable updates an existing receivable record.
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

// DeleteReceivable deletes a receivable by its ID.
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

// Package invoice_handlers provides the database implementation for invoice-related operations.
// It includes functionality for creating, retrieving, updating, and deleting invoices in the database.
package invoice_handlers

import (
	"database/sql"
	"erp/models"
	"errors"
)

// DBInvoiceStore is a struct to hold the database connection for invoice operations.
type DBInvoiceStore struct {
	DB *sql.DB
}

// CreateInvoice inserts a new invoice into the database.
func (store *DBInvoiceStore) CreateInvoice(invoice *models.Invoice) error {
	query := `
        INSERT INTO invoices (sales_order_id, customer_id, amount, status)
        VALUES ($1, $2, $3, $4)
        RETURNING id
    `
	err := store.DB.QueryRow(query, invoice.SalesOrderID, invoice.CustomerID, invoice.Amount, invoice.Status).Scan(&invoice.ID)
	if err != nil {
		return err
	}
	return nil
}

// GetInvoiceByID retrieves an invoice by its ID from the database.
func (store *DBInvoiceStore) GetInvoiceByID(id int) (*models.Invoice, error) {
	query := `
        SELECT id, sales_order_id, customer_id, amount, status
        FROM invoices
        WHERE id = $1
    `
	invoice := &models.Invoice{}
	err := store.DB.QueryRow(query, id).Scan(&invoice.ID, &invoice.SalesOrderID, &invoice.CustomerID, &invoice.Amount, &invoice.Status)
	if err == sql.ErrNoRows {
		return nil, errors.New("invoice not found")
	} else if err != nil {
		return nil, err
	}
	return invoice, nil
}

// UpdateInvoice updates an existing invoice's details in the database.
func (store *DBInvoiceStore) UpdateInvoice(invoice *models.Invoice) error {
	query := `
        UPDATE invoices
        SET sales_order_id = $1, customer_id = $2, amount = $3, status = $4
        WHERE id = $5
    `
	_, err := store.DB.Exec(query, invoice.SalesOrderID, invoice.CustomerID, invoice.Amount, invoice.Status, invoice.ID)
	if err != nil {
		return err
	}
	return nil
}

// DeleteInvoice deletes an invoice from the database by its ID.
func (store *DBInvoiceStore) DeleteInvoice(id int) error {
	query := `
        DELETE FROM invoices
        WHERE id = $1
    `
	_, err := store.DB.Exec(query, id)
	if err != nil {
		return err
	}
	return nil
}

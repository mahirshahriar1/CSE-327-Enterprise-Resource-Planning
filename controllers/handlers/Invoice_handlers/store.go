package invoice_handlers

import (
	"database/sql"
	"erp/models" // Adjust import path as needed
	"errors"
)

// DBInvoiceStore is a struct to hold the database connection for invoice operations
type DBInvoiceStore struct {
	DB *sql.DB
}

// CreateInvoice creates a new invoice in the database
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

// GetInvoiceByID retrieves an invoice by its ID from the database
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



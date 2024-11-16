// Package accounts_payable_handlers provides SQL-backed methods to manage payments
// for the accounts payable system. It includes functionalities to create, retrieve,
// update, and delete payments in the database.
package accounts_payable_handlers

import (
	"database/sql"
	"erp/models"
	"fmt"
)

// DBPaymentStore provides SQL-backed methods to manage payments.
// It implements the CRUD operations required for managing the `payments` table
// in the database.
type DBPaymentStore struct {
	DB *sql.DB // DB represents the database connection.
}

// CreatePayment inserts a new payment into the database.
// It generates a new ID for the payment and stores it in the provided `Payment` object.
//
// Parameters:
//   - payment: A pointer to the `Payment` object containing the payment details to be stored.
//
// Returns:
//   - error: An error if the query fails or the insertion is unsuccessful.
func (store *DBPaymentStore) CreatePayment(payment *models.Payment) error {
	return store.DB.QueryRow(
		"INSERT INTO payments (invoice_id, amount, payment_date, payment_method) VALUES ($1, $2, $3, $4) RETURNING id",
		payment.InvoiceID, payment.Amount, payment.PaymentDate, payment.PaymentMethod,
	).Scan(&payment.ID)
}

// GetPaymentByID retrieves a payment by its ID from the database.
//
// Parameters:
//   - id: The ID of the payment to retrieve.
//
// Returns:
//   - *Payment: A pointer to the `Payment` object containing the retrieved payment details.
//   - error: An error if the query fails or no payment is found with the provided ID.
func (store *DBPaymentStore) GetPaymentByID(id int) (*models.Payment, error) {
	row := store.DB.QueryRow("SELECT id, invoice_id, amount, payment_date, payment_method FROM payments WHERE id = $1", id)

	var payment models.Payment
	err := row.Scan(&payment.ID, &payment.InvoiceID, &payment.Amount, &payment.PaymentDate, &payment.PaymentMethod)
	if err != nil {
		return nil, err
	}
	return &payment, nil
}

// UpdatePayment updates an existing payment in the database.
//
// Parameters:
//   - payment: A pointer to the `Payment` object containing the updated payment details.
//
// Returns:
//   - error: An error if the query fails or if no payment exists with the provided ID.
func (store *DBPaymentStore) UpdatePayment(payment *models.Payment) error {
	result, err := store.DB.Exec(
		"UPDATE payments SET invoice_id = $1, amount = $2, payment_date = $3, payment_method = $4 WHERE id = $5",
		payment.InvoiceID, payment.Amount, payment.PaymentDate, payment.PaymentMethod, payment.ID,
	)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return fmt.Errorf("payment with ID %d does not exist", payment.ID)
	}

	return nil
}

// DeletePayment deletes a payment by its ID from the database.
//
// Parameters:
//   - id: The ID of the payment to delete.
//
// Returns:
//   - error: An error if the query fails or if no payment exists with the provided ID.
func (store *DBPaymentStore) DeletePayment(id int) error {
	result, err := store.DB.Exec("DELETE FROM payments WHERE id = $1", id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return fmt.Errorf("payment with ID %d does not exist", id)
	}

	return nil
}

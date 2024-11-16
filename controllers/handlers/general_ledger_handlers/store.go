package general_ledger_handlers

import (
	"database/sql"
	"erp/models"
	"fmt"
)

// DBFinancialTransactionStore provides SQL-backed methods to manage financial transactions.
// It acts as a store for interacting with the financial_transactions table in the database.
type DBFinancialTransactionStore struct {
	DB *sql.DB // DB represents the database connection.
}

// CreateTransaction inserts a new financial transaction into the database.
// It populates the ID of the transaction with the auto-generated ID from the database.
//
// Parameters:
//   - transaction: A pointer to the FinancialTransaction object containing transaction details.
//
// Returns:
//   - error: An error object if the transaction fails to be created, otherwise nil.
func (store *DBFinancialTransactionStore) CreateTransaction(transaction *models.FinancialTransaction) error {
	err := store.DB.QueryRow(
		"INSERT INTO financial_transactions (account_type, amount, transaction_date) VALUES ($1, $2, $3) RETURNING id",
		transaction.AccountType, transaction.Amount, transaction.TransactionDate,
	).Scan(&transaction.ID) // Scan the generated ID into the transaction.ID field

	return err
}

// GetTransactionByID retrieves a financial transaction from the database by its ID.
//
// Parameters:
//   - id: The ID of the transaction to retrieve.
//
// Returns:
//   - *FinancialTransaction: A pointer to the retrieved transaction object.
//   - error: An error object if the retrieval fails or if the transaction does not exist.
func (store *DBFinancialTransactionStore) GetTransactionByID(id int) (*models.FinancialTransaction, error) {
	row := store.DB.QueryRow("SELECT id, account_type, amount, transaction_date FROM financial_transactions WHERE id = $1", id)

	var transaction models.FinancialTransaction
	err := row.Scan(&transaction.ID, &transaction.AccountType, &transaction.Amount, &transaction.TransactionDate)
	if err != nil {
		return nil, err
	}
	return &transaction, nil
}

// UpdateTransaction updates an existing financial transaction in the database.
//
// Parameters:
//   - transaction: A pointer to the FinancialTransaction object containing updated transaction details.
//
// Returns:
//   - error: An error object if the update fails, or if the transaction ID does not exist.
func (store *DBFinancialTransactionStore) UpdateTransaction(transaction *models.FinancialTransaction) error {
	result, err := store.DB.Exec(
		"UPDATE financial_transactions SET account_type = $1, amount = $2, transaction_date = $3 WHERE id = $4",
		transaction.AccountType, transaction.Amount, transaction.TransactionDate, transaction.ID,
	)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return fmt.Errorf("transaction with ID %d does not exist", transaction.ID)
	}

	return nil
}

// DeleteTransaction deletes a financial transaction from the database by its ID.
//
// Parameters:
//   - id: The ID of the transaction to delete.
//
// Returns:
//   - error: An error object if the deletion fails, or if the transaction ID does not exist.
func (store *DBFinancialTransactionStore) DeleteTransaction(id int) error {
	result, err := store.DB.Exec("DELETE FROM financial_transactions WHERE id = $1", id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return fmt.Errorf("transaction with ID %d does not exist", id)
	}

	return nil
}

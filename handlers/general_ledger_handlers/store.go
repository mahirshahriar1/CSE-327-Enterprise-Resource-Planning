package general_ledger_handlers

import (
	"database/sql"
	"erp/models"
	"fmt"
)

// DBFinancialTransactionStore provides SQL-backed methods to manage financial transactions
type DBFinancialTransactionStore struct {
	DB *sql.DB
}

// CreateTransaction inserts a new transaction into the database and returns the generated ID
func (store *DBFinancialTransactionStore) CreateTransaction(transaction *models.FinancialTransaction) error {
	err := store.DB.QueryRow(
		"INSERT INTO financial_transactions (account_type, amount, transaction_date) VALUES ($1, $2, $3) RETURNING id",
		transaction.AccountType, transaction.Amount, transaction.TransactionDate,
	).Scan(&transaction.ID) // Scan the generated ID into the transaction.ID field

	return err
}

// GetTransactionByID retrieves a transaction by ID
func (store *DBFinancialTransactionStore) GetTransactionByID(id int) (*models.FinancialTransaction, error) {
	row := store.DB.QueryRow("SELECT id, account_type, amount, transaction_date FROM financial_transactions WHERE id = $1", id)

	var transaction models.FinancialTransaction
	err := row.Scan(&transaction.ID, &transaction.AccountType, &transaction.Amount, &transaction.TransactionDate)
	if err != nil {
		return nil, err
	}
	return &transaction, nil
}

// UpdateTransaction updates an existing transaction and returns an error if the transaction ID doesn't exist
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

// DeleteTransaction deletes a transaction by ID and returns an error if the ID does not exist
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

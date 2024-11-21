// Package financial_record_handlers provides SQL-backed methods to manage financial records.
package financial_record_handlers

import (
	"database/sql"
	"erp/models"
	"fmt"
)

// DBFinancialRecordStore provides SQL-backed methods to manage financial records.
// It implements the CRUD operations required for managing the `financial_records` table
// in the database.
type DBFinancialRecordStore struct {
	DB *sql.DB // DB represents the database connection.
}

// CreateFinancialRecord inserts a new financial record into the database.
// It generates a new ID for the financial record and stores it in the provided `FinancialRecord` object.
func (store *DBFinancialRecordStore) CreateFinancialRecord(financialRecord *models.FinancialRecord) error {
	return store.DB.QueryRow(
		"INSERT INTO financial_records (transaction_id, account_id, amount, transaction_date, transaction_type, description) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id",
		financialRecord.TransactionID, financialRecord.AccountID, financialRecord.Amount, financialRecord.TransactionDate, financialRecord.TransactionType, financialRecord.Description,
	).Scan(&financialRecord.ID)
}

// GetFinancialRecordByID retrieves a financial record by its ID from the database.
func (store *DBFinancialRecordStore) GetFinancialRecordByID(id int) (*models.FinancialRecord, error) {
	row := store.DB.QueryRow("SELECT id, transaction_id, account_id, amount, transaction_date, transaction_type, description FROM financial_records WHERE id = $1", id)

	var financialRecord models.FinancialRecord
	err := row.Scan(&financialRecord.ID, &financialRecord.TransactionID, &financialRecord.AccountID, &financialRecord.Amount, &financialRecord.TransactionDate, &financialRecord.TransactionType, &financialRecord.Description)
	if err != nil {
		return nil, err
	}
	return &financialRecord, nil
}

// UpdateFinancialRecord updates an existing financial record in the database.
func (store *DBFinancialRecordStore) UpdateFinancialRecord(financialRecord *models.FinancialRecord) error {
	result, err := store.DB.Exec(
		"UPDATE financial_records SET transaction_id = $1, account_id = $2, amount = $3, transaction_date = $4, transaction_type = $5, description = $6 WHERE id = $7",
		financialRecord.TransactionID, financialRecord.AccountID, financialRecord.Amount, financialRecord.TransactionDate, financialRecord.TransactionType, financialRecord.Description, financialRecord.ID,
	)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return fmt.Errorf("financial record with ID %d does not exist", financialRecord.ID)
	}

	return nil
}

// DeleteFinancialRecord deletes a financial record by its ID from the database.
func (store *DBFinancialRecordStore) DeleteFinancialRecord(id int) error {
	result, err := store.DB.Exec("DELETE FROM financial_records WHERE id = $1", id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return fmt.Errorf("financial record with ID %d does not exist", id)
	}

	return nil
}

// Package financial_record_handlers provides SQL-backed methods to manage financial records.
// It includes a DBFinancialRecordStore struct that implements the CRUD operations required
// for the `financial_records` table in the database.
package financial_record_handlers

import (
	"database/sql"
	"erp/models"
	"fmt"
)

// DBFinancialRecordStore provides SQL-backed methods to manage financial records.
// This struct interacts with the `financial_records` table in the database to perform
// create, read, update, and delete (CRUD) operations.
type DBFinancialRecordStore struct {
	// DB represents the database connection used for executing SQL queries.
	DB *sql.DB
}

// CreateFinancialRecord inserts a new financial record into the database and assigns
// the generated ID to the provided FinancialRecord object.
//
// Parameters:
//   - financialRecord: A pointer to the FinancialRecord object containing the details of the record to be created.
//
// Returns:
//   - An error if the operation fails, or nil if the record is successfully created.
func (store *DBFinancialRecordStore) CreateFinancialRecord(financialRecord *models.FinancialRecord) error {
	return store.DB.QueryRow(
		"INSERT INTO financial_records (transaction_id, account_id, amount, transaction_date, transaction_type, description) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id",
		financialRecord.TransactionID, financialRecord.AccountID, financialRecord.Amount, financialRecord.TransactionDate, financialRecord.TransactionType, financialRecord.Description,
	).Scan(&financialRecord.ID)
}

// GetFinancialRecordByID retrieves a financial record from the database by its ID.
//
// Parameters:
//   - id: The unique identifier of the financial record to be retrieved.
//
// Returns:
//   - A pointer to the FinancialRecord object if the record is found.
//   - An error if the record does not exist or if the operation fails.
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
//
// Parameters:
//   - financialRecord: A pointer to the FinancialRecord object containing updated details. The ID field must be set.
//
// Returns:
//   - An error if the operation fails, or if no rows are affected (indicating the record does not exist).
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

// DeleteFinancialRecord deletes a financial record from the database by its ID.
//
// Parameters:
//   - id: The unique identifier of the financial record to be deleted.
//
// Returns:
//   - An error if the operation fails, or if no rows are affected (indicating the record does not exist).
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

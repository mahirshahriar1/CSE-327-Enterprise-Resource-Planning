package general_ledger_handlers

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

// func TestCreateTransaction(t *testing.T) {
// 	// Set up mock database
// 	db, mock, err := sqlmock.New()
// 	if err != nil {
// 		t.Fatalf("could not create mock db: %v", err)
// 	}
// 	defer db.Close()

// 	store := &DBFinancialTransactionStore{DB: db}

// 	// Sample transaction data
// 	transaction := &models.FinancialTransaction{
// 		Amount:          1000.50,
// 		Description:     "Payment received",
// 		TransactionDate: time.Now(),
// 	}

// 	// Define expected behavior for mock (ensure types are correct)
// 	mock.ExpectQuery("INSERT INTO financial_transactions").
// 		WithArgs(transaction.Amount, transaction.Description, transaction.TransactionDate).
// 		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

// 	// Call the method
// 	err = store.CreateTransaction(transaction)

// 	// Assert that no error occurred
// 	assert.NoError(t, err)
// 	assert.Equal(t, 1, transaction.ID)

// 	// Assert that the expected query was executed
// 	if err := mock.ExpectationsWereMet(); err != nil {
// 		t.Errorf("there were unmet expectations: %v", err)
// 	}
// }

// func TestUpdateTransaction(t *testing.T) {
// 	// Set up mock database
// 	db, mock, err := sqlmock.New()
// 	if err != nil {
// 		t.Fatalf("could not create mock db: %v", err)
// 	}
// 	defer db.Close()

// 	store := &DBFinancialTransactionStore{DB: db}

// 	// Sample transaction data
// 	transaction := &models.FinancialTransaction{
// 		ID:              1,
// 		Amount:          1200.75,
// 		Description:     "Payment updated",
// 		TransactionDate: time.Now(),
// 	}

// 	// Define expected behavior for mock (ensure types are correct)
// 	mock.ExpectExec("UPDATE financial_transactions").
// 		WithArgs(transaction.Amount, transaction.Description, transaction.TransactionDate, transaction.ID).
// 		WillReturnResult(sqlmock.NewResult(1, 1))

// 	// Call the method
// 	err = store.UpdateTransaction(transaction)

// 	// Assert that no error occurred
// 	assert.NoError(t, err)

// 	// Assert that the expected query was executed
// 	if err := mock.ExpectationsWereMet(); err != nil {
// 		t.Errorf("there were unmet expectations: %v", err)
// 	}
// }

func TestDeleteTransaction(t *testing.T) {
	// Set up mock database
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("could not create mock db: %v", err)
	}
	defer db.Close()

	store := &DBFinancialTransactionStore{DB: db}

	// Define expected behavior for mock
	mock.ExpectExec("DELETE FROM financial_transactions").
		WithArgs(1).
		WillReturnResult(sqlmock.NewResult(1, 1))

	// Call the method
	err = store.DeleteTransaction(1)

	// Assert that no error occurred
	assert.NoError(t, err)

	// Assert that the expected query was executed
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unmet expectations: %v", err)
	}
}

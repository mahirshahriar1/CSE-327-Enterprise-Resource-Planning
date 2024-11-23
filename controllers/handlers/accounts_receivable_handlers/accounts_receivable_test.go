package accounts_receivable_handlers

import (
	"erp/models"
	"testing"
	"time"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestCreateReceivable(t *testing.T) {
	// Set up mock database
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("could not create mock db: %v", err)
	}
	defer db.Close()

	store := &DBReceivableStore{DB: db}

	// Sample receivable data
	receivable := &models.Receivable{
		CustomerName:  "Test Customer",
		Amount:        100.50,
		DueDate:       time.Now(),
		InvoiceNumber: "INV12345",
	}

	// Define expected behavior for mock
	mock.ExpectQuery("INSERT INTO receivables").
		WithArgs(receivable.CustomerName, receivable.Amount, receivable.DueDate, receivable.InvoiceNumber).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

	// Call the method
	err = store.CreateReceivable(receivable)

	// Assert that no error occurred
	assert.NoError(t, err)
	assert.Equal(t, 1, receivable.ID)

	// Assert that the expected query was executed
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unmet expectations: %v", err)
	}
}

func TestUpdateReceivable(t *testing.T) {
	// Set up mock database
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("could not create mock db: %v", err)
	}
	defer db.Close()

	store := &DBReceivableStore{DB: db}

	// Sample receivable data
	receivable := &models.Receivable{
		ID:            1,
		CustomerName:  "Test Customer",
		Amount:        150.75,
		DueDate:       time.Now(),
		InvoiceNumber: "INV12345",
	}

	// Define expected behavior for mock
	mock.ExpectExec("UPDATE receivables").
		WithArgs(receivable.CustomerName, receivable.Amount, receivable.DueDate, receivable.InvoiceNumber, receivable.ID).
		WillReturnResult(sqlmock.NewResult(1, 1))

	// Call the method
	err = store.UpdateReceivable(receivable)

	// Assert that no error occurred
	assert.NoError(t, err)

	// Assert that the expected query was executed
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unmet expectations: %v", err)
	}
}

func TestDeleteReceivable(t *testing.T) {
	// Set up mock database
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("could not create mock db: %v", err)
	}
	defer db.Close()

	store := &DBReceivableStore{DB: db}

	// Define expected behavior for mock
	mock.ExpectExec("DELETE FROM receivables").
		WithArgs(1).
		WillReturnResult(sqlmock.NewResult(1, 1))

	// Call the method
	err = store.DeleteReceivable(1)

	// Assert that no error occurred
	assert.NoError(t, err)

	// Assert that the expected query was executed
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unmet expectations: %v", err)
	}
}

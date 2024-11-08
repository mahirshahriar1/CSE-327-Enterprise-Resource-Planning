package models

import "time"

// FinancialTransaction represents a financial transaction in the system
type FinancialTransaction struct {
	ID            int       `json:"id"`
	AccountType   string    `json:"account_type"`
	Amount        float64   `json:"amount"`
	TransactionDate time.Time `json:"transaction_date"`
}

// FinancialTransactionStore defines an interface for financial transaction-related database operations
type FinancialTransactionStore interface {
	CreateTransaction(transaction *FinancialTransaction) error
	GetTransactionByID(id int) (*FinancialTransaction, error)
	UpdateTransaction(transaction *FinancialTransaction) error
	DeleteTransaction(id int) error
}

// Package models defines the structure for financial records and their store interface.
package models

import "time"

// FinancialRecord represents a financial record.
type FinancialRecord struct {
	ID              int       `json:"id"`
	TransactionID   int       `json:"transaction_id"`
	AccountID       int       `json:"account_id"`
	Amount          float64   `json:"amount"`
	TransactionDate time.Time `json:"transaction_date"`
	TransactionType string    `json:"transaction_type"`
	Description     string    `json:"description"`
}

// FinancialRecordStore is an interface that defines CRUD operations for financial records.
type FinancialRecordStore interface {
	CreateFinancialRecord(record *FinancialRecord) error
	GetFinancialRecordByID(id int) (*FinancialRecord, error)
	UpdateFinancialRecord(record *FinancialRecord) error
	DeleteFinancialRecord(id int) error
}

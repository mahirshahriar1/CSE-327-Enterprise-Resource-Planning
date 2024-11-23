// Package models defines the data models for various entities in the ERP system.
package models

import "time"

// Receivable represents a customer's receivable record, which tracks amounts due from customers.
type Receivable struct {
	ID            int       `json:"id"`             // Unique ID for the receivable entry
	CustomerName  string    `json:"customer_name"`  // Name of the customer who owes the amount
	Amount        float64   `json:"amount"`         // The amount the customer owes
	DueDate       time.Time `json:"due_date"`       // The date when the payment is due
	InvoiceNumber string    `json:"invoice_number"` // Unique invoice number for the receivable
	Status        string    `json:"status"`         // Current status of the receivable (e.g., "pending", "paid", "overdue")
	PaymentDate   time.Time `json:"payment_date"`   // Date when the payment was received (if applicable)
}

// ReceivableStore is the interface that wraps methods for managing receivable records.
type ReceivableStore interface {
	CreateReceivable(receivable *Receivable) error
	GetReceivableByID(id int) (*Receivable, error)
	UpdateReceivable(receivable *Receivable) error
	DeleteReceivable(id int) error
	GetAllReceivables() ([]Receivable, error)
}

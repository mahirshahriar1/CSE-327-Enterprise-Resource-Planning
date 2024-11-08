package models

import "time"

// Payment represents a payment in the system
type Payment struct {
	ID           int       `json:"id"`
	InvoiceID    int       `json:"invoice_id"`
	Amount       float64   `json:"amount"`
	PaymentDate  time.Time `json:"payment_date"`
	PaymentMethod string   `json:"payment_method"`
}

// PaymentStore defines an interface for payment-related database operations
type PaymentStore interface {
	CreatePayment(payment *Payment) error
	GetPaymentByID(id int) (*Payment, error)
	UpdatePayment(payment *Payment) error
	DeletePayment(id int) error
}

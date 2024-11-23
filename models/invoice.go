package models

// Invoice represents an invoice in the system
type Invoice struct {
	ID           int     `json:"id"`
	SalesOrderID int     `json:"sales_order_id"`
	CustomerID   int     `json:"customer_id"`
	Amount       float64 `json:"amount"`
	Status       string  `json:"status"`
}

// InvoiceStore defines an interface for invoice-related database operations
type InvoiceStore interface {
	CreateInvoice(invoice *Invoice) error
	GetInvoiceByID(id int) (*Invoice, error)
	//UpdateInvoice(invoice *Invoice) error
	//DeleteInvoice(id int) error
	// UpdateInvoice(invoice *Invoice) error
	// DeleteInvoice(id int) error
}

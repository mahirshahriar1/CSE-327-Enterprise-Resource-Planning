package models

import "time"

// SalesOrder represents a sales order in the system
type SalesOrder struct {
	ID         int       `json:"id"`
	CustomerID int       `json:"customer_id"`
	ProductID  int       `json:"product_id"`
	OrderDate  time.Time `json:"order_date"`
	Quantity   int       `json:"quantity"`
}

// SalesOrderStore defines an interface for sales order-related database operations
type SalesOrderStore interface {
	CreateSalesOrder(order *SalesOrder) error
	GetSalesOrderByID(id int) (*SalesOrder, error)
	UpdateSalesOrder(order *SalesOrder) error
	DeleteSalesOrder(id int) error
}

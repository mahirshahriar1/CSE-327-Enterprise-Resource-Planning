package models

// Customer represents a customer in the system
type Customer struct {
	ID           int    `json:"id"`
	Name         string `json:"name"`
	Contact      string `json:"contact"`
	OrderHistory string `json:"order_history"`
}

// CustomerStore defines an interface for customer-related database operations
type CustomerStore interface {
	CreateCustomer(customer *Customer) error
	GetCustomerByID(id int) (*Customer, error)
	UpdateCustomer(customer *Customer) error
	DeleteCustomer(id int) error
}

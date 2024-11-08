package models

// Product represents a product in the inventory
type Product struct {
	ID      int     `json:"id"`
	Name    string  `json:"name"`
	Brand   string  `json:"brand"`
	Season  string  `json:"season"`
	Price   float64 `json:"price"`
}

// ProductStore defines an interface for product-related database operations
type ProductStore interface {
	CreateProduct(product *Product) error
	GetProductByID(id int) (*Product, error)
	UpdateProduct(product *Product) error
	DeleteProduct(id int) error
}

package models

// Stock represents inventory stock information
type Stock struct {
	ID          int    `json:"id"`
	ProductID   int    `json:"product_id"`
	Quantity    int    `json:"quantity"`
	WarehouseID int    `json:"warehouse_id"`
	Location    string `json:"location"`
}

// StockStore defines an interface for stock-related database operations
type StockStore interface {
	CreateStock(stock *Stock) error
	GetStockByProductID(productID int) (*Stock, error)
	UpdateStock(stock *Stock) error
	DeleteStock(id int) error
}

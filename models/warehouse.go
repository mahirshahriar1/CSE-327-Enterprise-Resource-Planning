package models

// Warehouse represents warehouse details
type Warehouse struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Capacity int    `json:"capacity"`
	Location string `json:"location"`
}

// WarehouseStore defines an interface for warehouse-related database operations
type WarehouseStore interface {
	CreateWarehouse(warehouse *Warehouse) error
	GetWarehouseByID(id int) (*Warehouse, error)
	UpdateWarehouse(warehouse *Warehouse) error
	DeleteWarehouse(id int) error
}

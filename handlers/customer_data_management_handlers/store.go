package customer_data_management_handlers

import (
    "database/sql"
    "errors"
    "erp/models" // Adjust the import path if necessary
)

// DBStore is a struct to hold the database connection.
type DBStore struct {
    DB *sql.DB
}

// CreateCustomer inserts a new customer into the database.
func (store *DBStore) CreateCustomer(customer *models.Customer) error {
    query := `INSERT INTO customers (name, contact, order_history) VALUES ($1, $2, $3) RETURNING id`
    err := store.DB.QueryRow(query, customer.Name, customer.Contact, customer.OrderHistory).Scan(&customer.ID)
    if err != nil {
        return err
    }
    return nil
}

// GetCustomerByID retrieves a customer by their ID from the database.
func (store *DBStore) GetCustomerByID(id int) (*models.Customer, error) {
    query := `SELECT id, name, contact, order_history FROM customers WHERE id = $1`
    customer := &models.Customer{}
    err := store.DB.QueryRow(query, id).Scan(&customer.ID, &customer.Name, &customer.Contact, &customer.OrderHistory)
    if err == sql.ErrNoRows {
        return nil, errors.New("customer not found")
    } else if err != nil {
        return nil, err
    }
    return customer, nil
}

// UpdateCustomer updates an existing customer's details in the database.
func (store *DBStore) UpdateCustomer(customer *models.Customer) error {
	query := `UPDATE customers SET name = $1, contact = $2, order_history = $3 WHERE id = $4`
	_, err := store.DB.Exec(query, customer.Name, customer.Contact, customer.OrderHistory, customer.ID)
	if err != nil {
		return err
	}
	return nil
}

// DeleteCustomer deletes a customer from the database by their ID.
func (store *DBStore) DeleteCustomer(id int) error {
	query := `DELETE FROM customers WHERE id = $1`
	_, err := store.DB.Exec(query, id)
	if err != nil {
		return err
	}
	return nil
}


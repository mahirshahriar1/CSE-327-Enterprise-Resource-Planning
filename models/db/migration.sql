-- User Table
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    email VARCHAR(100) UNIQUE NOT NULL,
    password VARCHAR(255),
    role_id INT REFERENCES roles(id) ON DELETE SET NULL,
    department VARCHAR(100),
    needs_new_pass BOOLEAN DEFAULT FALSE
);

-- Role Table
CREATE TABLE roles (
    id SERIAL PRIMARY KEY,
    role_name VARCHAR(50) UNIQUE NOT NULL,
    permissions VARCHAR(255)
);

INSERT INTO roles (role_name, permissions)
VALUES 
    ('Admin', 'all_permissions'),
    ('Employee', 'basic_permissions'),
    ('Sales Group', 'sales_permissions'),
    ('Purchase Group', 'purchase_permissions'),
    ('Accountant', 'finance_permissions'),
    ('Corporate', 'corporate_permissions')
ON CONFLICT (role_name) DO NOTHING;

-- Attendance Table
CREATE TABLE attendance (
    id SERIAL PRIMARY KEY,
    user_id INT REFERENCES users(id) ON DELETE CASCADE,
    check_in TIMESTAMP,
    check_out TIMESTAMP,
    total_hours DECIMAL(5, 2)
);

-- Leave Table
CREATE TABLE leave (
    id SERIAL PRIMARY KEY,
    user_id INT REFERENCES users(id) ON DELETE CASCADE,
    leave_type VARCHAR(50) NOT NULL,
    start_date DATE NOT NULL,
    end_date DATE NOT NULL,
    status VARCHAR(20)
);

-- Product Table
CREATE TABLE products (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    brand VARCHAR(50),
    season VARCHAR(50),
    price DECIMAL(10, 2) NOT NULL
);

-- Stock Table
CREATE TABLE stock (
    id SERIAL PRIMARY KEY,
    product_id INT REFERENCES products(id) ON DELETE CASCADE,
    quantity INT NOT NULL,
    warehouse_id INT REFERENCES warehouses(id) ON DELETE SET NULL,
    location VARCHAR(100)
);

-- Warehouse Table
CREATE TABLE warehouses (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    capacity INT NOT NULL,
    location VARCHAR(100)
);

-- Customer Table
CREATE TABLE customers (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    contact VARCHAR(50),
    order_history TEXT
);

-- Sales Order Table
CREATE TABLE sales_orders (
    id SERIAL PRIMARY KEY,
    customer_id INT REFERENCES customers(id) ON DELETE CASCADE,
    product_id INT REFERENCES products(id) ON DELETE SET NULL,
    order_date DATE NOT NULL,
    quantity INT NOT NULL
);

-- Invoice Table
CREATE TABLE invoices (
    id SERIAL PRIMARY KEY,
    sales_order_id INT REFERENCES sales_orders(id) ON DELETE CASCADE,
    customer_id INT REFERENCES customers(id) ON DELETE SET NULL,
    amount DECIMAL(10, 2) NOT NULL,
    status VARCHAR(20)
);

-- Financial Transaction Table
CREATE TABLE financial_transactions (
    id SERIAL PRIMARY KEY,
    account_type VARCHAR(50) NOT NULL,
    amount DECIMAL(10, 2) NOT NULL,
    transaction_date DATE NOT NULL
);

-- Payment Table
CREATE TABLE payments (
    id SERIAL PRIMARY KEY,
    invoice_id INT REFERENCES invoices(id) ON DELETE SET NULL,
    amount DECIMAL(10, 2) NOT NULL,
    payment_date DATE NOT NULL,
    payment_method VARCHAR(50)
);

-- Financial Transaction Table with Foreign Keys
CREATE TABLE financial_transactions (
    id SERIAL PRIMARY KEY,
    account_type VARCHAR(50) NOT NULL,  -- 'accounts_receivable', 'revenue', 'expense', etc.
    amount DECIMAL(10, 2) NOT NULL,
    transaction_date DATE NOT NULL,
    transaction_type VARCHAR(50),  -- 'credit', 'debit' for tracking inflow and outflow
    invoice_id INT REFERENCES invoices(id) ON DELETE SET NULL,  -- Link to invoice if related
    payment_id INT REFERENCES payments(id) ON DELETE SET NULL,  -- Link to payment if related
    description TEXT   -- Optional, for further clarification (e.g., "Payment for invoice #123")
);

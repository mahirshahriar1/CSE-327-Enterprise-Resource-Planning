package routes

import (
	"database/sql"
	"erp/controllers/handlers/accounts_payable_handlers"
	"erp/controllers/handlers/auth_handlers"
	"erp/controllers/handlers/customer_data_management_handlers" // Import customer handlers package
	"erp/controllers/handlers/general_ledger_handlers"
	"erp/controllers/handlers/invoice_handlers"

	"github.com/gorilla/mux"
)

// InitRoutes initializes all routes in the application, mapping URL paths to handlers.
// It injects dependencies, like database connections, into handlers and stores.
func InitRoutes(db *sql.DB) *mux.Router {
	router := mux.NewRouter()

	// Initialize auth handlers and routes
	roleStore := &auth_handlers.DBRoleStore{DB: db}
	userStore := &auth_handlers.DBUserStore{
		DB:        db,
		RoleStore: roleStore,
	}
	authHandlers := &auth_handlers.AuthHandlers{UserStore: userStore}
	authRouter := router.PathPrefix("/auth").Subrouter()
	authHandlers.RegisterRoutes(authRouter)

	// Customer-related routes
	customerStore := &customer_data_management_handlers.DBStore{DB: db} // Assuming your customer store is in this package
	customerHandlers := &customer_data_management_handlers.CustomerHandlers{Store: customerStore}

	// Create a subrouter for customer routes
	customerRouter := router.PathPrefix("/customers").Subrouter()

	// Register customer routes
	customerRouter.HandleFunc("", customerHandlers.CreateCustomerHandler).Methods("POST")               // Create customer
	customerRouter.HandleFunc("/{id:[0-9]+}", customerHandlers.GetCustomerByIDHandler).Methods("GET")   // Get customer by ID
	customerRouter.HandleFunc("/{id:[0-9]+}", customerHandlers.UpdateCustomerHandler).Methods("PUT")    // Update customer
	customerRouter.HandleFunc("/{id:[0-9]+}", customerHandlers.DeleteCustomerHandler).Methods("DELETE") // Delete customer

	// Protected routes: requires JWT authentication (example)
	// router.Handle("/dashboard", middleware.JWTAuth(http.HandlerFunc(dashboard.Dashboard))).Methods("GET")
	// Initialize general ledger handlers and routes
	generalLedgerStore := &general_ledger_handlers.DBFinancialTransactionStore{DB: db}
	generalLedgerRouter := router.PathPrefix("/general_ledger").Subrouter()
	general_ledger_handlers.RegisterRoutes(generalLedgerRouter, generalLedgerStore)

	// Initialize accounts payable handlers and routes
	accountsPayableStore := &accounts_payable_handlers.DBPaymentStore{DB: db} // PaymentStore implementation
	accountsPayableRouter := router.PathPrefix("/accounts_payable").Subrouter()
	accounts_payable_handlers.RegisterRoutes(accountsPayableRouter, accountsPayableStore, generalLedgerStore)

	// Initialize invoice handlers and routes
	invoiceStore := &invoice_handlers.DBInvoiceStore{DB: db}
	invoiceHandlers := &invoice_handlers.InvoiceHandlers{Store: invoiceStore}

	// Create a subrouter for invoice routes
	invoiceRouter := router.PathPrefix("/invoices").Subrouter()

	// Register invoice routes
	invoiceRouter.HandleFunc("", invoiceHandlers.CreateInvoiceHandler).Methods("POST")             // Create invoice
	invoiceRouter.HandleFunc("/{id:[0-9]+}", invoiceHandlers.GetInvoiceByIDHandler).Methods("GET") // Get invoice by ID

	return router
}

// Protected routes: requires JWT authentication
// router.Handle("/dashboard", middleware.JWTAuth(http.HandlerFunc(dashboard.Dashboard))).Methods("GET")

// Additional routes can be added here
// router.HandleFunc("/profile", authHandlers.Profile).Methods("GET")

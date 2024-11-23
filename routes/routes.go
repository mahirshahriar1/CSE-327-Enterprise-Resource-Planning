package routes

import (
	"database/sql"
	"erp/handlers/accounts_payable_handlers"
	"erp/handlers/auth_handlers"
	"erp/handlers/customer_data_management_handlers" // Import customer handlers package
	"erp/handlers/general_ledger_handlers"

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
	customerRouter.HandleFunc("", customerHandlers.CreateCustomerHandler).Methods("POST")  // Create customer
	customerRouter.HandleFunc("/{id:[0-9]+}", customerHandlers.GetCustomerByIDHandler).Methods("GET")  // Get customer by ID
	customerRouter.HandleFunc("/{id:[0-9]+}", customerHandlers.UpdateCustomerHandler).Methods("PUT") // Update customer
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

	// Initialize accounts receivable handlers and routes
	accountReceivableStore := &accounts_payable_handlers.DBPaymentStore{DB: db} // PaymentStore implementation
	accountReceivableRouter := router.PathPrefix("/accounts_receivable").Subrouter()
	accounts_payable_handlers.RegisterRoutes(accountReceivableRouter, accountReceivableStore, generalLedgerStore)

	// initialize financial transaction handlers and routes
	// todo: implement financial transaction handlers

	return router
}

// Protected routes: requires JWT authentication
// router.Handle("/dashboard", middleware.JWTAuth(http.HandlerFunc(dashboard.Dashboard))).Methods("GET")

// Additional routes can be added here
// router.HandleFunc("/profile", authHandlers.Profile).Methods("GET")

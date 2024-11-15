package routes

import (
	"database/sql"
	"erp/handlers/auth_handlers"
	"erp/handlers/customer_data_management_handlers" // Import customer handlers package
	"github.com/gorilla/mux"
)

// InitRoutes sets up the application routes
func InitRoutes(db *sql.DB) *mux.Router {
	router := mux.NewRouter()

	// Initialize role store
	roleStore := &auth_handlers.DBRoleStore{DB: db}

	// Initialize user store with role store dependency
	userStore := &auth_handlers.DBUserStore{
		DB:        db,
		RoleStore: roleStore,
	}
	// Initialize auth handlers
	authHandlers := &auth_handlers.AuthHandlers{UserStore: userStore}

	// Create a subrouter for auth routes
	authRouter := router.PathPrefix("/auth").Subrouter()
	// Register auth routes
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

	// Additional routes can be added here
	// router.HandleFunc("/profile", authHandlers.Profile).Methods("GET")

	return router
}

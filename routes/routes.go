package routes

import (
	"database/sql"
	"erp/handlers/auth_handlers"
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

	// Initialize general ledger handlers and routes
	generalLedgerStore := &general_ledger_handlers.DBFinancialTransactionStore{DB: db}
	generalLedgerRouter := router.PathPrefix("/general_ledger").Subrouter()
	general_ledger_handlers.RegisterRoutes(generalLedgerRouter, generalLedgerStore)

	return router
}

// Protected routes: requires JWT authentication
// router.Handle("/dashboard", middleware.JWTAuth(http.HandlerFunc(dashboard.Dashboard))).Methods("GET")

// Additional routes can be added here
// router.HandleFunc("/profile", authHandlers.Profile).Methods("GET")

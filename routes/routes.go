package routes

import (
	"database/sql"
	"erp/handlers/auth_handlers"
	"erp/handlers/dashboard"
	"erp/middleware"
	"net/http"

	"github.com/gorilla/mux"
)

// InitRoutes sets up the application routes
func InitRoutes(db *sql.DB) *mux.Router {
	router := mux.NewRouter()

	// Initialize UserStore and AuthHandlers
	userStore := &auth_handlers.DBUserStore{DB: db}
	authHandlers := &auth_handlers.AuthHandlers{UserStore: userStore}

	// Create a subrouter for auth routes
	authRouter := router.PathPrefix("/auth").Subrouter()
	// Register auth routes
	authHandlers.RegisterRoutes(authRouter)

	// Protected routes: requires JWT authentication
	router.Handle("/dashboard", middleware.JWTAuth(http.HandlerFunc(dashboard.Dashboard))).Methods("GET")

	// Additional routes can be added here
	// router.HandleFunc("/profile", authHandlers.Profile).Methods("GET")

	return router
}

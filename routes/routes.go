package routes

import (
	"database/sql"
	dashboard "erp/handlers/Dashboard"
	"erp/handlers/auth_handlers"
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

	// Auth routes
	authRouter.HandleFunc("/signup", authHandlers.SignUp).Methods("POST")
	authRouter.HandleFunc("/check-user", authHandlers.CheckUser).Methods("POST")
	authRouter.HandleFunc("/set-new-password", authHandlers.SetNewPassword).Methods("POST")
	authRouter.HandleFunc("/login", authHandlers.Login).Methods("POST")

	// Protected routes: requires JWT authentication
	router.Handle("/dashboard", middleware.JWTAuth(http.HandlerFunc(dashboard.Dashboard))).Methods("GET")

	// Additional routes can be added here
	// router.HandleFunc("/profile", authHandlers.Profile).Methods("GET")

	return router
}

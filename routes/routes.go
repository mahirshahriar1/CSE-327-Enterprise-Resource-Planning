package routes

import (
	auth_handlers "erp/handlers/Auth"
	dashboard "erp/handlers/Dashboard"
	"erp/middleware"
	"net/http"

	"github.com/gorilla/mux"
)

// InitRoutes sets up the application routes
func InitRoutes() *mux.Router {
	router := mux.NewRouter()

	// Public routes (no authentication required)
	router.HandleFunc("/signup", auth_handlers.SignUp).Methods("POST")
	router.HandleFunc("/check-user", auth_handlers.CheckUser).Methods("POST")
	router.HandleFunc("/set-new-password", auth_handlers.SetNewPassword).Methods("POST")
	router.HandleFunc("/login", auth_handlers.Login).Methods("POST")

	// Protected routes: requires JWT authentication
	router.Handle("/dashboard", middleware.JWTAuth(http.HandlerFunc(dashboard.Dashboard))).Methods("GET")

	// You can add more routes here, for example:
	// router.HandleFunc("/profile", auth_handlers.Profile).Methods("GET")

	return router
}

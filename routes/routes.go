package routes

import (
	"erp/handlers"
	"erp/middleware"
	"net/http"

	"github.com/gorilla/mux"
)

// InitRoutes sets up the application routes
func InitRoutes() *mux.Router {
	router := mux.NewRouter()

	// Public routes (no authentication required)
	router.HandleFunc("/signup", handlers.SignUp).Methods("POST")
	router.HandleFunc("/check-user", handlers.CheckUser).Methods("POST")
	router.HandleFunc("/set-new-password", handlers.SetNewPassword).Methods("POST")
	router.HandleFunc("/login", handlers.Login).Methods("POST")

	// Protected routes: requires JWT authentication
	router.Handle("/dashboard", middleware.JWTAuth(http.HandlerFunc(handlers.Dashboard))).Methods("GET")

	// You can add more routes here, for example:
	// router.HandleFunc("/profile", handlers.Profile).Methods("GET")

	return router
}

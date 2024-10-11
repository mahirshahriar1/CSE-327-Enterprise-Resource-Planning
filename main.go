package main

import (
	"erp/db"
	"erp/handlers"
	"log"
	"net/http"
	"erp/middleware"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

func main() {
	// Initialize the database connection
	var err error
	db.DB, err = db.InitDB()
	if err != nil {
		log.Fatal("Failed to connect to the database:", err)
	}

	// Create the router
	router := mux.NewRouter()

	// Define the routes
	router.HandleFunc("/signup", handlers.SignUp).Methods("POST")
	router.HandleFunc("/check-user", handlers.CheckUser).Methods("POST")
	router.HandleFunc("/set-new-password", handlers.SetNewPassword).Methods("POST")
	router.HandleFunc("/login", handlers.Login).Methods("POST")

	// Protected route: requires JWT authentication
	router.Handle("/dashboard", middleware.JWTAuth(http.HandlerFunc(handlers.Dashboard))).Methods("GET")


	// Start the server
	log.Println("Server started on :8080")
	err = http.ListenAndServe(":8080", router)
	if err != nil {
		log.Fatal("Failed to start server:", err)
	}
}

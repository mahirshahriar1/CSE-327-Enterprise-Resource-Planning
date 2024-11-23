package main

import (
	"erp/models/db"
	"erp/controllers/routes"
	"log"
	"net/http"

	"github.com/gorilla/handlers"
	_ "github.com/lib/pq"
)

func main() {
	// Initialize the database connection
	var err error
	dbInstance, err := db.InitDB() // Use a local variable to avoid global state
	if err != nil {
		log.Fatal("Failed to connect to the database:", err)
	}
	defer dbInstance.Close()

	// Initialize the routes, passing the db instance
	router := routes.InitRoutes(dbInstance)

	// Set up CORS
	corsObj := handlers.AllowedOrigins([]string{"*"}) // You can replace "*" with your frontend URL
	corsHeaders := handlers.AllowedHeaders([]string{"Content-Type", "Authorization"})
	corsMethods := handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE", "OPTIONS"})

	// Start the server with CORS
	log.Println("Server started on :8080")
	err = http.ListenAndServe(":8080", handlers.CORS(corsObj, corsHeaders, corsMethods)(router))
	if err != nil {
		log.Fatal("Failed to start server:", err)
	}
}

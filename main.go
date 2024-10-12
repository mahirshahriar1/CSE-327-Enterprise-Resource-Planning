package main

import (
	"erp/db"
	"erp/routes"
	"log"
	"net/http"

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

	// Initialize the routes, pass the db instance
	router := routes.InitRoutes(dbInstance)

	// Start the server
	log.Println("Server started on :8080")
	err = http.ListenAndServe(":8080", router)
	if err != nil {
		log.Fatal("Failed to start server:", err)
	}
}

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
	db.DB, err = db.InitDB()
	if err != nil {
		log.Fatal("Failed to connect to the database:", err)
	}

	// Initialize the routes
	router := routes.InitRoutes()

	// Start the server
	log.Println("Server started on :8080")
	err = http.ListenAndServe(":8080", router)
	if err != nil {
		log.Fatal("Failed to start server:", err)
	}
}

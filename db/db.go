package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

var DB *sql.DB

func InitDB() (*sql.DB, error) {
	// Load environment variables from .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	// Get environment variables
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	sslMode := os.Getenv("SSL_MODE")

	// Create connection string
	connStr := fmt.Sprintf("user=%s password=%s dbname=%s host=%s port=%s sslmode=%s", dbUser, dbPassword, dbName, dbHost, dbPort, sslMode)

	// Open connection to the database
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	// Ping the database to test the connection
	err = db.Ping()
	if err != nil {
		return nil, err
	}

	log.Println("Connected to the database successfully!")
	return db, nil
}

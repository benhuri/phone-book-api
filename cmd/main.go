package main

import (
	"log"
	"net/http"
	"os"

	"github.com/benhuri/phone-book-api/internal/contacts"
	"github.com/benhuri/phone-book-api/internal/database"
	"github.com/benhuri/phone-book-api/internal/router"
)

func main() {
	// Get the database connection string from environment variables
	dbConnStr := os.Getenv("DB_CONNECTION_STRING")
	if dbConnStr == "" {
		log.Fatal("DB_CONNECTION_STRING environment variable is required")
	}

	// Initialize the database
	if err := database.InitDB(dbConnStr); err != nil {
		log.Fatalf("Could not connect to the database: %v", err)
	}
	defer database.DB.Close()

	// Initialize the contacts repository, service, and handler
	contactsRepo := contacts.NewRepository(database.DB)
	contactsService := contacts.NewService(contactsRepo)
	contactsHandler := contacts.NewHandler(contactsService)

	// Set up the router
	r := router.NewRouter(contactsHandler)

	// Start the HTTP server
	log.Println("Starting server on :8080")
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatalf("Could not start server: %v", err)
	}
}

package main

import (
	"context"
	"log"
	"net/http"

	"github.com/benhuri/phone-book-api/internal/config"
	"github.com/benhuri/phone-book-api/internal/contacts"
	"github.com/benhuri/phone-book-api/internal/database"
	"github.com/benhuri/phone-book-api/internal/metrics"
	"github.com/benhuri/phone-book-api/internal/router"
)

func main() {
	// Initialize the configuration
	config.InitConfig()

	// Initialize the database connection
	database.InitDB()

	// Create the contacts table if it doesn't exist
	createTableQuery := `
    CREATE TABLE IF NOT EXISTS contacts (
        id SERIAL PRIMARY KEY,
        first_name VARCHAR(50),
        last_name VARCHAR(50),
        phone_number VARCHAR(20),
        address VARCHAR(100)
    );`
	_, err := database.DB.ExecContext(context.Background(), createTableQuery)
	if err != nil {
		log.Fatalf("Error creating contacts table: %v", err)
	}

	// Initialize the contacts repository, service, and handler
	contactsRepo := contacts.NewRepository(database.DB)
	contactsService := contacts.NewService(contactsRepo)
	contactHandler := contacts.NewHandler(contactsService)

	// Initialize the router
	r := router.NewRouter(contactHandler)

	// Apply the metrics middleware
	r.Use(metrics.Middleware)

	// Start the server
	log.Fatal(http.ListenAndServe(":8080", r))
}

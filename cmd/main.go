package main

import (
    "log"
    "net/http"
    "phone-book-api/internal/database"
    "phone-book-api/internal/router"
)

func main() {
    // Initialize the database
    db, err := database.InitDB()
    if err != nil {
        log.Fatalf("Could not connect to the database: %v", err)
    }
    defer db.Close()

    // Set up the router
    r := router.NewRouter()

    // Start the HTTP server
    log.Println("Starting server on :8080")
    if err := http.ListenAndServe(":8080", r); err != nil {
        log.Fatalf("Could not start server: %v", err)
    }
}
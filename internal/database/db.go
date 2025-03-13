package database

import (
    "database/sql"
    "log"

    _ "github.com/lib/pq" // PostgreSQL driver
)

var db *sql.DB

// InitDB initializes the database connection
func InitDB(dataSourceName string) {
    var err error
    db, err = sql.Open("postgres", dataSourceName)
    if err != nil {
        log.Fatalf("Error opening database: %v", err)
    }

    if err = db.Ping(); err != nil {
        log.Fatalf("Error connecting to the database: %v", err)
    }

    log.Println("Database connection established")
}

// GetDB returns the database connection
func GetDB() *sql.DB {
    return db
}

// CloseDB closes the database connection
func CloseDB() {
    if err := db.Close(); err != nil {
        log.Fatalf("Error closing database: %v", err)
    }
    log.Println("Database connection closed")
}
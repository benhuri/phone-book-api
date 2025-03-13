package test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"testing"

	"github.com/benhuri/phone-book-api/internal/contacts"
	"github.com/benhuri/phone-book-api/internal/database"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

var contactHandler *contacts.Handler
var router *mux.Router

func setup() {
	// Set the database connection string for testing
	os.Setenv("DB_CONNECTION_STRING", "postgres://postgres:03051991@localhost/phonebook?sslmode=disable")
	database.InitDB(os.Getenv("DB_CONNECTION_STRING")) // Assuming this initializes the database connection

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
		panic(err)
	}

	// Initialize the contacts repository, service, and handler
	contactsRepo := contacts.NewRepository(database.DB)
	contactsService := contacts.NewService(contactsRepo)
	contactHandler = contacts.NewHandler(contactsService)

	// Initialize the router
	router = mux.NewRouter()
	router.HandleFunc("/contacts", contactHandler.AddContactHandler).Methods("POST")
	router.HandleFunc("/contacts", contactHandler.GetContactsHandler).Methods("GET")
	router.HandleFunc("/contacts/search", contactHandler.SearchContactHandler).Methods("GET")
	router.HandleFunc("/contacts/{id}", contactHandler.EditContactHandler).Methods("PUT")
	router.HandleFunc("/contacts/{id}", contactHandler.DeleteContactHandler).Methods("DELETE")
}

func TestAddContact(t *testing.T) {
	setup()
	contact := contacts.Contact{
		FirstName:   "John",
		LastName:    "Doe",
		PhoneNumber: "1234567890",
		Address:     "123 Main St",
	}

	body, _ := json.Marshal(contact)
	req, err := http.NewRequest("POST", "/contacts", bytes.NewBuffer(body))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusCreated, rr.Code)

	// Print the returned JSON
	fmt.Println("AddContact response:", rr.Body.String())
}

func TestGetContacts(t *testing.T) {
	setup()
	req, err := http.NewRequest("GET", "/contacts?page=1&limit=10", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
}

func TestEditContact(t *testing.T) {
	setup()

	// Add a contact to edit
	contact := contacts.Contact{
		FirstName:   "John",
		LastName:    "Doe",
		PhoneNumber: "1234567890",
		Address:     "123 Main St",
	}
	body, _ := json.Marshal(contact)
	req, err := http.NewRequest("POST", "/contacts", bytes.NewBuffer(body))
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusCreated, rr.Code)

	// Capture the ID of the newly created contact
	var createdContact contacts.Contact
	err = json.NewDecoder(rr.Body).Decode(&createdContact)
	if err != nil {
		t.Fatal(err)
	}

	// Print the captured contact
	fmt.Println("Created contact:", createdContact)

	// Edit the contact
	contact = contacts.Contact{
		FirstName:   "Jane",
		LastName:    "Doe",
		PhoneNumber: "0987654321",
		Address:     "456 Elm St",
	}
	body, _ = json.Marshal(contact)
	req, err = http.NewRequest("PUT", "/contacts/"+strconv.Itoa(createdContact.ID), bytes.NewBuffer(body))
	if err != nil {
		t.Fatal(err)
	}

	rr = httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
}

func TestDeleteContact(t *testing.T) {
	setup()

	// Add a contact to delete
	contact := contacts.Contact{
		FirstName:   "John",
		LastName:    "Doe",
		PhoneNumber: "1234567890",
		Address:     "123 Main St",
	}
	body, _ := json.Marshal(contact)
	req, err := http.NewRequest("POST", "/contacts", bytes.NewBuffer(body))
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusCreated, rr.Code)

	// Capture the ID of the newly created contact
	var createdContact contacts.Contact
	err = json.NewDecoder(rr.Body).Decode(&createdContact)
	if err != nil {
		t.Fatal(err)
	}

	// Print the captured contact
	fmt.Println("Created contact:", createdContact)

	// Delete the contact
	req, err = http.NewRequest("DELETE", "/contacts/"+strconv.Itoa(createdContact.ID), nil)
	if err != nil {
		t.Fatal(err)
	}

	rr = httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusNoContent, rr.Code)
}

func TestSearchContact(t *testing.T) {
	setup()
	req, err := http.NewRequest("GET", "/contacts/search?query=John", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
}

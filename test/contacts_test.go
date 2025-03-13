package test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/benhuri/phone-book-api/internal/contacts"
	"github.com/benhuri/phone-book-api/internal/database"
	"github.com/stretchr/testify/assert"
)

func setup() {
	// Set the database connection string for testing
	os.Setenv("DB_CONNECTION_STRING", "postgres://postgres:password@localhost/phonebook?sslmode=disable")
	database.InitDB(os.Getenv("DB_CONNECTION_STRING")) // Assuming this initializes the database connection
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
	handler := http.HandlerFunc((&contacts.Handler{}).AddContactHandler)
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusCreated, rr.Code)
}

func TestGetContacts(t *testing.T) {
	setup()
	req, err := http.NewRequest("GET", "/contacts?page=1&limit=10", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc((&contacts.Handler{}).GetContactsHandler)
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
}

func TestEditContact(t *testing.T) {
	setup()
	contact := contacts.Contact{
		FirstName:   "Jane",
		LastName:    "Doe",
		PhoneNumber: "0987654321",
		Address:     "456 Elm St",
	}

	body, _ := json.Marshal(contact)
	req, err := http.NewRequest("PUT", "/contacts/1", bytes.NewBuffer(body))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc((&contacts.Handler{}).EditContactHandler)
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
}

func TestDeleteContact(t *testing.T) {
	setup()
	req, err := http.NewRequest("DELETE", "/contacts/1", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc((&contacts.Handler{}).DeleteContactHandler)
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusNoContent, rr.Code)
}

func TestSearchContact(t *testing.T) {
	setup()
	req, err := http.NewRequest("GET", "/contacts/search?query=John", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc((&contacts.Handler{}).SearchContactHandler)
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
}

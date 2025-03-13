package test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"phone-book-api/internal/contacts"
	"phone-book-api/internal/database"
)

func setup() {
	database.Connect() // Assuming this initializes the database connection
}

func TestAddContact(t *testing.T) {
	setup()
	contact := contacts.Contact{
		FirstName: "John",
		LastName:  "Doe",
		Phone:     "1234567890",
		Address:   "123 Main St",
	}

	body, _ := json.Marshal(contact)
	req, err := http.NewRequest("POST", "/contacts", bytes.NewBuffer(body))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(contacts.AddContactHandler)
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
	handler := http.HandlerFunc(contacts.GetContactsHandler)
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
}

func TestEditContact(t *testing.T) {
	setup()
	contact := contacts.Contact{
		FirstName: "Jane",
		LastName:  "Doe",
		Phone:     "0987654321",
		Address:   "456 Elm St",
	}

	body, _ := json.Marshal(contact)
	req, err := http.NewRequest("PUT", "/contacts/1", bytes.NewBuffer(body))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(contacts.EditContactHandler)
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
	handler := http.HandlerFunc(contacts.DeleteContactHandler)
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
	handler := http.HandlerFunc(contacts.SearchContactHandler)
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
}
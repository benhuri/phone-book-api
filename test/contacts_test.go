package test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"testing"

	"github.com/benhuri/phone-book-api/internal/contacts"
	"github.com/benhuri/phone-book-api/internal/database"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

const (
	contentType         = "Content-Type"
	applicationJSON     = "application/json"
	basePath            = "/contacts"
	contactsPath        = basePath
	contactsSearchPath  = basePath + "/search"
	contactIDPath       = basePath + "/{id}"
	pageParam           = "page"
	limitParam          = "limit"
	queryParam          = "query"
	invalidRequestError = "Invalid request payload"
	invalidContactID    = "Invalid contact ID"
	internalServerError = "Internal Server Error"
)

var contactHandler *contacts.Handler
var router *mux.Router
var testContact contacts.Contact

func setup() {
	logrus.Info("Setting up the test environment")
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
		logrus.Fatalf("Failed to create contacts table: %v", err)
	}

	// Initialize the contacts repository, service, and handler
	contactsRepo := contacts.NewRepository(database.DB)
	contactsService := contacts.NewService(contactsRepo)
	contactHandler = contacts.NewHandler(contactsService)

	// Initialize the router
	router = mux.NewRouter()
	router.HandleFunc(contactsPath, contactHandler.AddContactHandler).Methods("POST")
	router.HandleFunc(contactsPath, contactHandler.GetContactsHandler).Methods("GET")
	router.HandleFunc(contactsSearchPath, contactHandler.SearchContactHandler).Methods("GET")
	router.HandleFunc(contactIDPath, contactHandler.EditContactHandler).Methods("PUT")
	router.HandleFunc(contactIDPath, contactHandler.DeleteContactHandler).Methods("DELETE")

	// Create a test contact
	testContact = contacts.Contact{
		FirstName:   "John",
		LastName:    "Doe",
		PhoneNumber: "1234567890",
		Address:     "123 Main St",
	}
	body, _ := json.Marshal(testContact)
	req, err := http.NewRequest("POST", contactsPath, bytes.NewBuffer(body))
	if err != nil {
		logrus.Fatalf("Failed to create test contact: %v", err)
	}
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)
	if rr.Code != http.StatusCreated {
		logrus.Fatalf("Failed to create test contact: %v", rr.Body.String())
	}
	err = json.NewDecoder(rr.Body).Decode(&testContact)
	if err != nil {
		logrus.Fatalf("Failed to decode test contact: %v", err)
	}
}

func teardown() {
	logrus.Info("Tearing down the test environment")
	// Delete all test data
	deleteQuery := `DELETE FROM contacts`
	_, err := database.DB.ExecContext(context.Background(), deleteQuery)
	if err != nil {
		logrus.Fatalf("Failed to delete test data: %v", err)
	}

	// Reset the ID sequence
	resetSequenceQuery := `ALTER SEQUENCE contacts_id_seq RESTART WITH 1`
	_, err = database.DB.ExecContext(context.Background(), resetSequenceQuery)
	if err != nil {
		logrus.Fatalf("Failed to reset ID sequence: %v", err)
	}
}

func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	teardown()
	os.Exit(code)
}

func TestAddContact(t *testing.T) {
	logrus.Info("Running TestAddContact")
	contact := contacts.Contact{
		FirstName:   "Jane",
		LastName:    "Smith",
		PhoneNumber: "0987654321",
		Address:     "456 Elm St",
	}

	body, _ := json.Marshal(contact)
	req, err := http.NewRequest("POST", contactsPath, bytes.NewBuffer(body))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusCreated, rr.Code)

	// Verify the contact was created correctly
	var createdContact contacts.Contact
	err = json.NewDecoder(rr.Body).Decode(&createdContact)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, contact.FirstName, createdContact.FirstName)
	assert.Equal(t, contact.LastName, createdContact.LastName)
	assert.Equal(t, contact.PhoneNumber, createdContact.PhoneNumber)
	assert.Equal(t, contact.Address, createdContact.Address)

	// Print the returned JSON
	logrus.Infof("AddContact response: %v", rr.Body.String())
}

func TestGetContacts(t *testing.T) {
	logrus.Info("Running TestGetContacts")
	req, err := http.NewRequest("GET", contactsPath+"?page=1&limit=10", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	// Verify the response contains the expected contacts
	var contacts []contacts.Contact
	err = json.NewDecoder(rr.Body).Decode(&contacts)
	if err != nil {
		t.Fatal(err)
	}

	assert.NotEmpty(t, contacts)
}

func TestSearchContact(t *testing.T) {
	logrus.Info("Running TestSearchContact")
	// Search for the test contact
	req, err := http.NewRequest("GET", contactsSearchPath+"?query=John", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	// Verify the search results
	var searchResults []contacts.Contact
	err = json.NewDecoder(rr.Body).Decode(&searchResults)
	if err != nil {
		t.Fatal(err)
	}

	assert.NotEmpty(t, searchResults)
	assert.Equal(t, testContact.FirstName, searchResults[0].FirstName)
	assert.Equal(t, testContact.LastName, searchResults[0].LastName)
	assert.Equal(t, testContact.PhoneNumber, searchResults[0].PhoneNumber)
	assert.Equal(t, testContact.Address, searchResults[0].Address)
}

func TestEditContact(t *testing.T) {
	logrus.Info("Running TestEditContact")
	// Edit the test contact
	contact := contacts.Contact{
		FirstName:   "Jane",
		LastName:    "Doe",
		PhoneNumber: "0987654321",
		Address:     "456 Elm St",
	}
	body, _ := json.Marshal(contact)
	req, err := http.NewRequest("PUT", contactsPath+"/"+strconv.Itoa(testContact.ID), bytes.NewBuffer(body))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	// Verify the contact was edited correctly
	var editedContact contacts.Contact
	err = json.NewDecoder(rr.Body).Decode(&editedContact)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, contact.FirstName, editedContact.FirstName)
	assert.Equal(t, contact.LastName, editedContact.LastName)
	assert.Equal(t, contact.PhoneNumber, editedContact.PhoneNumber)
	assert.Equal(t, contact.Address, editedContact.Address)
}

func TestDeleteContact(t *testing.T) {
	logrus.Info("Running TestDeleteContact")
	// Add a contact to delete
	contact := contacts.Contact{
		FirstName:   "Mark",
		LastName:    "Twain",
		PhoneNumber: "1122334455",
		Address:     "789 Oak St",
	}
	body, _ := json.Marshal(contact)
	req, err := http.NewRequest("POST", contactsPath, bytes.NewBuffer(body))
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

	// Delete the contact
	req, err = http.NewRequest("DELETE", contactsPath+"/"+strconv.Itoa(createdContact.ID), nil)
	if err != nil {
		t.Fatal(err)
	}

	rr = httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusNoContent, rr.Code)

	// Verify the contact was deleted by searching for it
	req, err = http.NewRequest("GET", contactsSearchPath+"?query="+createdContact.FirstName, nil)
	if err != nil {
		t.Fatal(err)
	}

	rr = httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	// Verify the search results do not contain the deleted contact
	var searchResults []contacts.Contact
	err = json.NewDecoder(rr.Body).Decode(&searchResults)
	if err != nil {
		t.Fatal(err)
	}

	for _, c := range searchResults {
		assert.NotEqual(t, createdContact.ID, c.ID)
	}
}

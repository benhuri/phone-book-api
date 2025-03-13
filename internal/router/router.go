package router

import (
	"net/http"

	"github.com/benhuri/phone-book-api/internal/contacts"
	"github.com/gorilla/mux"
)

func NewRouter(contactHandler *contacts.Handler) *mux.Router {
	router := mux.NewRouter()

	router.HandleFunc("/contacts", contactHandler.GetContacts).Methods(http.MethodGet)
	router.HandleFunc("/contacts/search", contactHandler.SearchContact).Methods(http.MethodGet)
	router.HandleFunc("/contacts", contactHandler.AddContact).Methods(http.MethodPost)
	router.HandleFunc("/contacts/{id}", contactHandler.EditContact).Methods(http.MethodPut)
	router.HandleFunc("/contacts/{id}", contactHandler.DeleteContact).Methods(http.MethodDelete)

	return router
}

package router

import (
	"net/http"

	"github.com/benhuri/phone-book-api/internal/contacts"
	"github.com/gorilla/mux"
)

func NewRouter(contactHandler *contacts.Handler) *mux.Router {
	router := mux.NewRouter()

	router.HandleFunc("/contacts", contactHandler.GetContactsHandler).Methods(http.MethodGet)
	router.HandleFunc("/contacts/search", contactHandler.SearchContactHandler).Methods(http.MethodGet)
	router.HandleFunc("/contacts", contactHandler.AddContactHandler).Methods(http.MethodPost)
	router.HandleFunc("/contacts/{id}", contactHandler.EditContactHandler).Methods(http.MethodPut)
	router.HandleFunc("/contacts/{id}", contactHandler.DeleteContactHandler).Methods(http.MethodDelete)

	return router
}

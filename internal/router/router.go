package router

import (
	"github.com/benhuri/phone-book-api/internal/contacts"
	"github.com/gorilla/mux"
)

const (
	basePath           = "/contacts"
	contactsPath       = basePath
	contactsSearchPath = basePath + "/search"
	contactIDPath      = basePath + "/{id}"
)

func NewRouter(handler *contacts.Handler) *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc(contactsPath, handler.AddContactHandler).Methods("POST")
	r.HandleFunc(contactsPath, handler.GetContactsHandler).Methods("GET")
	r.HandleFunc(contactsSearchPath, handler.SearchContactHandler).Methods("GET")
	r.HandleFunc(contactIDPath, handler.EditContactHandler).Methods("PUT")
	r.HandleFunc(contactIDPath, handler.DeleteContactHandler).Methods("DELETE")
	return r
}

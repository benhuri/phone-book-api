package contacts

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
)

const (
	contentType         = "Content-Type"
	applicationJSON     = "application/json"
	idParam             = "id"
	pageParam           = "page"
	limitParam          = "limit"
	queryParam          = "query"
	invalidRequestError = "Invalid request payload"
	invalidContactID    = "Invalid contact ID"
	internalServerError = "Internal Server Error"
)

type Handler struct {
	Service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{Service: service}
}

var validate *validator.Validate

func init() {
	validate = validator.New()
}

func (h *Handler) GetContactsHandler(w http.ResponseWriter, r *http.Request) {
	pageStr := r.URL.Query().Get(pageParam)
	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}

	limitStr := r.URL.Query().Get(limitParam)
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 1 {
		limit = 10 // default limit
	}

	contacts, err := h.Service.GetContacts(page, limit)
	if err != nil {
		log.Printf("Error getting contacts: %v", err)
		http.Error(w, internalServerError, http.StatusInternalServerError)
		return
	}

	w.Header().Set(contentType, applicationJSON)
	json.NewEncoder(w).Encode(contacts)
}

func (h *Handler) SearchContactHandler(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get(queryParam)
	contacts, err := h.Service.SearchContact(query)
	if err != nil {
		log.Printf("Error searching contacts: %v", err)
		http.Error(w, internalServerError, http.StatusInternalServerError)
		return
	}

	w.Header().Set(contentType, applicationJSON)
	json.NewEncoder(w).Encode(contacts)
}

func (h *Handler) AddContactHandler(w http.ResponseWriter, r *http.Request) {
	var contact Contact
	if err := json.NewDecoder(r.Body).Decode(&contact); err != nil {
		log.Printf("Error decoding contact: %v", err)
		http.Error(w, invalidRequestError, http.StatusBadRequest)
		return
	}

	if err := validate.Struct(contact); err != nil {
		log.Printf("Validation error: %v", err)
		http.Error(w, formatValidationError(err), http.StatusBadRequest)
		return
	}

	if err := h.Service.AddContact(&contact); err != nil {
		log.Printf("Error adding contact: %v", err)
		http.Error(w, internalServerError, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Header().Set(contentType, applicationJSON)
	json.NewEncoder(w).Encode(contact)
}

func (h *Handler) EditContactHandler(w http.ResponseWriter, r *http.Request) {
	var contact Contact
	if err := json.NewDecoder(r.Body).Decode(&contact); err != nil {
		log.Printf("Error decoding contact: %v", err)
		http.Error(w, invalidRequestError, http.StatusBadRequest)
		return
	}

	if err := validate.Struct(contact); err != nil {
		log.Printf("Validation error: %v", err)
		http.Error(w, formatValidationError(err), http.StatusBadRequest)
		return
	}

	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars[idParam])
	if err != nil {
		log.Printf("Invalid contact ID: %v", err)
		http.Error(w, invalidContactID, http.StatusBadRequest)
		return
	}
	contact.ID = id

	if err := h.Service.EditContact(contact); err != nil {
		log.Printf("Error editing contact: %v", err)
		http.Error(w, internalServerError, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set(contentType, applicationJSON)
	json.NewEncoder(w).Encode(contact)
}

func (h *Handler) DeleteContactHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars[idParam])
	if err != nil {
		log.Printf("Invalid contact ID: %v", err)
		http.Error(w, invalidContactID, http.StatusBadRequest)
		return
	}

	if err := h.Service.DeleteContact(id); err != nil {
		log.Printf("Error deleting contact: %v", err)
		http.Error(w, internalServerError, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func formatValidationError(err error) string {
	var errors []string
	for _, err := range err.(validator.ValidationErrors) {
		switch err.Tag() {
		case "required":
			errors = append(errors, err.Field()+" is required")
		case "min":
			errors = append(errors, err.Field()+" must be at least "+err.Param()+" characters")
		case "max":
			errors = append(errors, err.Field()+" must be at most "+err.Param()+" characters")
		default:
			errors = append(errors, err.Field()+" is invalid")
		}
	}
	return "Validation error: " + strings.Join(errors, ", ")
}

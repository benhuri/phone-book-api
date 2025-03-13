package contacts

import (
	"errors"
	"fmt"
)

// Contact represents a contact in the phone book.
type Contact struct {
	FirstName string
	LastName  string
	Phone     string
	Address   string
}

// Service provides methods for managing contacts.
type Service struct {
	repo Repository
}

// NewService creates a new contact service.
func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

// GetContacts retrieves a paginated list of contacts.
func (s *Service) GetContacts(page, limit int) ([]Contact, error) {
	if page < 1 || limit < 1 {
		return nil, errors.New("page and limit must be greater than 0")
	}
	return s.repo.FetchContacts(page, limit)
}

// SearchContact searches for a contact by name or phone number.
func (s *Service) SearchContact(query string) ([]Contact, error) {
	if query == "" {
		return nil, errors.New("query cannot be empty")
	}
	return s.repo.FindContact(query)
}

// AddContact adds a new contact to the phone book.
func (s *Service) AddContact(contact Contact) error {
	if contact.FirstName == "" || contact.LastName == "" || contact.Phone == "" {
		return errors.New("first name, last name, and phone are required")
	}
	return s.repo.CreateContact(contact)
}

// EditContact updates an existing contact.
func (s *Service) EditContact(id string, contact Contact) error {
	if id == "" {
		return errors.New("contact ID cannot be empty")
	}
	return s.repo.UpdateContact(id, contact)
}

// DeleteContact removes a contact from the phone book.
func (s *Service) DeleteContact(id string) error {
	if id == "" {
		return errors.New("contact ID cannot be empty")
	}
	return s.repo.RemoveContact(id)
}
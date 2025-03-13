package contacts

import (
	"context"
)

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) GetContacts(page, limit int) ([]Contact, error) {
	offset := (page - 1) * limit
	return s.repo.FetchContacts(context.Background(), limit, offset)
}

func (s *Service) SearchContact(query string) ([]Contact, error) {
	return s.repo.FindContact(context.Background(), query)
}

func (s *Service) AddContact(contact *Contact) error {
	return s.repo.CreateContact(context.Background(), contact)
}

func (s *Service) EditContact(contact Contact) error {
	return s.repo.UpdateContact(context.Background(), contact)
}

func (s *Service) DeleteContact(id int) error {
	return s.repo.RemoveContact(context.Background(), id)
}

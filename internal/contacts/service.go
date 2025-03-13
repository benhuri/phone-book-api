package contacts

import (
	"context"
)

type Service struct {
	Repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{Repo: repo}
}

func (s *Service) GetContacts(page, limit int) ([]Contact, error) {
	offset := (page - 1) * limit
	return s.Repo.FetchContacts(context.Background(), limit, offset)
}

func (s *Service) SearchContact(query string) ([]Contact, error) {
	return s.Repo.FindContact(context.Background(), query)
}

func (s *Service) AddContact(contact Contact) error {
	return s.Repo.CreateContact(context.Background(), contact)
}

func (s *Service) EditContact(contact Contact) error {
	return s.Repo.UpdateContact(context.Background(), contact)
}

func (s *Service) DeleteContact(id int) error {
	return s.Repo.RemoveContact(context.Background(), id)
}

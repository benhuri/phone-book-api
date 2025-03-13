package contacts

import (
	"context"
	"database/sql"
	"errors"
)

type Repository interface {
	FetchContacts(ctx context.Context, limit, offset int) ([]Contact, error)
	FindContact(ctx context.Context, query string) ([]Contact, error)
	CreateContact(ctx context.Context, contact *Contact) error
	UpdateContact(ctx context.Context, contact Contact) error
	RemoveContact(ctx context.Context, id int) error
}

type contactRepository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) Repository {
	return &contactRepository{db: db}
}

func (r *contactRepository) FetchContacts(ctx context.Context, limit, offset int) ([]Contact, error) {
	rows, err := r.db.QueryContext(ctx, "SELECT id, first_name, last_name, phone_number, address FROM contacts LIMIT $1 OFFSET $2", limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var contacts []Contact
	for rows.Next() {
		var contact Contact
		if err := rows.Scan(&contact.ID, &contact.FirstName, &contact.LastName, &contact.PhoneNumber, &contact.Address); err != nil {
			return nil, err
		}
		contacts = append(contacts, contact)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return contacts, nil
}

func (r *contactRepository) FindContact(ctx context.Context, query string) ([]Contact, error) {
	rows, err := r.db.QueryContext(ctx, "SELECT id, first_name, last_name, phone_number, address FROM contacts WHERE first_name LIKE $1 OR last_name LIKE $2", "%"+query+"%", "%"+query+"%")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var contacts []Contact
	for rows.Next() {
		var contact Contact
		if err := rows.Scan(&contact.ID, &contact.FirstName, &contact.LastName, &contact.PhoneNumber, &contact.Address); err != nil {
			return nil, err
		}
		contacts = append(contacts, contact)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return contacts, nil
}

func (r *contactRepository) CreateContact(ctx context.Context, contact *Contact) error {
	err := r.db.QueryRowContext(ctx, "INSERT INTO contacts (first_name, last_name, phone_number, address) VALUES ($1, $2, $3, $4) RETURNING id", contact.FirstName, contact.LastName, contact.PhoneNumber, contact.Address).Scan(&contact.ID)
	return err
}

func (r *contactRepository) UpdateContact(ctx context.Context, contact Contact) error {
	result, err := r.db.ExecContext(ctx, "UPDATE contacts SET first_name = $1, last_name = $2, phone_number = $3, address = $4 WHERE id = $5", contact.FirstName, contact.LastName, contact.PhoneNumber, contact.Address, contact.ID)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return errors.New("contact not found")
	}
	return nil
}

func (r *contactRepository) RemoveContact(ctx context.Context, id int) error {
	result, err := r.db.ExecContext(ctx, "DELETE FROM contacts WHERE id = $1", id)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return errors.New("contact not found")
	}
	return nil
}

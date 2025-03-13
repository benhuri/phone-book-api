package contacts

import (
	"context"
	"database/sql"
	"errors"
)

type Repository interface {
	GetContacts(ctx context.Context, limit, offset int) ([]Contact, error)
	SearchContact(ctx context.Context, query string) ([]Contact, error)
	AddContact(ctx context.Context, contact Contact) error
	UpdateContact(ctx context.Context, contact Contact) error
	DeleteContact(ctx context.Context, id int) error
}

type repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) Repository {
	return &repository{db: db}
}

func (r *repository) GetContacts(ctx context.Context, limit, offset int) ([]Contact, error) {
	rows, err := r.db.QueryContext(ctx, "SELECT id, first_name, last_name, phone_number, address FROM contacts LIMIT ? OFFSET ?", limit, offset)
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

func (r *repository) SearchContact(ctx context.Context, query string) ([]Contact, error) {
	rows, err := r.db.QueryContext(ctx, "SELECT id, first_name, last_name, phone_number, address FROM contacts WHERE first_name LIKE ? OR last_name LIKE ?", "%"+query+"%", "%"+query+"%")
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

func (r *repository) AddContact(ctx context.Context, contact Contact) error {
	_, err := r.db.ExecContext(ctx, "INSERT INTO contacts (first_name, last_name, phone_number, address) VALUES (?, ?, ?, ?)", contact.FirstName, contact.LastName, contact.PhoneNumber, contact.Address)
	return err
}

func (r *repository) UpdateContact(ctx context.Context, contact Contact) error {
	result, err := r.db.ExecContext(ctx, "UPDATE contacts SET first_name = ?, last_name = ?, phone_number = ?, address = ? WHERE id = ?", contact.FirstName, contact.LastName, contact.PhoneNumber, contact.Address, contact.ID)
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

func (r *repository) DeleteContact(ctx context.Context, id int) error {
	result, err := r.db.ExecContext(ctx, "DELETE FROM contacts WHERE id = ?", id)
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
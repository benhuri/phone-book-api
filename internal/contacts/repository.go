package contacts

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
)

const (
	selectContactsQuery  = "SELECT id, first_name, last_name, phone_number, address FROM contacts"
	selectContactByQuery = "SELECT id, first_name, last_name, phone_number, address FROM contacts WHERE first_name LIKE $1 OR last_name LIKE $2 OR phone_number LIKE $3"
	insertContactQuery   = "INSERT INTO contacts (first_name, last_name, phone_number, address) VALUES ($1, $2, $3, $4) RETURNING id"
	updateContactQuery   = "UPDATE contacts SET first_name = $1, last_name = $2, phone_number = $3, address = $4 WHERE id = $5"
	deleteContactQuery   = "DELETE FROM contacts WHERE id = $1"
	fetchContactsError   = "failed to fetch contacts: %w"
	scanContactError     = "failed to scan contact: %w"
	rowsError            = "rows error: %w"
	findContactError     = "failed to find contact: %w"
	createContactError   = "failed to create contact: %w"
	updateContactError   = "failed to update contact: %w"
	getRowsAffectedError = "failed to get rows affected: %w"
	contactNotFoundError = "contact not found"
	removeContactError   = "failed to remove contact: %w"
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
	rows, err := r.db.QueryContext(ctx, selectContactsQuery+" LIMIT $1 OFFSET $2", limit, offset)
	if err != nil {
		return nil, fmt.Errorf(fetchContactsError, err)
	}
	defer rows.Close()

	var contacts []Contact
	for rows.Next() {
		var contact Contact
		if err := rows.Scan(&contact.ID, &contact.FirstName, &contact.LastName, &contact.PhoneNumber, &contact.Address); err != nil {
			return nil, fmt.Errorf(scanContactError, err)
		}
		contacts = append(contacts, contact)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf(rowsError, err)
	}

	return contacts, nil
}

func (r *contactRepository) FindContact(ctx context.Context, query string) ([]Contact, error) {
	rows, err := r.db.QueryContext(ctx, selectContactByQuery, "%"+query+"%", "%"+query+"%", "%"+query+"%")
	if err != nil {
		return nil, fmt.Errorf(findContactError, err)
	}
	defer rows.Close()

	var contacts []Contact
	for rows.Next() {
		var contact Contact
		if err := rows.Scan(&contact.ID, &contact.FirstName, &contact.LastName, &contact.PhoneNumber, &contact.Address); err != nil {
			return nil, fmt.Errorf(scanContactError, err)
		}
		contacts = append(contacts, contact)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf(rowsError, err)
	}

	return contacts, nil
}

func (r *contactRepository) CreateContact(ctx context.Context, contact *Contact) error {
	err := r.db.QueryRowContext(ctx, insertContactQuery, contact.FirstName, contact.LastName, contact.PhoneNumber, contact.Address).Scan(&contact.ID)
	if err != nil {
		return fmt.Errorf(createContactError, err)
	}
	return nil
}

func (r *contactRepository) UpdateContact(ctx context.Context, contact Contact) error {
	result, err := r.db.ExecContext(ctx, updateContactQuery, contact.FirstName, contact.LastName, contact.PhoneNumber, contact.Address, contact.ID)
	if err != nil {
		return fmt.Errorf(updateContactError, err)
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf(getRowsAffectedError, err)
	}
	if rowsAffected == 0 {
		return errors.New(contactNotFoundError)
	}
	return nil
}

func (r *contactRepository) RemoveContact(ctx context.Context, id int) error {
	result, err := r.db.ExecContext(ctx, deleteContactQuery, id)
	if err != nil {
		return fmt.Errorf(removeContactError, err)
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf(getRowsAffectedError, err)
	}
	if rowsAffected == 0 {
		return errors.New(contactNotFoundError)
	}
	return nil
}

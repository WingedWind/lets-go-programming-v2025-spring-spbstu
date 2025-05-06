package repository

import (
	"context"
	"errors"

	"task-10/internal/model"
)

var (
	// ErrContactNotFound is returned when a contact with the given ID is not found.
	ErrContactNotFound = errors.New("contact not found")
)

// ContactRepository defines the interface for contact data operations.
type ContactRepository interface {
	// Create creates a new contact.
	Create(ctx context.Context, contact *model.Contact) error

	// Get retrieves a contact by ID.
	Get(ctx context.Context, id string) (*model.Contact, error)

	// List retrieves all contacts.
	List(ctx context.Context) ([]*model.Contact, error)

	// Update updates an existing contact.
	Update(ctx context.Context, contact *model.Contact) error

	// Delete deletes a contact by ID.
	Delete(ctx context.Context, id string) error
}

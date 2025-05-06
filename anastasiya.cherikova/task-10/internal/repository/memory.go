package repository

import (
	"context"
	"sync"

	"task-10/internal/model"

	"github.com/google/uuid"
)

// MemoryContactRepository implements ContactRepository using in-memory storage.
type MemoryContactRepository struct {
	contacts map[string]*model.Contact
	mu       sync.RWMutex
}

// NewMemoryContactRepository creates a new MemoryContactRepository.
func NewMemoryContactRepository() *MemoryContactRepository {
	return &MemoryContactRepository{
		contacts: make(map[string]*model.Contact),
	}
}

// Create creates a new contact.
func (r *MemoryContactRepository) Create(ctx context.Context, contact *model.Contact) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if contact.ID == "" {
		contact.ID = uuid.New().String()
	}

	r.contacts[contact.ID] = contact
	return nil
}

// Get retrieves a contact by ID.
func (r *MemoryContactRepository) Get(ctx context.Context, id string) (*model.Contact, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	contact, ok := r.contacts[id]
	if !ok {
		return nil, ErrContactNotFound
	}

	return contact, nil
}

// List retrieves all contacts.
func (r *MemoryContactRepository) List(ctx context.Context) ([]*model.Contact, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	contacts := make([]*model.Contact, 0, len(r.contacts))
	for _, contact := range r.contacts {
		contacts = append(contacts, contact)
	}

	return contacts, nil
}

// Update updates an existing contact.
func (r *MemoryContactRepository) Update(ctx context.Context, contact *model.Contact) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.contacts[contact.ID]; !ok {
		return ErrContactNotFound
	}

	r.contacts[contact.ID] = contact
	return nil
}

// Delete deletes a contact by ID.
func (r *MemoryContactRepository) Delete(ctx context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.contacts[id]; !ok {
		return ErrContactNotFound
	}

	delete(r.contacts, id)
	return nil
}

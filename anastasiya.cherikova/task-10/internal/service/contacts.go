package service

import (
	"context"
	"errors"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"task-10/internal/model"
	"task-10/internal/repository"
	contactsv1 "task-10/pkg/api/contacts/v1"
)

// ContactsService implements the gRPC ContactsService.
type ContactsService struct {
	contactsv1.UnimplementedContactsServiceServer

	repo repository.ContactRepository
}

// NewContactsService creates a new ContactsService.
func NewContactsService(repo repository.ContactRepository) *ContactsService {
	return &ContactsService{
		repo: repo,
	}
}

// CreateContact implements contactsv1.ContactsServiceServer.CreateContact.
func (s *ContactsService) CreateContact(ctx context.Context, req *contactsv1.CreateContactRequest) (*contactsv1.CreateContactResponse, error) {
	contact := &model.Contact{
		FirstName:   req.Contact.FirstName,
		LastName:    req.Contact.LastName,
		PhoneNumber: req.Contact.PhoneNumber,
		Email:       req.Contact.Email,
	}

	if err := s.repo.Create(ctx, contact); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create contact: %v", err)
	}

	return &contactsv1.CreateContactResponse{
		Contact: &contactsv1.Contact{
			Id:          contact.ID,
			FirstName:   contact.FirstName,
			LastName:    contact.LastName,
			PhoneNumber: contact.PhoneNumber,
			Email:       contact.Email,
		},
	}, nil
}

// GetContact implements contactsv1.ContactsServiceServer.GetContact.
func (s *ContactsService) GetContact(ctx context.Context, req *contactsv1.GetContactRequest) (*contactsv1.GetContactResponse, error) {
	contact, err := s.repo.Get(ctx, req.Id)
	if err != nil {
		if errors.Is(err, repository.ErrContactNotFound) {
			return nil, status.Errorf(codes.NotFound, "contact not found: %v", err)
		}
		return nil, status.Errorf(codes.Internal, "failed to get contact: %v", err)
	}

	return &contactsv1.GetContactResponse{
		Contact: &contactsv1.Contact{
			Id:          contact.ID,
			FirstName:   contact.FirstName,
			LastName:    contact.LastName,
			PhoneNumber: contact.PhoneNumber,
			Email:       contact.Email,
		},
	}, nil
}

// ListContacts implements contactsv1.ContactsServiceServer.ListContacts.
func (s *ContactsService) ListContacts(ctx context.Context, req *contactsv1.ListContactsRequest) (*contactsv1.ListContactsResponse, error) {
	contacts, err := s.repo.List(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to list contacts: %v", err)
	}

	protoContacts := make([]*contactsv1.Contact, 0, len(contacts))
	for _, contact := range contacts {
		protoContacts = append(protoContacts, &contactsv1.Contact{
			Id:          contact.ID,
			FirstName:   contact.FirstName,
			LastName:    contact.LastName,
			PhoneNumber: contact.PhoneNumber,
			Email:       contact.Email,
		})
	}

	return &contactsv1.ListContactsResponse{
		Contacts: protoContacts,
	}, nil
}

// UpdateContact implements contactsv1.ContactsServiceServer.UpdateContact.
func (s *ContactsService) UpdateContact(ctx context.Context, req *contactsv1.UpdateContactRequest) (*contactsv1.UpdateContactResponse, error) {
	contact := &model.Contact{
		ID:          req.Contact.Id,
		FirstName:   req.Contact.FirstName,
		LastName:    req.Contact.LastName,
		PhoneNumber: req.Contact.PhoneNumber,
		Email:       req.Contact.Email,
	}

	if err := s.repo.Update(ctx, contact); err != nil {
		if errors.Is(err, repository.ErrContactNotFound) {
			return nil, status.Errorf(codes.NotFound, "contact not found: %v", err)
		}
		return nil, status.Errorf(codes.Internal, "failed to update contact: %v", err)
	}

	return &contactsv1.UpdateContactResponse{
		Contact: &contactsv1.Contact{
			Id:          contact.ID,
			FirstName:   contact.FirstName,
			LastName:    contact.LastName,
			PhoneNumber: contact.PhoneNumber,
			Email:       contact.Email,
		},
	}, nil
}

// DeleteContact implements contactsv1.ContactsServiceServer.DeleteContact.
func (s *ContactsService) DeleteContact(ctx context.Context, req *contactsv1.DeleteContactRequest) (*contactsv1.DeleteContactResponse, error) {
	if err := s.repo.Delete(ctx, req.Id); err != nil {
		if errors.Is(err, repository.ErrContactNotFound) {
			return nil, status.Errorf(codes.NotFound, "contact not found: %v", err)
		}
		return nil, status.Errorf(codes.Internal, "failed to delete contact: %v", err)
	}

	return &contactsv1.DeleteContactResponse{
		Success: true,
	}, nil
}

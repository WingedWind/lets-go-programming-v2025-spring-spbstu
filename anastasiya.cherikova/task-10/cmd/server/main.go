package main

import (
	"context"
	"errors"
	"flag"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"

	"task-10/internal/model"
	"task-10/internal/repository"
	contactsv1 "task-10/pkg/api/contacts/v1"
)

var (
	grpcPort = flag.String("grpc-port", "50051", "gRPC server port")
)

// ContactsServiceServer реализует gRPC сервис для контактов
type ContactsServiceServer struct {
	contactsv1.UnimplementedContactsServiceServer
	repo repository.ContactRepository
}

// NewContactsServiceServer создает новый экземпляр сервера
func NewContactsServiceServer(repo repository.ContactRepository) *ContactsServiceServer {
	return &ContactsServiceServer{
		repo: repo,
	}
}

// CreateContact создает новый контакт
func (s *ContactsServiceServer) CreateContact(ctx context.Context, req *contactsv1.CreateContactRequest) (*contactsv1.CreateContactResponse, error) {
	contact := &model.Contact{
		FirstName:   req.GetContact().GetFirstName(),
		LastName:    req.GetContact().GetLastName(),
		PhoneNumber: req.GetContact().GetPhoneNumber(),
		Email:       req.GetContact().GetEmail(),
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

// GetContact получает контакт по ID
func (s *ContactsServiceServer) GetContact(ctx context.Context, req *contactsv1.GetContactRequest) (*contactsv1.GetContactResponse, error) {
	contact, err := s.repo.Get(ctx, req.GetId())
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

// ListContacts получает список всех контактов
func (s *ContactsServiceServer) ListContacts(ctx context.Context, req *contactsv1.ListContactsRequest) (*contactsv1.ListContactsResponse, error) {
	contacts, err := s.repo.List(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to list contacts: %v", err)
	}

	result := make([]*contactsv1.Contact, 0, len(contacts))
	for _, c := range contacts {
		result = append(result, &contactsv1.Contact{
			Id:          c.ID,
			FirstName:   c.FirstName,
			LastName:    c.LastName,
			PhoneNumber: c.PhoneNumber,
			Email:       c.Email,
		})
	}

	return &contactsv1.ListContactsResponse{
		Contacts: result,
	}, nil
}

// UpdateContact обновляет существующий контакт
func (s *ContactsServiceServer) UpdateContact(ctx context.Context, req *contactsv1.UpdateContactRequest) (*contactsv1.UpdateContactResponse, error) {
	contact := &model.Contact{
		ID:          req.GetContact().GetId(),
		FirstName:   req.GetContact().GetFirstName(),
		LastName:    req.GetContact().GetLastName(),
		PhoneNumber: req.GetContact().GetPhoneNumber(),
		Email:       req.GetContact().GetEmail(),
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

// DeleteContact удаляет контакт по ID
func (s *ContactsServiceServer) DeleteContact(ctx context.Context, req *contactsv1.DeleteContactRequest) (*contactsv1.DeleteContactResponse, error) {
	if err := s.repo.Delete(ctx, req.GetId()); err != nil {
		if errors.Is(err, repository.ErrContactNotFound) {
			return nil, status.Errorf(codes.NotFound, "contact not found: %v", err)
		}
		return nil, status.Errorf(codes.Internal, "failed to delete contact: %v", err)
	}

	return &contactsv1.DeleteContactResponse{
		Success: true,
	}, nil
}

func main() {
	flag.Parse()

	// Создаем репозиторий для контактов
	repo := repository.NewMemoryContactRepository()

	// Инициализируем gRPC сервер
	grpcServer := grpc.NewServer()
	contactsServer := NewContactsServiceServer(repo)

	// Регистрируем наш сервис
	contactsv1.RegisterContactsServiceServer(grpcServer, contactsServer)

	// Включаем рефлексию для grpcurl
	reflection.Register(grpcServer)

	// Запускаем gRPC сервер
	lis, err := net.Listen("tcp", ":"+*grpcPort)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	log.Printf("Starting gRPC server on port %s", *grpcPort)

	// Запускаем сервер в отдельной горутине
	go func() {
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("Failed to serve: %v", err)
		}
	}()

	// Ожидаем сигнала для грациозного завершения
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down gRPC server...")
	grpcServer.GracefulStop()
	log.Println("Server stopped")
}

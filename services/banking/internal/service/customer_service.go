package service

import (
	"context"

	"github.com/jonbelaire/repotown/packages/go-core/logging"
	"github.com/jonbelaire/repotown/services/banking/internal/domain"
	"github.com/jonbelaire/repotown/services/banking/internal/repository"
)

// CustomerService provides business logic for customers
type CustomerService interface {
	GetCustomer(ctx context.Context, id string) (*domain.Customer, error)
	GetCustomerByEmail(ctx context.Context, email string) (*domain.Customer, error)
	ListCustomers(ctx context.Context, limit, offset int) ([]*domain.Customer, error)
	CreateCustomer(ctx context.Context, firstName, lastName, email, phoneNumber string, address domain.Address) (*domain.Customer, error)
	UpdateCustomer(ctx context.Context, id, firstName, lastName, email, phoneNumber string, address domain.Address) (*domain.Customer, error)
	UpdateCustomerStatus(ctx context.Context, id string, status domain.CustomerStatus) (*domain.Customer, error)
	DeleteCustomer(ctx context.Context, id string) error
}

// customerService implements CustomerService
type customerService struct {
	customerRepo repository.CustomerRepository
	logger       logging.Logger
}

// NewCustomerService creates a new customer service
func NewCustomerService(customerRepo repository.CustomerRepository, logger logging.Logger) CustomerService {
	return &customerService{
		customerRepo: customerRepo,
		logger:       logger,
	}
}

// GetCustomer retrieves a customer by ID
func (s *customerService) GetCustomer(ctx context.Context, id string) (*domain.Customer, error) {
	return s.customerRepo.GetByID(ctx, id)
}

// GetCustomerByEmail retrieves a customer by email
func (s *customerService) GetCustomerByEmail(ctx context.Context, email string) (*domain.Customer, error) {
	return s.customerRepo.GetByEmail(ctx, email)
}

// ListCustomers retrieves customers with pagination
func (s *customerService) ListCustomers(ctx context.Context, limit, offset int) ([]*domain.Customer, error) {
	return s.customerRepo.List(ctx, limit, offset)
}

// CreateCustomer creates a new customer
func (s *customerService) CreateCustomer(ctx context.Context, firstName, lastName, email, phoneNumber string, address domain.Address) (*domain.Customer, error) {
	// Check if customer with email already exists
	existingCustomer, err := s.customerRepo.GetByEmail(ctx, email)
	if err == nil && existingCustomer != nil {
		return nil, domain.ErrCustomerExists
	}

	// Create new customer
	customer := domain.NewCustomer(firstName, lastName, email, phoneNumber, address)
	if err := s.customerRepo.Create(ctx, customer); err != nil {
		return nil, err
	}

	return customer, nil
}

// UpdateCustomer updates customer details
func (s *customerService) UpdateCustomer(ctx context.Context, id, firstName, lastName, email, phoneNumber string, address domain.Address) (*domain.Customer, error) {
	// Get existing customer
	customer, err := s.customerRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Check if email is being changed and new email already exists
	if email != customer.Email {
		existingCustomer, err := s.customerRepo.GetByEmail(ctx, email)
		if err == nil && existingCustomer != nil && existingCustomer.ID != id {
			return nil, domain.ErrCustomerExists
		}
	}

	// Update customer
	customer.FirstName = firstName
	customer.LastName = lastName
	customer.Email = email
	customer.PhoneNumber = phoneNumber
	customer.Address = address

	if err := s.customerRepo.Update(ctx, customer); err != nil {
		return nil, err
	}

	return customer, nil
}

// UpdateCustomerStatus updates a customer's status
func (s *customerService) UpdateCustomerStatus(ctx context.Context, id string, status domain.CustomerStatus) (*domain.Customer, error) {
	// Get existing customer
	customer, err := s.customerRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Update status
	customer.UpdateStatus(status)

	if err := s.customerRepo.Update(ctx, customer); err != nil {
		return nil, err
	}

	return customer, nil
}

// DeleteCustomer deletes a customer
func (s *customerService) DeleteCustomer(ctx context.Context, id string) error {
	return s.customerRepo.Delete(ctx, id)
}

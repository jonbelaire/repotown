package service

import (
	"context"

	"github.com/jonbelaire/repotown/packages/go-core/logging"
	"github.com/jonbelaire/repotown/services/treasury/internal/domain"
	"github.com/jonbelaire/repotown/services/treasury/internal/repository"
)

// TaxpayerService provides business logic for taxpayer management
type TaxpayerService interface {
	GetTaxpayer(ctx context.Context, id string) (*domain.Taxpayer, error)
	GetTaxpayerByTaxIdentifier(ctx context.Context, identifier string) (*domain.Taxpayer, error)
	ListTaxpayers(ctx context.Context, limit, offset int) ([]*domain.Taxpayer, error)
	ListTaxpayersByType(ctx context.Context, taxpayerType domain.TaxpayerType) ([]*domain.Taxpayer, error)
	ListTaxpayersByStatus(ctx context.Context, status domain.TaxpayerStatus) ([]*domain.Taxpayer, error)
	ListBusinessesByIndustry(ctx context.Context, industry string) ([]*domain.Taxpayer, error)
	CreateTaxpayer(ctx context.Context, taxpayerType domain.TaxpayerType, name, taxIdentifier, contactEmail, contactPhone string, address domain.Address, exemptionCodes []string, annualRevenue int64, businessType, industry string) (*domain.Taxpayer, error)
	UpdateTaxpayerStatus(ctx context.Context, id string, status domain.TaxpayerStatus) (*domain.Taxpayer, error)
	UpdateTaxpayerContact(ctx context.Context, id, email, phone string) (*domain.Taxpayer, error)
	UpdateTaxpayerAddress(ctx context.Context, id string, address domain.Address) (*domain.Taxpayer, error)
	UpdateBusinessInfo(ctx context.Context, id string, annualRevenue int64, businessType, industry string) (*domain.Taxpayer, error)
	AddExemption(ctx context.Context, id, code string) (*domain.Taxpayer, error)
	RemoveExemption(ctx context.Context, id, code string) (*domain.Taxpayer, error)
	SearchTaxpayers(ctx context.Context, query string, limit int) ([]*domain.Taxpayer, error)
}

// taxpayerService implements TaxpayerService
type taxpayerService struct {
	taxpayerRepo repository.TaxpayerRepository
	logger       logging.Logger
}

// NewTaxpayerService creates a new taxpayer service
func NewTaxpayerService(taxpayerRepo repository.TaxpayerRepository, logger logging.Logger) TaxpayerService {
	return &taxpayerService{
		taxpayerRepo: taxpayerRepo,
		logger:       logger,
	}
}

// GetTaxpayer retrieves a taxpayer by ID
func (s *taxpayerService) GetTaxpayer(ctx context.Context, id string) (*domain.Taxpayer, error) {
	return s.taxpayerRepo.GetByID(ctx, id)
}

// GetTaxpayerByTaxIdentifier retrieves a taxpayer by tax identifier
func (s *taxpayerService) GetTaxpayerByTaxIdentifier(ctx context.Context, identifier string) (*domain.Taxpayer, error) {
	return s.taxpayerRepo.GetByTaxIdentifier(ctx, identifier)
}

// ListTaxpayers retrieves taxpayers with pagination
func (s *taxpayerService) ListTaxpayers(ctx context.Context, limit, offset int) ([]*domain.Taxpayer, error) {
	return s.taxpayerRepo.List(ctx, limit, offset)
}

// ListTaxpayersByType retrieves taxpayers by type
func (s *taxpayerService) ListTaxpayersByType(ctx context.Context, taxpayerType domain.TaxpayerType) ([]*domain.Taxpayer, error) {
	return s.taxpayerRepo.ListByType(ctx, taxpayerType)
}

// ListTaxpayersByStatus retrieves taxpayers by status
func (s *taxpayerService) ListTaxpayersByStatus(ctx context.Context, status domain.TaxpayerStatus) ([]*domain.Taxpayer, error) {
	return s.taxpayerRepo.ListByStatus(ctx, status)
}

// ListBusinessesByIndustry retrieves business taxpayers by industry
func (s *taxpayerService) ListBusinessesByIndustry(ctx context.Context, industry string) ([]*domain.Taxpayer, error) {
	return s.taxpayerRepo.ListBusinessesByIndustry(ctx, industry)
}

// CreateTaxpayer creates a new taxpayer
func (s *taxpayerService) CreateTaxpayer(
	ctx context.Context,
	taxpayerType domain.TaxpayerType,
	name, taxIdentifier, contactEmail, contactPhone string,
	address domain.Address,
	exemptionCodes []string,
	annualRevenue int64,
	businessType, industry string,
) (*domain.Taxpayer, error) {
	// Check if taxpayer with this identifier already exists
	existingTaxpayer, err := s.taxpayerRepo.GetByTaxIdentifier(ctx, taxIdentifier)
	if err == nil && existingTaxpayer != nil {
		return nil, domain.ErrTaxpayerExists
	}

	taxpayer := domain.NewTaxpayer(
		taxpayerType,
		name,
		taxIdentifier,
		contactEmail,
		contactPhone,
		address,
		exemptionCodes,
		annualRevenue,
		businessType,
		industry,
	)

	if err := s.taxpayerRepo.Create(ctx, taxpayer); err != nil {
		return nil, err
	}

	return taxpayer, nil
}

// UpdateTaxpayerStatus updates a taxpayer's status
func (s *taxpayerService) UpdateTaxpayerStatus(ctx context.Context, id string, status domain.TaxpayerStatus) (*domain.Taxpayer, error) {
	taxpayer, err := s.taxpayerRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	taxpayer.UpdateStatus(status)

	if err := s.taxpayerRepo.Update(ctx, taxpayer); err != nil {
		return nil, err
	}

	return taxpayer, nil
}

// UpdateTaxpayerContact updates a taxpayer's contact information
func (s *taxpayerService) UpdateTaxpayerContact(ctx context.Context, id, email, phone string) (*domain.Taxpayer, error) {
	taxpayer, err := s.taxpayerRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	taxpayer.UpdateContact(email, phone)

	if err := s.taxpayerRepo.Update(ctx, taxpayer); err != nil {
		return nil, err
	}

	return taxpayer, nil
}

// UpdateTaxpayerAddress updates a taxpayer's address
func (s *taxpayerService) UpdateTaxpayerAddress(ctx context.Context, id string, address domain.Address) (*domain.Taxpayer, error) {
	taxpayer, err := s.taxpayerRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	taxpayer.UpdateAddress(address)

	if err := s.taxpayerRepo.Update(ctx, taxpayer); err != nil {
		return nil, err
	}

	return taxpayer, nil
}

// UpdateBusinessInfo updates business-specific information
func (s *taxpayerService) UpdateBusinessInfo(ctx context.Context, id string, annualRevenue int64, businessType, industry string) (*domain.Taxpayer, error) {
	taxpayer, err := s.taxpayerRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	taxpayer.UpdateBusinessInfo(annualRevenue, businessType, industry)

	if err := s.taxpayerRepo.Update(ctx, taxpayer); err != nil {
		return nil, err
	}

	return taxpayer, nil
}

// AddExemption adds a tax exemption code to a taxpayer
func (s *taxpayerService) AddExemption(ctx context.Context, id, code string) (*domain.Taxpayer, error) {
	taxpayer, err := s.taxpayerRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	taxpayer.AddExemption(code)

	if err := s.taxpayerRepo.Update(ctx, taxpayer); err != nil {
		return nil, err
	}

	return taxpayer, nil
}

// RemoveExemption removes a tax exemption code from a taxpayer
func (s *taxpayerService) RemoveExemption(ctx context.Context, id, code string) (*domain.Taxpayer, error) {
	taxpayer, err := s.taxpayerRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	taxpayer.RemoveExemption(code)

	if err := s.taxpayerRepo.Update(ctx, taxpayer); err != nil {
		return nil, err
	}

	return taxpayer, nil
}

// SearchTaxpayers searches for taxpayers by name, tax identifier, or other identifiable information
func (s *taxpayerService) SearchTaxpayers(ctx context.Context, query string, limit int) ([]*domain.Taxpayer, error) {
	return s.taxpayerRepo.Search(ctx, query, limit)
}
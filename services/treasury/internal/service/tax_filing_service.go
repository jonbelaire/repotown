package service

import (
	"context"
	"time"

	"github.com/jonbelaire/repotown/packages/go-core/logging"
	"github.com/jonbelaire/repotown/services/treasury/internal/domain"
	"github.com/jonbelaire/repotown/services/treasury/internal/repository"
)

// TaxFilingService provides business logic for tax filing management
type TaxFilingService interface {
	GetTaxFiling(ctx context.Context, id string) (*domain.TaxFiling, error)
	ListTaxFilings(ctx context.Context, limit, offset int) ([]*domain.TaxFiling, error)
	ListTaxFilingsByTaxpayer(ctx context.Context, taxpayerID string, limit, offset int) ([]*domain.TaxFiling, error)
	ListTaxFilingsByStatus(ctx context.Context, status domain.FilingStatus, limit, offset int) ([]*domain.TaxFiling, error)
	ListTaxFilingsByPeriod(ctx context.Context, year int, period domain.FilingPeriod) ([]*domain.TaxFiling, error)
	ListOverdueTaxFilings(ctx context.Context) ([]*domain.TaxFiling, error)
	ListRecentlySubmittedFilings(ctx context.Context, days int) ([]*domain.TaxFiling, error)
	CreateTaxFiling(ctx context.Context, taxpayerID string, taxYear int, period domain.FilingPeriod, periodStart, periodEnd time.Time, filingType domain.TaxType, dueDate time.Time) (*domain.TaxFiling, error)
	UpdateTaxFilingAmounts(ctx context.Context, id string, grossIncome, taxableIncome, totalSales, taxableAmount int64) (*domain.TaxFiling, error)
	AddDeduction(ctx context.Context, id, code, description string, amount int64) (*domain.TaxFiling, error)
	AddCredit(ctx context.Context, id, code, description string, amount int64) (*domain.TaxFiling, error)
	SubmitTaxFiling(ctx context.Context, id string) (*domain.TaxFiling, error)
	ProcessTaxFiling(ctx context.Context, id string, taxCalculated int64) (*domain.TaxFiling, error)
	AcceptTaxFiling(ctx context.Context, id string) (*domain.TaxFiling, error)
	RejectTaxFiling(ctx context.Context, id, reason string) (*domain.TaxFiling, error)
	AmendTaxFiling(ctx context.Context, id string) (*domain.TaxFiling, error)
}

// taxFilingService implements TaxFilingService
type taxFilingService struct {
	taxFilingRepo  repository.TaxFilingRepository
	taxpayerRepo   repository.TaxpayerRepository
	taxRateRepo    repository.TaxRateRepository
	paymentService TaxPaymentService
	logger         logging.Logger
}

// NewTaxFilingService creates a new tax filing service
func NewTaxFilingService(
	taxFilingRepo repository.TaxFilingRepository,
	taxpayerRepo repository.TaxpayerRepository,
	taxRateRepo repository.TaxRateRepository,
	paymentService TaxPaymentService,
	logger logging.Logger,
) TaxFilingService {
	return &taxFilingService{
		taxFilingRepo:  taxFilingRepo,
		taxpayerRepo:   taxpayerRepo,
		taxRateRepo:    taxRateRepo,
		paymentService: paymentService,
		logger:         logger,
	}
}

// GetTaxFiling retrieves a tax filing by ID
func (s *taxFilingService) GetTaxFiling(ctx context.Context, id string) (*domain.TaxFiling, error) {
	return s.taxFilingRepo.GetByID(ctx, id)
}

// ListTaxFilings retrieves tax filings with pagination
func (s *taxFilingService) ListTaxFilings(ctx context.Context, limit, offset int) ([]*domain.TaxFiling, error) {
	return s.taxFilingRepo.List(ctx, limit, offset)
}

// ListTaxFilingsByTaxpayer retrieves tax filings for a specific taxpayer
func (s *taxFilingService) ListTaxFilingsByTaxpayer(ctx context.Context, taxpayerID string, limit, offset int) ([]*domain.TaxFiling, error) {
	return s.taxFilingRepo.ListByTaxpayer(ctx, taxpayerID, limit, offset)
}

// ListTaxFilingsByStatus retrieves tax filings by status
func (s *taxFilingService) ListTaxFilingsByStatus(ctx context.Context, status domain.FilingStatus, limit, offset int) ([]*domain.TaxFiling, error) {
	return s.taxFilingRepo.ListByStatus(ctx, status, limit, offset)
}

// ListTaxFilingsByPeriod retrieves tax filings for a specific period
func (s *taxFilingService) ListTaxFilingsByPeriod(ctx context.Context, year int, period domain.FilingPeriod) ([]*domain.TaxFiling, error) {
	return s.taxFilingRepo.ListByPeriod(ctx, year, period)
}

// ListOverdueTaxFilings retrieves overdue tax filings
func (s *taxFilingService) ListOverdueTaxFilings(ctx context.Context) ([]*domain.TaxFiling, error) {
	return s.taxFilingRepo.ListOverdue(ctx)
}

// ListRecentlySubmittedFilings retrieves recently submitted tax filings
func (s *taxFilingService) ListRecentlySubmittedFilings(ctx context.Context, days int) ([]*domain.TaxFiling, error) {
	return s.taxFilingRepo.ListRecentlySubmitted(ctx, days)
}

// CreateTaxFiling creates a new tax filing
func (s *taxFilingService) CreateTaxFiling(
	ctx context.Context,
	taxpayerID string,
	taxYear int,
	period domain.FilingPeriod,
	periodStart, periodEnd time.Time,
	filingType domain.TaxType,
	dueDate time.Time,
) (*domain.TaxFiling, error) {
	// Verify taxpayer exists
	if _, err := s.taxpayerRepo.GetByID(ctx, taxpayerID); err != nil {
		return nil, err
	}

	filing, err := domain.NewTaxFiling(
		taxpayerID,
		taxYear,
		period,
		periodStart,
		periodEnd,
		filingType,
		dueDate,
	)
	if err != nil {
		return nil, err
	}

	if err := s.taxFilingRepo.Create(ctx, filing); err != nil {
		return nil, err
	}

	return filing, nil
}

// UpdateTaxFilingAmounts updates the amounts in a tax filing
func (s *taxFilingService) UpdateTaxFilingAmounts(
	ctx context.Context,
	id string,
	grossIncome, taxableIncome, totalSales, taxableAmount int64,
) (*domain.TaxFiling, error) {
	filing, err := s.taxFilingRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Only update if the filing is in draft status
	if filing.Status != domain.FilingStatusDraft {
		return nil, domain.ErrInvalidFilingStatus
	}

	filing.GrossIncome = grossIncome
	filing.TaxableIncome = taxableIncome
	filing.TotalSales = totalSales
	filing.TaxableAmount = taxableAmount
	filing.UpdatedAt = time.Now()

	if err := s.taxFilingRepo.Update(ctx, filing); err != nil {
		return nil, err
	}

	return filing, nil
}

// AddDeduction adds a deduction to a tax filing
func (s *taxFilingService) AddDeduction(ctx context.Context, id, code, description string, amount int64) (*domain.TaxFiling, error) {
	filing, err := s.taxFilingRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Only add deduction if the filing is in draft status
	if filing.Status != domain.FilingStatusDraft {
		return nil, domain.ErrInvalidFilingStatus
	}

	filing.AddDeduction(code, description, amount)

	if err := s.taxFilingRepo.Update(ctx, filing); err != nil {
		return nil, err
	}

	return filing, nil
}

// AddCredit adds a credit to a tax filing
func (s *taxFilingService) AddCredit(ctx context.Context, id, code, description string, amount int64) (*domain.TaxFiling, error) {
	filing, err := s.taxFilingRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Only add credit if the filing is in draft status
	if filing.Status != domain.FilingStatusDraft {
		return nil, domain.ErrInvalidFilingStatus
	}

	filing.AddCredit(code, description, amount)

	if err := s.taxFilingRepo.Update(ctx, filing); err != nil {
		return nil, err
	}

	return filing, nil
}

// SubmitTaxFiling submits a tax filing
func (s *taxFilingService) SubmitTaxFiling(ctx context.Context, id string) (*domain.TaxFiling, error) {
	filing, err := s.taxFilingRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if err := filing.SubmitFiling(); err != nil {
		return nil, err
	}

	if err := s.taxFilingRepo.Update(ctx, filing); err != nil {
		return nil, err
	}

	return filing, nil
}

// ProcessTaxFiling processes a tax filing and calculates tax
func (s *taxFilingService) ProcessTaxFiling(ctx context.Context, id string, taxCalculated int64) (*domain.TaxFiling, error) {
	filing, err := s.taxFilingRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Only process if filing is submitted
	if filing.Status != domain.FilingStatusSubmitted {
		return nil, domain.ErrInvalidFilingStatus
	}

	// Update processing status
	filing.Status = domain.FilingStatusProcessing
	filing.TaxCalculated = taxCalculated
	filing.UpdatedAt = time.Now()

	if err := s.taxFilingRepo.Update(ctx, filing); err != nil {
		return nil, err
	}

	return filing, nil
}

// AcceptTaxFiling accepts a tax filing
func (s *taxFilingService) AcceptTaxFiling(ctx context.Context, id string) (*domain.TaxFiling, error) {
	filing, err := s.taxFilingRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if err := filing.AcceptFiling(); err != nil {
		return nil, err
	}

	if err := s.taxFilingRepo.Update(ctx, filing); err != nil {
		return nil, err
	}

	return filing, nil
}

// RejectTaxFiling rejects a tax filing
func (s *taxFilingService) RejectTaxFiling(ctx context.Context, id, reason string) (*domain.TaxFiling, error) {
	filing, err := s.taxFilingRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if err := filing.RejectFiling(reason); err != nil {
		return nil, err
	}

	if err := s.taxFilingRepo.Update(ctx, filing); err != nil {
		return nil, err
	}

	return filing, nil
}

// AmendTaxFiling creates an amended tax filing based on an existing filing
func (s *taxFilingService) AmendTaxFiling(ctx context.Context, id string) (*domain.TaxFiling, error) {
	filing, err := s.taxFilingRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	amendedFiling := filing.AmendFiling()

	if err := s.taxFilingRepo.Create(ctx, amendedFiling); err != nil {
		return nil, err
	}

	return amendedFiling, nil
}
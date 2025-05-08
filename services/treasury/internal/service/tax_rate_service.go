package service

import (
	"context"
	"time"

	"github.com/jonbelaire/repotown/packages/go-core/logging"
	"github.com/jonbelaire/repotown/services/treasury/internal/domain"
	"github.com/jonbelaire/repotown/services/treasury/internal/repository"
)

// TaxRateService provides business logic for tax rate management
type TaxRateService interface {
	GetTaxRate(ctx context.Context, id string) (*domain.TaxRate, error)
	ListTaxRates(ctx context.Context, limit, offset int) ([]*domain.TaxRate, error)
	ListTaxRatesByType(ctx context.Context, taxType domain.TaxType) ([]*domain.TaxRate, error)
	ListActiveTaxRates(ctx context.Context) ([]*domain.TaxRate, error)
	ListTaxRatesByJurisdiction(ctx context.Context, jurisdictionCode string) ([]*domain.TaxRate, error)
	GetTaxRatesForIncome(ctx context.Context, amount int64, jurisdictionCode string) ([]*domain.TaxRate, error)
	CalculateIncomeTax(ctx context.Context, amount int64, jurisdictionCode string) (int64, error)
	CalculateSalesTax(ctx context.Context, amount int64, category string, jurisdictionCode string) (int64, error)
	CreateTaxRate(ctx context.Context, taxType domain.TaxType, name, description string, rate float64, bracketType domain.TaxBracketType, minAmount, maxAmount int64, category, jurisdictionCode string, effectiveDate time.Time, expirationDate *time.Time) (*domain.TaxRate, error)
	UpdateTaxRate(ctx context.Context, id, name, description string, rate float64, category string, effectiveDate time.Time, expirationDate *time.Time) (*domain.TaxRate, error)
	ActivateTaxRate(ctx context.Context, id string) (*domain.TaxRate, error)
	DeactivateTaxRate(ctx context.Context, id string) (*domain.TaxRate, error)
	ArchiveTaxRate(ctx context.Context, id string) (*domain.TaxRate, error)
}

// taxRateService implements TaxRateService
type taxRateService struct {
	taxRateRepo repository.TaxRateRepository
	logger      logging.Logger
}

// NewTaxRateService creates a new tax rate service
func NewTaxRateService(taxRateRepo repository.TaxRateRepository, logger logging.Logger) TaxRateService {
	return &taxRateService{
		taxRateRepo: taxRateRepo,
		logger:      logger,
	}
}

// GetTaxRate retrieves a tax rate by ID
func (s *taxRateService) GetTaxRate(ctx context.Context, id string) (*domain.TaxRate, error) {
	return s.taxRateRepo.GetByID(ctx, id)
}

// ListTaxRates retrieves tax rates with pagination
func (s *taxRateService) ListTaxRates(ctx context.Context, limit, offset int) ([]*domain.TaxRate, error) {
	return s.taxRateRepo.List(ctx, limit, offset)
}

// ListTaxRatesByType retrieves tax rates by type
func (s *taxRateService) ListTaxRatesByType(ctx context.Context, taxType domain.TaxType) ([]*domain.TaxRate, error) {
	return s.taxRateRepo.ListByType(ctx, taxType)
}

// ListActiveTaxRates retrieves active tax rates
func (s *taxRateService) ListActiveTaxRates(ctx context.Context) ([]*domain.TaxRate, error) {
	return s.taxRateRepo.ListActive(ctx)
}

// ListTaxRatesByJurisdiction retrieves tax rates by jurisdiction
func (s *taxRateService) ListTaxRatesByJurisdiction(ctx context.Context, jurisdictionCode string) ([]*domain.TaxRate, error) {
	return s.taxRateRepo.ListByJurisdiction(ctx, jurisdictionCode)
}

// GetTaxRatesForIncome retrieves applicable tax rates for income
func (s *taxRateService) GetTaxRatesForIncome(ctx context.Context, amount int64, jurisdictionCode string) ([]*domain.TaxRate, error) {
	return s.taxRateRepo.GetRatesForIncome(ctx, amount, jurisdictionCode)
}

// CalculateIncomeTax calculates income tax based on amount and jurisdiction
func (s *taxRateService) CalculateIncomeTax(ctx context.Context, amount int64, jurisdictionCode string) (int64, error) {
	rates, err := s.taxRateRepo.GetRatesForIncome(ctx, amount, jurisdictionCode)
	if err != nil {
		return 0, err
	}

	var totalTax int64 = 0
	for _, rate := range rates {
		if rate.Type == domain.TaxTypeIncome && rate.IsActive() {
			totalTax += rate.CalculateTax(amount)
		}
	}

	return totalTax, nil
}

// CalculateSalesTax calculates sales tax based on amount, category, and jurisdiction
func (s *taxRateService) CalculateSalesTax(ctx context.Context, amount int64, category string, jurisdictionCode string) (int64, error) {
	rates, err := s.taxRateRepo.ListByJurisdiction(ctx, jurisdictionCode)
	if err != nil {
		return 0, err
	}

	var salesTaxRate *domain.TaxRate
	// Find the appropriate sales tax rate based on category
	for _, rate := range rates {
		if rate.Type == domain.TaxTypeSales && rate.IsActive() {
			if category != "" && rate.Category == category {
				salesTaxRate = rate
				break
			} else if rate.Category == "" && salesTaxRate == nil {
				// Use the default sales tax if no category-specific tax is found
				salesTaxRate = rate
			}
		}
	}

	if salesTaxRate == nil {
		return 0, nil // No applicable sales tax found
	}

	return salesTaxRate.CalculateTax(amount), nil
}

// CreateTaxRate creates a new tax rate
func (s *taxRateService) CreateTaxRate(
	ctx context.Context,
	taxType domain.TaxType,
	name, description string,
	rate float64,
	bracketType domain.TaxBracketType,
	minAmount, maxAmount int64,
	category, jurisdictionCode string,
	effectiveDate time.Time,
	expirationDate *time.Time,
) (*domain.TaxRate, error) {
	taxRate, err := domain.NewTaxRate(
		taxType,
		name,
		description,
		rate,
		bracketType,
		minAmount,
		maxAmount,
		category,
		jurisdictionCode,
		effectiveDate,
		expirationDate,
	)
	if err != nil {
		return nil, err
	}

	if err := s.taxRateRepo.Create(ctx, taxRate); err != nil {
		return nil, err
	}

	return taxRate, nil
}

// UpdateTaxRate updates an existing tax rate
func (s *taxRateService) UpdateTaxRate(
	ctx context.Context,
	id, name, description string,
	rate float64,
	category string,
	effectiveDate time.Time,
	expirationDate *time.Time,
) (*domain.TaxRate, error) {
	taxRate, err := s.taxRateRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if err := taxRate.UpdateTaxRate(
		name,
		description,
		rate,
		category,
		effectiveDate,
		expirationDate,
	); err != nil {
		return nil, err
	}

	if err := s.taxRateRepo.Update(ctx, taxRate); err != nil {
		return nil, err
	}

	return taxRate, nil
}

// ActivateTaxRate activates a tax rate
func (s *taxRateService) ActivateTaxRate(ctx context.Context, id string) (*domain.TaxRate, error) {
	taxRate, err := s.taxRateRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	taxRate.ActivateTaxRate()

	if err := s.taxRateRepo.Update(ctx, taxRate); err != nil {
		return nil, err
	}

	return taxRate, nil
}

// DeactivateTaxRate deactivates a tax rate
func (s *taxRateService) DeactivateTaxRate(ctx context.Context, id string) (*domain.TaxRate, error) {
	taxRate, err := s.taxRateRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	taxRate.DeactivateTaxRate()

	if err := s.taxRateRepo.Update(ctx, taxRate); err != nil {
		return nil, err
	}

	return taxRate, nil
}

// ArchiveTaxRate archives a tax rate
func (s *taxRateService) ArchiveTaxRate(ctx context.Context, id string) (*domain.TaxRate, error) {
	taxRate, err := s.taxRateRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	taxRate.ArchiveTaxRate()

	if err := s.taxRateRepo.Update(ctx, taxRate); err != nil {
		return nil, err
	}

	return taxRate, nil
}
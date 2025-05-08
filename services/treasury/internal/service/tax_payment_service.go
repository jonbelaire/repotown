package service

import (
	"context"
	"time"

	"github.com/jonbelaire/repotown/packages/go-core/logging"
	"github.com/jonbelaire/repotown/services/treasury/internal/domain"
	"github.com/jonbelaire/repotown/services/treasury/internal/repository"
)

// TaxPaymentService provides business logic for tax payment management
type TaxPaymentService interface {
	GetTaxPayment(ctx context.Context, id string) (*domain.TaxPayment, error)
	GetTaxPaymentByConfirmationCode(ctx context.Context, code string) (*domain.TaxPayment, error)
	ListTaxPayments(ctx context.Context, limit, offset int) ([]*domain.TaxPayment, error)
	ListTaxPaymentsByTaxpayer(ctx context.Context, taxpayerID string, limit, offset int) ([]*domain.TaxPayment, error)
	ListTaxPaymentsByFiling(ctx context.Context, filingID string) ([]*domain.TaxPayment, error)
	ListTaxPaymentsByStatus(ctx context.Context, status domain.PaymentStatus) ([]*domain.TaxPayment, error)
	ListTaxPaymentsByDateRange(ctx context.Context, startDate, endDate time.Time) ([]*domain.TaxPayment, error)
	ListRecentTaxPayments(ctx context.Context, days int) ([]*domain.TaxPayment, error)
	GetTotalPaymentsByTaxType(ctx context.Context, taxType domain.TaxType, startDate, endDate time.Time) (int64, error)
	CreateTaxPayment(ctx context.Context, taxpayerID, filingID string, taxType domain.TaxType, amount int64, paymentMethod domain.PaymentMethod, paymentDate time.Time, notes string) (*domain.TaxPayment, error)
	ProcessTaxPayment(ctx context.Context, id string) (*domain.TaxPayment, error)
	MarkTaxPaymentAsFailed(ctx context.Context, id, reason string) (*domain.TaxPayment, error)
	RefundTaxPayment(ctx context.Context, id, reason string) (*domain.TaxPayment, error)
	VoidTaxPayment(ctx context.Context, id, reason string) (*domain.TaxPayment, error)
	UpdateTaxPaymentAmount(ctx context.Context, id string, amount int64) (*domain.TaxPayment, error)
}

// taxPaymentService implements TaxPaymentService
type taxPaymentService struct {
	taxPaymentRepo repository.TaxPaymentRepository
	taxFilingRepo  repository.TaxFilingRepository
	logger         logging.Logger
}

// NewTaxPaymentService creates a new tax payment service
func NewTaxPaymentService(
	taxPaymentRepo repository.TaxPaymentRepository,
	taxFilingRepo repository.TaxFilingRepository,
	logger logging.Logger,
) TaxPaymentService {
	return &taxPaymentService{
		taxPaymentRepo: taxPaymentRepo,
		taxFilingRepo:  taxFilingRepo,
		logger:         logger,
	}
}

// GetTaxPayment retrieves a tax payment by ID
func (s *taxPaymentService) GetTaxPayment(ctx context.Context, id string) (*domain.TaxPayment, error) {
	return s.taxPaymentRepo.GetByID(ctx, id)
}

// GetTaxPaymentByConfirmationCode retrieves a tax payment by confirmation code
func (s *taxPaymentService) GetTaxPaymentByConfirmationCode(ctx context.Context, code string) (*domain.TaxPayment, error) {
	return s.taxPaymentRepo.GetByConfirmationCode(ctx, code)
}

// ListTaxPayments retrieves tax payments with pagination
func (s *taxPaymentService) ListTaxPayments(ctx context.Context, limit, offset int) ([]*domain.TaxPayment, error) {
	return s.taxPaymentRepo.List(ctx, limit, offset)
}

// ListTaxPaymentsByTaxpayer retrieves tax payments for a specific taxpayer
func (s *taxPaymentService) ListTaxPaymentsByTaxpayer(ctx context.Context, taxpayerID string, limit, offset int) ([]*domain.TaxPayment, error) {
	return s.taxPaymentRepo.ListByTaxpayer(ctx, taxpayerID, limit, offset)
}

// ListTaxPaymentsByFiling retrieves tax payments for a specific filing
func (s *taxPaymentService) ListTaxPaymentsByFiling(ctx context.Context, filingID string) ([]*domain.TaxPayment, error) {
	return s.taxPaymentRepo.ListByFiling(ctx, filingID)
}

// ListTaxPaymentsByStatus retrieves tax payments by status
func (s *taxPaymentService) ListTaxPaymentsByStatus(ctx context.Context, status domain.PaymentStatus) ([]*domain.TaxPayment, error) {
	return s.taxPaymentRepo.ListByStatus(ctx, status)
}

// ListTaxPaymentsByDateRange retrieves tax payments within a date range
func (s *taxPaymentService) ListTaxPaymentsByDateRange(ctx context.Context, startDate, endDate time.Time) ([]*domain.TaxPayment, error) {
	return s.taxPaymentRepo.ListByDateRange(ctx, startDate, endDate)
}

// ListRecentTaxPayments retrieves recent tax payments
func (s *taxPaymentService) ListRecentTaxPayments(ctx context.Context, days int) ([]*domain.TaxPayment, error) {
	return s.taxPaymentRepo.ListRecentPayments(ctx, days)
}

// GetTotalPaymentsByTaxType calculates the total payment amount by tax type within a date range
func (s *taxPaymentService) GetTotalPaymentsByTaxType(ctx context.Context, taxType domain.TaxType, startDate, endDate time.Time) (int64, error) {
	return s.taxPaymentRepo.GetTotalByTaxType(ctx, taxType, startDate, endDate)
}

// CreateTaxPayment creates a new tax payment
func (s *taxPaymentService) CreateTaxPayment(
	ctx context.Context,
	taxpayerID, filingID string,
	taxType domain.TaxType,
	amount int64,
	paymentMethod domain.PaymentMethod,
	paymentDate time.Time,
	notes string,
) (*domain.TaxPayment, error) {
	// Verify filing if provided
	if filingID != "" {
		if _, err := s.taxFilingRepo.GetByID(ctx, filingID); err != nil {
			return nil, err
		}
	}

	payment, err := domain.NewTaxPayment(
		taxpayerID,
		filingID,
		taxType,
		amount,
		paymentMethod,
		paymentDate,
		notes,
	)
	if err != nil {
		return nil, err
	}

	if err := s.taxPaymentRepo.Create(ctx, payment); err != nil {
		return nil, err
	}

	return payment, nil
}

// ProcessTaxPayment marks a tax payment as completed (processed)
func (s *taxPaymentService) ProcessTaxPayment(ctx context.Context, id string) (*domain.TaxPayment, error) {
	payment, err := s.taxPaymentRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	payment.MarkAsCompleted()

	if err := s.taxPaymentRepo.Update(ctx, payment); err != nil {
		return nil, err
	}

	// If payment is for a filing, update the filing's tax paid amount
	if payment.FilingID != "" {
		filing, err := s.taxFilingRepo.GetByID(ctx, payment.FilingID)
		if err != nil {
			s.logger.Error("Failed to get filing for payment", "filing_id", payment.FilingID, "error", err)
		} else {
			filing.TaxPaid += payment.Amount
			filing.UpdatedAt = time.Now()
			if err := s.taxFilingRepo.Update(ctx, filing); err != nil {
				s.logger.Error("Failed to update filing's paid amount", "filing_id", payment.FilingID, "error", err)
			}
		}
	}

	return payment, nil
}

// MarkTaxPaymentAsFailed marks a tax payment as failed
func (s *taxPaymentService) MarkTaxPaymentAsFailed(ctx context.Context, id, reason string) (*domain.TaxPayment, error) {
	payment, err := s.taxPaymentRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	payment.MarkAsFailed(reason)

	if err := s.taxPaymentRepo.Update(ctx, payment); err != nil {
		return nil, err
	}

	return payment, nil
}

// RefundTaxPayment refunds a tax payment
func (s *taxPaymentService) RefundTaxPayment(ctx context.Context, id, reason string) (*domain.TaxPayment, error) {
	payment, err := s.taxPaymentRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if err := payment.Refund(reason); err != nil {
		return nil, err
	}

	if err := s.taxPaymentRepo.Update(ctx, payment); err != nil {
		return nil, err
	}

	// If payment is for a filing, update the filing's tax paid amount
	if payment.FilingID != "" {
		filing, err := s.taxFilingRepo.GetByID(ctx, payment.FilingID)
		if err != nil {
			s.logger.Error("Failed to get filing for refunded payment", "filing_id", payment.FilingID, "error", err)
		} else {
			filing.TaxPaid -= payment.Amount
			if filing.TaxPaid < 0 {
				filing.TaxPaid = 0
			}
			filing.UpdatedAt = time.Now()
			if err := s.taxFilingRepo.Update(ctx, filing); err != nil {
				s.logger.Error("Failed to update filing's paid amount after refund", "filing_id", payment.FilingID, "error", err)
			}
		}
	}

	return payment, nil
}

// VoidTaxPayment voids a tax payment
func (s *taxPaymentService) VoidTaxPayment(ctx context.Context, id, reason string) (*domain.TaxPayment, error) {
	payment, err := s.taxPaymentRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if err := payment.Void(reason); err != nil {
		return nil, err
	}

	if err := s.taxPaymentRepo.Update(ctx, payment); err != nil {
		return nil, err
	}

	return payment, nil
}

// UpdateTaxPaymentAmount updates the payment amount
func (s *taxPaymentService) UpdateTaxPaymentAmount(ctx context.Context, id string, amount int64) (*domain.TaxPayment, error) {
	payment, err := s.taxPaymentRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if err := payment.UpdateAmount(amount); err != nil {
		return nil, err
	}

	if err := s.taxPaymentRepo.Update(ctx, payment); err != nil {
		return nil, err
	}

	return payment, nil
}
package service

import (
	"github.com/jonbelaire/repotown/packages/go-core/logging"
	"github.com/jonbelaire/repotown/services/treasury/internal/repository"
)

// Services holds all service implementations
type Services struct {
	TaxRate     TaxRateService
	Taxpayer    TaxpayerService
	TaxFiling   TaxFilingService
	TaxPayment  TaxPaymentService
	TaxReport   TaxReportService
}

// NewServices creates all services
func NewServices(repos *repository.Repositories, logger logging.Logger) *Services {
	paymentService := NewTaxPaymentService(repos.TaxPayment, repos.TaxFiling, logger)
	filingService := NewTaxFilingService(repos.TaxFiling, repos.Taxpayer, repos.TaxRate, paymentService, logger)

	return &Services{
		TaxRate:    NewTaxRateService(repos.TaxRate, logger),
		Taxpayer:   NewTaxpayerService(repos.Taxpayer, logger),
		TaxFiling:  filingService,
		TaxPayment: paymentService,
		TaxReport:  NewTaxReportService(repos.TaxPayment, repos.TaxFiling, repos.Taxpayer, logger),
	}
}
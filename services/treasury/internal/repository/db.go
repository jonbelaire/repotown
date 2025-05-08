package repository

import (
	"github.com/jonbelaire/repotown/packages/go-core/database"
	"github.com/jonbelaire/repotown/packages/go-core/logging"
	"github.com/jonbelaire/repotown/services/treasury/internal/config"
)

// Connect establishes a connection to the database
func Connect(cfg config.DatabaseConfig, logger logging.Logger) (*database.DB, error) {
	dbCfg := database.Config{
		DSN:             cfg.URL,
		MaxOpenConns:    cfg.MaxOpenConns,
		MaxIdleConns:    cfg.MaxIdleConns,
		ConnMaxLifetime: cfg.MaxLifetime,
		ConnMaxIdleTime: cfg.MaxIdleTime,
	}

	return database.New(dbCfg, logger)
}

// Repositories holds all repositories
type Repositories struct {
	TaxRate     TaxRateRepository
	Taxpayer    TaxpayerRepository
	TaxFiling   TaxFilingRepository
	TaxPayment  TaxPaymentRepository
}

// NewRepositories creates all repositories
func NewRepositories(db *database.DB, logger logging.Logger) *Repositories {
	return &Repositories{
		TaxRate:     NewTaxRateRepository(db, logger),
		Taxpayer:    NewTaxpayerRepository(db, logger),
		TaxFiling:   NewTaxFilingRepository(db, logger),
		TaxPayment:  NewTaxPaymentRepository(db, logger),
	}
}
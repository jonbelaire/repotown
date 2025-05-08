package repository

import (
	"github.com/jonbelaire/repotown/packages/go-core/database"
	"github.com/jonbelaire/repotown/packages/go-core/logging"
	"github.com/jonbelaire/repotown/services/banking/internal/config"
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
	Account     AccountRepository
	Transaction TransactionRepository
	Customer    CustomerRepository
}

// NewRepositories creates all repositories
func NewRepositories(db *database.DB, logger logging.Logger) *Repositories {
	return &Repositories{
		Account:     NewAccountRepository(db, logger),
		Transaction: NewTransactionRepository(db, logger),
		Customer:    NewCustomerRepository(db, logger),
	}
}
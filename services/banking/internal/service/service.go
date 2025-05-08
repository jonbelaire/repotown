package service

import (
	"github.com/jonbelaire/repotown/packages/go-core/logging"
	"github.com/jonbelaire/repotown/services/banking/internal/repository"
)

// Services holds all services
type Services struct {
	Account     AccountService
	Transaction TransactionService
	Customer    CustomerService
}

// NewServices creates all services
func NewServices(repos *repository.Repositories, logger logging.Logger) *Services {
	return &Services{
		Account:     NewAccountService(repos.Account, repos.Transaction, logger),
		Transaction: NewTransactionService(repos.Transaction, repos.Account, logger),
		Customer:    NewCustomerService(repos.Customer, logger),
	}
}

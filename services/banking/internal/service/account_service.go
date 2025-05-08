package service

import (
	"context"
	"errors"

	"github.com/jonbelaire/repotown/packages/go-core/logging"
	"github.com/jonbelaire/repotown/services/banking/internal/domain"
	"github.com/jonbelaire/repotown/services/banking/internal/repository"
)

// AccountService provides business logic for accounts
type AccountService interface {
	GetAccount(ctx context.Context, id string) (*domain.Account, error)
	ListAccounts(ctx context.Context, limit, offset int) ([]*domain.Account, error)
	ListAccountsByCustomer(ctx context.Context, customerID string) ([]*domain.Account, error)
	CreateAccount(ctx context.Context, customerID string, accType domain.AccountType, name, currencyCode string) (*domain.Account, error)
	UpdateAccount(ctx context.Context, id string, name string) (*domain.Account, error)
	CloseAccount(ctx context.Context, id string) error
	Deposit(ctx context.Context, accountID string, amount int64, description string) (*domain.Transaction, error)
	Withdraw(ctx context.Context, accountID string, amount int64, description string) (*domain.Transaction, error)
	Transfer(ctx context.Context, sourceID, targetID string, amount int64, description string) (*domain.Transaction, error)
}

// accountService implements AccountService
type accountService struct {
	accountRepo     repository.AccountRepository
	transactionRepo repository.TransactionRepository
	logger          logging.Logger
}

// NewAccountService creates a new account service
func NewAccountService(accountRepo repository.AccountRepository, transactionRepo repository.TransactionRepository, logger logging.Logger) AccountService {
	return &accountService{
		accountRepo:     accountRepo,
		transactionRepo: transactionRepo,
		logger:          logger,
	}
}

// GetAccount retrieves an account by ID
func (s *accountService) GetAccount(ctx context.Context, id string) (*domain.Account, error) {
	return s.accountRepo.GetByID(ctx, id)
}

// ListAccounts retrieves accounts with pagination
func (s *accountService) ListAccounts(ctx context.Context, limit, offset int) ([]*domain.Account, error) {
	return s.accountRepo.List(ctx, limit, offset)
}

// ListAccountsByCustomer retrieves accounts for a specific customer
func (s *accountService) ListAccountsByCustomer(ctx context.Context, customerID string) ([]*domain.Account, error) {
	return s.accountRepo.ListByCustomer(ctx, customerID)
}

// CreateAccount creates a new account
func (s *accountService) CreateAccount(ctx context.Context, customerID string, accType domain.AccountType, name, currencyCode string) (*domain.Account, error) {
	account := domain.NewAccount(customerID, accType, name, currencyCode)
	if err := s.accountRepo.Create(ctx, account); err != nil {
		return nil, err
	}
	return account, nil
}

// UpdateAccount updates an account's details
func (s *accountService) UpdateAccount(ctx context.Context, id string, name string) (*domain.Account, error) {
	account, err := s.accountRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if !account.IsActive() {
		return nil, domain.ErrAccountClosed
	}

	account.Name = name
	if err := s.accountRepo.Update(ctx, account); err != nil {
		return nil, err
	}

	return account, nil
}

// CloseAccount closes an account
func (s *accountService) CloseAccount(ctx context.Context, id string) error {
	account, err := s.accountRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	if !account.IsActive() {
		return domain.ErrAccountClosed
	}

	account.Close()
	return s.accountRepo.Update(ctx, account)
}

// Deposit adds funds to an account
func (s *accountService) Deposit(ctx context.Context, accountID string, amount int64, description string) (*domain.Transaction, error) {
	// In a real implementation, this would be wrapped in a transaction
	account, err := s.accountRepo.GetByID(ctx, accountID)
	if err != nil {
		return nil, err
	}

	if err := account.Deposit(amount); err != nil {
		return nil, err
	}

	if err := s.accountRepo.Update(ctx, account); err != nil {
		return nil, err
	}

	// Create transaction record
	tx := domain.NewTransaction(domain.TransactionTypeDeposit, accountID, amount, account.CurrencyCode, description)
	tx.Complete()

	// In a real implementation, we would save the transaction
	// This is just a placeholder
	return tx, errors.New("transaction repository not implemented")
}

// Withdraw removes funds from an account
func (s *accountService) Withdraw(ctx context.Context, accountID string, amount int64, description string) (*domain.Transaction, error) {
	// In a real implementation, this would be wrapped in a transaction
	account, err := s.accountRepo.GetByID(ctx, accountID)
	if err != nil {
		return nil, err
	}

	if err := account.Withdraw(amount); err != nil {
		return nil, err
	}

	if err := s.accountRepo.Update(ctx, account); err != nil {
		return nil, err
	}

	// Create transaction record
	tx := domain.NewTransaction(domain.TransactionTypeWithdrawal, accountID, amount, account.CurrencyCode, description)
	tx.Complete()

	// In a real implementation, we would save the transaction
	// This is just a placeholder
	return tx, errors.New("transaction repository not implemented")
}

// Transfer moves funds between accounts
func (s *accountService) Transfer(ctx context.Context, sourceID, targetID string, amount int64, description string) (*domain.Transaction, error) {
	// In a real implementation, this would be wrapped in a database transaction
	// Get source account
	sourceAccount, err := s.accountRepo.GetByID(ctx, sourceID)
	if err != nil {
		return nil, err
	}

	// Get target account
	targetAccount, err := s.accountRepo.GetByID(ctx, targetID)
	if err != nil {
		return nil, err
	}

	// Ensure accounts are active
	if !sourceAccount.IsActive() {
		return nil, domain.ErrAccountClosed
	}
	if !targetAccount.IsActive() {
		return nil, domain.ErrAccountClosed
	}

	// Withdraw from source
	if err := sourceAccount.Withdraw(amount); err != nil {
		return nil, err
	}

	// Deposit to target
	if err := targetAccount.Deposit(amount); err != nil {
		return nil, err
	}

	// Update accounts
	if err := s.accountRepo.Update(ctx, sourceAccount); err != nil {
		return nil, err
	}
	if err := s.accountRepo.Update(ctx, targetAccount); err != nil {
		return nil, err
	}

	// Create transaction record
	tx := domain.NewTransferTransaction(sourceID, targetID, amount, sourceAccount.CurrencyCode, description)
	tx.Complete()

	// In a real implementation, we would save the transaction
	// This is just a placeholder
	return tx, errors.New("transaction repository not implemented")
}
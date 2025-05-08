package service

import (
	"context"

	"github.com/jonbelaire/repotown/packages/go-core/logging"
	"github.com/jonbelaire/repotown/services/banking/internal/domain"
	"github.com/jonbelaire/repotown/services/banking/internal/repository"
)

// TransactionService provides business logic for transactions
type TransactionService interface {
	GetTransaction(ctx context.Context, id string) (*domain.Transaction, error)
	ListTransactions(ctx context.Context, limit, offset int) ([]*domain.Transaction, error)
	ListAccountTransactions(ctx context.Context, accountID string, limit, offset int) ([]*domain.Transaction, error)
	CreateTransaction(ctx context.Context, txType domain.TransactionType, accountID string, amount int64, currencyCode, description string) (*domain.Transaction, error)
}

// transactionService implements TransactionService
type transactionService struct {
	transactionRepo repository.TransactionRepository
	accountRepo     repository.AccountRepository
	logger          logging.Logger
}

// NewTransactionService creates a new transaction service
func NewTransactionService(
	transactionRepo repository.TransactionRepository,
	accountRepo repository.AccountRepository,
	logger logging.Logger,
) TransactionService {
	return &transactionService{
		transactionRepo: transactionRepo,
		accountRepo:     accountRepo,
		logger:          logger,
	}
}

// GetTransaction retrieves a transaction by ID
func (s *transactionService) GetTransaction(ctx context.Context, id string) (*domain.Transaction, error) {
	return s.transactionRepo.GetByID(ctx, id)
}

// ListTransactions retrieves transactions with pagination
func (s *transactionService) ListTransactions(ctx context.Context, limit, offset int) ([]*domain.Transaction, error) {
	return s.transactionRepo.List(ctx, limit, offset)
}

// ListAccountTransactions retrieves transactions for a specific account
func (s *transactionService) ListAccountTransactions(ctx context.Context, accountID string, limit, offset int) ([]*domain.Transaction, error) {
	return s.transactionRepo.ListByAccount(ctx, accountID, limit, offset)
}

// CreateTransaction creates a new transaction
func (s *transactionService) CreateTransaction(ctx context.Context, txType domain.TransactionType, accountID string, amount int64, currencyCode, description string) (*domain.Transaction, error) {
	// Verify account exists
	account, err := s.accountRepo.GetByID(ctx, accountID)
	if err != nil {
		return nil, err
	}

	// Create transaction
	transaction := domain.NewTransaction(txType, accountID, amount, currencyCode, description)
	
	// Save transaction (pending)
	if err := s.transactionRepo.Create(ctx, transaction); err != nil {
		return nil, err
	}

	// Process transaction based on type
	switch txType {
	case domain.TransactionTypeDeposit:
		if err := account.Deposit(amount); err != nil {
			transaction.Fail()
			s.transactionRepo.Update(ctx, transaction)
			return nil, err
		}
	case domain.TransactionTypeWithdrawal:
		if err := account.Withdraw(amount); err != nil {
			transaction.Fail()
			s.transactionRepo.Update(ctx, transaction)
			return nil, err
		}
	}

	// Update account and complete transaction
	if err := s.accountRepo.Update(ctx, account); err != nil {
		transaction.Fail()
		s.transactionRepo.Update(ctx, transaction)
		return nil, err
	}

	// Mark transaction as completed
	transaction.Complete()
	if err := s.transactionRepo.Update(ctx, transaction); err != nil {
		return nil, err
	}

	return transaction, nil
}
package domain

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

// Common errors
var (
	ErrInsufficientFunds = errors.New("insufficient funds")
	ErrAccountNotFound   = errors.New("account not found")
	ErrAccountClosed     = errors.New("account is closed")
)

// AccountType represents the type of account
type AccountType string

// AccountStatus represents the status of an account
type AccountStatus string

// Account types
const (
	AccountTypeSavings    AccountType = "savings"
	AccountTypeChecking   AccountType = "checking"
	AccountTypeCreditCard AccountType = "credit_card"
)

// Account statuses
const (
	AccountStatusActive   AccountStatus = "active"
	AccountStatusInactive AccountStatus = "inactive"
	AccountStatusClosed   AccountStatus = "closed"
)

// Account represents a bank account
type Account struct {
	ID           string        `json:"id"`
	CustomerID   string        `json:"customer_id"`
	Type         AccountType   `json:"type"`
	Status       AccountStatus `json:"status"`
	Balance      int64         `json:"balance"` // Stored in cents
	CurrencyCode string        `json:"currency_code"`
	Name         string        `json:"name"`
	Number       string        `json:"number"`
	CreatedAt    time.Time     `json:"created_at"`
	UpdatedAt    time.Time     `json:"updated_at"`
	ClosedAt     *time.Time    `json:"closed_at,omitempty"`
}

// NewAccount creates a new account with default values
func NewAccount(customerID string, accType AccountType, name, currencyCode string) *Account {
	now := time.Now()
	return &Account{
		ID:           uuid.New().String(),
		CustomerID:   customerID,
		Type:         accType,
		Status:       AccountStatusActive,
		Balance:      0,
		CurrencyCode: currencyCode,
		Name:         name,
		Number:       generateAccountNumber(),
		CreatedAt:    now,
		UpdatedAt:    now,
	}
}

// IsActive checks if the account is active
func (a *Account) IsActive() bool {
	return a.Status == AccountStatusActive
}

// Close closes the account
func (a *Account) Close() {
	now := time.Now()
	a.Status = AccountStatusClosed
	a.ClosedAt = &now
	a.UpdatedAt = now
}

// Deposit adds funds to the account
func (a *Account) Deposit(amount int64) error {
	if !a.IsActive() {
		return ErrAccountClosed
	}

	a.Balance += amount
	a.UpdatedAt = time.Now()
	return nil
}

// Withdraw removes funds from the account
func (a *Account) Withdraw(amount int64) error {
	if !a.IsActive() {
		return ErrAccountClosed
	}

	if a.Balance < amount {
		return ErrInsufficientFunds
	}

	a.Balance -= amount
	a.UpdatedAt = time.Now()
	return nil
}

// generateAccountNumber creates a random account number
// In a real system, this would follow specific bank rules
func generateAccountNumber() string {
	// Simple implementation - would be more complex in real system
	return uuid.New().String()[:8]
}

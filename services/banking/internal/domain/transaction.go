package domain

import (
	"time"

	"github.com/google/uuid"
)

// TransactionType represents the type of transaction
type TransactionType string

// TransactionStatus represents the status of a transaction
type TransactionStatus string

// Transaction types
const (
	TransactionTypeDeposit    TransactionType = "deposit"
	TransactionTypeWithdrawal TransactionType = "withdrawal"
	TransactionTypeTransfer   TransactionType = "transfer"
	TransactionTypeFee        TransactionType = "fee"
	TransactionTypeInterest   TransactionType = "interest"
)

// Transaction statuses
const (
	TransactionStatusPending   TransactionStatus = "pending"
	TransactionStatusCompleted TransactionStatus = "completed"
	TransactionStatusFailed    TransactionStatus = "failed"
	TransactionStatusReversed  TransactionStatus = "reversed"
)

// Transaction represents a financial transaction
type Transaction struct {
	ID              string            `json:"id"`
	Type            TransactionType   `json:"type"`
	Status          TransactionStatus `json:"status"`
	AccountID       string            `json:"account_id"`
	SourceAccountID *string           `json:"source_account_id,omitempty"`
	TargetAccountID *string           `json:"target_account_id,omitempty"`
	Amount          int64             `json:"amount"` // Stored in cents
	CurrencyCode    string            `json:"currency_code"`
	Description     string            `json:"description"`
	Reference       string            `json:"reference,omitempty"`
	Metadata        map[string]string `json:"metadata,omitempty"`
	CreatedAt       time.Time         `json:"created_at"`
	UpdatedAt       time.Time         `json:"updated_at"`
	CompletedAt     *time.Time        `json:"completed_at,omitempty"`
}

// NewTransaction creates a new transaction
func NewTransaction(txType TransactionType, accountID string, amount int64, currencyCode string, description string) *Transaction {
	now := time.Now()
	return &Transaction{
		ID:           uuid.New().String(),
		Type:         txType,
		Status:       TransactionStatusPending,
		AccountID:    accountID,
		Amount:       amount,
		CurrencyCode: currencyCode,
		Description:  description,
		Reference:    generateReference(),
		Metadata:     make(map[string]string),
		CreatedAt:    now,
		UpdatedAt:    now,
	}
}

// NewTransferTransaction creates a new transfer transaction
func NewTransferTransaction(sourceAccountID, targetAccountID string, amount int64, currencyCode string, description string) *Transaction {
	tx := NewTransaction(TransactionTypeTransfer, sourceAccountID, amount, currencyCode, description)
	tx.SourceAccountID = &sourceAccountID
	tx.TargetAccountID = &targetAccountID
	return tx
}

// Complete marks the transaction as completed
func (t *Transaction) Complete() {
	now := time.Now()
	t.Status = TransactionStatusCompleted
	t.CompletedAt = &now
	t.UpdatedAt = now
}

// Fail marks the transaction as failed
func (t *Transaction) Fail() {
	t.Status = TransactionStatusFailed
	t.UpdatedAt = time.Now()
}

// Reverse marks the transaction as reversed
func (t *Transaction) Reverse() {
	t.Status = TransactionStatusReversed
	t.UpdatedAt = time.Now()
}

// AddMetadata adds metadata to the transaction
func (t *Transaction) AddMetadata(key, value string) {
	if t.Metadata == nil {
		t.Metadata = make(map[string]string)
	}
	t.Metadata[key] = value
	t.UpdatedAt = time.Now()
}

// generateReference creates a unique reference number for the transaction
func generateReference() string {
	return "TX-" + uuid.New().String()[:8]
}

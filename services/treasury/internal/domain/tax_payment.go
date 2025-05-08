package domain

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

// Common errors
var (
	ErrTaxPaymentNotFound  = errors.New("tax payment not found")
	ErrInvalidPaymentAmount = errors.New("invalid payment amount")
)

// PaymentMethod represents the method of tax payment
type PaymentMethod string

// Payment methods
const (
	PaymentMethodElectronic PaymentMethod = "electronic"
	PaymentMethodCheck      PaymentMethod = "check"
	PaymentMethodCredit     PaymentMethod = "credit"
	PaymentMethodDebit      PaymentMethod = "debit"
	PaymentMethodCash       PaymentMethod = "cash"
	PaymentMethodWire       PaymentMethod = "wire"
)

// PaymentStatus represents the status of a tax payment
type PaymentStatus string

// Payment statuses
const (
	PaymentStatusPending   PaymentStatus = "pending"
	PaymentStatusCompleted PaymentStatus = "completed"
	PaymentStatusFailed    PaymentStatus = "failed"
	PaymentStatusRefunded  PaymentStatus = "refunded"
	PaymentStatusVoided    PaymentStatus = "voided"
)

// TaxPayment represents a payment for taxes
type TaxPayment struct {
	ID              string        `json:"id"`
	TaxpayerID      string        `json:"taxpayer_id"`
	FilingID        string        `json:"filing_id,omitempty"` // Optional, if payment is for a specific filing
	TaxType         TaxType       `json:"tax_type"`
	Amount          int64         `json:"amount"` // In cents
	PaymentMethod   PaymentMethod `json:"payment_method"`
	Status          PaymentStatus `json:"status"`
	PaymentDate     time.Time     `json:"payment_date"`
	ConfirmationCode string        `json:"confirmation_code"`
	Notes           string        `json:"notes,omitempty"`
	ProcessedAt     *time.Time    `json:"processed_at,omitempty"` // When the payment was processed
	RefundedAt      *time.Time    `json:"refunded_at,omitempty"`  // When the payment was refunded, if applicable
	CreatedAt       time.Time     `json:"created_at"`
	UpdatedAt       time.Time     `json:"updated_at"`
}

// NewTaxPayment creates a new tax payment
func NewTaxPayment(
	taxpayerID string,
	filingID string,
	taxType TaxType,
	amount int64,
	paymentMethod PaymentMethod,
	paymentDate time.Time,
	notes string,
) (*TaxPayment, error) {
	// Validate payment amount
	if amount <= 0 {
		return nil, ErrInvalidPaymentAmount
	}

	now := time.Now()
	return &TaxPayment{
		ID:               uuid.New().String(),
		TaxpayerID:       taxpayerID,
		FilingID:         filingID,
		TaxType:          taxType,
		Amount:           amount,
		PaymentMethod:    paymentMethod,
		Status:           PaymentStatusPending,
		PaymentDate:      paymentDate,
		ConfirmationCode: generateConfirmationCode(),
		Notes:            notes,
		CreatedAt:        now,
		UpdatedAt:        now,
	}, nil
}

// IsCompleted checks if the payment is completed
func (tp *TaxPayment) IsCompleted() bool {
	return tp.Status == PaymentStatusCompleted
}

// MarkAsCompleted marks the payment as completed
func (tp *TaxPayment) MarkAsCompleted() {
	now := time.Now()
	tp.Status = PaymentStatusCompleted
	tp.ProcessedAt = &now
	tp.UpdatedAt = now
}

// MarkAsFailed marks the payment as failed
func (tp *TaxPayment) MarkAsFailed(reason string) {
	tp.Status = PaymentStatusFailed
	if reason != "" {
		tp.Notes += "\nFailure reason: " + reason
	}
	tp.UpdatedAt = time.Now()
}

// Refund refunds the payment
func (tp *TaxPayment) Refund(reason string) error {
	if tp.Status != PaymentStatusCompleted {
		return errors.New("only completed payments can be refunded")
	}
	
	now := time.Now()
	tp.Status = PaymentStatusRefunded
	tp.RefundedAt = &now
	tp.Notes += "\nRefund reason: " + reason
	tp.UpdatedAt = now
	return nil
}

// Void voids the payment
func (tp *TaxPayment) Void(reason string) error {
	if tp.Status != PaymentStatusPending {
		return errors.New("only pending payments can be voided")
	}
	
	tp.Status = PaymentStatusVoided
	tp.Notes += "\nVoid reason: " + reason
	tp.UpdatedAt = time.Now()
	return nil
}

// UpdateAmount updates the payment amount
func (tp *TaxPayment) UpdateAmount(amount int64) error {
	if tp.Status != PaymentStatusPending {
		return errors.New("only pending payments can be updated")
	}
	
	if amount <= 0 {
		return ErrInvalidPaymentAmount
	}
	
	tp.Amount = amount
	tp.UpdatedAt = time.Now()
	return nil
}

// generateConfirmationCode creates a unique confirmation code for the payment
func generateConfirmationCode() string {
	return "PAY-" + uuid.New().String()[:8]
}
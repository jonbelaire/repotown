package domain

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

// Common errors
var (
	ErrTaxFilingNotFound   = errors.New("tax filing not found")
	ErrInvalidFilingStatus = errors.New("invalid filing status")
	ErrInvalidFilingPeriod = errors.New("invalid filing period")
)

// FilingStatus represents the status of a tax filing
type FilingStatus string

// Filing statuses
const (
	FilingStatusDraft     FilingStatus = "draft"
	FilingStatusSubmitted FilingStatus = "submitted"
	FilingStatusProcessing FilingStatus = "processing"
	FilingStatusAccepted  FilingStatus = "accepted"
	FilingStatusRejected  FilingStatus = "rejected"
	FilingStatusAmended   FilingStatus = "amended"
	FilingStatusAudited   FilingStatus = "audited"
)

// FilingPeriod represents the time period for a tax filing
type FilingPeriod string

// Filing periods
const (
	FilingPeriodMonthly    FilingPeriod = "monthly"
	FilingPeriodQuarterly  FilingPeriod = "quarterly"
	FilingPeriodSemiAnnual FilingPeriod = "semi_annual"
	FilingPeriodAnnual     FilingPeriod = "annual"
)

// TaxFiling represents a tax filing submission
type TaxFiling struct {
	ID             string       `json:"id"`
	TaxpayerID     string       `json:"taxpayer_id"`
	TaxYear        int          `json:"tax_year"`
	Period         FilingPeriod `json:"period"`
	PeriodStart    time.Time    `json:"period_start"`
	PeriodEnd      time.Time    `json:"period_end"`
	FilingType     TaxType      `json:"filing_type"` // Income, Sales, etc.
	Status         FilingStatus `json:"status"`
	GrossIncome    int64        `json:"gross_income,omitempty"`    // For income tax
	TaxableIncome  int64        `json:"taxable_income,omitempty"`  // For income tax
	TotalSales     int64        `json:"total_sales,omitempty"`     // For sales tax
	TaxableAmount  int64        `json:"taxable_amount"`            // Amount subject to tax
	TaxCalculated  int64        `json:"tax_calculated"`            // Tax owed based on calculation
	TaxPaid        int64        `json:"tax_paid"`                  // Amount actually paid
	SubmissionDate *time.Time   `json:"submission_date,omitempty"` // When the filing was submitted
	AcceptanceDate *time.Time   `json:"acceptance_date,omitempty"` // When the filing was accepted
	DueDate        time.Time    `json:"due_date"`                  // When the filing is due
	Deductions     []Deduction  `json:"deductions,omitempty"`      // List of tax deductions
	Credits        []Credit     `json:"credits,omitempty"`         // List of tax credits
	Notes          string       `json:"notes,omitempty"`
	CreatedAt      time.Time    `json:"created_at"`
	UpdatedAt      time.Time    `json:"updated_at"`
}

// Deduction represents a tax deduction in a filing
type Deduction struct {
	Code        string `json:"code"`
	Description string `json:"description"`
	Amount      int64  `json:"amount"`
}

// Credit represents a tax credit in a filing
type Credit struct {
	Code        string `json:"code"`
	Description string `json:"description"`
	Amount      int64  `json:"amount"`
}

// NewTaxFiling creates a new tax filing
func NewTaxFiling(
	taxpayerID string,
	taxYear int,
	period FilingPeriod,
	periodStart time.Time,
	periodEnd time.Time,
	filingType TaxType,
	dueDate time.Time,
) (*TaxFiling, error) {
	// Validate period
	if periodEnd.Before(periodStart) {
		return nil, ErrInvalidFilingPeriod
	}

	now := time.Now()
	return &TaxFiling{
		ID:            uuid.New().String(),
		TaxpayerID:    taxpayerID,
		TaxYear:       taxYear,
		Period:        period,
		PeriodStart:   periodStart,
		PeriodEnd:     periodEnd,
		FilingType:    filingType,
		Status:        FilingStatusDraft,
		Deductions:    []Deduction{},
		Credits:       []Credit{},
		DueDate:       dueDate,
		CreatedAt:     now,
		UpdatedAt:     now,
	}, nil
}

// IsOverdue checks if the filing is past its due date
func (tf *TaxFiling) IsOverdue() bool {
	return time.Now().After(tf.DueDate) && tf.Status != FilingStatusAccepted
}

// IsSubmitted checks if the filing has been submitted
func (tf *TaxFiling) IsSubmitted() bool {
	return tf.Status != FilingStatusDraft
}

// SubmitFiling submits the filing for processing
func (tf *TaxFiling) SubmitFiling() error {
	if tf.Status != FilingStatusDraft {
		return ErrInvalidFilingStatus
	}
	
	now := time.Now()
	tf.Status = FilingStatusSubmitted
	tf.SubmissionDate = &now
	tf.UpdatedAt = now
	return nil
}

// AcceptFiling marks the filing as accepted
func (tf *TaxFiling) AcceptFiling() error {
	if tf.Status != FilingStatusSubmitted && tf.Status != FilingStatusProcessing {
		return ErrInvalidFilingStatus
	}
	
	now := time.Now()
	tf.Status = FilingStatusAccepted
	tf.AcceptanceDate = &now
	tf.UpdatedAt = now
	return nil
}

// RejectFiling marks the filing as rejected
func (tf *TaxFiling) RejectFiling(reason string) error {
	if tf.Status != FilingStatusSubmitted && tf.Status != FilingStatusProcessing {
		return ErrInvalidFilingStatus
	}
	
	tf.Status = FilingStatusRejected
	tf.Notes += "\nRejection reason: " + reason
	tf.UpdatedAt = time.Now()
	return nil
}

// AmendFiling creates a new amended filing based on an existing one
func (tf *TaxFiling) AmendFiling() *TaxFiling {
	now := time.Time{}
	amended := *tf
	amended.ID = uuid.New().String()
	amended.Status = FilingStatusAmended
	amended.SubmissionDate = nil
	amended.AcceptanceDate = nil
	amended.Notes = "Amended from filing ID: " + tf.ID
	amended.CreatedAt = now
	amended.UpdatedAt = now
	return &amended
}

// AddDeduction adds a deduction to the filing
func (tf *TaxFiling) AddDeduction(code, description string, amount int64) {
	deduction := Deduction{
		Code:        code,
		Description: description,
		Amount:      amount,
	}
	tf.Deductions = append(tf.Deductions, deduction)
	tf.UpdatedAt = time.Now()
}

// AddCredit adds a credit to the filing
func (tf *TaxFiling) AddCredit(code, description string, amount int64) {
	credit := Credit{
		Code:        code,
		Description: description,
		Amount:      amount,
	}
	tf.Credits = append(tf.Credits, credit)
	tf.UpdatedAt = time.Now()
}

// CalculateTotalDeductions calculates the total deductions
func (tf *TaxFiling) CalculateTotalDeductions() int64 {
	total := int64(0)
	for _, deduction := range tf.Deductions {
		total += deduction.Amount
	}
	return total
}

// CalculateTotalCredits calculates the total credits
func (tf *TaxFiling) CalculateTotalCredits() int64 {
	total := int64(0)
	for _, credit := range tf.Credits {
		total += credit.Amount
	}
	return total
}
package domain

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

// Common errors
var (
	ErrTaxRateNotFound = errors.New("tax rate not found")
	ErrInvalidTaxRate  = errors.New("invalid tax rate")
)

// TaxType represents the type of tax
type TaxType string

// Tax types
const (
	TaxTypeIncome     TaxType = "income"
	TaxTypeSales      TaxType = "sales"
	TaxTypeProperty   TaxType = "property"
	TaxTypeBusiness   TaxType = "business"
	TaxTypeExcise     TaxType = "excise"
)

// TaxStatus represents the status of a tax rate
type TaxStatus string

// Tax statuses
const (
	TaxStatusActive    TaxStatus = "active"
	TaxStatusInactive  TaxStatus = "inactive"
	TaxStatusProposed  TaxStatus = "proposed"
	TaxStatusArchived  TaxStatus = "archived"
)

// TaxBracketType represents how the tax rate is applied
type TaxBracketType string

// Tax bracket types
const (
	TaxBracketFlat       TaxBracketType = "flat"      // Same percentage for all
	TaxBracketProgressive TaxBracketType = "progressive" // Different rates for different income levels
	TaxBracketTiered     TaxBracketType = "tiered"    // Different rates for different tiers
)

// TaxRate represents a tax rate for a specific type of tax
type TaxRate struct {
	ID                string         `json:"id"`
	Type              TaxType        `json:"type"`
	Name              string         `json:"name"`
	Description       string         `json:"description"`
	Rate              float64        `json:"rate"` // Stored as decimal (e.g., 0.07 for 7%)
	BracketType       TaxBracketType `json:"bracket_type"`
	MinAmount         int64          `json:"min_amount,omitempty"` // For progressive taxes, minimum amount this rate applies to
	MaxAmount         int64          `json:"max_amount,omitempty"` // For progressive taxes, maximum amount this rate applies to
	Status            TaxStatus      `json:"status"`
	Category          string         `json:"category,omitempty"` // Optional category for specific taxes (e.g., "Luxury" for sales tax)
	JurisdictionCode  string         `json:"jurisdiction_code"` // Geographic jurisdiction code
	EffectiveDate     time.Time      `json:"effective_date"`
	ExpirationDate    *time.Time     `json:"expiration_date,omitempty"`
	CreatedAt         time.Time      `json:"created_at"`
	UpdatedAt         time.Time      `json:"updated_at"`
}

// NewTaxRate creates a new tax rate
func NewTaxRate(
	taxType TaxType,
	name string,
	description string,
	rate float64,
	bracketType TaxBracketType,
	minAmount int64,
	maxAmount int64,
	category string,
	jurisdictionCode string,
	effectiveDate time.Time,
	expirationDate *time.Time,
) (*TaxRate, error) {
	// Validate tax rate
	if rate < 0 || rate > 1 {
		return nil, ErrInvalidTaxRate
	}

	now := time.Now()
	return &TaxRate{
		ID:               uuid.New().String(),
		Type:             taxType,
		Name:             name,
		Description:      description,
		Rate:             rate,
		BracketType:      bracketType,
		MinAmount:        minAmount,
		MaxAmount:        maxAmount,
		Status:           TaxStatusProposed, // New tax rates start as proposed
		Category:         category,
		JurisdictionCode: jurisdictionCode,
		EffectiveDate:    effectiveDate,
		ExpirationDate:   expirationDate,
		CreatedAt:        now,
		UpdatedAt:        now,
	}, nil
}

// IsActive checks if the tax rate is active
func (tr *TaxRate) IsActive() bool {
	return tr.Status == TaxStatusActive
}

// ActivateTaxRate sets the tax rate status to active
func (tr *TaxRate) ActivateTaxRate() {
	tr.Status = TaxStatusActive
	tr.UpdatedAt = time.Now()
}

// DeactivateTaxRate sets the tax rate status to inactive
func (tr *TaxRate) DeactivateTaxRate() {
	tr.Status = TaxStatusInactive
	tr.UpdatedAt = time.Now()
}

// ArchiveTaxRate sets the tax rate status to archived
func (tr *TaxRate) ArchiveTaxRate() {
	tr.Status = TaxStatusArchived
	tr.UpdatedAt = time.Now()
}

// UpdateTaxRate updates the tax rate details
func (tr *TaxRate) UpdateTaxRate(
	name string,
	description string,
	rate float64,
	category string,
	effectiveDate time.Time,
	expirationDate *time.Time,
) error {
	// Validate tax rate
	if rate < 0 || rate > 1 {
		return ErrInvalidTaxRate
	}

	tr.Name = name
	tr.Description = description
	tr.Rate = rate
	tr.Category = category
	tr.EffectiveDate = effectiveDate
	tr.ExpirationDate = expirationDate
	tr.UpdatedAt = time.Now()
	return nil
}

// IsApplicable checks if the tax rate is applicable to a given amount
func (tr *TaxRate) IsApplicable(amount int64) bool {
	if tr.BracketType == TaxBracketFlat {
		return true
	}
	
	if tr.MinAmount > 0 && amount < tr.MinAmount {
		return false
	}
	
	if tr.MaxAmount > 0 && amount > tr.MaxAmount {
		return false
	}
	
	return true
}

// CalculateTax calculates the tax amount for a given amount
func (tr *TaxRate) CalculateTax(amount int64) int64 {
	if !tr.IsApplicable(amount) {
		return 0
	}
	
	// For progressive, we only tax the amount within the bracket
	if tr.BracketType == TaxBracketProgressive {
		taxableAmount := amount
		
		if tr.MinAmount > 0 {
			taxableAmount = amount - tr.MinAmount
			if taxableAmount < 0 {
				taxableAmount = 0
			}
		}
		
		if tr.MaxAmount > 0 && amount > tr.MaxAmount {
			taxableAmount = tr.MaxAmount - tr.MinAmount
		}
		
		return int64(float64(taxableAmount) * tr.Rate)
	}
	
	// For flat and tiered rates
	return int64(float64(amount) * tr.Rate)
}
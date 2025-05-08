package domain

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

// Common errors
var (
	ErrTaxpayerNotFound = errors.New("taxpayer not found")
	ErrTaxpayerExists   = errors.New("taxpayer with identifier already exists")
)

// TaxpayerType represents the type of taxpayer
type TaxpayerType string

// Taxpayer types
const (
	TaxpayerTypeIndividual TaxpayerType = "individual"
	TaxpayerTypeBusiness   TaxpayerType = "business"
	TaxpayerTypeNonProfit  TaxpayerType = "non_profit"
	TaxpayerTypeGovernment TaxpayerType = "government"
)

// TaxpayerStatus represents the status of a taxpayer
type TaxpayerStatus string

// Taxpayer statuses
const (
	TaxpayerStatusActive   TaxpayerStatus = "active"
	TaxpayerStatusInactive TaxpayerStatus = "inactive"
	TaxpayerStatusExempt   TaxpayerStatus = "exempt"
	TaxpayerStatusDelinquent TaxpayerStatus = "delinquent"
)

// Taxpayer represents an entity that pays taxes
type Taxpayer struct {
	ID             string         `json:"id"`
	Type           TaxpayerType   `json:"type"`
	Status         TaxpayerStatus `json:"status"`
	Name           string         `json:"name"`
	TaxIdentifier  string         `json:"tax_identifier"` // SSN, EIN, etc.
	ContactEmail   string         `json:"contact_email"`
	ContactPhone   string         `json:"contact_phone"`
	Address        Address        `json:"address"`
	ExemptionCodes []string       `json:"exemption_codes,omitempty"` // Any tax exemption codes
	AnnualRevenue  int64          `json:"annual_revenue,omitempty"`  // For businesses
	BusinessType   string         `json:"business_type,omitempty"`   // For businesses
	Industry       string         `json:"industry,omitempty"`        // For businesses
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
}

// Address represents a physical address
type Address struct {
	Street     string `json:"street"`
	City       string `json:"city"`
	State      string `json:"state"`
	Country    string `json:"country"`
	PostalCode string `json:"postal_code"`
}

// NewTaxpayer creates a new taxpayer
func NewTaxpayer(
	taxpayerType TaxpayerType,
	name string,
	taxIdentifier string,
	contactEmail string,
	contactPhone string,
	address Address,
	exemptionCodes []string,
	annualRevenue int64,
	businessType string,
	industry string,
) *Taxpayer {
	now := time.Now()
	return &Taxpayer{
		ID:             uuid.New().String(),
		Type:           taxpayerType,
		Status:         TaxpayerStatusActive,
		Name:           name,
		TaxIdentifier:  taxIdentifier,
		ContactEmail:   contactEmail,
		ContactPhone:   contactPhone,
		Address:        address,
		ExemptionCodes: exemptionCodes,
		AnnualRevenue:  annualRevenue,
		BusinessType:   businessType,
		Industry:       industry,
		CreatedAt:      now,
		UpdatedAt:      now,
	}
}

// IsActive checks if the taxpayer is active
func (tp *Taxpayer) IsActive() bool {
	return tp.Status == TaxpayerStatusActive
}

// IsExempt checks if the taxpayer is exempt from taxes
func (tp *Taxpayer) IsExempt() bool {
	return tp.Status == TaxpayerStatusExempt
}

// HasExemption checks if the taxpayer has a specific exemption
func (tp *Taxpayer) HasExemption(code string) bool {
	for _, c := range tp.ExemptionCodes {
		if c == code {
			return true
		}
	}
	return false
}

// UpdateStatus updates the taxpayer status
func (tp *Taxpayer) UpdateStatus(status TaxpayerStatus) {
	tp.Status = status
	tp.UpdatedAt = time.Now()
}

// UpdateContact updates the taxpayer's contact information
func (tp *Taxpayer) UpdateContact(email, phone string) {
	tp.ContactEmail = email
	tp.ContactPhone = phone
	tp.UpdatedAt = time.Now()
}

// UpdateAddress updates the taxpayer's address
func (tp *Taxpayer) UpdateAddress(address Address) {
	tp.Address = address
	tp.UpdatedAt = time.Now()
}

// UpdateBusinessInfo updates business-specific information
func (tp *Taxpayer) UpdateBusinessInfo(annualRevenue int64, businessType, industry string) {
	tp.AnnualRevenue = annualRevenue
	tp.BusinessType = businessType
	tp.Industry = industry
	tp.UpdatedAt = time.Now()
}

// AddExemption adds a tax exemption code
func (tp *Taxpayer) AddExemption(code string) {
	// Check if exemption already exists
	for _, c := range tp.ExemptionCodes {
		if c == code {
			return // Already has this exemption
		}
	}
	
	tp.ExemptionCodes = append(tp.ExemptionCodes, code)
	tp.UpdatedAt = time.Now()
}

// RemoveExemption removes a tax exemption code
func (tp *Taxpayer) RemoveExemption(code string) {
	for i, c := range tp.ExemptionCodes {
		if c == code {
			// Remove the exemption by replacing it with the last element and truncating the slice
			tp.ExemptionCodes[i] = tp.ExemptionCodes[len(tp.ExemptionCodes)-1]
			tp.ExemptionCodes = tp.ExemptionCodes[:len(tp.ExemptionCodes)-1]
			tp.UpdatedAt = time.Now()
			return
		}
	}
}
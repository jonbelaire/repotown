package domain

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

// Common errors
var (
	ErrGarageNotFound    = errors.New("garage not found")
	ErrGarageCapacityFull = errors.New("garage capacity is full")
)

// GarageStatus represents the status of a garage
type GarageStatus string

// Garage statuses
const (
	GarageStatusOperational GarageStatus = "operational"
	GarageStatusMaintenance GarageStatus = "maintenance"
	GarageStatusClosed      GarageStatus = "closed"
)

// Garage represents a parking garage
type Garage struct {
	ID               string       `json:"id"`
	Name             string       `json:"name"`
	Address          string       `json:"address"`
	Status           GarageStatus `json:"status"`
	TotalSpaces      int          `json:"total_spaces"`
	AvailableSpaces  int          `json:"available_spaces"`
	HourlyRate       int64        `json:"hourly_rate"` // Stored in cents
	DailyRate        int64        `json:"daily_rate"`  // Stored in cents
	OperatingHours   string       `json:"operating_hours"`
	HasElectricCharging bool      `json:"has_electric_charging"`
	CreatedAt        time.Time    `json:"created_at"`
	UpdatedAt        time.Time    `json:"updated_at"`
}

// NewGarage creates a new garage with default values
func NewGarage(name, address string, totalSpaces int, hourlyRate, dailyRate int64, operatingHours string, hasElectricCharging bool) *Garage {
	now := time.Now()
	return &Garage{
		ID:                 uuid.New().String(),
		Name:               name,
		Address:            address,
		Status:             GarageStatusOperational,
		TotalSpaces:        totalSpaces,
		AvailableSpaces:    totalSpaces, // Initially all spaces are available
		HourlyRate:         hourlyRate,
		DailyRate:          dailyRate,
		OperatingHours:     operatingHours,
		HasElectricCharging: hasElectricCharging,
		CreatedAt:          now,
		UpdatedAt:          now,
	}
}

// IsOperational checks if the garage is operational
func (g *Garage) IsOperational() bool {
	return g.Status == GarageStatusOperational
}

// HasAvailableSpace checks if the garage has available parking spaces
func (g *Garage) HasAvailableSpace() bool {
	return g.AvailableSpaces > 0
}

// UpdateStatus updates the garage status
func (g *Garage) UpdateStatus(status GarageStatus) {
	g.Status = status
	g.UpdatedAt = time.Now()
}

// UpdateRates updates the garage pricing rates
func (g *Garage) UpdateRates(hourlyRate, dailyRate int64) {
	g.HourlyRate = hourlyRate
	g.DailyRate = dailyRate
	g.UpdatedAt = time.Now()
}

// UpdateOperatingHours updates the garage operating hours
func (g *Garage) UpdateOperatingHours(operatingHours string) {
	g.OperatingHours = operatingHours
	g.UpdatedAt = time.Now()
}

// DecrementAvailableSpaces decreases the available spaces count
func (g *Garage) DecrementAvailableSpaces() error {
	if g.AvailableSpaces <= 0 {
		return ErrGarageCapacityFull
	}
	
	g.AvailableSpaces--
	g.UpdatedAt = time.Now()
	return nil
}

// IncrementAvailableSpaces increases the available spaces count
func (g *Garage) IncrementAvailableSpaces() {
	if g.AvailableSpaces < g.TotalSpaces {
		g.AvailableSpaces++
		g.UpdatedAt = time.Now()
	}
}
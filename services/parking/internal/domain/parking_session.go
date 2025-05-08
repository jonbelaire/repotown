package domain

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

// Common errors
var (
	ErrSessionNotFound     = errors.New("parking session not found")
	ErrSessionAlreadyEnded = errors.New("parking session already ended")
)

// ParkingSessionStatus represents the status of a parking session
type ParkingSessionStatus string

// Parking session statuses
const (
	SessionStatusActive   ParkingSessionStatus = "active"
	SessionStatusCompleted ParkingSessionStatus = "completed"
	SessionStatusCancelled ParkingSessionStatus = "cancelled"
)

// PaymentStatus represents the payment status of a parking session
type PaymentStatus string

// Payment statuses
const (
	PaymentStatusPending   PaymentStatus = "pending"
	PaymentStatusPaid      PaymentStatus = "paid"
	PaymentStatusFailed    PaymentStatus = "failed"
	PaymentStatusRefunded  PaymentStatus = "refunded"
)

// ParkingSession represents a vehicle parking session
type ParkingSession struct {
	ID              string               `json:"id"`
	GarageID        string               `json:"garage_id"`
	VehicleID       string               `json:"vehicle_id"`
	SpotNumber      string               `json:"spot_number,omitempty"`
	Status          ParkingSessionStatus `json:"status"`
	PaymentStatus   PaymentStatus        `json:"payment_status"`
	StartTime       time.Time            `json:"start_time"`
	EndTime         *time.Time           `json:"end_time,omitempty"`
	Duration        *time.Duration       `json:"duration,omitempty"`
	AmountCharged   int64                `json:"amount_charged,omitempty"` // Stored in cents
	IsPrepaid       bool                 `json:"is_prepaid"`
	Notes           string               `json:"notes,omitempty"`
	CreatedAt       time.Time            `json:"created_at"`
	UpdatedAt       time.Time            `json:"updated_at"`
}

// NewParkingSession creates a new parking session
func NewParkingSession(garageID, vehicleID, spotNumber string, isPrepaid bool) *ParkingSession {
	now := time.Now()
	return &ParkingSession{
		ID:            uuid.New().String(),
		GarageID:      garageID,
		VehicleID:     vehicleID,
		SpotNumber:    spotNumber,
		Status:        SessionStatusActive,
		PaymentStatus: isPrepaid ? PaymentStatusPaid : PaymentStatusPending,
		StartTime:     now,
		IsPrepaid:     isPrepaid,
		CreatedAt:     now,
		UpdatedAt:     now,
	}
}

// EndSession completes the parking session and calculates duration and amount
func (s *ParkingSession) EndSession(garage *Garage, endTime time.Time) error {
	if s.Status != SessionStatusActive {
		return ErrSessionAlreadyEnded
	}

	s.EndTime = &endTime
	
	// Calculate duration
	duration := endTime.Sub(s.StartTime)
	s.Duration = &duration
	
	// Calculate amount charged - simple implementation
	hours := int64(duration.Hours())
	if hours < 1 {
		hours = 1 // Minimum one hour charge
	}
	
	// If more than 8 hours, use daily rate
	if hours > 8 {
		days := (hours / 24) + 1
		s.AmountCharged = days * garage.DailyRate
	} else {
		s.AmountCharged = hours * garage.HourlyRate
	}
	
	s.Status = SessionStatusCompleted
	if !s.IsPrepaid {
		s.PaymentStatus = PaymentStatusPending
	}
	
	s.UpdatedAt = endTime
	return nil
}

// MarkAsPaid marks the session as paid
func (s *ParkingSession) MarkAsPaid() {
	s.PaymentStatus = PaymentStatusPaid
	s.UpdatedAt = time.Time{}
}

// CancelSession cancels the parking session
func (s *ParkingSession) CancelSession() error {
	if s.Status != SessionStatusActive {
		return ErrSessionAlreadyEnded
	}
	
	s.Status = SessionStatusCancelled
	now := time.Now()
	s.EndTime = &now
	
	// Calculate duration up to cancellation
	duration := now.Sub(s.StartTime)
	s.Duration = &duration
	
	s.UpdatedAt = now
	return nil
}

// UpdateSpotNumber updates the parking spot number
func (s *ParkingSession) UpdateSpotNumber(spotNumber string) {
	s.SpotNumber = spotNumber
	s.UpdatedAt = time.Now()
}

// AddNotes adds notes to the parking session
func (s *ParkingSession) AddNotes(notes string) {
	if s.Notes != "" {
		s.Notes = s.Notes + "; " + notes
	} else {
		s.Notes = notes
	}
	s.UpdatedAt = time.Now()
}
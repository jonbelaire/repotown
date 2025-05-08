package domain

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

// Common errors
var (
	ErrReservationNotFound  = errors.New("reservation not found")
	ErrReservationCancelled = errors.New("reservation already cancelled")
	ErrReservationExpired   = errors.New("reservation has expired")
)

// ReservationStatus represents the status of a reservation
type ReservationStatus string

// Reservation statuses
const (
	ReservationStatusActive    ReservationStatus = "active"
	ReservationStatusUsed      ReservationStatus = "used"
	ReservationStatusCancelled ReservationStatus = "cancelled"
	ReservationStatusExpired   ReservationStatus = "expired"
)

// Reservation represents a parking spot reservation
type Reservation struct {
	ID            string            `json:"id"`
	GarageID      string            `json:"garage_id"`
	VehicleID     string            `json:"vehicle_id"`
	Status        ReservationStatus `json:"status"`
	StartTime     time.Time         `json:"start_time"`
	EndTime       time.Time         `json:"end_time"`
	AmountPaid    int64             `json:"amount_paid"` // Stored in cents
	ConfirmationCode string         `json:"confirmation_code"`
	CreatedAt     time.Time         `json:"created_at"`
	UpdatedAt     time.Time         `json:"updated_at"`
	UsedAt        *time.Time        `json:"used_at,omitempty"`
	CancelledAt   *time.Time        `json:"cancelled_at,omitempty"`
}

// NewReservation creates a new parking reservation
func NewReservation(garageID, vehicleID string, startTime, endTime time.Time, amountPaid int64) *Reservation {
	now := time.Now()
	return &Reservation{
		ID:              uuid.New().String(),
		GarageID:        garageID,
		VehicleID:       vehicleID,
		Status:          ReservationStatusActive,
		StartTime:       startTime,
		EndTime:         endTime,
		AmountPaid:      amountPaid,
		ConfirmationCode: generateConfirmationCode(),
		CreatedAt:       now,
		UpdatedAt:       now,
	}
}

// IsActive checks if the reservation is active
func (r *Reservation) IsActive() bool {
	return r.Status == ReservationStatusActive
}

// IsExpired checks if the reservation has expired
func (r *Reservation) IsExpired() bool {
	return time.Now().After(r.EndTime) && r.Status == ReservationStatusActive
}

// MarkAsUsed marks the reservation as used
func (r *Reservation) MarkAsUsed() error {
	if r.Status != ReservationStatusActive {
		return ErrReservationCancelled
	}
	
	now := time.Now()
	r.Status = ReservationStatusUsed
	r.UsedAt = &now
	r.UpdatedAt = now
	return nil
}

// CancelReservation cancels the reservation
func (r *Reservation) CancelReservation() error {
	if r.Status != ReservationStatusActive {
		return ErrReservationCancelled
	}
	
	now := time.Now()
	r.Status = ReservationStatusCancelled
	r.CancelledAt = &now
	r.UpdatedAt = now
	return nil
}

// MarkAsExpired marks the reservation as expired
func (r *Reservation) MarkAsExpired() error {
	if r.Status != ReservationStatusActive {
		return ErrReservationCancelled
	}
	
	if !r.IsExpired() {
		return ErrReservationExpired
	}
	
	now := time.Now()
	r.Status = ReservationStatusExpired
	r.UpdatedAt = now
	return nil
}

// generateConfirmationCode creates a unique confirmation code for the reservation
func generateConfirmationCode() string {
	return "RSV-" + uuid.New().String()[:8]
}
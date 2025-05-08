package service

import (
	"context"
	"time"

	"github.com/jonbelaire/repotown/packages/go-core/logging"
	"github.com/jonbelaire/repotown/services/parking/internal/domain"
	"github.com/jonbelaire/repotown/services/parking/internal/repository"
)

// ReservationService provides business logic for reservation management
type ReservationService interface {
	GetReservation(ctx context.Context, id string) (*domain.Reservation, error)
	GetReservationByConfirmationCode(ctx context.Context, code string) (*domain.Reservation, error)
	ListReservations(ctx context.Context, limit, offset int) ([]*domain.Reservation, error)
	ListActiveReservations(ctx context.Context) ([]*domain.Reservation, error)
	ListReservationsByGarage(ctx context.Context, garageID string) ([]*domain.Reservation, error)
	ListReservationsByVehicle(ctx context.Context, vehicleID string) ([]*domain.Reservation, error)
	CreateReservation(ctx context.Context, garageID, vehicleID string, startTime, endTime time.Time, amountPaid int64) (*domain.Reservation, error)
	UseReservation(ctx context.Context, id string) (*domain.Reservation, error)
	CancelReservation(ctx context.Context, id string) (*domain.Reservation, error)
	CheckExpiredReservations(ctx context.Context) (int, error)
}

// reservationService implements ReservationService
type reservationService struct {
	reservationRepo repository.ReservationRepository
	garageRepo      repository.GarageRepository
	logger          logging.Logger
}

// NewReservationService creates a new reservation service
func NewReservationService(
	reservationRepo repository.ReservationRepository,
	garageRepo repository.GarageRepository,
	logger logging.Logger,
) ReservationService {
	return &reservationService{
		reservationRepo: reservationRepo,
		garageRepo:      garageRepo,
		logger:          logger,
	}
}

// GetReservation retrieves a reservation by ID
func (s *reservationService) GetReservation(ctx context.Context, id string) (*domain.Reservation, error) {
	return s.reservationRepo.GetByID(ctx, id)
}

// GetReservationByConfirmationCode retrieves a reservation by confirmation code
func (s *reservationService) GetReservationByConfirmationCode(ctx context.Context, code string) (*domain.Reservation, error) {
	return s.reservationRepo.GetByConfirmationCode(ctx, code)
}

// ListReservations retrieves reservations with pagination
func (s *reservationService) ListReservations(ctx context.Context, limit, offset int) ([]*domain.Reservation, error) {
	return s.reservationRepo.List(ctx, limit, offset)
}

// ListActiveReservations retrieves active reservations
func (s *reservationService) ListActiveReservations(ctx context.Context) ([]*domain.Reservation, error) {
	return s.reservationRepo.ListActive(ctx)
}

// ListReservationsByGarage retrieves reservations for a specific garage
func (s *reservationService) ListReservationsByGarage(ctx context.Context, garageID string) ([]*domain.Reservation, error) {
	return s.reservationRepo.ListByGarage(ctx, garageID)
}

// ListReservationsByVehicle retrieves reservations for a specific vehicle
func (s *reservationService) ListReservationsByVehicle(ctx context.Context, vehicleID string) ([]*domain.Reservation, error) {
	return s.reservationRepo.ListByVehicle(ctx, vehicleID)
}

// CreateReservation creates a new reservation
func (s *reservationService) CreateReservation(ctx context.Context, garageID, vehicleID string, startTime, endTime time.Time, amountPaid int64) (*domain.Reservation, error) {
	// Verify garage exists
	garage, err := s.garageRepo.GetByID(ctx, garageID)
	if err != nil {
		return nil, err
	}

	if !garage.IsOperational() {
		return nil, domain.ErrGarageNotOperational
	}

	// Check if the garage has available space
	if !garage.HasAvailableSpace() {
		return nil, domain.ErrGarageCapacityFull
	}

	// Create reservation
	reservation := domain.NewReservation(garageID, vehicleID, startTime, endTime, amountPaid)
	if err := s.reservationRepo.Create(ctx, reservation); err != nil {
		return nil, err
	}

	return reservation, nil
}

// UseReservation marks a reservation as used
func (s *reservationService) UseReservation(ctx context.Context, id string) (*domain.Reservation, error) {
	// Get reservation
	reservation, err := s.reservationRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Mark as used
	if err := reservation.MarkAsUsed(); err != nil {
		return nil, err
	}

	// Update reservation
	if err := s.reservationRepo.Update(ctx, reservation); err != nil {
		return nil, err
	}

	return reservation, nil
}

// CancelReservation cancels a reservation
func (s *reservationService) CancelReservation(ctx context.Context, id string) (*domain.Reservation, error) {
	// Get reservation
	reservation, err := s.reservationRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Cancel reservation
	if err := reservation.CancelReservation(); err != nil {
		return nil, err
	}

	// Update reservation
	if err := s.reservationRepo.Update(ctx, reservation); err != nil {
		return nil, err
	}

	return reservation, nil
}

// CheckExpiredReservations checks for and marks expired reservations
func (s *reservationService) CheckExpiredReservations(ctx context.Context) (int, error) {
	// Get active reservations
	activeReservations, err := s.reservationRepo.ListActive(ctx)
	if err != nil {
		return 0, err
	}

	expiredCount := 0
	for _, reservation := range activeReservations {
		if reservation.IsExpired() {
			if err := reservation.MarkAsExpired(); err != nil {
				s.logger.Error("Failed to mark reservation as expired", "reservation_id", reservation.ID, "error", err)
				continue
			}

			if err := s.reservationRepo.Update(ctx, reservation); err != nil {
				s.logger.Error("Failed to update expired reservation", "reservation_id", reservation.ID, "error", err)
				continue
			}

			expiredCount++
		}
	}

	return expiredCount, nil
}
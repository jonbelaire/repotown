package service

import (
	"context"
	"time"

	"github.com/jonbelaire/repotown/packages/go-core/logging"
	"github.com/jonbelaire/repotown/services/parking/internal/domain"
	"github.com/jonbelaire/repotown/services/parking/internal/repository"
)

// ParkingSessionService provides business logic for parking session management
type ParkingSessionService interface {
	GetSession(ctx context.Context, id string) (*domain.ParkingSession, error)
	ListSessions(ctx context.Context, limit, offset int) ([]*domain.ParkingSession, error)
	ListActiveSessions(ctx context.Context) ([]*domain.ParkingSession, error)
	ListSessionsByGarage(ctx context.Context, garageID string) ([]*domain.ParkingSession, error)
	ListSessionsByVehicle(ctx context.Context, vehicleID string) ([]*domain.ParkingSession, error)
	StartSession(ctx context.Context, garageID, vehicleID, spotNumber string, isPrepaid bool) (*domain.ParkingSession, error)
	EndSession(ctx context.Context, id string) (*domain.ParkingSession, error)
	EndSessionWithCustomTime(ctx context.Context, id string, endTime time.Time) (*domain.ParkingSession, error)
	CancelSession(ctx context.Context, id string) (*domain.ParkingSession, error)
	UpdateSessionSpot(ctx context.Context, id, spotNumber string) (*domain.ParkingSession, error)
	AddSessionNotes(ctx context.Context, id, notes string) (*domain.ParkingSession, error)
	MarkSessionAsPaid(ctx context.Context, id string) (*domain.ParkingSession, error)
}

// parkingSessionService implements ParkingSessionService
type parkingSessionService struct {
	sessionRepo repository.ParkingSessionRepository
	garageRepo  repository.GarageRepository
	vehicleRepo repository.VehicleRepository
	logger      logging.Logger
}

// NewParkingSessionService creates a new parking session service
func NewParkingSessionService(
	sessionRepo repository.ParkingSessionRepository,
	garageRepo repository.GarageRepository,
	vehicleRepo repository.VehicleRepository,
	logger logging.Logger,
) ParkingSessionService {
	return &parkingSessionService{
		sessionRepo: sessionRepo,
		garageRepo:  garageRepo,
		vehicleRepo: vehicleRepo,
		logger:      logger,
	}
}

// GetSession retrieves a parking session by ID
func (s *parkingSessionService) GetSession(ctx context.Context, id string) (*domain.ParkingSession, error) {
	return s.sessionRepo.GetByID(ctx, id)
}

// ListSessions retrieves parking sessions with pagination
func (s *parkingSessionService) ListSessions(ctx context.Context, limit, offset int) ([]*domain.ParkingSession, error) {
	return s.sessionRepo.List(ctx, limit, offset)
}

// ListActiveSessions retrieves active parking sessions
func (s *parkingSessionService) ListActiveSessions(ctx context.Context) ([]*domain.ParkingSession, error) {
	return s.sessionRepo.ListActive(ctx)
}

// ListSessionsByGarage retrieves parking sessions for a specific garage
func (s *parkingSessionService) ListSessionsByGarage(ctx context.Context, garageID string) ([]*domain.ParkingSession, error) {
	return s.sessionRepo.ListByGarage(ctx, garageID)
}

// ListSessionsByVehicle retrieves parking sessions for a specific vehicle
func (s *parkingSessionService) ListSessionsByVehicle(ctx context.Context, vehicleID string) ([]*domain.ParkingSession, error) {
	return s.sessionRepo.ListByVehicle(ctx, vehicleID)
}

// StartSession starts a new parking session
func (s *parkingSessionService) StartSession(ctx context.Context, garageID, vehicleID, spotNumber string, isPrepaid bool) (*domain.ParkingSession, error) {
	// Verify garage exists and has available space
	garage, err := s.garageRepo.GetByID(ctx, garageID)
	if err != nil {
		return nil, err
	}

	if !garage.IsOperational() {
		return nil, domain.ErrGarageNotOperational
	}

	if err := garage.DecrementAvailableSpaces(); err != nil {
		return nil, err
	}

	// Update garage available spaces
	if err := s.garageRepo.Update(ctx, garage); err != nil {
		return nil, err
	}

	// Verify vehicle exists
	if _, err := s.vehicleRepo.GetByID(ctx, vehicleID); err != nil {
		// Restore garage space on error
		_ = s.restoreGarageSpace(ctx, garageID)
		return nil, err
	}

	// Create parking session
	session := domain.NewParkingSession(garageID, vehicleID, spotNumber, isPrepaid)
	if err := s.sessionRepo.Create(ctx, session); err != nil {
		// Restore garage space on error
		_ = s.restoreGarageSpace(ctx, garageID)
		return nil, err
	}

	return session, nil
}

// EndSession ends a parking session
func (s *parkingSessionService) EndSession(ctx context.Context, id string) (*domain.ParkingSession, error) {
	return s.EndSessionWithCustomTime(ctx, id, time.Now())
}

// EndSessionWithCustomTime ends a parking session with a custom end time
func (s *parkingSessionService) EndSessionWithCustomTime(ctx context.Context, id string, endTime time.Time) (*domain.ParkingSession, error) {
	// Get session
	session, err := s.sessionRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Get garage for pricing info
	garage, err := s.garageRepo.GetByID(ctx, session.GarageID)
	if err != nil {
		return nil, err
	}

	// End session
	if err := session.EndSession(garage, endTime); err != nil {
		return nil, err
	}

	// Update session
	if err := s.sessionRepo.Update(ctx, session); err != nil {
		return nil, err
	}

	// Restore garage space
	if err := s.restoreGarageSpace(ctx, session.GarageID); err != nil {
		s.logger.Error("Failed to restore garage space", "error", err)
	}

	return session, nil
}

// CancelSession cancels a parking session
func (s *parkingSessionService) CancelSession(ctx context.Context, id string) (*domain.ParkingSession, error) {
	// Get session
	session, err := s.sessionRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Cancel session
	if err := session.CancelSession(); err != nil {
		return nil, err
	}

	// Update session
	if err := s.sessionRepo.Update(ctx, session); err != nil {
		return nil, err
	}

	// Restore garage space
	if err := s.restoreGarageSpace(ctx, session.GarageID); err != nil {
		s.logger.Error("Failed to restore garage space", "error", err)
	}

	return session, nil
}

// UpdateSessionSpot updates the parking spot number for a session
func (s *parkingSessionService) UpdateSessionSpot(ctx context.Context, id, spotNumber string) (*domain.ParkingSession, error) {
	// Get session
	session, err := s.sessionRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Update spot number
	session.UpdateSpotNumber(spotNumber)

	// Update session
	if err := s.sessionRepo.Update(ctx, session); err != nil {
		return nil, err
	}

	return session, nil
}

// AddSessionNotes adds notes to a parking session
func (s *parkingSessionService) AddSessionNotes(ctx context.Context, id, notes string) (*domain.ParkingSession, error) {
	// Get session
	session, err := s.sessionRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Add notes
	session.AddNotes(notes)

	// Update session
	if err := s.sessionRepo.Update(ctx, session); err != nil {
		return nil, err
	}

	return session, nil
}

// MarkSessionAsPaid marks a session as paid
func (s *parkingSessionService) MarkSessionAsPaid(ctx context.Context, id string) (*domain.ParkingSession, error) {
	// Get session
	session, err := s.sessionRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Mark as paid
	session.MarkAsPaid()

	// Update session
	if err := s.sessionRepo.Update(ctx, session); err != nil {
		return nil, err
	}

	return session, nil
}

// Helper method to restore garage space
func (s *parkingSessionService) restoreGarageSpace(ctx context.Context, garageID string) error {
	garage, err := s.garageRepo.GetByID(ctx, garageID)
	if err != nil {
		return err
	}

	garage.IncrementAvailableSpaces()
	return s.garageRepo.Update(ctx, garage)
}
package service

import (
	"context"
	"time"

	"github.com/jonbelaire/repotown/packages/go-core/logging"
	"github.com/jonbelaire/repotown/services/parking/internal/domain"
	"github.com/jonbelaire/repotown/services/parking/internal/repository"
)

// GarageService provides business logic for garage management
type GarageService interface {
	GetGarage(ctx context.Context, id string) (*domain.Garage, error)
	ListGarages(ctx context.Context, limit, offset int) ([]*domain.Garage, error)
	ListAvailableGarages(ctx context.Context) ([]*domain.Garage, error)
	CreateGarage(ctx context.Context, name, address string, totalSpaces int, hourlyRate, dailyRate int64, operatingHours string, hasElectricCharging bool) (*domain.Garage, error)
	UpdateGarage(ctx context.Context, id, name, address, operatingHours string, hourlyRate, dailyRate int64) (*domain.Garage, error)
	UpdateGarageStatus(ctx context.Context, id string, status domain.GarageStatus) (*domain.Garage, error)
	UpdateAvailableSpaces(ctx context.Context, id string, increment bool) (*domain.Garage, error)
}

// garageService implements GarageService
type garageService struct {
	garageRepo repository.GarageRepository
	logger     logging.Logger
}

// NewGarageService creates a new garage service
func NewGarageService(garageRepo repository.GarageRepository, logger logging.Logger) GarageService {
	return &garageService{
		garageRepo: garageRepo,
		logger:     logger,
	}
}

// GetGarage retrieves a garage by ID
func (s *garageService) GetGarage(ctx context.Context, id string) (*domain.Garage, error) {
	return s.garageRepo.GetByID(ctx, id)
}

// ListGarages retrieves garages with pagination
func (s *garageService) ListGarages(ctx context.Context, limit, offset int) ([]*domain.Garage, error) {
	return s.garageRepo.List(ctx, limit, offset)
}

// ListAvailableGarages retrieves garages with available spaces
func (s *garageService) ListAvailableGarages(ctx context.Context) ([]*domain.Garage, error) {
	return s.garageRepo.ListAvailable(ctx)
}

// CreateGarage creates a new garage
func (s *garageService) CreateGarage(ctx context.Context, name, address string, totalSpaces int, hourlyRate, dailyRate int64, operatingHours string, hasElectricCharging bool) (*domain.Garage, error) {
	garage := domain.NewGarage(name, address, totalSpaces, hourlyRate, dailyRate, operatingHours, hasElectricCharging)
	if err := s.garageRepo.Create(ctx, garage); err != nil {
		return nil, err
	}
	return garage, nil
}

// UpdateGarage updates a garage's details
func (s *garageService) UpdateGarage(ctx context.Context, id, name, address, operatingHours string, hourlyRate, dailyRate int64) (*domain.Garage, error) {
	garage, err := s.garageRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	garage.Name = name
	garage.Address = address
	garage.OperatingHours = operatingHours
	garage.HourlyRate = hourlyRate
	garage.DailyRate = dailyRate
	garage.UpdatedAt = time.Now()

	if err := s.garageRepo.Update(ctx, garage); err != nil {
		return nil, err
	}

	return garage, nil
}

// UpdateGarageStatus updates a garage's operational status
func (s *garageService) UpdateGarageStatus(ctx context.Context, id string, status domain.GarageStatus) (*domain.Garage, error) {
	garage, err := s.garageRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	garage.UpdateStatus(status)

	if err := s.garageRepo.Update(ctx, garage); err != nil {
		return nil, err
	}

	return garage, nil
}

// UpdateAvailableSpaces increments or decrements the available spaces
func (s *garageService) UpdateAvailableSpaces(ctx context.Context, id string, increment bool) (*domain.Garage, error) {
	garage, err := s.garageRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if increment {
		garage.IncrementAvailableSpaces()
	} else {
		if err := garage.DecrementAvailableSpaces(); err != nil {
			return nil, err
		}
	}

	if err := s.garageRepo.Update(ctx, garage); err != nil {
		return nil, err
	}

	return garage, nil
}

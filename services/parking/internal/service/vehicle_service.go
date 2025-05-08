package service

import (
	"context"

	"github.com/jonbelaire/repotown/packages/go-core/logging"
	"github.com/jonbelaire/repotown/services/parking/internal/domain"
	"github.com/jonbelaire/repotown/services/parking/internal/repository"
)

// VehicleService provides business logic for vehicle management
type VehicleService interface {
	GetVehicle(ctx context.Context, id string) (*domain.Vehicle, error)
	GetVehicleByLicensePlate(ctx context.Context, licensePlate string) (*domain.Vehicle, error)
	ListVehicles(ctx context.Context, limit, offset int) ([]*domain.Vehicle, error)
	ListVehiclesByOwner(ctx context.Context, ownerID string) ([]*domain.Vehicle, error)
	CreateVehicle(ctx context.Context, licensePlate string, vehicleType domain.VehicleType, make, model, color, ownerID string) (*domain.Vehicle, error)
	UpdateVehicle(ctx context.Context, id, make, model, color string) (*domain.Vehicle, error)
	UpdateVehicleType(ctx context.Context, id string, vehicleType domain.VehicleType) (*domain.Vehicle, error)
	ChangeVehicleOwner(ctx context.Context, id, ownerID string) (*domain.Vehicle, error)
}

// vehicleService implements VehicleService
type vehicleService struct {
	vehicleRepo repository.VehicleRepository
	logger      logging.Logger
}

// NewVehicleService creates a new vehicle service
func NewVehicleService(vehicleRepo repository.VehicleRepository, logger logging.Logger) VehicleService {
	return &vehicleService{
		vehicleRepo: vehicleRepo,
		logger:      logger,
	}
}

// GetVehicle retrieves a vehicle by ID
func (s *vehicleService) GetVehicle(ctx context.Context, id string) (*domain.Vehicle, error) {
	return s.vehicleRepo.GetByID(ctx, id)
}

// GetVehicleByLicensePlate retrieves a vehicle by license plate
func (s *vehicleService) GetVehicleByLicensePlate(ctx context.Context, licensePlate string) (*domain.Vehicle, error) {
	return s.vehicleRepo.GetByLicensePlate(ctx, licensePlate)
}

// ListVehicles retrieves vehicles with pagination
func (s *vehicleService) ListVehicles(ctx context.Context, limit, offset int) ([]*domain.Vehicle, error) {
	return s.vehicleRepo.List(ctx, limit, offset)
}

// ListVehiclesByOwner retrieves vehicles for a specific owner
func (s *vehicleService) ListVehiclesByOwner(ctx context.Context, ownerID string) ([]*domain.Vehicle, error) {
	return s.vehicleRepo.ListByOwner(ctx, ownerID)
}

// CreateVehicle creates a new vehicle
func (s *vehicleService) CreateVehicle(ctx context.Context, licensePlate string, vehicleType domain.VehicleType, make, model, color, ownerID string) (*domain.Vehicle, error) {
	// Check if vehicle with this license plate already exists
	existingVehicle, err := s.vehicleRepo.GetByLicensePlate(ctx, licensePlate)
	if err == nil && existingVehicle != nil {
		return nil, domain.ErrVehicleExists
	}

	vehicle := domain.NewVehicle(licensePlate, vehicleType, make, model, color, ownerID)
	if err := s.vehicleRepo.Create(ctx, vehicle); err != nil {
		return nil, err
	}
	return vehicle, nil
}

// UpdateVehicle updates a vehicle's details
func (s *vehicleService) UpdateVehicle(ctx context.Context, id, make, model, color string) (*domain.Vehicle, error) {
	vehicle, err := s.vehicleRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	vehicle.UpdateDetails(make, model, color)

	if err := s.vehicleRepo.Update(ctx, vehicle); err != nil {
		return nil, err
	}

	return vehicle, nil
}

// UpdateVehicleType updates a vehicle's type
func (s *vehicleService) UpdateVehicleType(ctx context.Context, id string, vehicleType domain.VehicleType) (*domain.Vehicle, error) {
	vehicle, err := s.vehicleRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	vehicle.UpdateType(vehicleType)

	if err := s.vehicleRepo.Update(ctx, vehicle); err != nil {
		return nil, err
	}

	return vehicle, nil
}

// ChangeVehicleOwner changes the owner of a vehicle
func (s *vehicleService) ChangeVehicleOwner(ctx context.Context, id, ownerID string) (*domain.Vehicle, error) {
	vehicle, err := s.vehicleRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	vehicle.ChangeOwner(ownerID)

	if err := s.vehicleRepo.Update(ctx, vehicle); err != nil {
		return nil, err
	}

	return vehicle, nil
}
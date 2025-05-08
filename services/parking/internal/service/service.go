package service

import (
	"github.com/jonbelaire/repotown/packages/go-core/logging"
	"github.com/jonbelaire/repotown/services/parking/internal/repository"
)

// Services holds all service implementations
type Services struct {
	Garage         GarageService
	Vehicle        VehicleService
	ParkingSession ParkingSessionService
	Reservation    ReservationService
}

// NewServices creates all services
func NewServices(repos *repository.Repositories, logger logging.Logger) *Services {
	return &Services{
		Garage:         NewGarageService(repos.Garage, logger),
		Vehicle:        NewVehicleService(repos.Vehicle, logger),
		ParkingSession: NewParkingSessionService(repos.ParkingSession, repos.Garage, repos.Vehicle, logger),
		Reservation:    NewReservationService(repos.Reservation, repos.Garage, logger),
	}
}
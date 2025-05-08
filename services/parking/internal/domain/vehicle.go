package domain

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

// Common errors
var (
	ErrVehicleNotFound = errors.New("vehicle not found")
	ErrVehicleExists   = errors.New("vehicle with license plate already exists")
)

// VehicleType represents the type of vehicle
type VehicleType string

// Vehicle types
const (
	VehicleTypeCar         VehicleType = "car"
	VehicleTypeMotorcycle  VehicleType = "motorcycle"
	VehicleTypeTruck       VehicleType = "truck"
	VehicleTypeElectric    VehicleType = "electric"
)

// Vehicle represents a vehicle that can park in garages
type Vehicle struct {
	ID              string      `json:"id"`
	LicensePlate    string      `json:"license_plate"`
	Type            VehicleType `json:"type"`
	Make            string      `json:"make"`
	Model           string      `json:"model"`
	Color           string      `json:"color"`
	OwnerID         string      `json:"owner_id,omitempty"` // Optional reference to a customer
	CreatedAt       time.Time   `json:"created_at"`
	UpdatedAt       time.Time   `json:"updated_at"`
}

// NewVehicle creates a new vehicle
func NewVehicle(licensePlate string, vehicleType VehicleType, make, model, color string, ownerID string) *Vehicle {
	now := time.Now()
	return &Vehicle{
		ID:              uuid.New().String(),
		LicensePlate:    licensePlate,
		Type:            vehicleType,
		Make:            make,
		Model:           model,
		Color:           color,
		OwnerID:         ownerID,
		CreatedAt:       now,
		UpdatedAt:       now,
	}
}

// IsElectric checks if the vehicle is electric
func (v *Vehicle) IsElectric() bool {
	return v.Type == VehicleTypeElectric
}

// UpdateDetails updates the vehicle details
func (v *Vehicle) UpdateDetails(make, model, color string) {
	v.Make = make
	v.Model = model
	v.Color = color
	v.UpdatedAt = time.Now()
}

// UpdateType updates the vehicle type
func (v *Vehicle) UpdateType(vehicleType VehicleType) {
	v.Type = vehicleType
	v.UpdatedAt = time.Now()
}

// ChangeOwner changes the owner of the vehicle
func (v *Vehicle) ChangeOwner(ownerID string) {
	v.OwnerID = ownerID
	v.UpdatedAt = time.Now()
}
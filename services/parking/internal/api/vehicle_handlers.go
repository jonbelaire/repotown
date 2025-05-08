package api

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/jonbelaire/repotown/packages/go-core/httputils"
	"github.com/jonbelaire/repotown/packages/go-core/logging"
	"github.com/jonbelaire/repotown/services/parking/internal/domain"
	"github.com/jonbelaire/repotown/services/parking/internal/service"
)

// VehicleHandler manages vehicle-related HTTP requests
type VehicleHandler struct {
	vehicleService service.VehicleService
	logger         logging.Logger
}

// NewVehicleHandler creates a new vehicle handler
func NewVehicleHandler(vehicleService service.VehicleService, logger logging.Logger) *VehicleHandler {
	return &VehicleHandler{
		vehicleService: vehicleService,
		logger:         logger,
	}
}

// CreateVehicleRequest defines the request body for creating a vehicle
type CreateVehicleRequest struct {
	LicensePlate string            `json:"license_plate" validate:"required"`
	Type         domain.VehicleType `json:"type" validate:"required"`
	Make         string            `json:"make" validate:"required"`
	Model        string            `json:"model" validate:"required"`
	Color        string            `json:"color" validate:"required"`
	OwnerID      string            `json:"owner_id,omitempty"`
}

// UpdateVehicleRequest defines the request body for updating a vehicle
type UpdateVehicleRequest struct {
	Make         string `json:"make" validate:"required"`
	Model        string `json:"model" validate:"required"`
	Color        string `json:"color" validate:"required"`
}

// UpdateVehicleTypeRequest defines the request body for updating a vehicle type
type UpdateVehicleTypeRequest struct {
	Type domain.VehicleType `json:"type" validate:"required"`
}

// ChangeVehicleOwnerRequest defines the request body for changing a vehicle owner
type ChangeVehicleOwnerRequest struct {
	OwnerID string `json:"owner_id" validate:"required"`
}

// RegisterRoutes registers vehicle routes
func (h *VehicleHandler) RegisterRoutes(r chi.Router) {
	r.Route("/vehicles", func(r chi.Router) {
		r.Get("/", h.listVehicles)
		r.Post("/", h.createVehicle)
		r.Get("/license/{plate}", h.getVehicleByLicensePlate)
		r.Get("/owner/{ownerID}", h.listVehiclesByOwner)
		r.Get("/{id}", h.getVehicle)
		r.Put("/{id}", h.updateVehicle)
		r.Patch("/{id}/type", h.updateVehicleType)
		r.Patch("/{id}/owner", h.changeVehicleOwner)
	})
}

// listVehicles handles GET /vehicles
func (h *VehicleHandler) listVehicles(w http.ResponseWriter, r *http.Request) {
	// Parse pagination parameters
	limit, offset := getPaginationParams(r)

	// Get vehicles
	vehicles, err := h.vehicleService.ListVehicles(r.Context(), limit, offset)
	if err != nil {
		h.logger.Error("Failed to list vehicles", "error", err)
		httputils.ErrorJSON(w, httputils.ErrInternal)
		return
	}

	httputils.JSON(w, http.StatusOK, vehicles)
}

// createVehicle handles POST /vehicles
func (h *VehicleHandler) createVehicle(w http.ResponseWriter, r *http.Request) {
	// Parse request body
	var req CreateVehicleRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httputils.ErrorJSON(w, httputils.ErrBadRequest)
		return
	}

	// Create vehicle
	vehicle, err := h.vehicleService.CreateVehicle(
		r.Context(),
		req.LicensePlate,
		req.Type,
		req.Make,
		req.Model,
		req.Color,
		req.OwnerID,
	)
	if err != nil {
		if err == domain.ErrVehicleExists {
			httputils.ErrorJSON(w, httputils.NewError(http.StatusBadRequest, "VEHICLE_EXISTS", "Vehicle with this license plate already exists", nil))
			return
		}
		h.logger.Error("Failed to create vehicle", "error", err)
		httputils.ErrorJSON(w, httputils.ErrInternal)
		return
	}

	httputils.JSON(w, http.StatusCreated, vehicle)
}

// getVehicle handles GET /vehicles/{id}
func (h *VehicleHandler) getVehicle(w http.ResponseWriter, r *http.Request) {
	// Get vehicle ID from path
	id := chi.URLParam(r, "id")

	// Get vehicle
	vehicle, err := h.vehicleService.GetVehicle(r.Context(), id)
	if err != nil {
		if err == domain.ErrVehicleNotFound {
			httputils.ErrorJSON(w, httputils.ErrNotFound)
			return
		}
		h.logger.Error("Failed to get vehicle", "error", err)
		httputils.ErrorJSON(w, httputils.ErrInternal)
		return
	}

	httputils.JSON(w, http.StatusOK, vehicle)
}

// getVehicleByLicensePlate handles GET /vehicles/license/{plate}
func (h *VehicleHandler) getVehicleByLicensePlate(w http.ResponseWriter, r *http.Request) {
	// Get license plate from path
	plate := chi.URLParam(r, "plate")

	// Get vehicle
	vehicle, err := h.vehicleService.GetVehicleByLicensePlate(r.Context(), plate)
	if err != nil {
		if err == domain.ErrVehicleNotFound {
			httputils.ErrorJSON(w, httputils.ErrNotFound)
			return
		}
		h.logger.Error("Failed to get vehicle by license plate", "error", err)
		httputils.ErrorJSON(w, httputils.ErrInternal)
		return
	}

	httputils.JSON(w, http.StatusOK, vehicle)
}

// listVehiclesByOwner handles GET /vehicles/owner/{ownerID}
func (h *VehicleHandler) listVehiclesByOwner(w http.ResponseWriter, r *http.Request) {
	// Get owner ID from path
	ownerID := chi.URLParam(r, "ownerID")

	// Get vehicles
	vehicles, err := h.vehicleService.ListVehiclesByOwner(r.Context(), ownerID)
	if err != nil {
		h.logger.Error("Failed to list vehicles by owner", "error", err)
		httputils.ErrorJSON(w, httputils.ErrInternal)
		return
	}

	httputils.JSON(w, http.StatusOK, vehicles)
}

// updateVehicle handles PUT /vehicles/{id}
func (h *VehicleHandler) updateVehicle(w http.ResponseWriter, r *http.Request) {
	// Get vehicle ID from path
	id := chi.URLParam(r, "id")

	// Parse request body
	var req UpdateVehicleRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httputils.ErrorJSON(w, httputils.ErrBadRequest)
		return
	}

	// Update vehicle
	vehicle, err := h.vehicleService.UpdateVehicle(
		r.Context(),
		id,
		req.Make,
		req.Model,
		req.Color,
	)
	if err != nil {
		if err == domain.ErrVehicleNotFound {
			httputils.ErrorJSON(w, httputils.ErrNotFound)
			return
		}
		h.logger.Error("Failed to update vehicle", "error", err)
		httputils.ErrorJSON(w, httputils.ErrInternal)
		return
	}

	httputils.JSON(w, http.StatusOK, vehicle)
}

// updateVehicleType handles PATCH /vehicles/{id}/type
func (h *VehicleHandler) updateVehicleType(w http.ResponseWriter, r *http.Request) {
	// Get vehicle ID from path
	id := chi.URLParam(r, "id")

	// Parse request body
	var req UpdateVehicleTypeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httputils.ErrorJSON(w, httputils.ErrBadRequest)
		return
	}

	// Update vehicle type
	vehicle, err := h.vehicleService.UpdateVehicleType(r.Context(), id, req.Type)
	if err != nil {
		if err == domain.ErrVehicleNotFound {
			httputils.ErrorJSON(w, httputils.ErrNotFound)
			return
		}
		h.logger.Error("Failed to update vehicle type", "error", err)
		httputils.ErrorJSON(w, httputils.ErrInternal)
		return
	}

	httputils.JSON(w, http.StatusOK, vehicle)
}

// changeVehicleOwner handles PATCH /vehicles/{id}/owner
func (h *VehicleHandler) changeVehicleOwner(w http.ResponseWriter, r *http.Request) {
	// Get vehicle ID from path
	id := chi.URLParam(r, "id")

	// Parse request body
	var req ChangeVehicleOwnerRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httputils.ErrorJSON(w, httputils.ErrBadRequest)
		return
	}

	// Change vehicle owner
	vehicle, err := h.vehicleService.ChangeVehicleOwner(r.Context(), id, req.OwnerID)
	if err != nil {
		if err == domain.ErrVehicleNotFound {
			httputils.ErrorJSON(w, httputils.ErrNotFound)
			return
		}
		h.logger.Error("Failed to change vehicle owner", "error", err)
		httputils.ErrorJSON(w, httputils.ErrInternal)
		return
	}

	httputils.JSON(w, http.StatusOK, vehicle)
}
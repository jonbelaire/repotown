package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/jonbelaire/repotown/packages/go-core/httputils"
	"github.com/jonbelaire/repotown/packages/go-core/logging"
	"github.com/jonbelaire/repotown/services/parking/internal/domain"
	"github.com/jonbelaire/repotown/services/parking/internal/service"
)

// GarageHandler manages garage-related HTTP requests
type GarageHandler struct {
	garageService service.GarageService
	logger        logging.Logger
}

// NewGarageHandler creates a new garage handler
func NewGarageHandler(garageService service.GarageService, logger logging.Logger) *GarageHandler {
	return &GarageHandler{
		garageService: garageService,
		logger:        logger,
	}
}

// CreateGarageRequest defines the request body for creating a garage
type CreateGarageRequest struct {
	Name               string  `json:"name" validate:"required"`
	Address            string  `json:"address" validate:"required"`
	TotalSpaces        int     `json:"total_spaces" validate:"required,gt=0"`
	HourlyRate         int64   `json:"hourly_rate" validate:"required,gte=0"`
	DailyRate          int64   `json:"daily_rate" validate:"required,gte=0"`
	OperatingHours     string  `json:"operating_hours" validate:"required"`
	HasElectricCharging bool    `json:"has_electric_charging"`
}

// UpdateGarageRequest defines the request body for updating a garage
type UpdateGarageRequest struct {
	Name           string  `json:"name" validate:"required"`
	Address        string  `json:"address" validate:"required"`
	HourlyRate     int64   `json:"hourly_rate" validate:"required,gte=0"`
	DailyRate      int64   `json:"daily_rate" validate:"required,gte=0"`
	OperatingHours string  `json:"operating_hours" validate:"required"`
}

// UpdateGarageStatusRequest defines the request body for updating a garage's status
type UpdateGarageStatusRequest struct {
	Status domain.GarageStatus `json:"status" validate:"required"`
}

// RegisterRoutes registers garage routes
func (h *GarageHandler) RegisterRoutes(r chi.Router) {
	r.Route("/garages", func(r chi.Router) {
		r.Get("/", h.listGarages)
		r.Post("/", h.createGarage)
		r.Get("/available", h.listAvailableGarages)
		r.Get("/{id}", h.getGarage)
		r.Put("/{id}", h.updateGarage)
		r.Patch("/{id}/status", h.updateGarageStatus)
	})
}

// listGarages handles GET /garages
func (h *GarageHandler) listGarages(w http.ResponseWriter, r *http.Request) {
	// Parse pagination parameters
	limit, offset := getPaginationParams(r)

	// Get garages
	garages, err := h.garageService.ListGarages(r.Context(), limit, offset)
	if err != nil {
		h.logger.Error("Failed to list garages", "error", err)
		httputils.ErrorJSON(w, httputils.ErrInternal)
		return
	}

	httputils.JSON(w, http.StatusOK, garages)
}

// createGarage handles POST /garages
func (h *GarageHandler) createGarage(w http.ResponseWriter, r *http.Request) {
	// Parse request body
	var req CreateGarageRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httputils.ErrorJSON(w, httputils.ErrBadRequest)
		return
	}

	// Create garage
	garage, err := h.garageService.CreateGarage(
		r.Context(),
		req.Name,
		req.Address,
		req.TotalSpaces,
		req.HourlyRate,
		req.DailyRate,
		req.OperatingHours,
		req.HasElectricCharging,
	)
	if err != nil {
		h.logger.Error("Failed to create garage", "error", err)
		httputils.ErrorJSON(w, httputils.ErrInternal)
		return
	}

	httputils.JSON(w, http.StatusCreated, garage)
}

// listAvailableGarages handles GET /garages/available
func (h *GarageHandler) listAvailableGarages(w http.ResponseWriter, r *http.Request) {
	// Get available garages
	garages, err := h.garageService.ListAvailableGarages(r.Context())
	if err != nil {
		h.logger.Error("Failed to list available garages", "error", err)
		httputils.ErrorJSON(w, httputils.ErrInternal)
		return
	}

	httputils.JSON(w, http.StatusOK, garages)
}

// getGarage handles GET /garages/{id}
func (h *GarageHandler) getGarage(w http.ResponseWriter, r *http.Request) {
	// Get garage ID from path
	id := chi.URLParam(r, "id")

	// Get garage
	garage, err := h.garageService.GetGarage(r.Context(), id)
	if err != nil {
		if err == domain.ErrGarageNotFound {
			httputils.ErrorJSON(w, httputils.ErrNotFound)
			return
		}
		h.logger.Error("Failed to get garage", "error", err)
		httputils.ErrorJSON(w, httputils.ErrInternal)
		return
	}

	httputils.JSON(w, http.StatusOK, garage)
}

// updateGarage handles PUT /garages/{id}
func (h *GarageHandler) updateGarage(w http.ResponseWriter, r *http.Request) {
	// Get garage ID from path
	id := chi.URLParam(r, "id")

	// Parse request body
	var req UpdateGarageRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httputils.ErrorJSON(w, httputils.ErrBadRequest)
		return
	}

	// Update garage
	garage, err := h.garageService.UpdateGarage(
		r.Context(),
		id,
		req.Name,
		req.Address,
		req.OperatingHours,
		req.HourlyRate,
		req.DailyRate,
	)
	if err != nil {
		if err == domain.ErrGarageNotFound {
			httputils.ErrorJSON(w, httputils.ErrNotFound)
			return
		}
		h.logger.Error("Failed to update garage", "error", err)
		httputils.ErrorJSON(w, httputils.ErrInternal)
		return
	}

	httputils.JSON(w, http.StatusOK, garage)
}

// updateGarageStatus handles PATCH /garages/{id}/status
func (h *GarageHandler) updateGarageStatus(w http.ResponseWriter, r *http.Request) {
	// Get garage ID from path
	id := chi.URLParam(r, "id")

	// Parse request body
	var req UpdateGarageStatusRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httputils.ErrorJSON(w, httputils.ErrBadRequest)
		return
	}

	// Update garage status
	garage, err := h.garageService.UpdateGarageStatus(r.Context(), id, req.Status)
	if err != nil {
		if err == domain.ErrGarageNotFound {
			httputils.ErrorJSON(w, httputils.ErrNotFound)
			return
		}
		h.logger.Error("Failed to update garage status", "error", err)
		httputils.ErrorJSON(w, httputils.ErrInternal)
		return
	}

	httputils.JSON(w, http.StatusOK, garage)
}

// Helper functions

// getPaginationParams extracts pagination parameters from request
func getPaginationParams(r *http.Request) (int, int) {
	// Parse limit
	limitStr := r.URL.Query().Get("limit")
	limit := 10 // Default limit
	if limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
			limit = l
		}
	}

	// Parse offset
	offsetStr := r.URL.Query().Get("offset")
	offset := 0 // Default offset
	if offsetStr != "" {
		if o, err := strconv.Atoi(offsetStr); err == nil && o >= 0 {
			offset = o
		}
	}

	return limit, offset
}
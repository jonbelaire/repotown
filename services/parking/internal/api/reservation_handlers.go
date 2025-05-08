package api

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/jonbelaire/repotown/packages/go-core/httputils"
	"github.com/jonbelaire/repotown/packages/go-core/logging"
	"github.com/jonbelaire/repotown/services/parking/internal/domain"
	"github.com/jonbelaire/repotown/services/parking/internal/service"
)

// ReservationHandler manages reservation-related HTTP requests
type ReservationHandler struct {
	reservationService service.ReservationService
	logger             logging.Logger
}

// NewReservationHandler creates a new reservation handler
func NewReservationHandler(reservationService service.ReservationService, logger logging.Logger) *ReservationHandler {
	return &ReservationHandler{
		reservationService: reservationService,
		logger:             logger,
	}
}

// CreateReservationRequest defines the request body for creating a reservation
type CreateReservationRequest struct {
	GarageID    string    `json:"garage_id" validate:"required"`
	VehicleID   string    `json:"vehicle_id" validate:"required"`
	StartTime   time.Time `json:"start_time" validate:"required"`
	EndTime     time.Time `json:"end_time" validate:"required"`
	AmountPaid  int64     `json:"amount_paid" validate:"required,gte=0"`
}

// RegisterRoutes registers reservation routes
func (h *ReservationHandler) RegisterRoutes(r chi.Router) {
	r.Route("/reservations", func(r chi.Router) {
		r.Get("/", h.listReservations)
		r.Post("/", h.createReservation)
		r.Get("/active", h.listActiveReservations)
		r.Get("/garage/{garageID}", h.listReservationsByGarage)
		r.Get("/vehicle/{vehicleID}", h.listReservationsByVehicle)
		r.Get("/code/{code}", h.getReservationByCode)
		r.Get("/{id}", h.getReservation)
		r.Put("/{id}/use", h.useReservation)
		r.Put("/{id}/cancel", h.cancelReservation)
		r.Post("/check-expired", h.checkExpiredReservations)
	})
}

// listReservations handles GET /reservations
func (h *ReservationHandler) listReservations(w http.ResponseWriter, r *http.Request) {
	// Parse pagination parameters
	limit, offset := getPaginationParams(r)

	// Get reservations
	reservations, err := h.reservationService.ListReservations(r.Context(), limit, offset)
	if err != nil {
		h.logger.Error("Failed to list reservations", "error", err)
		httputils.ErrorJSON(w, httputils.ErrInternal)
		return
	}

	httputils.JSON(w, http.StatusOK, reservations)
}

// createReservation handles POST /reservations
func (h *ReservationHandler) createReservation(w http.ResponseWriter, r *http.Request) {
	// Parse request body
	var req CreateReservationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httputils.ErrorJSON(w, httputils.ErrBadRequest)
		return
	}

	// Validate time range
	if req.StartTime.After(req.EndTime) {
		httputils.ErrorJSON(w, httputils.NewError(http.StatusBadRequest, "INVALID_TIME_RANGE", "Start time must be before end time", nil))
		return
	}

	// Create reservation
	reservation, err := h.reservationService.CreateReservation(
		r.Context(),
		req.GarageID,
		req.VehicleID,
		req.StartTime,
		req.EndTime,
		req.AmountPaid,
	)
	if err != nil {
		h.handleReservationError(w, err, "Failed to create reservation")
		return
	}

	httputils.JSON(w, http.StatusCreated, reservation)
}

// listActiveReservations handles GET /reservations/active
func (h *ReservationHandler) listActiveReservations(w http.ResponseWriter, r *http.Request) {
	// Get active reservations
	reservations, err := h.reservationService.ListActiveReservations(r.Context())
	if err != nil {
		h.logger.Error("Failed to list active reservations", "error", err)
		httputils.ErrorJSON(w, httputils.ErrInternal)
		return
	}

	httputils.JSON(w, http.StatusOK, reservations)
}

// listReservationsByGarage handles GET /reservations/garage/{garageID}
func (h *ReservationHandler) listReservationsByGarage(w http.ResponseWriter, r *http.Request) {
	// Get garage ID from path
	garageID := chi.URLParam(r, "garageID")

	// Get reservations
	reservations, err := h.reservationService.ListReservationsByGarage(r.Context(), garageID)
	if err != nil {
		h.logger.Error("Failed to list reservations by garage", "error", err)
		httputils.ErrorJSON(w, httputils.ErrInternal)
		return
	}

	httputils.JSON(w, http.StatusOK, reservations)
}

// listReservationsByVehicle handles GET /reservations/vehicle/{vehicleID}
func (h *ReservationHandler) listReservationsByVehicle(w http.ResponseWriter, r *http.Request) {
	// Get vehicle ID from path
	vehicleID := chi.URLParam(r, "vehicleID")

	// Get reservations
	reservations, err := h.reservationService.ListReservationsByVehicle(r.Context(), vehicleID)
	if err != nil {
		h.logger.Error("Failed to list reservations by vehicle", "error", err)
		httputils.ErrorJSON(w, httputils.ErrInternal)
		return
	}

	httputils.JSON(w, http.StatusOK, reservations)
}

// getReservation handles GET /reservations/{id}
func (h *ReservationHandler) getReservation(w http.ResponseWriter, r *http.Request) {
	// Get reservation ID from path
	id := chi.URLParam(r, "id")

	// Get reservation
	reservation, err := h.reservationService.GetReservation(r.Context(), id)
	if err != nil {
		if err == domain.ErrReservationNotFound {
			httputils.ErrorJSON(w, httputils.ErrNotFound)
			return
		}
		h.logger.Error("Failed to get reservation", "error", err)
		httputils.ErrorJSON(w, httputils.ErrInternal)
		return
	}

	httputils.JSON(w, http.StatusOK, reservation)
}

// getReservationByCode handles GET /reservations/code/{code}
func (h *ReservationHandler) getReservationByCode(w http.ResponseWriter, r *http.Request) {
	// Get confirmation code from path
	code := chi.URLParam(r, "code")

	// Get reservation
	reservation, err := h.reservationService.GetReservationByConfirmationCode(r.Context(), code)
	if err != nil {
		if err == domain.ErrReservationNotFound {
			httputils.ErrorJSON(w, httputils.ErrNotFound)
			return
		}
		h.logger.Error("Failed to get reservation by code", "error", err)
		httputils.ErrorJSON(w, httputils.ErrInternal)
		return
	}

	httputils.JSON(w, http.StatusOK, reservation)
}

// useReservation handles PUT /reservations/{id}/use
func (h *ReservationHandler) useReservation(w http.ResponseWriter, r *http.Request) {
	// Get reservation ID from path
	id := chi.URLParam(r, "id")

	// Use reservation
	reservation, err := h.reservationService.UseReservation(r.Context(), id)
	if err != nil {
		h.handleReservationError(w, err, "Failed to use reservation")
		return
	}

	httputils.JSON(w, http.StatusOK, reservation)
}

// cancelReservation handles PUT /reservations/{id}/cancel
func (h *ReservationHandler) cancelReservation(w http.ResponseWriter, r *http.Request) {
	// Get reservation ID from path
	id := chi.URLParam(r, "id")

	// Cancel reservation
	reservation, err := h.reservationService.CancelReservation(r.Context(), id)
	if err != nil {
		h.handleReservationError(w, err, "Failed to cancel reservation")
		return
	}

	httputils.JSON(w, http.StatusOK, reservation)
}

// checkExpiredReservations handles POST /reservations/check-expired
func (h *ReservationHandler) checkExpiredReservations(w http.ResponseWriter, r *http.Request) {
	// Check for expired reservations
	count, err := h.reservationService.CheckExpiredReservations(r.Context())
	if err != nil {
		h.logger.Error("Failed to check expired reservations", "error", err)
		httputils.ErrorJSON(w, httputils.ErrInternal)
		return
	}

	httputils.JSON(w, http.StatusOK, map[string]int{"expired_count": count})
}

// handleReservationError handles common reservation errors
func (h *ReservationHandler) handleReservationError(w http.ResponseWriter, err error, logMessage string) {
	switch err {
	case domain.ErrReservationNotFound:
		httputils.ErrorJSON(w, httputils.ErrNotFound)
	case domain.ErrReservationCancelled:
		httputils.ErrorJSON(w, httputils.NewError(http.StatusBadRequest, "RESERVATION_CANCELLED", "Reservation is already cancelled or used", nil))
	case domain.ErrReservationExpired:
		httputils.ErrorJSON(w, httputils.NewError(http.StatusBadRequest, "RESERVATION_EXPIRED", "Reservation has expired", nil))
	case domain.ErrGarageNotFound:
		httputils.ErrorJSON(w, httputils.NewError(http.StatusBadRequest, "GARAGE_NOT_FOUND", "Garage not found", nil))
	case domain.ErrGarageCapacityFull:
		httputils.ErrorJSON(w, httputils.NewError(http.StatusBadRequest, "GARAGE_FULL", "Garage is at full capacity", nil))
	case domain.ErrVehicleNotFound:
		httputils.ErrorJSON(w, httputils.NewError(http.StatusBadRequest, "VEHICLE_NOT_FOUND", "Vehicle not found", nil))
	default:
		h.logger.Error(logMessage, "error", err)
		httputils.ErrorJSON(w, httputils.ErrInternal)
	}
}
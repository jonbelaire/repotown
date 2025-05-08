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

// ParkingSessionHandler manages parking session-related HTTP requests
type ParkingSessionHandler struct {
	sessionService service.ParkingSessionService
	logger         logging.Logger
}

// NewParkingSessionHandler creates a new parking session handler
func NewParkingSessionHandler(sessionService service.ParkingSessionService, logger logging.Logger) *ParkingSessionHandler {
	return &ParkingSessionHandler{
		sessionService: sessionService,
		logger:         logger,
	}
}

// StartSessionRequest defines the request body for starting a parking session
type StartSessionRequest struct {
	GarageID    string `json:"garage_id" validate:"required"`
	VehicleID   string `json:"vehicle_id" validate:"required"`
	SpotNumber  string `json:"spot_number,omitempty"`
	IsPrepaid   bool   `json:"is_prepaid"`
}

// EndSessionRequest defines the request body for ending a parking session with a custom end time
type EndSessionRequest struct {
	EndTime time.Time `json:"end_time,omitempty"`
}

// UpdateSpotRequest defines the request body for updating the parking spot
type UpdateSpotRequest struct {
	SpotNumber string `json:"spot_number" validate:"required"`
}

// AddNotesRequest defines the request body for adding notes to a session
type AddNotesRequest struct {
	Notes string `json:"notes" validate:"required"`
}

// RegisterRoutes registers parking session routes
func (h *ParkingSessionHandler) RegisterRoutes(r chi.Router) {
	r.Route("/sessions", func(r chi.Router) {
		r.Get("/", h.listSessions)
		r.Post("/", h.startSession)
		r.Get("/active", h.listActiveSessions)
		r.Get("/garage/{garageID}", h.listSessionsByGarage)
		r.Get("/vehicle/{vehicleID}", h.listSessionsByVehicle)
		r.Get("/{id}", h.getSession)
		r.Put("/{id}/end", h.endSession)
		r.Put("/{id}/cancel", h.cancelSession)
		r.Patch("/{id}/spot", h.updateSessionSpot)
		r.Patch("/{id}/notes", h.addSessionNotes)
		r.Patch("/{id}/pay", h.markSessionAsPaid)
	})
}

// listSessions handles GET /sessions
func (h *ParkingSessionHandler) listSessions(w http.ResponseWriter, r *http.Request) {
	// Parse pagination parameters
	limit, offset := getPaginationParams(r)

	// Get sessions
	sessions, err := h.sessionService.ListSessions(r.Context(), limit, offset)
	if err != nil {
		h.logger.Error("Failed to list sessions", "error", err)
		httputils.ErrorJSON(w, httputils.ErrInternal)
		return
	}

	httputils.JSON(w, http.StatusOK, sessions)
}

// startSession handles POST /sessions
func (h *ParkingSessionHandler) startSession(w http.ResponseWriter, r *http.Request) {
	// Parse request body
	var req StartSessionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httputils.ErrorJSON(w, httputils.ErrBadRequest)
		return
	}

	// Start session
	session, err := h.sessionService.StartSession(
		r.Context(),
		req.GarageID,
		req.VehicleID,
		req.SpotNumber,
		req.IsPrepaid,
	)
	if err != nil {
		h.handleSessionError(w, err, "Failed to start session")
		return
	}

	httputils.JSON(w, http.StatusCreated, session)
}

// listActiveSessions handles GET /sessions/active
func (h *ParkingSessionHandler) listActiveSessions(w http.ResponseWriter, r *http.Request) {
	// Get active sessions
	sessions, err := h.sessionService.ListActiveSessions(r.Context())
	if err != nil {
		h.logger.Error("Failed to list active sessions", "error", err)
		httputils.ErrorJSON(w, httputils.ErrInternal)
		return
	}

	httputils.JSON(w, http.StatusOK, sessions)
}

// listSessionsByGarage handles GET /sessions/garage/{garageID}
func (h *ParkingSessionHandler) listSessionsByGarage(w http.ResponseWriter, r *http.Request) {
	// Get garage ID from path
	garageID := chi.URLParam(r, "garageID")

	// Get sessions
	sessions, err := h.sessionService.ListSessionsByGarage(r.Context(), garageID)
	if err != nil {
		h.logger.Error("Failed to list sessions by garage", "error", err)
		httputils.ErrorJSON(w, httputils.ErrInternal)
		return
	}

	httputils.JSON(w, http.StatusOK, sessions)
}

// listSessionsByVehicle handles GET /sessions/vehicle/{vehicleID}
func (h *ParkingSessionHandler) listSessionsByVehicle(w http.ResponseWriter, r *http.Request) {
	// Get vehicle ID from path
	vehicleID := chi.URLParam(r, "vehicleID")

	// Get sessions
	sessions, err := h.sessionService.ListSessionsByVehicle(r.Context(), vehicleID)
	if err != nil {
		h.logger.Error("Failed to list sessions by vehicle", "error", err)
		httputils.ErrorJSON(w, httputils.ErrInternal)
		return
	}

	httputils.JSON(w, http.StatusOK, sessions)
}

// getSession handles GET /sessions/{id}
func (h *ParkingSessionHandler) getSession(w http.ResponseWriter, r *http.Request) {
	// Get session ID from path
	id := chi.URLParam(r, "id")

	// Get session
	session, err := h.sessionService.GetSession(r.Context(), id)
	if err != nil {
		if err == domain.ErrSessionNotFound {
			httputils.ErrorJSON(w, httputils.ErrNotFound)
			return
		}
		h.logger.Error("Failed to get session", "error", err)
		httputils.ErrorJSON(w, httputils.ErrInternal)
		return
	}

	httputils.JSON(w, http.StatusOK, session)
}

// endSession handles PUT /sessions/{id}/end
func (h *ParkingSessionHandler) endSession(w http.ResponseWriter, r *http.Request) {
	// Get session ID from path
	id := chi.URLParam(r, "id")

	// Parse request body to check if there's a custom end time
	var req EndSessionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.EndTime.IsZero() {
		// If no valid end time provided, use current time
		session, err := h.sessionService.EndSession(r.Context(), id)
		if err != nil {
			h.handleSessionError(w, err, "Failed to end session")
			return
		}
		httputils.JSON(w, http.StatusOK, session)
		return
	}

	// Use custom end time
	session, err := h.sessionService.EndSessionWithCustomTime(r.Context(), id, req.EndTime)
	if err != nil {
		h.handleSessionError(w, err, "Failed to end session with custom time")
		return
	}

	httputils.JSON(w, http.StatusOK, session)
}

// cancelSession handles PUT /sessions/{id}/cancel
func (h *ParkingSessionHandler) cancelSession(w http.ResponseWriter, r *http.Request) {
	// Get session ID from path
	id := chi.URLParam(r, "id")

	// Cancel session
	session, err := h.sessionService.CancelSession(r.Context(), id)
	if err != nil {
		h.handleSessionError(w, err, "Failed to cancel session")
		return
	}

	httputils.JSON(w, http.StatusOK, session)
}

// updateSessionSpot handles PATCH /sessions/{id}/spot
func (h *ParkingSessionHandler) updateSessionSpot(w http.ResponseWriter, r *http.Request) {
	// Get session ID from path
	id := chi.URLParam(r, "id")

	// Parse request body
	var req UpdateSpotRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httputils.ErrorJSON(w, httputils.ErrBadRequest)
		return
	}

	// Update spot
	session, err := h.sessionService.UpdateSessionSpot(r.Context(), id, req.SpotNumber)
	if err != nil {
		h.handleSessionError(w, err, "Failed to update session spot")
		return
	}

	httputils.JSON(w, http.StatusOK, session)
}

// addSessionNotes handles PATCH /sessions/{id}/notes
func (h *ParkingSessionHandler) addSessionNotes(w http.ResponseWriter, r *http.Request) {
	// Get session ID from path
	id := chi.URLParam(r, "id")

	// Parse request body
	var req AddNotesRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httputils.ErrorJSON(w, httputils.ErrBadRequest)
		return
	}

	// Add notes
	session, err := h.sessionService.AddSessionNotes(r.Context(), id, req.Notes)
	if err != nil {
		h.handleSessionError(w, err, "Failed to add session notes")
		return
	}

	httputils.JSON(w, http.StatusOK, session)
}

// markSessionAsPaid handles PATCH /sessions/{id}/pay
func (h *ParkingSessionHandler) markSessionAsPaid(w http.ResponseWriter, r *http.Request) {
	// Get session ID from path
	id := chi.URLParam(r, "id")

	// Mark as paid
	session, err := h.sessionService.MarkSessionAsPaid(r.Context(), id)
	if err != nil {
		h.handleSessionError(w, err, "Failed to mark session as paid")
		return
	}

	httputils.JSON(w, http.StatusOK, session)
}

// handleSessionError handles common session errors
func (h *ParkingSessionHandler) handleSessionError(w http.ResponseWriter, err error, logMessage string) {
	switch err {
	case domain.ErrSessionNotFound:
		httputils.ErrorJSON(w, httputils.ErrNotFound)
	case domain.ErrSessionAlreadyEnded:
		httputils.ErrorJSON(w, httputils.NewError(http.StatusBadRequest, "SESSION_ALREADY_ENDED", "Session is already ended or cancelled", nil))
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
package api

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/jonbelaire/repotown/packages/go-core/httputils"
	"github.com/jonbelaire/repotown/packages/go-core/logging"
	"github.com/jonbelaire/repotown/services/treasury/internal/domain"
	"github.com/jonbelaire/repotown/services/treasury/internal/service"
)

// TaxpayerHandler manages taxpayer-related HTTP requests
type TaxpayerHandler struct {
	taxpayerService service.TaxpayerService
	logger          logging.Logger
}

// NewTaxpayerHandler creates a new taxpayer handler
func NewTaxpayerHandler(taxpayerService service.TaxpayerService, logger logging.Logger) *TaxpayerHandler {
	return &TaxpayerHandler{
		taxpayerService: taxpayerService,
		logger:          logger,
	}
}

// RegisterRoutes registers taxpayer routes
func (h *TaxpayerHandler) RegisterRoutes(r chi.Router) {
	r.Route("/taxpayers", func(r chi.Router) {
		r.Get("/", h.listTaxpayers)
		r.Post("/", h.createTaxpayer)
		r.Get("/search", h.searchTaxpayers)
		r.Get("/type/{type}", h.listTaxpayersByType)
		r.Get("/status/{status}", h.listTaxpayersByStatus)
		r.Get("/industry/{industry}", h.listBusinessesByIndustry)
		r.Get("/identifier/{identifier}", h.getTaxpayerByIdentifier)
		r.Get("/{id}", h.getTaxpayer)
		r.Patch("/{id}/status", h.updateTaxpayerStatus)
		r.Patch("/{id}/contact", h.updateTaxpayerContact)
		r.Patch("/{id}/address", h.updateTaxpayerAddress)
		r.Patch("/{id}/business", h.updateBusinessInfo)
		r.Post("/{id}/exemptions", h.addExemption)
		r.Delete("/{id}/exemptions/{code}", h.removeExemption)
	})
}

// CreateTaxpayerRequest defines the request body for creating a taxpayer
type CreateTaxpayerRequest struct {
	Type           domain.TaxpayerType `json:"type" validate:"required"`
	Name           string              `json:"name" validate:"required"`
	TaxIdentifier  string              `json:"tax_identifier" validate:"required"`
	ContactEmail   string              `json:"contact_email" validate:"required,email"`
	ContactPhone   string              `json:"contact_phone" validate:"required"`
	Address        domain.Address      `json:"address" validate:"required"`
	ExemptionCodes []string            `json:"exemption_codes,omitempty"`
	AnnualRevenue  int64               `json:"annual_revenue,omitempty"`
	BusinessType   string              `json:"business_type,omitempty"`
	Industry       string              `json:"industry,omitempty"`
}

// UpdateTaxpayerStatusRequest defines the request for updating taxpayer status
type UpdateTaxpayerStatusRequest struct {
	Status domain.TaxpayerStatus `json:"status" validate:"required"`
}

// UpdateTaxpayerContactRequest defines the request for updating taxpayer contact
type UpdateTaxpayerContactRequest struct {
	Email string `json:"email" validate:"required,email"`
	Phone string `json:"phone" validate:"required"`
}

// UpdateBusinessInfoRequest defines the request for updating business info
type UpdateBusinessInfoRequest struct {
	AnnualRevenue int64  `json:"annual_revenue" validate:"required,gte=0"`
	BusinessType  string `json:"business_type" validate:"required"`
	Industry      string `json:"industry" validate:"required"`
}

// AddExemptionRequest defines the request for adding an exemption
type AddExemptionRequest struct {
	Code string `json:"code" validate:"required"`
}

// listTaxpayers handles GET /taxpayers
func (h *TaxpayerHandler) listTaxpayers(w http.ResponseWriter, r *http.Request) {
	// Parse pagination parameters
	limit, offset := getPaginationParams(r)

	// Get taxpayers
	taxpayers, err := h.taxpayerService.ListTaxpayers(r.Context(), limit, offset)
	if err != nil {
		h.logger.Error("Failed to list taxpayers", "error", err)
		httputils.ErrorJSON(w, httputils.ErrInternal)
		return
	}

	httputils.JSON(w, http.StatusOK, taxpayers)
}

// createTaxpayer handles POST /taxpayers
func (h *TaxpayerHandler) createTaxpayer(w http.ResponseWriter, r *http.Request) {
	// Parse request body
	var req CreateTaxpayerRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httputils.ErrorJSON(w, httputils.ErrBadRequest)
		return
	}

	// Create taxpayer
	taxpayer, err := h.taxpayerService.CreateTaxpayer(
		r.Context(),
		req.Type,
		req.Name,
		req.TaxIdentifier,
		req.ContactEmail,
		req.ContactPhone,
		req.Address,
		req.ExemptionCodes,
		req.AnnualRevenue,
		req.BusinessType,
		req.Industry,
	)
	if err != nil {
		if err == domain.ErrTaxpayerExists {
			httputils.ErrorJSON(w, httputils.NewError(http.StatusBadRequest, "TAXPAYER_EXISTS", "Taxpayer with this identifier already exists", nil))
			return
		}
		h.logger.Error("Failed to create taxpayer", "error", err)
		httputils.ErrorJSON(w, httputils.ErrInternal)
		return
	}

	httputils.JSON(w, http.StatusCreated, taxpayer)
}

// getTaxpayer handles GET /taxpayers/{id}
func (h *TaxpayerHandler) getTaxpayer(w http.ResponseWriter, r *http.Request) {
	// Get taxpayer ID from path
	id := chi.URLParam(r, "id")

	// Get taxpayer
	taxpayer, err := h.taxpayerService.GetTaxpayer(r.Context(), id)
	if err != nil {
		if err == domain.ErrTaxpayerNotFound {
			httputils.ErrorJSON(w, httputils.ErrNotFound)
			return
		}
		h.logger.Error("Failed to get taxpayer", "error", err)
		httputils.ErrorJSON(w, httputils.ErrInternal)
		return
	}

	httputils.JSON(w, http.StatusOK, taxpayer)
}

// getTaxpayerByIdentifier handles GET /taxpayers/identifier/{identifier}
func (h *TaxpayerHandler) getTaxpayerByIdentifier(w http.ResponseWriter, r *http.Request) {
	// Get tax identifier from path
	identifier := chi.URLParam(r, "identifier")

	// Get taxpayer
	taxpayer, err := h.taxpayerService.GetTaxpayerByTaxIdentifier(r.Context(), identifier)
	if err != nil {
		if err == domain.ErrTaxpayerNotFound {
			httputils.ErrorJSON(w, httputils.ErrNotFound)
			return
		}
		h.logger.Error("Failed to get taxpayer by identifier", "error", err)
		httputils.ErrorJSON(w, httputils.ErrInternal)
		return
	}

	httputils.JSON(w, http.StatusOK, taxpayer)
}

// listTaxpayersByType handles GET /taxpayers/type/{type}
func (h *TaxpayerHandler) listTaxpayersByType(w http.ResponseWriter, r *http.Request) {
	// Get taxpayer type from path
	typeStr := chi.URLParam(r, "type")
	taxpayerType := domain.TaxpayerType(typeStr)

	// Validate taxpayer type
	validTypes := map[domain.TaxpayerType]bool{
		domain.TaxpayerTypeIndividual: true,
		domain.TaxpayerTypeBusiness:   true,
		domain.TaxpayerTypeNonProfit:  true,
		domain.TaxpayerTypeGovernment: true,
	}
	if !validTypes[taxpayerType] {
		httputils.ErrorJSON(w, httputils.NewError(http.StatusBadRequest, "INVALID_TAXPAYER_TYPE", "Invalid taxpayer type", nil))
		return
	}

	// Get taxpayers by type
	taxpayers, err := h.taxpayerService.ListTaxpayersByType(r.Context(), taxpayerType)
	if err != nil {
		h.logger.Error("Failed to list taxpayers by type", "error", err)
		httputils.ErrorJSON(w, httputils.ErrInternal)
		return
	}

	httputils.JSON(w, http.StatusOK, taxpayers)
}

// listTaxpayersByStatus handles GET /taxpayers/status/{status}
func (h *TaxpayerHandler) listTaxpayersByStatus(w http.ResponseWriter, r *http.Request) {
	// Get taxpayer status from path
	statusStr := chi.URLParam(r, "status")
	taxpayerStatus := domain.TaxpayerStatus(statusStr)

	// Validate taxpayer status
	validStatuses := map[domain.TaxpayerStatus]bool{
		domain.TaxpayerStatusActive:     true,
		domain.TaxpayerStatusInactive:   true,
		domain.TaxpayerStatusExempt:     true,
		domain.TaxpayerStatusDelinquent: true,
	}
	if !validStatuses[taxpayerStatus] {
		httputils.ErrorJSON(w, httputils.NewError(http.StatusBadRequest, "INVALID_TAXPAYER_STATUS", "Invalid taxpayer status", nil))
		return
	}

	// Get taxpayers by status
	taxpayers, err := h.taxpayerService.ListTaxpayersByStatus(r.Context(), taxpayerStatus)
	if err != nil {
		h.logger.Error("Failed to list taxpayers by status", "error", err)
		httputils.ErrorJSON(w, httputils.ErrInternal)
		return
	}

	httputils.JSON(w, http.StatusOK, taxpayers)
}

// listBusinessesByIndustry handles GET /taxpayers/industry/{industry}
func (h *TaxpayerHandler) listBusinessesByIndustry(w http.ResponseWriter, r *http.Request) {
	// Get industry from path
	industry := chi.URLParam(r, "industry")

	// Get businesses by industry
	taxpayers, err := h.taxpayerService.ListBusinessesByIndustry(r.Context(), industry)
	if err != nil {
		h.logger.Error("Failed to list businesses by industry", "error", err)
		httputils.ErrorJSON(w, httputils.ErrInternal)
		return
	}

	httputils.JSON(w, http.StatusOK, taxpayers)
}

// searchTaxpayers handles GET /taxpayers/search
func (h *TaxpayerHandler) searchTaxpayers(w http.ResponseWriter, r *http.Request) {
	// Get search query from query parameter
	query := r.URL.Query().Get("q")
	if query == "" {
		httputils.ErrorJSON(w, httputils.NewError(http.StatusBadRequest, "MISSING_QUERY", "Search query parameter is required", nil))
		return
	}

	// Parse limit
	limit, _ := getPaginationParams(r)

	// Search taxpayers
	taxpayers, err := h.taxpayerService.SearchTaxpayers(r.Context(), query, limit)
	if err != nil {
		h.logger.Error("Failed to search taxpayers", "error", err)
		httputils.ErrorJSON(w, httputils.ErrInternal)
		return
	}

	httputils.JSON(w, http.StatusOK, taxpayers)
}

// updateTaxpayerStatus handles PATCH /taxpayers/{id}/status
func (h *TaxpayerHandler) updateTaxpayerStatus(w http.ResponseWriter, r *http.Request) {
	// Get taxpayer ID from path
	id := chi.URLParam(r, "id")

	// Parse request body
	var req UpdateTaxpayerStatusRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httputils.ErrorJSON(w, httputils.ErrBadRequest)
		return
	}

	// Validate taxpayer status
	validStatuses := map[domain.TaxpayerStatus]bool{
		domain.TaxpayerStatusActive:     true,
		domain.TaxpayerStatusInactive:   true,
		domain.TaxpayerStatusExempt:     true,
		domain.TaxpayerStatusDelinquent: true,
	}
	if !validStatuses[req.Status] {
		httputils.ErrorJSON(w, httputils.NewError(http.StatusBadRequest, "INVALID_TAXPAYER_STATUS", "Invalid taxpayer status", nil))
		return
	}

	// Update taxpayer status
	taxpayer, err := h.taxpayerService.UpdateTaxpayerStatus(r.Context(), id, req.Status)
	if err != nil {
		if err == domain.ErrTaxpayerNotFound {
			httputils.ErrorJSON(w, httputils.ErrNotFound)
			return
		}
		h.logger.Error("Failed to update taxpayer status", "error", err)
		httputils.ErrorJSON(w, httputils.ErrInternal)
		return
	}

	httputils.JSON(w, http.StatusOK, taxpayer)
}

// updateTaxpayerContact handles PATCH /taxpayers/{id}/contact
func (h *TaxpayerHandler) updateTaxpayerContact(w http.ResponseWriter, r *http.Request) {
	// Get taxpayer ID from path
	id := chi.URLParam(r, "id")

	// Parse request body
	var req UpdateTaxpayerContactRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httputils.ErrorJSON(w, httputils.ErrBadRequest)
		return
	}

	// Update taxpayer contact
	taxpayer, err := h.taxpayerService.UpdateTaxpayerContact(r.Context(), id, req.Email, req.Phone)
	if err != nil {
		if err == domain.ErrTaxpayerNotFound {
			httputils.ErrorJSON(w, httputils.ErrNotFound)
			return
		}
		h.logger.Error("Failed to update taxpayer contact", "error", err)
		httputils.ErrorJSON(w, httputils.ErrInternal)
		return
	}

	httputils.JSON(w, http.StatusOK, taxpayer)
}

// updateTaxpayerAddress handles PATCH /taxpayers/{id}/address
func (h *TaxpayerHandler) updateTaxpayerAddress(w http.ResponseWriter, r *http.Request) {
	// Get taxpayer ID from path
	id := chi.URLParam(r, "id")

	// Parse request body
	var req domain.Address
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httputils.ErrorJSON(w, httputils.ErrBadRequest)
		return
	}

	// Update taxpayer address
	taxpayer, err := h.taxpayerService.UpdateTaxpayerAddress(r.Context(), id, req)
	if err != nil {
		if err == domain.ErrTaxpayerNotFound {
			httputils.ErrorJSON(w, httputils.ErrNotFound)
			return
		}
		h.logger.Error("Failed to update taxpayer address", "error", err)
		httputils.ErrorJSON(w, httputils.ErrInternal)
		return
	}

	httputils.JSON(w, http.StatusOK, taxpayer)
}

// updateBusinessInfo handles PATCH /taxpayers/{id}/business
func (h *TaxpayerHandler) updateBusinessInfo(w http.ResponseWriter, r *http.Request) {
	// Get taxpayer ID from path
	id := chi.URLParam(r, "id")

	// Parse request body
	var req UpdateBusinessInfoRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httputils.ErrorJSON(w, httputils.ErrBadRequest)
		return
	}

	// Update business info
	taxpayer, err := h.taxpayerService.UpdateBusinessInfo(r.Context(), id, req.AnnualRevenue, req.BusinessType, req.Industry)
	if err != nil {
		if err == domain.ErrTaxpayerNotFound {
			httputils.ErrorJSON(w, httputils.ErrNotFound)
			return
		}
		h.logger.Error("Failed to update business info", "error", err)
		httputils.ErrorJSON(w, httputils.ErrInternal)
		return
	}

	httputils.JSON(w, http.StatusOK, taxpayer)
}

// addExemption handles POST /taxpayers/{id}/exemptions
func (h *TaxpayerHandler) addExemption(w http.ResponseWriter, r *http.Request) {
	// Get taxpayer ID from path
	id := chi.URLParam(r, "id")

	// Parse request body
	var req AddExemptionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httputils.ErrorJSON(w, httputils.ErrBadRequest)
		return
	}

	// Add exemption
	taxpayer, err := h.taxpayerService.AddExemption(r.Context(), id, req.Code)
	if err != nil {
		if err == domain.ErrTaxpayerNotFound {
			httputils.ErrorJSON(w, httputils.ErrNotFound)
			return
		}
		h.logger.Error("Failed to add exemption", "error", err)
		httputils.ErrorJSON(w, httputils.ErrInternal)
		return
	}

	httputils.JSON(w, http.StatusOK, taxpayer)
}

// removeExemption handles DELETE /taxpayers/{id}/exemptions/{code}
func (h *TaxpayerHandler) removeExemption(w http.ResponseWriter, r *http.Request) {
	// Get taxpayer ID and exemption code from path
	id := chi.URLParam(r, "id")
	code := chi.URLParam(r, "code")

	// Remove exemption
	taxpayer, err := h.taxpayerService.RemoveExemption(r.Context(), id, code)
	if err != nil {
		if err == domain.ErrTaxpayerNotFound {
			httputils.ErrorJSON(w, httputils.ErrNotFound)
			return
		}
		h.logger.Error("Failed to remove exemption", "error", err)
		httputils.ErrorJSON(w, httputils.ErrInternal)
		return
	}

	httputils.JSON(w, http.StatusOK, taxpayer)
}
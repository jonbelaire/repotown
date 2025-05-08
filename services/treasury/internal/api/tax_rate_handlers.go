package api

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/jonbelaire/repotown/packages/go-core/httputils"
	"github.com/jonbelaire/repotown/packages/go-core/logging"
	"github.com/jonbelaire/repotown/services/treasury/internal/domain"
	"github.com/jonbelaire/repotown/services/treasury/internal/service"
)

// TaxRateHandler manages tax rate-related HTTP requests
type TaxRateHandler struct {
	taxRateService service.TaxRateService
	logger         logging.Logger
}

// NewTaxRateHandler creates a new tax rate handler
func NewTaxRateHandler(taxRateService service.TaxRateService, logger logging.Logger) *TaxRateHandler {
	return &TaxRateHandler{
		taxRateService: taxRateService,
		logger:         logger,
	}
}

// RegisterRoutes registers tax rate routes
func (h *TaxRateHandler) RegisterRoutes(r chi.Router) {
	r.Route("/tax-rates", func(r chi.Router) {
		r.Get("/", h.listTaxRates)
		r.Post("/", h.createTaxRate)
		r.Get("/active", h.listActiveTaxRates)
		r.Get("/type/{type}", h.listTaxRatesByType)
		r.Get("/jurisdiction/{jurisdictionCode}", h.listTaxRatesByJurisdiction)
		r.Get("/calculate/income", h.calculateIncomeTax)
		r.Get("/calculate/sales", h.calculateSalesTax)
		r.Get("/{id}", h.getTaxRate)
		r.Put("/{id}", h.updateTaxRate)
		r.Post("/{id}/activate", h.activateTaxRate)
		r.Post("/{id}/deactivate", h.deactivateTaxRate)
		r.Post("/{id}/archive", h.archiveTaxRate)
	})
}

// CreateTaxRateRequest defines the request body for creating a tax rate
type CreateTaxRateRequest struct {
	TaxType           domain.TaxType        `json:"tax_type" validate:"required"`
	Name              string                `json:"name" validate:"required"`
	Description       string                `json:"description" validate:"required"`
	Rate              float64               `json:"rate" validate:"required,gte=0,lte=1"`
	BracketType       domain.TaxBracketType `json:"bracket_type" validate:"required"`
	MinAmount         int64                 `json:"min_amount"`
	MaxAmount         int64                 `json:"max_amount"`
	Category          string                `json:"category,omitempty"`
	JurisdictionCode  string                `json:"jurisdiction_code" validate:"required"`
	EffectiveDate     time.Time             `json:"effective_date" validate:"required"`
	ExpirationDate    *time.Time            `json:"expiration_date,omitempty"`
}

// UpdateTaxRateRequest defines the request body for updating a tax rate
type UpdateTaxRateRequest struct {
	Name           string     `json:"name" validate:"required"`
	Description    string     `json:"description" validate:"required"`
	Rate           float64    `json:"rate" validate:"required,gte=0,lte=1"`
	Category       string     `json:"category,omitempty"`
	EffectiveDate  time.Time  `json:"effective_date" validate:"required"`
	ExpirationDate *time.Time `json:"expiration_date,omitempty"`
}

// listTaxRates handles GET /tax-rates
func (h *TaxRateHandler) listTaxRates(w http.ResponseWriter, r *http.Request) {
	// Parse pagination parameters
	limit, offset := getPaginationParams(r)

	// Get tax rates
	taxRates, err := h.taxRateService.ListTaxRates(r.Context(), limit, offset)
	if err != nil {
		h.logger.Error("Failed to list tax rates", "error", err)
		httputils.ErrorJSON(w, httputils.ErrInternal)
		return
	}

	httputils.JSON(w, http.StatusOK, taxRates)
}

// createTaxRate handles POST /tax-rates
func (h *TaxRateHandler) createTaxRate(w http.ResponseWriter, r *http.Request) {
	// Parse request body
	var req CreateTaxRateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httputils.ErrorJSON(w, httputils.ErrBadRequest)
		return
	}

	// Create tax rate
	taxRate, err := h.taxRateService.CreateTaxRate(
		r.Context(),
		req.TaxType,
		req.Name,
		req.Description,
		req.Rate,
		req.BracketType,
		req.MinAmount,
		req.MaxAmount,
		req.Category,
		req.JurisdictionCode,
		req.EffectiveDate,
		req.ExpirationDate,
	)
	if err != nil {
		if err == domain.ErrInvalidTaxRate {
			httputils.ErrorJSON(w, httputils.NewError(http.StatusBadRequest, "INVALID_TAX_RATE", "Invalid tax rate value", nil))
			return
		}
		h.logger.Error("Failed to create tax rate", "error", err)
		httputils.ErrorJSON(w, httputils.ErrInternal)
		return
	}

	httputils.JSON(w, http.StatusCreated, taxRate)
}

// getTaxRate handles GET /tax-rates/{id}
func (h *TaxRateHandler) getTaxRate(w http.ResponseWriter, r *http.Request) {
	// Get tax rate ID from path
	id := chi.URLParam(r, "id")

	// Get tax rate
	taxRate, err := h.taxRateService.GetTaxRate(r.Context(), id)
	if err != nil {
		if err == domain.ErrTaxRateNotFound {
			httputils.ErrorJSON(w, httputils.ErrNotFound)
			return
		}
		h.logger.Error("Failed to get tax rate", "error", err)
		httputils.ErrorJSON(w, httputils.ErrInternal)
		return
	}

	httputils.JSON(w, http.StatusOK, taxRate)
}

// listActiveTaxRates handles GET /tax-rates/active
func (h *TaxRateHandler) listActiveTaxRates(w http.ResponseWriter, r *http.Request) {
	// Get active tax rates
	taxRates, err := h.taxRateService.ListActiveTaxRates(r.Context())
	if err != nil {
		h.logger.Error("Failed to list active tax rates", "error", err)
		httputils.ErrorJSON(w, httputils.ErrInternal)
		return
	}

	httputils.JSON(w, http.StatusOK, taxRates)
}

// listTaxRatesByType handles GET /tax-rates/type/{type}
func (h *TaxRateHandler) listTaxRatesByType(w http.ResponseWriter, r *http.Request) {
	// Get tax type from path
	taxTypeStr := chi.URLParam(r, "type")
	taxType := domain.TaxType(taxTypeStr)

	// Validate tax type
	validTypes := map[domain.TaxType]bool{
		domain.TaxTypeIncome:   true,
		domain.TaxTypeSales:    true,
		domain.TaxTypeProperty: true,
		domain.TaxTypeBusiness: true,
		domain.TaxTypeExcise:   true,
	}
	if !validTypes[taxType] {
		httputils.ErrorJSON(w, httputils.NewError(http.StatusBadRequest, "INVALID_TAX_TYPE", "Invalid tax type", nil))
		return
	}

	// Get tax rates by type
	taxRates, err := h.taxRateService.ListTaxRatesByType(r.Context(), taxType)
	if err != nil {
		h.logger.Error("Failed to list tax rates by type", "error", err)
		httputils.ErrorJSON(w, httputils.ErrInternal)
		return
	}

	httputils.JSON(w, http.StatusOK, taxRates)
}

// listTaxRatesByJurisdiction handles GET /tax-rates/jurisdiction/{jurisdictionCode}
func (h *TaxRateHandler) listTaxRatesByJurisdiction(w http.ResponseWriter, r *http.Request) {
	// Get jurisdiction code from path
	jurisdictionCode := chi.URLParam(r, "jurisdictionCode")

	// Get tax rates by jurisdiction
	taxRates, err := h.taxRateService.ListTaxRatesByJurisdiction(r.Context(), jurisdictionCode)
	if err != nil {
		h.logger.Error("Failed to list tax rates by jurisdiction", "error", err)
		httputils.ErrorJSON(w, httputils.ErrInternal)
		return
	}

	httputils.JSON(w, http.StatusOK, taxRates)
}

// calculateIncomeTax handles GET /tax-rates/calculate/income
func (h *TaxRateHandler) calculateIncomeTax(w http.ResponseWriter, r *http.Request) {
	// Parse amount
	amountStr := r.URL.Query().Get("amount")
	if amountStr == "" {
		httputils.ErrorJSON(w, httputils.NewError(http.StatusBadRequest, "MISSING_AMOUNT", "Amount parameter is required", nil))
		return
	}
	amount, err := strconv.ParseInt(amountStr, 10, 64)
	if err != nil || amount < 0 {
		httputils.ErrorJSON(w, httputils.NewError(http.StatusBadRequest, "INVALID_AMOUNT", "Invalid amount value", nil))
		return
	}

	// Get jurisdiction code
	jurisdictionCode := r.URL.Query().Get("jurisdiction")
	if jurisdictionCode == "" {
		jurisdictionCode = "CITYTOWN" // Default jurisdiction
	}

	// Calculate income tax
	taxAmount, err := h.taxRateService.CalculateIncomeTax(r.Context(), amount, jurisdictionCode)
	if err != nil {
		h.logger.Error("Failed to calculate income tax", "error", err)
		httputils.ErrorJSON(w, httputils.ErrInternal)
		return
	}

	httputils.JSON(w, http.StatusOK, map[string]interface{}{
		"amount":           amount,
		"tax_amount":       taxAmount,
		"jurisdiction_code": jurisdictionCode,
	})
}

// calculateSalesTax handles GET /tax-rates/calculate/sales
func (h *TaxRateHandler) calculateSalesTax(w http.ResponseWriter, r *http.Request) {
	// Parse amount
	amountStr := r.URL.Query().Get("amount")
	if amountStr == "" {
		httputils.ErrorJSON(w, httputils.NewError(http.StatusBadRequest, "MISSING_AMOUNT", "Amount parameter is required", nil))
		return
	}
	amount, err := strconv.ParseInt(amountStr, 10, 64)
	if err != nil || amount < 0 {
		httputils.ErrorJSON(w, httputils.NewError(http.StatusBadRequest, "INVALID_AMOUNT", "Invalid amount value", nil))
		return
	}

	// Get optional parameters
	category := r.URL.Query().Get("category")
	jurisdictionCode := r.URL.Query().Get("jurisdiction")
	if jurisdictionCode == "" {
		jurisdictionCode = "CITYTOWN" // Default jurisdiction
	}

	// Calculate sales tax
	taxAmount, err := h.taxRateService.CalculateSalesTax(r.Context(), amount, category, jurisdictionCode)
	if err != nil {
		h.logger.Error("Failed to calculate sales tax", "error", err)
		httputils.ErrorJSON(w, httputils.ErrInternal)
		return
	}

	httputils.JSON(w, http.StatusOK, map[string]interface{}{
		"amount":           amount,
		"tax_amount":       taxAmount,
		"category":         category,
		"jurisdiction_code": jurisdictionCode,
	})
}

// updateTaxRate handles PUT /tax-rates/{id}
func (h *TaxRateHandler) updateTaxRate(w http.ResponseWriter, r *http.Request) {
	// Get tax rate ID from path
	id := chi.URLParam(r, "id")

	// Parse request body
	var req UpdateTaxRateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httputils.ErrorJSON(w, httputils.ErrBadRequest)
		return
	}

	// Update tax rate
	taxRate, err := h.taxRateService.UpdateTaxRate(
		r.Context(),
		id,
		req.Name,
		req.Description,
		req.Rate,
		req.Category,
		req.EffectiveDate,
		req.ExpirationDate,
	)
	if err != nil {
		if err == domain.ErrTaxRateNotFound {
			httputils.ErrorJSON(w, httputils.ErrNotFound)
			return
		}
		if err == domain.ErrInvalidTaxRate {
			httputils.ErrorJSON(w, httputils.NewError(http.StatusBadRequest, "INVALID_TAX_RATE", "Invalid tax rate value", nil))
			return
		}
		h.logger.Error("Failed to update tax rate", "error", err)
		httputils.ErrorJSON(w, httputils.ErrInternal)
		return
	}

	httputils.JSON(w, http.StatusOK, taxRate)
}

// activateTaxRate handles POST /tax-rates/{id}/activate
func (h *TaxRateHandler) activateTaxRate(w http.ResponseWriter, r *http.Request) {
	// Get tax rate ID from path
	id := chi.URLParam(r, "id")

	// Activate tax rate
	taxRate, err := h.taxRateService.ActivateTaxRate(r.Context(), id)
	if err != nil {
		if err == domain.ErrTaxRateNotFound {
			httputils.ErrorJSON(w, httputils.ErrNotFound)
			return
		}
		h.logger.Error("Failed to activate tax rate", "error", err)
		httputils.ErrorJSON(w, httputils.ErrInternal)
		return
	}

	httputils.JSON(w, http.StatusOK, taxRate)
}

// deactivateTaxRate handles POST /tax-rates/{id}/deactivate
func (h *TaxRateHandler) deactivateTaxRate(w http.ResponseWriter, r *http.Request) {
	// Get tax rate ID from path
	id := chi.URLParam(r, "id")

	// Deactivate tax rate
	taxRate, err := h.taxRateService.DeactivateTaxRate(r.Context(), id)
	if err != nil {
		if err == domain.ErrTaxRateNotFound {
			httputils.ErrorJSON(w, httputils.ErrNotFound)
			return
		}
		h.logger.Error("Failed to deactivate tax rate", "error", err)
		httputils.ErrorJSON(w, httputils.ErrInternal)
		return
	}

	httputils.JSON(w, http.StatusOK, taxRate)
}

// archiveTaxRate handles POST /tax-rates/{id}/archive
func (h *TaxRateHandler) archiveTaxRate(w http.ResponseWriter, r *http.Request) {
	// Get tax rate ID from path
	id := chi.URLParam(r, "id")

	// Archive tax rate
	taxRate, err := h.taxRateService.ArchiveTaxRate(r.Context(), id)
	if err != nil {
		if err == domain.ErrTaxRateNotFound {
			httputils.ErrorJSON(w, httputils.ErrNotFound)
			return
		}
		h.logger.Error("Failed to archive tax rate", "error", err)
		httputils.ErrorJSON(w, httputils.ErrInternal)
		return
	}

	httputils.JSON(w, http.StatusOK, taxRate)
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
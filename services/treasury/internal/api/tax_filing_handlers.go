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

// TaxFilingHandler manages tax filing-related HTTP requests
type TaxFilingHandler struct {
	taxFilingService service.TaxFilingService
	logger           logging.Logger
}

// NewTaxFilingHandler creates a new tax filing handler
func NewTaxFilingHandler(taxFilingService service.TaxFilingService, logger logging.Logger) *TaxFilingHandler {
	return &TaxFilingHandler{
		taxFilingService: taxFilingService,
		logger:           logger,
	}
}

// RegisterRoutes registers tax filing routes
func (h *TaxFilingHandler) RegisterRoutes(r chi.Router) {
	r.Route("/filings", func(r chi.Router) {
		r.Get("/", h.listTaxFilings)
		r.Post("/", h.createTaxFiling)
		r.Get("/status/{status}", h.listTaxFilingsByStatus)
		r.Get("/taxpayer/{taxpayerID}", h.listTaxFilingsByTaxpayer)
		r.Get("/period/{year}/{period}", h.listTaxFilingsByPeriod)
		r.Get("/overdue", h.listOverdueTaxFilings)
		r.Get("/recent", h.listRecentlySubmittedFilings)
		r.Get("/{id}", h.getTaxFiling)
		r.Patch("/{id}/amounts", h.updateTaxFilingAmounts)
		r.Post("/{id}/deductions", h.addDeduction)
		r.Post("/{id}/credits", h.addCredit)
		r.Post("/{id}/submit", h.submitTaxFiling)
		r.Post("/{id}/process", h.processTaxFiling)
		r.Post("/{id}/accept", h.acceptTaxFiling)
		r.Post("/{id}/reject", h.rejectTaxFiling)
		r.Post("/{id}/amend", h.amendTaxFiling)
	})
}

// CreateTaxFilingRequest defines the request body for creating a tax filing
type CreateTaxFilingRequest struct {
	TaxpayerID   string             `json:"taxpayer_id" validate:"required"`
	TaxYear      int                `json:"tax_year" validate:"required"`
	Period       domain.FilingPeriod `json:"period" validate:"required"`
	PeriodStart  time.Time          `json:"period_start" validate:"required"`
	PeriodEnd    time.Time          `json:"period_end" validate:"required"`
	FilingType   domain.TaxType     `json:"filing_type" validate:"required"`
	DueDate      time.Time          `json:"due_date" validate:"required"`
}

// UpdateTaxFilingAmountsRequest defines the request for updating filing amounts
type UpdateTaxFilingAmountsRequest struct {
	GrossIncome   int64 `json:"gross_income,omitempty"`
	TaxableIncome int64 `json:"taxable_income,omitempty"`
	TotalSales    int64 `json:"total_sales,omitempty"`
	TaxableAmount int64 `json:"taxable_amount" validate:"required"`
}

// AddDeductionRequest defines the request for adding a deduction
type AddDeductionRequest struct {
	Code        string `json:"code" validate:"required"`
	Description string `json:"description" validate:"required"`
	Amount      int64  `json:"amount" validate:"required,gt=0"`
}

// AddCreditRequest defines the request for adding a credit
type AddCreditRequest struct {
	Code        string `json:"code" validate:"required"`
	Description string `json:"description" validate:"required"`
	Amount      int64  `json:"amount" validate:"required,gt=0"`
}

// ProcessTaxFilingRequest defines the request for processing a tax filing
type ProcessTaxFilingRequest struct {
	TaxCalculated int64 `json:"tax_calculated" validate:"required,gte=0"`
}

// RejectTaxFilingRequest defines the request for rejecting a tax filing
type RejectTaxFilingRequest struct {
	Reason string `json:"reason" validate:"required"`
}

// listTaxFilings handles GET /filings
func (h *TaxFilingHandler) listTaxFilings(w http.ResponseWriter, r *http.Request) {
	// Parse pagination parameters
	limit, offset := getPaginationParams(r)

	// Get tax filings
	filings, err := h.taxFilingService.ListTaxFilings(r.Context(), limit, offset)
	if err != nil {
		h.logger.Error("Failed to list tax filings", "error", err)
		httputils.ErrorJSON(w, httputils.ErrInternal)
		return
	}

	httputils.JSON(w, http.StatusOK, filings)
}

// createTaxFiling handles POST /filings
func (h *TaxFilingHandler) createTaxFiling(w http.ResponseWriter, r *http.Request) {
	// Parse request body
	var req CreateTaxFilingRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httputils.ErrorJSON(w, httputils.ErrBadRequest)
		return
	}

	// Validate time period
	if req.PeriodEnd.Before(req.PeriodStart) {
		httputils.ErrorJSON(w, httputils.NewError(http.StatusBadRequest, "INVALID_PERIOD", "Period end must be after period start", nil))
		return
	}

	// Create tax filing
	filing, err := h.taxFilingService.CreateTaxFiling(
		r.Context(),
		req.TaxpayerID,
		req.TaxYear,
		req.Period,
		req.PeriodStart,
		req.PeriodEnd,
		req.FilingType,
		req.DueDate,
	)
	if err != nil {
		if err == domain.ErrInvalidFilingPeriod {
			httputils.ErrorJSON(w, httputils.NewError(http.StatusBadRequest, "INVALID_PERIOD", "Invalid filing period", nil))
			return
		}
		if err == domain.ErrTaxpayerNotFound {
			httputils.ErrorJSON(w, httputils.NewError(http.StatusBadRequest, "TAXPAYER_NOT_FOUND", "Taxpayer not found", nil))
			return
		}
		h.logger.Error("Failed to create tax filing", "error", err)
		httputils.ErrorJSON(w, httputils.ErrInternal)
		return
	}

	httputils.JSON(w, http.StatusCreated, filing)
}

// getTaxFiling handles GET /filings/{id}
func (h *TaxFilingHandler) getTaxFiling(w http.ResponseWriter, r *http.Request) {
	// Get filing ID from path
	id := chi.URLParam(r, "id")

	// Get tax filing
	filing, err := h.taxFilingService.GetTaxFiling(r.Context(), id)
	if err != nil {
		if err == domain.ErrTaxFilingNotFound {
			httputils.ErrorJSON(w, httputils.ErrNotFound)
			return
		}
		h.logger.Error("Failed to get tax filing", "error", err)
		httputils.ErrorJSON(w, httputils.ErrInternal)
		return
	}

	httputils.JSON(w, http.StatusOK, filing)
}

// listTaxFilingsByStatus handles GET /filings/status/{status}
func (h *TaxFilingHandler) listTaxFilingsByStatus(w http.ResponseWriter, r *http.Request) {
	// Get status from path
	statusStr := chi.URLParam(r, "status")
	status := domain.FilingStatus(statusStr)

	// Validate status
	validStatuses := map[domain.FilingStatus]bool{
		domain.FilingStatusDraft:      true,
		domain.FilingStatusSubmitted:  true,
		domain.FilingStatusProcessing: true,
		domain.FilingStatusAccepted:   true,
		domain.FilingStatusRejected:   true,
		domain.FilingStatusAmended:    true,
		domain.FilingStatusAudited:    true,
	}
	if !validStatuses[status] {
		httputils.ErrorJSON(w, httputils.NewError(http.StatusBadRequest, "INVALID_STATUS", "Invalid filing status", nil))
		return
	}

	// Parse pagination parameters
	limit, offset := getPaginationParams(r)

	// Get tax filings by status
	filings, err := h.taxFilingService.ListTaxFilingsByStatus(r.Context(), status, limit, offset)
	if err != nil {
		h.logger.Error("Failed to list tax filings by status", "error", err)
		httputils.ErrorJSON(w, httputils.ErrInternal)
		return
	}

	httputils.JSON(w, http.StatusOK, filings)
}

// listTaxFilingsByTaxpayer handles GET /filings/taxpayer/{taxpayerID}
func (h *TaxFilingHandler) listTaxFilingsByTaxpayer(w http.ResponseWriter, r *http.Request) {
	// Get taxpayer ID from path
	taxpayerID := chi.URLParam(r, "taxpayerID")

	// Parse pagination parameters
	limit, offset := getPaginationParams(r)

	// Get tax filings by taxpayer
	filings, err := h.taxFilingService.ListTaxFilingsByTaxpayer(r.Context(), taxpayerID, limit, offset)
	if err != nil {
		h.logger.Error("Failed to list tax filings by taxpayer", "error", err)
		httputils.ErrorJSON(w, httputils.ErrInternal)
		return
	}

	httputils.JSON(w, http.StatusOK, filings)
}

// listTaxFilingsByPeriod handles GET /filings/period/{year}/{period}
func (h *TaxFilingHandler) listTaxFilingsByPeriod(w http.ResponseWriter, r *http.Request) {
	// Get year and period from path
	yearStr := chi.URLParam(r, "year")
	periodStr := chi.URLParam(r, "period")

	// Parse year
	year, err := strconv.Atoi(yearStr)
	if err != nil {
		httputils.ErrorJSON(w, httputils.NewError(http.StatusBadRequest, "INVALID_YEAR", "Invalid year format", nil))
		return
	}

	// Validate period
	period := domain.FilingPeriod(periodStr)
	validPeriods := map[domain.FilingPeriod]bool{
		domain.FilingPeriodMonthly:    true,
		domain.FilingPeriodQuarterly:  true,
		domain.FilingPeriodSemiAnnual: true,
		domain.FilingPeriodAnnual:     true,
	}
	if !validPeriods[period] {
		httputils.ErrorJSON(w, httputils.NewError(http.StatusBadRequest, "INVALID_PERIOD", "Invalid filing period", nil))
		return
	}

	// Get tax filings by period
	filings, err := h.taxFilingService.ListTaxFilingsByPeriod(r.Context(), year, period)
	if err != nil {
		h.logger.Error("Failed to list tax filings by period", "error", err)
		httputils.ErrorJSON(w, httputils.ErrInternal)
		return
	}

	httputils.JSON(w, http.StatusOK, filings)
}

// listOverdueTaxFilings handles GET /filings/overdue
func (h *TaxFilingHandler) listOverdueTaxFilings(w http.ResponseWriter, r *http.Request) {
	// Get overdue tax filings
	filings, err := h.taxFilingService.ListOverdueTaxFilings(r.Context())
	if err != nil {
		h.logger.Error("Failed to list overdue tax filings", "error", err)
		httputils.ErrorJSON(w, httputils.ErrInternal)
		return
	}

	httputils.JSON(w, http.StatusOK, filings)
}

// listRecentlySubmittedFilings handles GET /filings/recent
func (h *TaxFilingHandler) listRecentlySubmittedFilings(w http.ResponseWriter, r *http.Request) {
	// Parse days parameter
	daysStr := r.URL.Query().Get("days")
	days := 7 // Default to 7 days
	if daysStr != "" {
		if d, err := strconv.Atoi(daysStr); err == nil && d > 0 {
			days = d
		}
	}

	// Get recently submitted tax filings
	filings, err := h.taxFilingService.ListRecentlySubmittedFilings(r.Context(), days)
	if err != nil {
		h.logger.Error("Failed to list recently submitted tax filings", "error", err)
		httputils.ErrorJSON(w, httputils.ErrInternal)
		return
	}

	httputils.JSON(w, http.StatusOK, filings)
}

// updateTaxFilingAmounts handles PATCH /filings/{id}/amounts
func (h *TaxFilingHandler) updateTaxFilingAmounts(w http.ResponseWriter, r *http.Request) {
	// Get filing ID from path
	id := chi.URLParam(r, "id")

	// Parse request body
	var req UpdateTaxFilingAmountsRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httputils.ErrorJSON(w, httputils.ErrBadRequest)
		return
	}

	// Update tax filing amounts
	filing, err := h.taxFilingService.UpdateTaxFilingAmounts(
		r.Context(),
		id,
		req.GrossIncome,
		req.TaxableIncome,
		req.TotalSales,
		req.TaxableAmount,
	)
	if err != nil {
		if err == domain.ErrTaxFilingNotFound {
			httputils.ErrorJSON(w, httputils.ErrNotFound)
			return
		}
		if err == domain.ErrInvalidFilingStatus {
			httputils.ErrorJSON(w, httputils.NewError(http.StatusBadRequest, "INVALID_STATUS", "Filing status does not allow updates", nil))
			return
		}
		h.logger.Error("Failed to update tax filing amounts", "error", err)
		httputils.ErrorJSON(w, httputils.ErrInternal)
		return
	}

	httputils.JSON(w, http.StatusOK, filing)
}

// addDeduction handles POST /filings/{id}/deductions
func (h *TaxFilingHandler) addDeduction(w http.ResponseWriter, r *http.Request) {
	// Get filing ID from path
	id := chi.URLParam(r, "id")

	// Parse request body
	var req AddDeductionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httputils.ErrorJSON(w, httputils.ErrBadRequest)
		return
	}

	// Add deduction
	filing, err := h.taxFilingService.AddDeduction(r.Context(), id, req.Code, req.Description, req.Amount)
	if err != nil {
		if err == domain.ErrTaxFilingNotFound {
			httputils.ErrorJSON(w, httputils.ErrNotFound)
			return
		}
		if err == domain.ErrInvalidFilingStatus {
			httputils.ErrorJSON(w, httputils.NewError(http.StatusBadRequest, "INVALID_STATUS", "Filing status does not allow updates", nil))
			return
		}
		h.logger.Error("Failed to add deduction", "error", err)
		httputils.ErrorJSON(w, httputils.ErrInternal)
		return
	}

	httputils.JSON(w, http.StatusOK, filing)
}

// addCredit handles POST /filings/{id}/credits
func (h *TaxFilingHandler) addCredit(w http.ResponseWriter, r *http.Request) {
	// Get filing ID from path
	id := chi.URLParam(r, "id")

	// Parse request body
	var req AddCreditRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httputils.ErrorJSON(w, httputils.ErrBadRequest)
		return
	}

	// Add credit
	filing, err := h.taxFilingService.AddCredit(r.Context(), id, req.Code, req.Description, req.Amount)
	if err != nil {
		if err == domain.ErrTaxFilingNotFound {
			httputils.ErrorJSON(w, httputils.ErrNotFound)
			return
		}
		if err == domain.ErrInvalidFilingStatus {
			httputils.ErrorJSON(w, httputils.NewError(http.StatusBadRequest, "INVALID_STATUS", "Filing status does not allow updates", nil))
			return
		}
		h.logger.Error("Failed to add credit", "error", err)
		httputils.ErrorJSON(w, httputils.ErrInternal)
		return
	}

	httputils.JSON(w, http.StatusOK, filing)
}

// submitTaxFiling handles POST /filings/{id}/submit
func (h *TaxFilingHandler) submitTaxFiling(w http.ResponseWriter, r *http.Request) {
	// Get filing ID from path
	id := chi.URLParam(r, "id")

	// Submit tax filing
	filing, err := h.taxFilingService.SubmitTaxFiling(r.Context(), id)
	if err != nil {
		if err == domain.ErrTaxFilingNotFound {
			httputils.ErrorJSON(w, httputils.ErrNotFound)
			return
		}
		if err == domain.ErrInvalidFilingStatus {
			httputils.ErrorJSON(w, httputils.NewError(http.StatusBadRequest, "INVALID_STATUS", "Filing has already been submitted", nil))
			return
		}
		h.logger.Error("Failed to submit tax filing", "error", err)
		httputils.ErrorJSON(w, httputils.ErrInternal)
		return
	}

	httputils.JSON(w, http.StatusOK, filing)
}

// processTaxFiling handles POST /filings/{id}/process
func (h *TaxFilingHandler) processTaxFiling(w http.ResponseWriter, r *http.Request) {
	// Get filing ID from path
	id := chi.URLParam(r, "id")

	// Parse request body
	var req ProcessTaxFilingRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httputils.ErrorJSON(w, httputils.ErrBadRequest)
		return
	}

	// Process tax filing
	filing, err := h.taxFilingService.ProcessTaxFiling(r.Context(), id, req.TaxCalculated)
	if err != nil {
		if err == domain.ErrTaxFilingNotFound {
			httputils.ErrorJSON(w, httputils.ErrNotFound)
			return
		}
		if err == domain.ErrInvalidFilingStatus {
			httputils.ErrorJSON(w, httputils.NewError(http.StatusBadRequest, "INVALID_STATUS", "Filing status does not allow processing", nil))
			return
		}
		h.logger.Error("Failed to process tax filing", "error", err)
		httputils.ErrorJSON(w, httputils.ErrInternal)
		return
	}

	httputils.JSON(w, http.StatusOK, filing)
}

// acceptTaxFiling handles POST /filings/{id}/accept
func (h *TaxFilingHandler) acceptTaxFiling(w http.ResponseWriter, r *http.Request) {
	// Get filing ID from path
	id := chi.URLParam(r, "id")

	// Accept tax filing
	filing, err := h.taxFilingService.AcceptTaxFiling(r.Context(), id)
	if err != nil {
		if err == domain.ErrTaxFilingNotFound {
			httputils.ErrorJSON(w, httputils.ErrNotFound)
			return
		}
		if err == domain.ErrInvalidFilingStatus {
			httputils.ErrorJSON(w, httputils.NewError(http.StatusBadRequest, "INVALID_STATUS", "Filing status does not allow acceptance", nil))
			return
		}
		h.logger.Error("Failed to accept tax filing", "error", err)
		httputils.ErrorJSON(w, httputils.ErrInternal)
		return
	}

	httputils.JSON(w, http.StatusOK, filing)
}

// rejectTaxFiling handles POST /filings/{id}/reject
func (h *TaxFilingHandler) rejectTaxFiling(w http.ResponseWriter, r *http.Request) {
	// Get filing ID from path
	id := chi.URLParam(r, "id")

	// Parse request body
	var req RejectTaxFilingRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httputils.ErrorJSON(w, httputils.ErrBadRequest)
		return
	}

	// Reject tax filing
	filing, err := h.taxFilingService.RejectTaxFiling(r.Context(), id, req.Reason)
	if err != nil {
		if err == domain.ErrTaxFilingNotFound {
			httputils.ErrorJSON(w, httputils.ErrNotFound)
			return
		}
		if err == domain.ErrInvalidFilingStatus {
			httputils.ErrorJSON(w, httputils.NewError(http.StatusBadRequest, "INVALID_STATUS", "Filing status does not allow rejection", nil))
			return
		}
		h.logger.Error("Failed to reject tax filing", "error", err)
		httputils.ErrorJSON(w, httputils.ErrInternal)
		return
	}

	httputils.JSON(w, http.StatusOK, filing)
}

// amendTaxFiling handles POST /filings/{id}/amend
func (h *TaxFilingHandler) amendTaxFiling(w http.ResponseWriter, r *http.Request) {
	// Get filing ID from path
	id := chi.URLParam(r, "id")

	// Amend tax filing
	filing, err := h.taxFilingService.AmendTaxFiling(r.Context(), id)
	if err != nil {
		if err == domain.ErrTaxFilingNotFound {
			httputils.ErrorJSON(w, httputils.ErrNotFound)
			return
		}
		h.logger.Error("Failed to amend tax filing", "error", err)
		httputils.ErrorJSON(w, httputils.ErrInternal)
		return
	}

	httputils.JSON(w, http.StatusCreated, filing)
}
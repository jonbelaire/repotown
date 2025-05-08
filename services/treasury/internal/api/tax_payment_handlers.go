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

// TaxPaymentHandler manages tax payment-related HTTP requests
type TaxPaymentHandler struct {
	taxPaymentService service.TaxPaymentService
	logger            logging.Logger
}

// NewTaxPaymentHandler creates a new tax payment handler
func NewTaxPaymentHandler(taxPaymentService service.TaxPaymentService, logger logging.Logger) *TaxPaymentHandler {
	return &TaxPaymentHandler{
		taxPaymentService: taxPaymentService,
		logger:            logger,
	}
}

// RegisterRoutes registers tax payment routes
func (h *TaxPaymentHandler) RegisterRoutes(r chi.Router) {
	r.Route("/payments", func(r chi.Router) {
		r.Get("/", h.listTaxPayments)
		r.Post("/", h.createTaxPayment)
		r.Get("/code/{code}", h.getTaxPaymentByConfirmationCode)
		r.Get("/taxpayer/{taxpayerID}", h.listTaxPaymentsByTaxpayer)
		r.Get("/filing/{filingID}", h.listTaxPaymentsByFiling)
		r.Get("/status/{status}", h.listTaxPaymentsByStatus)
		r.Get("/date-range", h.listTaxPaymentsByDateRange)
		r.Get("/recent", h.listRecentTaxPayments)
		r.Get("/totals/by-type", h.getTotalPaymentsByTaxType)
		r.Get("/{id}", h.getTaxPayment)
		r.Post("/{id}/process", h.processTaxPayment)
		r.Post("/{id}/fail", h.markTaxPaymentAsFailed)
		r.Post("/{id}/refund", h.refundTaxPayment)
		r.Post("/{id}/void", h.voidTaxPayment)
		r.Patch("/{id}/amount", h.updateTaxPaymentAmount)
	})
}

// CreateTaxPaymentRequest defines the request body for creating a tax payment
type CreateTaxPaymentRequest struct {
	TaxpayerID     string               `json:"taxpayer_id" validate:"required"`
	FilingID       string               `json:"filing_id,omitempty"`
	TaxType        domain.TaxType       `json:"tax_type" validate:"required"`
	Amount         int64                `json:"amount" validate:"required,gt=0"`
	PaymentMethod  domain.PaymentMethod `json:"payment_method" validate:"required"`
	PaymentDate    time.Time            `json:"payment_date" validate:"required"`
	Notes          string               `json:"notes,omitempty"`
}

// UpdateTaxPaymentAmountRequest defines the request for updating payment amount
type UpdateTaxPaymentAmountRequest struct {
	Amount int64 `json:"amount" validate:"required,gt=0"`
}

// ReasonRequest defines a request that includes a reason
type ReasonRequest struct {
	Reason string `json:"reason" validate:"required"`
}

// listTaxPayments handles GET /payments
func (h *TaxPaymentHandler) listTaxPayments(w http.ResponseWriter, r *http.Request) {
	// Parse pagination parameters
	limit, offset := getPaginationParams(r)

	// Get tax payments
	payments, err := h.taxPaymentService.ListTaxPayments(r.Context(), limit, offset)
	if err != nil {
		h.logger.Error("Failed to list tax payments", "error", err)
		httputils.ErrorJSON(w, httputils.ErrInternal)
		return
	}

	httputils.JSON(w, http.StatusOK, payments)
}

// createTaxPayment handles POST /payments
func (h *TaxPaymentHandler) createTaxPayment(w http.ResponseWriter, r *http.Request) {
	// Parse request body
	var req CreateTaxPaymentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httputils.ErrorJSON(w, httputils.ErrBadRequest)
		return
	}

	// Create tax payment
	payment, err := h.taxPaymentService.CreateTaxPayment(
		r.Context(),
		req.TaxpayerID,
		req.FilingID,
		req.TaxType,
		req.Amount,
		req.PaymentMethod,
		req.PaymentDate,
		req.Notes,
	)
	if err != nil {
		if err == domain.ErrInvalidPaymentAmount {
			httputils.ErrorJSON(w, httputils.NewError(http.StatusBadRequest, "INVALID_AMOUNT", "Invalid payment amount", nil))
			return
		}
		if err == domain.ErrTaxFilingNotFound {
			httputils.ErrorJSON(w, httputils.NewError(http.StatusBadRequest, "FILING_NOT_FOUND", "Tax filing not found", nil))
			return
		}
		h.logger.Error("Failed to create tax payment", "error", err)
		httputils.ErrorJSON(w, httputils.ErrInternal)
		return
	}

	httputils.JSON(w, http.StatusCreated, payment)
}

// getTaxPayment handles GET /payments/{id}
func (h *TaxPaymentHandler) getTaxPayment(w http.ResponseWriter, r *http.Request) {
	// Get payment ID from path
	id := chi.URLParam(r, "id")

	// Get tax payment
	payment, err := h.taxPaymentService.GetTaxPayment(r.Context(), id)
	if err != nil {
		if err == domain.ErrTaxPaymentNotFound {
			httputils.ErrorJSON(w, httputils.ErrNotFound)
			return
		}
		h.logger.Error("Failed to get tax payment", "error", err)
		httputils.ErrorJSON(w, httputils.ErrInternal)
		return
	}

	httputils.JSON(w, http.StatusOK, payment)
}

// getTaxPaymentByConfirmationCode handles GET /payments/code/{code}
func (h *TaxPaymentHandler) getTaxPaymentByConfirmationCode(w http.ResponseWriter, r *http.Request) {
	// Get confirmation code from path
	code := chi.URLParam(r, "code")

	// Get tax payment by confirmation code
	payment, err := h.taxPaymentService.GetTaxPaymentByConfirmationCode(r.Context(), code)
	if err != nil {
		if err == domain.ErrTaxPaymentNotFound {
			httputils.ErrorJSON(w, httputils.ErrNotFound)
			return
		}
		h.logger.Error("Failed to get tax payment by confirmation code", "error", err)
		httputils.ErrorJSON(w, httputils.ErrInternal)
		return
	}

	httputils.JSON(w, http.StatusOK, payment)
}

// listTaxPaymentsByTaxpayer handles GET /payments/taxpayer/{taxpayerID}
func (h *TaxPaymentHandler) listTaxPaymentsByTaxpayer(w http.ResponseWriter, r *http.Request) {
	// Get taxpayer ID from path
	taxpayerID := chi.URLParam(r, "taxpayerID")

	// Parse pagination parameters
	limit, offset := getPaginationParams(r)

	// Get tax payments by taxpayer
	payments, err := h.taxPaymentService.ListTaxPaymentsByTaxpayer(r.Context(), taxpayerID, limit, offset)
	if err != nil {
		h.logger.Error("Failed to list tax payments by taxpayer", "error", err)
		httputils.ErrorJSON(w, httputils.ErrInternal)
		return
	}

	httputils.JSON(w, http.StatusOK, payments)
}

// listTaxPaymentsByFiling handles GET /payments/filing/{filingID}
func (h *TaxPaymentHandler) listTaxPaymentsByFiling(w http.ResponseWriter, r *http.Request) {
	// Get filing ID from path
	filingID := chi.URLParam(r, "filingID")

	// Get tax payments by filing
	payments, err := h.taxPaymentService.ListTaxPaymentsByFiling(r.Context(), filingID)
	if err != nil {
		h.logger.Error("Failed to list tax payments by filing", "error", err)
		httputils.ErrorJSON(w, httputils.ErrInternal)
		return
	}

	httputils.JSON(w, http.StatusOK, payments)
}

// listTaxPaymentsByStatus handles GET /payments/status/{status}
func (h *TaxPaymentHandler) listTaxPaymentsByStatus(w http.ResponseWriter, r *http.Request) {
	// Get status from path
	statusStr := chi.URLParam(r, "status")
	status := domain.PaymentStatus(statusStr)

	// Validate status
	validStatuses := map[domain.PaymentStatus]bool{
		domain.PaymentStatusPending:   true,
		domain.PaymentStatusCompleted: true,
		domain.PaymentStatusFailed:    true,
		domain.PaymentStatusRefunded:  true,
		domain.PaymentStatusVoided:    true,
	}
	if !validStatuses[status] {
		httputils.ErrorJSON(w, httputils.NewError(http.StatusBadRequest, "INVALID_STATUS", "Invalid payment status", nil))
		return
	}

	// Get tax payments by status
	payments, err := h.taxPaymentService.ListTaxPaymentsByStatus(r.Context(), status)
	if err != nil {
		h.logger.Error("Failed to list tax payments by status", "error", err)
		httputils.ErrorJSON(w, httputils.ErrInternal)
		return
	}

	httputils.JSON(w, http.StatusOK, payments)
}

// listTaxPaymentsByDateRange handles GET /payments/date-range
func (h *TaxPaymentHandler) listTaxPaymentsByDateRange(w http.ResponseWriter, r *http.Request) {
	// Parse start date
	startDateStr := r.URL.Query().Get("start_date")
	if startDateStr == "" {
		httputils.ErrorJSON(w, httputils.NewError(http.StatusBadRequest, "MISSING_START_DATE", "Start date parameter is required", nil))
		return
	}
	startDate, err := time.Parse(time.RFC3339, startDateStr)
	if err != nil {
		httputils.ErrorJSON(w, httputils.NewError(http.StatusBadRequest, "INVALID_START_DATE", "Invalid start date format", nil))
		return
	}

	// Parse end date
	endDateStr := r.URL.Query().Get("end_date")
	if endDateStr == "" {
		httputils.ErrorJSON(w, httputils.NewError(http.StatusBadRequest, "MISSING_END_DATE", "End date parameter is required", nil))
		return
	}
	endDate, err := time.Parse(time.RFC3339, endDateStr)
	if err != nil {
		httputils.ErrorJSON(w, httputils.NewError(http.StatusBadRequest, "INVALID_END_DATE", "Invalid end date format", nil))
		return
	}

	// Validate date range
	if endDate.Before(startDate) {
		httputils.ErrorJSON(w, httputils.NewError(http.StatusBadRequest, "INVALID_DATE_RANGE", "End date must be after start date", nil))
		return
	}

	// Get tax payments by date range
	payments, err := h.taxPaymentService.ListTaxPaymentsByDateRange(r.Context(), startDate, endDate)
	if err != nil {
		h.logger.Error("Failed to list tax payments by date range", "error", err)
		httputils.ErrorJSON(w, httputils.ErrInternal)
		return
	}

	httputils.JSON(w, http.StatusOK, payments)
}

// listRecentTaxPayments handles GET /payments/recent
func (h *TaxPaymentHandler) listRecentTaxPayments(w http.ResponseWriter, r *http.Request) {
	// Parse days parameter
	daysStr := r.URL.Query().Get("days")
	days := 7 // Default to 7 days
	if daysStr != "" {
		if d, err := strconv.Atoi(daysStr); err == nil && d > 0 {
			days = d
		}
	}

	// Get recent tax payments
	payments, err := h.taxPaymentService.ListRecentTaxPayments(r.Context(), days)
	if err != nil {
		h.logger.Error("Failed to list recent tax payments", "error", err)
		httputils.ErrorJSON(w, httputils.ErrInternal)
		return
	}

	httputils.JSON(w, http.StatusOK, payments)
}

// getTotalPaymentsByTaxType handles GET /payments/totals/by-type
func (h *TaxPaymentHandler) getTotalPaymentsByTaxType(w http.ResponseWriter, r *http.Request) {
	// Parse tax type
	taxTypeStr := r.URL.Query().Get("tax_type")
	if taxTypeStr == "" {
		httputils.ErrorJSON(w, httputils.NewError(http.StatusBadRequest, "MISSING_TAX_TYPE", "Tax type parameter is required", nil))
		return
	}
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

	// Parse start date
	startDateStr := r.URL.Query().Get("start_date")
	if startDateStr == "" {
		httputils.ErrorJSON(w, httputils.NewError(http.StatusBadRequest, "MISSING_START_DATE", "Start date parameter is required", nil))
		return
	}
	startDate, err := time.Parse(time.RFC3339, startDateStr)
	if err != nil {
		httputils.ErrorJSON(w, httputils.NewError(http.StatusBadRequest, "INVALID_START_DATE", "Invalid start date format", nil))
		return
	}

	// Parse end date
	endDateStr := r.URL.Query().Get("end_date")
	if endDateStr == "" {
		httputils.ErrorJSON(w, httputils.NewError(http.StatusBadRequest, "MISSING_END_DATE", "End date parameter is required", nil))
		return
	}
	endDate, err := time.Parse(time.RFC3339, endDateStr)
	if err != nil {
		httputils.ErrorJSON(w, httputils.NewError(http.StatusBadRequest, "INVALID_END_DATE", "Invalid end date format", nil))
		return
	}

	// Get total payments by tax type
	total, err := h.taxPaymentService.GetTotalPaymentsByTaxType(r.Context(), taxType, startDate, endDate)
	if err != nil {
		h.logger.Error("Failed to get total payments by tax type", "error", err)
		httputils.ErrorJSON(w, httputils.ErrInternal)
		return
	}

	httputils.JSON(w, http.StatusOK, map[string]interface{}{
		"tax_type":    taxType,
		"start_date":  startDate,
		"end_date":    endDate,
		"total":       total,
	})
}

// processTaxPayment handles POST /payments/{id}/process
func (h *TaxPaymentHandler) processTaxPayment(w http.ResponseWriter, r *http.Request) {
	// Get payment ID from path
	id := chi.URLParam(r, "id")

	// Process tax payment
	payment, err := h.taxPaymentService.ProcessTaxPayment(r.Context(), id)
	if err != nil {
		if err == domain.ErrTaxPaymentNotFound {
			httputils.ErrorJSON(w, httputils.ErrNotFound)
			return
		}
		h.logger.Error("Failed to process tax payment", "error", err)
		httputils.ErrorJSON(w, httputils.ErrInternal)
		return
	}

	httputils.JSON(w, http.StatusOK, payment)
}

// markTaxPaymentAsFailed handles POST /payments/{id}/fail
func (h *TaxPaymentHandler) markTaxPaymentAsFailed(w http.ResponseWriter, r *http.Request) {
	// Get payment ID from path
	id := chi.URLParam(r, "id")

	// Parse request body
	var req ReasonRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httputils.ErrorJSON(w, httputils.ErrBadRequest)
		return
	}

	// Mark tax payment as failed
	payment, err := h.taxPaymentService.MarkTaxPaymentAsFailed(r.Context(), id, req.Reason)
	if err != nil {
		if err == domain.ErrTaxPaymentNotFound {
			httputils.ErrorJSON(w, httputils.ErrNotFound)
			return
		}
		h.logger.Error("Failed to mark tax payment as failed", "error", err)
		httputils.ErrorJSON(w, httputils.ErrInternal)
		return
	}

	httputils.JSON(w, http.StatusOK, payment)
}

// refundTaxPayment handles POST /payments/{id}/refund
func (h *TaxPaymentHandler) refundTaxPayment(w http.ResponseWriter, r *http.Request) {
	// Get payment ID from path
	id := chi.URLParam(r, "id")

	// Parse request body
	var req ReasonRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httputils.ErrorJSON(w, httputils.ErrBadRequest)
		return
	}

	// Refund tax payment
	payment, err := h.taxPaymentService.RefundTaxPayment(r.Context(), id, req.Reason)
	if err != nil {
		if err == domain.ErrTaxPaymentNotFound {
			httputils.ErrorJSON(w, httputils.ErrNotFound)
			return
		}
		h.logger.Error("Failed to refund tax payment", "error", err)
		httputils.ErrorJSON(w, httputils.ErrInternal)
		return
	}

	httputils.JSON(w, http.StatusOK, payment)
}

// voidTaxPayment handles POST /payments/{id}/void
func (h *TaxPaymentHandler) voidTaxPayment(w http.ResponseWriter, r *http.Request) {
	// Get payment ID from path
	id := chi.URLParam(r, "id")

	// Parse request body
	var req ReasonRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httputils.ErrorJSON(w, httputils.ErrBadRequest)
		return
	}

	// Void tax payment
	payment, err := h.taxPaymentService.VoidTaxPayment(r.Context(), id, req.Reason)
	if err != nil {
		if err == domain.ErrTaxPaymentNotFound {
			httputils.ErrorJSON(w, httputils.ErrNotFound)
			return
		}
		h.logger.Error("Failed to void tax payment", "error", err)
		httputils.ErrorJSON(w, httputils.ErrInternal)
		return
	}

	httputils.JSON(w, http.StatusOK, payment)
}

// updateTaxPaymentAmount handles PATCH /payments/{id}/amount
func (h *TaxPaymentHandler) updateTaxPaymentAmount(w http.ResponseWriter, r *http.Request) {
	// Get payment ID from path
	id := chi.URLParam(r, "id")

	// Parse request body
	var req UpdateTaxPaymentAmountRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httputils.ErrorJSON(w, httputils.ErrBadRequest)
		return
	}

	// Update tax payment amount
	payment, err := h.taxPaymentService.UpdateTaxPaymentAmount(r.Context(), id, req.Amount)
	if err != nil {
		if err == domain.ErrTaxPaymentNotFound {
			httputils.ErrorJSON(w, httputils.ErrNotFound)
			return
		}
		if err == domain.ErrInvalidPaymentAmount {
			httputils.ErrorJSON(w, httputils.NewError(http.StatusBadRequest, "INVALID_AMOUNT", "Invalid payment amount", nil))
			return
		}
		h.logger.Error("Failed to update tax payment amount", "error", err)
		httputils.ErrorJSON(w, httputils.ErrInternal)
		return
	}

	httputils.JSON(w, http.StatusOK, payment)
}
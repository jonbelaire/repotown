package api

import (
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/jonbelaire/repotown/packages/go-core/httputils"
	"github.com/jonbelaire/repotown/packages/go-core/logging"
	"github.com/jonbelaire/repotown/services/treasury/internal/service"
)

// TaxReportHandler manages tax report-related HTTP requests
type TaxReportHandler struct {
	taxReportService service.TaxReportService
	logger           logging.Logger
}

// NewTaxReportHandler creates a new tax report handler
func NewTaxReportHandler(taxReportService service.TaxReportService, logger logging.Logger) *TaxReportHandler {
	return &TaxReportHandler{
		taxReportService: taxReportService,
		logger:           logger,
	}
}

// RegisterRoutes registers tax report routes
func (h *TaxReportHandler) RegisterRoutes(r chi.Router) {
	r.Route("/reports", func(r chi.Router) {
		r.Get("/revenue", h.getRevenueReport)
		r.Get("/filing-status", h.getFilingStatusReport)
		r.Get("/taxpayer-compliance", h.getTaxpayerComplianceReport)
		r.Get("/tax-type-breakdown", h.getTaxTypeBreakdownReport)
	})
}

// getRevenueReport handles GET /reports/revenue
func (h *TaxReportHandler) getRevenueReport(w http.ResponseWriter, r *http.Request) {
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

	// Generate revenue report
	report, err := h.taxReportService.GenerateRevenueReport(r.Context(), startDate, endDate)
	if err != nil {
		h.logger.Error("Failed to generate revenue report", "error", err)
		httputils.ErrorJSON(w, httputils.ErrInternal)
		return
	}

	httputils.JSON(w, http.StatusOK, report)
}

// getFilingStatusReport handles GET /reports/filing-status
func (h *TaxReportHandler) getFilingStatusReport(w http.ResponseWriter, r *http.Request) {
	// Parse tax year
	yearStr := r.URL.Query().Get("year")
	if yearStr == "" {
		httputils.ErrorJSON(w, httputils.NewError(http.StatusBadRequest, "MISSING_YEAR", "Year parameter is required", nil))
		return
	}
	year, err := strconv.Atoi(yearStr)
	if err != nil || year < 1900 || year > 2100 {
		httputils.ErrorJSON(w, httputils.NewError(http.StatusBadRequest, "INVALID_YEAR", "Invalid year format", nil))
		return
	}

	// Generate filing status report
	report, err := h.taxReportService.GenerateFilingStatusReport(r.Context(), year)
	if err != nil {
		h.logger.Error("Failed to generate filing status report", "error", err)
		httputils.ErrorJSON(w, httputils.ErrInternal)
		return
	}

	httputils.JSON(w, http.StatusOK, report)
}

// getTaxpayerComplianceReport handles GET /reports/taxpayer-compliance
func (h *TaxReportHandler) getTaxpayerComplianceReport(w http.ResponseWriter, r *http.Request) {
	// Generate taxpayer compliance report
	report, err := h.taxReportService.GenerateTaxpayerComplianceReport(r.Context())
	if err != nil {
		h.logger.Error("Failed to generate taxpayer compliance report", "error", err)
		httputils.ErrorJSON(w, httputils.ErrInternal)
		return
	}

	httputils.JSON(w, http.StatusOK, report)
}

// getTaxTypeBreakdownReport handles GET /reports/tax-type-breakdown
func (h *TaxReportHandler) getTaxTypeBreakdownReport(w http.ResponseWriter, r *http.Request) {
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

	// Generate tax type breakdown report
	report, err := h.taxReportService.GenerateTaxTypeBreakdownReport(r.Context(), startDate, endDate)
	if err != nil {
		h.logger.Error("Failed to generate tax type breakdown report", "error", err)
		httputils.ErrorJSON(w, httputils.ErrInternal)
		return
	}

	httputils.JSON(w, http.StatusOK, report)
}
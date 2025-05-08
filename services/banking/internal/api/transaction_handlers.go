package api

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/jonbelaire/repotown/packages/go-core/httputils"
	"github.com/jonbelaire/repotown/packages/go-core/logging"
	"github.com/jonbelaire/repotown/services/banking/internal/domain"
	"github.com/jonbelaire/repotown/services/banking/internal/service"
)

// TransactionHandler manages transaction-related HTTP requests
type TransactionHandler struct {
	transactionService service.TransactionService
	logger            logging.Logger
}

// NewTransactionHandler creates a new transaction handler
func NewTransactionHandler(transactionService service.TransactionService, logger logging.Logger) *TransactionHandler {
	return &TransactionHandler{
		transactionService: transactionService,
		logger:            logger,
	}
}

// RegisterRoutes registers transaction routes
func (h *TransactionHandler) RegisterRoutes(r chi.Router) {
	r.Route("/transactions", func(r chi.Router) {
		r.Get("/", h.listTransactions)
		r.Post("/", h.createTransaction)
		r.Get("/{id}", h.getTransaction)
	})
}

// listTransactions handles GET /transactions
func (h *TransactionHandler) listTransactions(w http.ResponseWriter, r *http.Request) {
	// Parse pagination parameters
	limit, offset := getPaginationParams(r)
	
	// Check if filtering by account ID
	accountID := r.URL.Query().Get("account_id")
	
	var transactions []*domain.Transaction
	var err error
	
	// Get transactions based on filters
	if accountID != "" {
		transactions, err = h.transactionService.ListAccountTransactions(r.Context(), accountID, limit, offset)
	} else {
		transactions, err = h.transactionService.ListTransactions(r.Context(), limit, offset)
	}
	
	if err != nil {
		h.logger.Error("Failed to list transactions", "error", err)
		httputils.ErrorJSON(w, httputils.ErrInternal)
		return
	}
	
	httputils.JSON(w, http.StatusOK, transactions)
}

// createTransaction handles POST /transactions
func (h *TransactionHandler) createTransaction(w http.ResponseWriter, r *http.Request) {
	// Parse request body
	var req struct {
		Type        domain.TransactionType `json:"type" validate:"required"`
		AccountID   string                `json:"account_id" validate:"required"`
		Amount      int64                 `json:"amount" validate:"required,gt=0"`
		CurrencyCode string                `json:"currency_code" validate:"required"`
		Description string                `json:"description" validate:"required"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httputils.ErrorJSON(w, httputils.ErrBadRequest)
		return
	}
	
	// Validate transaction type
	if req.Type != domain.TransactionTypeDeposit && req.Type != domain.TransactionTypeWithdrawal {
		httputils.ErrorJSON(w, httputils.NewError(
			http.StatusBadRequest,
			"INVALID_TRANSACTION_TYPE",
			"Transaction type must be 'deposit' or 'withdrawal'",
			nil,
		))
		return
	}
	
	// Create transaction
	transaction, err := h.transactionService.CreateTransaction(
		r.Context(),
		req.Type,
		req.AccountID,
		req.Amount,
		req.CurrencyCode,
		req.Description,
	)
	if err != nil {
		if err == domain.ErrAccountNotFound {
			httputils.ErrorJSON(w, httputils.ErrNotFound)
			return
		}
		if err == domain.ErrAccountClosed {
			httputils.ErrorJSON(w, httputils.NewError(
				http.StatusBadRequest,
				"ACCOUNT_CLOSED",
				"Account is closed",
				nil,
			))
			return
		}
		if err == domain.ErrInsufficientFunds {
			httputils.ErrorJSON(w, httputils.NewError(
				http.StatusBadRequest,
				"INSUFFICIENT_FUNDS",
				"Insufficient funds for withdrawal",
				nil,
			))
			return
		}
		h.logger.Error("Failed to create transaction", "error", err)
		httputils.ErrorJSON(w, httputils.ErrInternal)
		return
	}
	
	httputils.JSON(w, http.StatusCreated, transaction)
}

// getTransaction handles GET /transactions/{id}
func (h *TransactionHandler) getTransaction(w http.ResponseWriter, r *http.Request) {
	// Get transaction ID from path
	id := chi.URLParam(r, "id")
	
	// Get transaction
	transaction, err := h.transactionService.GetTransaction(r.Context(), id)
	if err != nil {
		h.logger.Error("Failed to get transaction", "error", err)
		httputils.ErrorJSON(w, httputils.ErrInternal)
		return
	}
	
	if transaction == nil {
		httputils.ErrorJSON(w, httputils.ErrNotFound)
		return
	}
	
	httputils.JSON(w, http.StatusOK, transaction)
}
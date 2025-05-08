package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/jonbelaire/repotown/packages/go-core/httputils"
	"github.com/jonbelaire/repotown/packages/go-core/logging"
	"github.com/jonbelaire/repotown/services/banking/internal/domain"
	"github.com/jonbelaire/repotown/services/banking/internal/service"
)

// AccountHandler manages account-related HTTP requests
type AccountHandler struct {
	accountService service.AccountService
	logger         logging.Logger
}

// NewAccountHandler creates a new account handler
func NewAccountHandler(accountService service.AccountService, logger logging.Logger) *AccountHandler {
	return &AccountHandler{
		accountService: accountService,
		logger:         logger,
	}
}

// CreateAccountRequest defines the request body for creating an account
type CreateAccountRequest struct {
	CustomerID   string             `json:"customer_id" validate:"required"`
	AccountType  domain.AccountType `json:"account_type" validate:"required"`
	Name         string             `json:"name" validate:"required"`
	CurrencyCode string             `json:"currency_code" validate:"required"`
}

// RegisterRoutes registers account routes
func (h *AccountHandler) RegisterRoutes(r chi.Router) {
	r.Route("/accounts", func(r chi.Router) {
		r.Get("/", h.listAccounts)
		r.Post("/", h.createAccount)
		r.Get("/{id}", h.getAccount)
		r.Put("/{id}", h.updateAccount)
		r.Delete("/{id}", h.closeAccount)

		// Transaction-related endpoints
		r.Post("/{id}/deposit", h.deposit)
		r.Post("/{id}/withdraw", h.withdraw)
		r.Post("/{id}/transfer", h.transfer)
	})
}

// listAccounts handles GET /accounts
func (h *AccountHandler) listAccounts(w http.ResponseWriter, r *http.Request) {
	// Parse pagination parameters
	limit, offset := getPaginationParams(r)

	// Get accounts
	accounts, err := h.accountService.ListAccounts(r.Context(), limit, offset)
	if err != nil {
		h.logger.Error("Failed to list accounts", "error", err)
		httputils.ErrorJSON(w, httputils.ErrInternal)
		return
	}

	httputils.JSON(w, http.StatusOK, accounts)
}

// createAccount handles POST /accounts
func (h *AccountHandler) createAccount(w http.ResponseWriter, r *http.Request) {
	// Parse request body
	var req CreateAccountRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httputils.ErrorJSON(w, httputils.ErrBadRequest)
		return
	}

	// Create account
	account, err := h.accountService.CreateAccount(
		r.Context(),
		req.CustomerID,
		req.AccountType,
		req.Name,
		req.CurrencyCode,
	)
	if err != nil {
		h.logger.Error("Failed to create account", "error", err)
		httputils.ErrorJSON(w, httputils.ErrInternal)
		return
	}

	httputils.JSON(w, http.StatusCreated, account)
}

// getAccount handles GET /accounts/{id}
func (h *AccountHandler) getAccount(w http.ResponseWriter, r *http.Request) {
	// Get account ID from path
	id := chi.URLParam(r, "id")

	// Get account
	account, err := h.accountService.GetAccount(r.Context(), id)
	if err != nil {
		if err == domain.ErrAccountNotFound {
			httputils.ErrorJSON(w, httputils.ErrNotFound)
			return
		}
		h.logger.Error("Failed to get account", "error", err)
		httputils.ErrorJSON(w, httputils.ErrInternal)
		return
	}

	httputils.JSON(w, http.StatusOK, account)
}

// updateAccount handles PUT /accounts/{id}
func (h *AccountHandler) updateAccount(w http.ResponseWriter, r *http.Request) {
	// Get account ID from path
	id := chi.URLParam(r, "id")

	// Parse request body
	var req struct {
		Name string `json:"name" validate:"required"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httputils.ErrorJSON(w, httputils.ErrBadRequest)
		return
	}

	// Update account
	account, err := h.accountService.UpdateAccount(r.Context(), id, req.Name)
	if err != nil {
		if err == domain.ErrAccountNotFound {
			httputils.ErrorJSON(w, httputils.ErrNotFound)
			return
		}
		if err == domain.ErrAccountClosed {
			httputils.ErrorJSON(w, httputils.NewError(http.StatusBadRequest, "ACCOUNT_CLOSED", "Account is closed", nil))
			return
		}
		h.logger.Error("Failed to update account", "error", err)
		httputils.ErrorJSON(w, httputils.ErrInternal)
		return
	}

	httputils.JSON(w, http.StatusOK, account)
}

// closeAccount handles DELETE /accounts/{id}
func (h *AccountHandler) closeAccount(w http.ResponseWriter, r *http.Request) {
	// Get account ID from path
	id := chi.URLParam(r, "id")

	// Close account
	if err := h.accountService.CloseAccount(r.Context(), id); err != nil {
		if err == domain.ErrAccountNotFound {
			httputils.ErrorJSON(w, httputils.ErrNotFound)
			return
		}
		if err == domain.ErrAccountClosed {
			httputils.ErrorJSON(w, httputils.NewError(http.StatusBadRequest, "ACCOUNT_CLOSED", "Account is already closed", nil))
			return
		}
		h.logger.Error("Failed to close account", "error", err)
		httputils.ErrorJSON(w, httputils.ErrInternal)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// deposit handles POST /accounts/{id}/deposit
func (h *AccountHandler) deposit(w http.ResponseWriter, r *http.Request) {
	// Get account ID from path
	id := chi.URLParam(r, "id")

	// Parse request body
	var req struct {
		Amount      int64  `json:"amount" validate:"required,gt=0"`
		Description string `json:"description" validate:"required"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httputils.ErrorJSON(w, httputils.ErrBadRequest)
		return
	}

	// Deposit
	tx, err := h.accountService.Deposit(r.Context(), id, req.Amount, req.Description)
	if err != nil {
		if err == domain.ErrAccountNotFound {
			httputils.ErrorJSON(w, httputils.ErrNotFound)
			return
		}
		if err == domain.ErrAccountClosed {
			httputils.ErrorJSON(w, httputils.NewError(http.StatusBadRequest, "ACCOUNT_CLOSED", "Account is closed", nil))
			return
		}
		h.logger.Error("Failed to deposit", "error", err)
		httputils.ErrorJSON(w, httputils.ErrInternal)
		return
	}

	httputils.JSON(w, http.StatusOK, tx)
}

// withdraw handles POST /accounts/{id}/withdraw
func (h *AccountHandler) withdraw(w http.ResponseWriter, r *http.Request) {
	// Get account ID from path
	id := chi.URLParam(r, "id")

	// Parse request body
	var req struct {
		Amount      int64  `json:"amount" validate:"required,gt=0"`
		Description string `json:"description" validate:"required"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httputils.ErrorJSON(w, httputils.ErrBadRequest)
		return
	}

	// Withdraw
	tx, err := h.accountService.Withdraw(r.Context(), id, req.Amount, req.Description)
	if err != nil {
		if err == domain.ErrAccountNotFound {
			httputils.ErrorJSON(w, httputils.ErrNotFound)
			return
		}
		if err == domain.ErrAccountClosed {
			httputils.ErrorJSON(w, httputils.NewError(http.StatusBadRequest, "ACCOUNT_CLOSED", "Account is closed", nil))
			return
		}
		if err == domain.ErrInsufficientFunds {
			httputils.ErrorJSON(w, httputils.NewError(http.StatusBadRequest, "INSUFFICIENT_FUNDS", "Insufficient funds", nil))
			return
		}
		h.logger.Error("Failed to withdraw", "error", err)
		httputils.ErrorJSON(w, httputils.ErrInternal)
		return
	}

	httputils.JSON(w, http.StatusOK, tx)
}

// transfer handles POST /accounts/{id}/transfer
func (h *AccountHandler) transfer(w http.ResponseWriter, r *http.Request) {
	// Get source account ID from path
	sourceID := chi.URLParam(r, "id")

	// Parse request body
	var req struct {
		TargetID    string `json:"target_id" validate:"required"`
		Amount      int64  `json:"amount" validate:"required,gt=0"`
		Description string `json:"description" validate:"required"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httputils.ErrorJSON(w, httputils.ErrBadRequest)
		return
	}

	// Transfer
	tx, err := h.accountService.Transfer(r.Context(), sourceID, req.TargetID, req.Amount, req.Description)
	if err != nil {
		if err == domain.ErrAccountNotFound {
			httputils.ErrorJSON(w, httputils.ErrNotFound)
			return
		}
		if err == domain.ErrAccountClosed {
			httputils.ErrorJSON(w, httputils.NewError(http.StatusBadRequest, "ACCOUNT_CLOSED", "Account is closed", nil))
			return
		}
		if err == domain.ErrInsufficientFunds {
			httputils.ErrorJSON(w, httputils.NewError(http.StatusBadRequest, "INSUFFICIENT_FUNDS", "Insufficient funds", nil))
			return
		}
		h.logger.Error("Failed to transfer", "error", err)
		httputils.ErrorJSON(w, httputils.ErrInternal)
		return
	}

	httputils.JSON(w, http.StatusOK, tx)
}

// Helpers

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

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

// CustomerHandler manages customer-related HTTP requests
type CustomerHandler struct {
	customerService service.CustomerService
	accountService  service.AccountService
	logger         logging.Logger
}

// NewCustomerHandler creates a new customer handler
func NewCustomerHandler(
	customerService service.CustomerService,
	accountService service.AccountService,
	logger logging.Logger,
) *CustomerHandler {
	return &CustomerHandler{
		customerService: customerService,
		accountService:  accountService,
		logger:         logger,
	}
}

// RegisterRoutes registers customer routes
func (h *CustomerHandler) RegisterRoutes(r chi.Router) {
	r.Route("/customers", func(r chi.Router) {
		r.Get("/", h.listCustomers)
		r.Post("/", h.createCustomer)
		r.Get("/{id}", h.getCustomer)
		r.Put("/{id}", h.updateCustomer)
		r.Delete("/{id}", h.deleteCustomer)
		r.Get("/{id}/accounts", h.getCustomerAccounts)
	})
}

// CustomerRequest defines the request body for creating/updating a customer
type CustomerRequest struct {
	FirstName   string         `json:"first_name" validate:"required"`
	LastName    string         `json:"last_name" validate:"required"`
	Email       string         `json:"email" validate:"required,email"`
	PhoneNumber string         `json:"phone_number"`
	Address     domain.Address `json:"address" validate:"required"`
}

// listCustomers handles GET /customers
func (h *CustomerHandler) listCustomers(w http.ResponseWriter, r *http.Request) {
	// Parse pagination parameters
	limit, offset := getPaginationParams(r)
	
	// Get customers
	customers, err := h.customerService.ListCustomers(r.Context(), limit, offset)
	if err != nil {
		h.logger.Error("Failed to list customers", "error", err)
		httputils.ErrorJSON(w, httputils.ErrInternal)
		return
	}
	
	httputils.JSON(w, http.StatusOK, customers)
}

// createCustomer handles POST /customers
func (h *CustomerHandler) createCustomer(w http.ResponseWriter, r *http.Request) {
	// Parse request body
	var req CustomerRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httputils.ErrorJSON(w, httputils.ErrBadRequest)
		return
	}
	
	// Create customer
	customer, err := h.customerService.CreateCustomer(
		r.Context(),
		req.FirstName,
		req.LastName,
		req.Email,
		req.PhoneNumber,
		req.Address,
	)
	if err != nil {
		if err == domain.ErrCustomerExists {
			httputils.ErrorJSON(w, httputils.NewError(
				http.StatusConflict,
				"CUSTOMER_EXISTS",
				"A customer with this email already exists",
				nil,
			))
			return
		}
		h.logger.Error("Failed to create customer", "error", err)
		httputils.ErrorJSON(w, httputils.ErrInternal)
		return
	}
	
	httputils.JSON(w, http.StatusCreated, customer)
}

// getCustomer handles GET /customers/{id}
func (h *CustomerHandler) getCustomer(w http.ResponseWriter, r *http.Request) {
	// Get customer ID from path
	id := chi.URLParam(r, "id")
	
	// Get customer
	customer, err := h.customerService.GetCustomer(r.Context(), id)
	if err != nil {
		if err == domain.ErrCustomerNotFound {
			httputils.ErrorJSON(w, httputils.ErrNotFound)
			return
		}
		h.logger.Error("Failed to get customer", "error", err)
		httputils.ErrorJSON(w, httputils.ErrInternal)
		return
	}
	
	httputils.JSON(w, http.StatusOK, customer)
}

// updateCustomer handles PUT /customers/{id}
func (h *CustomerHandler) updateCustomer(w http.ResponseWriter, r *http.Request) {
	// Get customer ID from path
	id := chi.URLParam(r, "id")
	
	// Parse request body
	var req CustomerRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httputils.ErrorJSON(w, httputils.ErrBadRequest)
		return
	}
	
	// Update customer
	customer, err := h.customerService.UpdateCustomer(
		r.Context(),
		id,
		req.FirstName,
		req.LastName,
		req.Email,
		req.PhoneNumber,
		req.Address,
	)
	if err != nil {
		if err == domain.ErrCustomerNotFound {
			httputils.ErrorJSON(w, httputils.ErrNotFound)
			return
		}
		if err == domain.ErrCustomerExists {
			httputils.ErrorJSON(w, httputils.NewError(
				http.StatusConflict,
				"CUSTOMER_EXISTS",
				"A customer with this email already exists",
				nil,
			))
			return
		}
		h.logger.Error("Failed to update customer", "error", err)
		httputils.ErrorJSON(w, httputils.ErrInternal)
		return
	}
	
	httputils.JSON(w, http.StatusOK, customer)
}

// deleteCustomer handles DELETE /customers/{id}
func (h *CustomerHandler) deleteCustomer(w http.ResponseWriter, r *http.Request) {
	// Get customer ID from path
	id := chi.URLParam(r, "id")
	
	// Delete customer
	if err := h.customerService.DeleteCustomer(r.Context(), id); err != nil {
		if err == domain.ErrCustomerNotFound {
			httputils.ErrorJSON(w, httputils.ErrNotFound)
			return
		}
		h.logger.Error("Failed to delete customer", "error", err)
		httputils.ErrorJSON(w, httputils.ErrInternal)
		return
	}
	
	w.WriteHeader(http.StatusNoContent)
}

// getCustomerAccounts handles GET /customers/{id}/accounts
func (h *CustomerHandler) getCustomerAccounts(w http.ResponseWriter, r *http.Request) {
	// Get customer ID from path
	id := chi.URLParam(r, "id")
	
	// Get customer to verify it exists
	_, err := h.customerService.GetCustomer(r.Context(), id)
	if err != nil {
		if err == domain.ErrCustomerNotFound {
			httputils.ErrorJSON(w, httputils.ErrNotFound)
			return
		}
		h.logger.Error("Failed to get customer", "error", err)
		httputils.ErrorJSON(w, httputils.ErrInternal)
		return
	}
	
	// Get customer accounts
	accounts, err := h.accountService.ListAccountsByCustomer(r.Context(), id)
	if err != nil {
		h.logger.Error("Failed to get customer accounts", "error", err)
		httputils.ErrorJSON(w, httputils.ErrInternal)
		return
	}
	
	httputils.JSON(w, http.StatusOK, accounts)
}
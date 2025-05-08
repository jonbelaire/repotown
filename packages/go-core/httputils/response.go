package httputils

import (
	"encoding/json"
	"net/http"
)

// ResponseError represents an error response
type ResponseError struct {
	Status  int         `json:"-"`
	Code    string      `json:"code"`
	Message string      `json:"message"`
	Details interface{} `json:"details,omitempty"`
}

// Response represents the standard API response format
type Response struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   interface{} `json:"error,omitempty"`
	Meta    interface{} `json:"meta,omitempty"`
}

// JSON sends a JSON response with the provided status code and data
func JSON(w http.ResponseWriter, status int, data interface{}) {
	response := Response{
		Success: status >= 200 && status < 400,
		Data:    data,
	}

	JSONResponse(w, status, response)
}

// ErrorJSON sends an error JSON response with the provided status code and error
func ErrorJSON(w http.ResponseWriter, err ResponseError) {
	response := Response{
		Success: false,
		Error:   err,
	}

	JSONResponse(w, err.Status, response)
}

// JSONResponse writes the response as JSON
func JSONResponse(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	
	if err := json.NewEncoder(w).Encode(data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// NewError creates a new ResponseError
func NewError(status int, code, message string, details interface{}) ResponseError {
	return ResponseError{
		Status:  status,
		Code:    code,
		Message: message,
		Details: details,
	}
}

// CommonErrors for API responses
var (
	ErrBadRequest = NewError(http.StatusBadRequest, "BAD_REQUEST", "Invalid request", nil)
	ErrNotFound = NewError(http.StatusNotFound, "NOT_FOUND", "Resource not found", nil)
	ErrUnauthorized = NewError(http.StatusUnauthorized, "UNAUTHORIZED", "Unauthorized access", nil)
	ErrForbidden = NewError(http.StatusForbidden, "FORBIDDEN", "Access forbidden", nil)
	ErrInternal = NewError(http.StatusInternalServerError, "INTERNAL_ERROR", "Internal server error", nil)
)
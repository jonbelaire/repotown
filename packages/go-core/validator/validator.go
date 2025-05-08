package validator

import (
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/jonbelaire/repotown/packages/go-core/httputils"
)

// Validator represents a struct validator
type Validator struct {
	validate *validator.Validate
}

// NewValidator creates a new validator
func NewValidator() *Validator {
	v := validator.New()
	
	// Register function to get json tag name instead of struct field name
	v.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})

	return &Validator{validate: v}
}

// Validate validates a struct and returns validation errors
func (v *Validator) Validate(i interface{}) error {
	return v.validate.Struct(i)
}

// ValidateJSON validates a struct and returns formatted validation errors
func (v *Validator) ValidateJSON(i interface{}) (bool, *httputils.ResponseError) {
	err := v.validate.Struct(i)
	if err == nil {
		return true, nil
	}

	var validationErrors validator.ValidationErrors
	if errors.As(err, &validationErrors) {
		errorDetails := make(map[string]string)
		
		for _, e := range validationErrors {
			field := e.Field()
			errorDetails[field] = formattedValidationError(e)
		}

		responseErr := httputils.NewError(
			400,
			"VALIDATION_ERROR",
			"Validation failed",
			errorDetails,
		)
		
		return false, &responseErr
	}

	// If we got here, it's not a validation error
	responseErr := httputils.ErrBadRequest
	return false, &responseErr
}

// Helper to format validation error messages
func formattedValidationError(err validator.FieldError) string {
	switch err.Tag() {
	case "required":
		return fmt.Sprintf("The %s field is required", err.Field())
	case "email":
		return fmt.Sprintf("The %s field must be a valid email", err.Field())
	case "min":
		return fmt.Sprintf("The %s field must be at least %s characters", err.Field(), err.Param())
	case "max":
		return fmt.Sprintf("The %s field must not exceed %s characters", err.Field(), err.Param())
	default:
		return fmt.Sprintf("The %s field failed validation: %s", err.Field(), err.Tag())
	}
}
package domain

import (
	"fmt"
	"strings"
)

// ValidationError represents a domain validation failure with detailed messages
type ValidationError struct {
	Message string
	Details []string
}

func (e ValidationError) Error() string {
	if len(e.Details) == 0 {
		return fmt.Sprintf("validation failed: %s", e.Message)
	}
	return fmt.Sprintf("validation failed: %s [%s]", e.Message, strings.Join(e.Details, ", "))
}

// NewValidationError creates a validation error with details
func NewValidationError(message string, details ...string) ValidationError {
	return ValidationError{
		Message: message,
		Details: details,
	}
}

// NotFoundError represents a resource not found
type NotFoundError struct {
	Resource string
}

func (e NotFoundError) Error() string {
	return fmt.Sprintf("%s not found", e.Resource)
}

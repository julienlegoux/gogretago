package errors

import "fmt"

// ApplicationError is the base error type for application-level errors
type ApplicationError struct {
	Message string
	Code    string
}

func (e *ApplicationError) Error() string {
	return e.Message
}

// ValidationError is returned when input validation fails
type ValidationError struct {
	ApplicationError
	Details map[string][]string
}

func NewValidationError(message string, details map[string][]string) *ValidationError {
	return &ValidationError{
		ApplicationError: ApplicationError{
			Message: message,
			Code:    "VALIDATION_ERROR",
		},
		Details: details,
	}
}

// NotFoundError is returned when a resource is not found
type NotFoundError struct {
	ApplicationError
}

func NewNotFoundError(resource string, identifier string) *NotFoundError {
	return &NotFoundError{
		ApplicationError: ApplicationError{
			Message: fmt.Sprintf("%s not found: %s", resource, identifier),
			Code:    "NOT_FOUND",
		},
	}
}

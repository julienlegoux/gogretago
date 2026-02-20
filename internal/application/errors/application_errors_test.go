package errors

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidationError_HasCorrectCode(t *testing.T) {
	details := map[string][]string{
		"Email": {"Invalid email format"},
	}
	err := NewValidationError("Validation failed", details)
	assert.Equal(t, "VALIDATION_ERROR", err.Code)
}

func TestValidationError_HasDetails(t *testing.T) {
	details := map[string][]string{
		"Email":    {"Invalid email format"},
		"Password": {"Password is required", "Password must contain at least one lowercase, one uppercase, and one number"},
	}
	err := NewValidationError("Validation failed", details)

	assert.Equal(t, details, err.Details)
	assert.Len(t, err.Details["Email"], 1)
	assert.Len(t, err.Details["Password"], 2)
}

func TestValidationError_ErrorMethod(t *testing.T) {
	err := NewValidationError("Validation failed", nil)
	assert.Equal(t, "Validation failed", err.Error())
}

func TestNotFoundError_HasCorrectFormat(t *testing.T) {
	err := NewNotFoundError("User", "abc-123")
	assert.Equal(t, "NOT_FOUND", err.Code)
	assert.Contains(t, err.Message, "User")
	assert.Contains(t, err.Message, "abc-123")
	assert.Equal(t, "User not found: abc-123", err.Message)
}

func TestApplicationError_ImplementsErrorInterface(t *testing.T) {
	tests := []struct {
		name string
		err  error
	}{
		{"ApplicationError", &ApplicationError{Message: "test", Code: "TEST"}},
		{"ValidationError", NewValidationError("test", nil)},
		{"NotFoundError", NewNotFoundError("User", "1")},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Implements(t, (*error)(nil), tt.err)
			assert.NotEmpty(t, tt.err.Error())
		})
	}
}

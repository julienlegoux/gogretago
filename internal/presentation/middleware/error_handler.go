package middleware

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	apperrors "github.com/lgxju/gogretago/internal/application/errors"
	domainerrors "github.com/lgxju/gogretago/internal/domain/errors"
)

// ErrorResponse represents the standard error response format
type ErrorResponse struct {
	Success bool `json:"success"`
	Error   struct {
		Code    string              `json:"code"`
		Message string              `json:"message"`
		Details map[string][]string `json:"details,omitempty"`
	} `json:"error"`
}

// ErrorHandler middleware catches errors and returns standardized responses
func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		// Check if there are any errors
		if len(c.Errors) > 0 {
			err := c.Errors.Last().Err
			response := buildErrorResponse(err)
			status := getStatusCode(err)
			c.JSON(status, response)
		}
	}
}

func buildErrorResponse(err error) ErrorResponse {
	response := ErrorResponse{Success: false}

	// Check for validation errors
	var validationErr *apperrors.ValidationError
	if errors.As(err, &validationErr) {
		response.Error.Code = validationErr.Code
		response.Error.Message = validationErr.Message
		response.Error.Details = validationErr.Details
		return response
	}

	// Check for domain errors
	var domainErr *domainerrors.DomainError
	if errors.As(err, &domainErr) {
		response.Error.Code = domainErr.Code
		response.Error.Message = domainErr.Message
		return response
	}

	// Check for UserAlreadyExistsError
	var userExistsErr *domainerrors.UserAlreadyExistsError
	if errors.As(err, &userExistsErr) {
		response.Error.Code = userExistsErr.Code
		response.Error.Message = userExistsErr.Message
		return response
	}

	// Check for InvalidCredentialsError
	var invalidCredsErr *domainerrors.InvalidCredentialsError
	if errors.As(err, &invalidCredsErr) {
		response.Error.Code = invalidCredsErr.Code
		response.Error.Message = invalidCredsErr.Message
		return response
	}

	// Check for UserNotFoundError
	var userNotFoundErr *domainerrors.UserNotFoundError
	if errors.As(err, &userNotFoundErr) {
		response.Error.Code = userNotFoundErr.Code
		response.Error.Message = userNotFoundErr.Message
		return response
	}

	// Check for application errors
	var appErr *apperrors.ApplicationError
	if errors.As(err, &appErr) {
		response.Error.Code = appErr.Code
		response.Error.Message = appErr.Message
		return response
	}

	// Default internal error
	response.Error.Code = "INTERNAL_ERROR"
	response.Error.Message = "An unexpected error occurred"
	return response
}

func getStatusCode(err error) int {
	// Validation errors
	var validationErr *apperrors.ValidationError
	if errors.As(err, &validationErr) {
		return http.StatusBadRequest
	}

	// User already exists
	var userExistsErr *domainerrors.UserAlreadyExistsError
	if errors.As(err, &userExistsErr) {
		return http.StatusBadRequest
	}

	// Invalid credentials
	var invalidCredsErr *domainerrors.InvalidCredentialsError
	if errors.As(err, &invalidCredsErr) {
		return http.StatusUnauthorized
	}

	// User not found
	var userNotFoundErr *domainerrors.UserNotFoundError
	if errors.As(err, &userNotFoundErr) {
		return http.StatusNotFound
	}

	// Not found errors
	var notFoundErr *apperrors.NotFoundError
	if errors.As(err, &notFoundErr) {
		return http.StatusNotFound
	}

	return http.StatusInternalServerError
}

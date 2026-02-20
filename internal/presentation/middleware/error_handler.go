package middleware

import (
	"errors"

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

	var validationErr *apperrors.ValidationError
	if errors.As(err, &validationErr) {
		response.Error.Code = validationErr.Code
		response.Error.Message = validationErr.Message
		response.Error.Details = validationErr.Details
		return response
	}

	var domainErr *domainerrors.DomainError
	if errors.As(err, &domainErr) {
		response.Error.Code = domainErr.Code
		response.Error.Message = domainErr.Message
		return response
	}

	var appErr *apperrors.ApplicationError
	if errors.As(err, &appErr) {
		response.Error.Code = appErr.Code
		response.Error.Message = appErr.Message
		return response
	}

	response.Error.Code = "INTERNAL_ERROR"
	response.Error.Message = "An unexpected error occurred"
	return response
}

func getStatusCode(err error) int {
	var domainErr *domainerrors.DomainError
	if errors.As(err, &domainErr) {
		return domainerrors.GetHTTPStatus(domainErr.Code)
	}

	var validationErr *apperrors.ValidationError
	if errors.As(err, &validationErr) {
		return domainerrors.GetHTTPStatus(validationErr.Code)
	}

	var appErr *apperrors.ApplicationError
	if errors.As(err, &appErr) {
		return domainerrors.GetHTTPStatus(appErr.Code)
	}

	return 500
}

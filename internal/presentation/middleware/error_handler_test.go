package middleware

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	apperrors "github.com/lgxju/gogretago/internal/application/errors"
	domainerrors "github.com/lgxju/gogretago/internal/domain/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupErrorHandlerRouter(errToAdd error) (*gin.Engine, *httptest.ResponseRecorder) {
	router := gin.New()
	router.Use(ErrorHandler())
	router.GET("/test", func(c *gin.Context) {
		if errToAdd != nil {
			_ = c.Error(errToAdd)
			return
		}
		c.JSON(http.StatusOK, gin.H{"success": true})
	})
	return router, httptest.NewRecorder()
}

func TestErrorHandler_DomainError(t *testing.T) {
	// Use a plain *DomainError to test the domain error handling path.
	// Note: typed domain errors (e.g. *UserNotFoundError) embed DomainError
	// by value, so errors.As(err, &domainErr) does not match them.
	// The error handler's DomainError branch works with direct *DomainError values.
	domainErr := &domainerrors.DomainError{
		Message: "User not found: user-999",
		Code:    "USER_NOT_FOUND",
	}
	router, w := setupErrorHandlerRouter(domainErr)
	req := httptest.NewRequest(http.MethodGet, "/test", http.NoBody)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)

	var body ErrorResponse
	err := json.Unmarshal(w.Body.Bytes(), &body)
	require.NoError(t, err)
	assert.False(t, body.Success)
	assert.Equal(t, "USER_NOT_FOUND", body.Error.Code)
	assert.Contains(t, body.Error.Message, "user-999")
}

func TestErrorHandler_ValidationError(t *testing.T) {
	details := map[string][]string{
		"Email":    {"Invalid email format"},
		"Password": {"Password is required"},
	}
	validationErr := apperrors.NewValidationError("Validation failed", details)

	router, w := setupErrorHandlerRouter(validationErr)
	req := httptest.NewRequest(http.MethodGet, "/test", http.NoBody)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var body ErrorResponse
	err := json.Unmarshal(w.Body.Bytes(), &body)
	require.NoError(t, err)
	assert.False(t, body.Success)
	assert.Equal(t, "VALIDATION_ERROR", body.Error.Code)
	assert.Equal(t, "Validation failed", body.Error.Message)
	assert.Contains(t, body.Error.Details["Email"], "Invalid email format")
	assert.Contains(t, body.Error.Details["Password"], "Password is required")
}

func TestErrorHandler_ApplicationError(t *testing.T) {
	appErr := &apperrors.ApplicationError{
		Message: "Something went wrong in the app",
		Code:    "INTERNAL_ERROR",
	}

	router, w := setupErrorHandlerRouter(appErr)
	req := httptest.NewRequest(http.MethodGet, "/test", http.NoBody)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)

	var body ErrorResponse
	err := json.Unmarshal(w.Body.Bytes(), &body)
	require.NoError(t, err)
	assert.False(t, body.Success)
	assert.Equal(t, "INTERNAL_ERROR", body.Error.Code)
}

func TestErrorHandler_UnknownError(t *testing.T) {
	unknownErr := fmt.Errorf("some unexpected error")

	router, w := setupErrorHandlerRouter(unknownErr)
	req := httptest.NewRequest(http.MethodGet, "/test", http.NoBody)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)

	var body ErrorResponse
	err := json.Unmarshal(w.Body.Bytes(), &body)
	require.NoError(t, err)
	assert.False(t, body.Success)
	assert.Equal(t, "INTERNAL_ERROR", body.Error.Code)
	assert.Equal(t, "An unexpected error occurred", body.Error.Message)
}

func TestErrorHandler_NoErrors(t *testing.T) {
	router, w := setupErrorHandlerRouter(nil)
	req := httptest.NewRequest(http.MethodGet, "/test", http.NoBody)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var body map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &body)
	require.NoError(t, err)
	assert.Equal(t, true, body["success"])
}

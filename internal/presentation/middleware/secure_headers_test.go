package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestSecureHeaders_SetsAllHeaders(t *testing.T) {
	router := gin.New()
	router.Use(SecureHeaders())
	router.GET("/test", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "nosniff", w.Header().Get("X-Content-Type-Options"))
	assert.Equal(t, "SAMEORIGIN", w.Header().Get("X-Frame-Options"))
	assert.Equal(t, "max-age=63072000; includeSubDomains", w.Header().Get("Strict-Transport-Security"))
	assert.Equal(t, "0", w.Header().Get("X-XSS-Protection"))
	assert.Equal(t, "default-src 'self'", w.Header().Get("Content-Security-Policy"))
	assert.Equal(t, "no-referrer", w.Header().Get("Referrer-Policy"))
	assert.Equal(t, "none", w.Header().Get("X-Permitted-Cross-Domain-Policies"))
	assert.Equal(t, "noopen", w.Header().Get("X-Download-Options"))
}

func TestSecureHeaders_HeadersPresentOnErrorResponses(t *testing.T) {
	router := gin.New()
	router.Use(SecureHeaders())
	router.GET("/error", func(c *gin.Context) {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "fail"})
	})

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/error", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Equal(t, "nosniff", w.Header().Get("X-Content-Type-Options"))
	assert.Equal(t, "SAMEORIGIN", w.Header().Get("X-Frame-Options"))
}

func TestSecureHeaders_HeadersPresentOnDifferentMethods(t *testing.T) {
	router := gin.New()
	router.Use(SecureHeaders())
	router.POST("/test", func(c *gin.Context) {
		c.Status(http.StatusCreated)
	})

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/test", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
	assert.Equal(t, "nosniff", w.Header().Get("X-Content-Type-Options"))
	assert.NotEmpty(t, w.Header().Get("Strict-Transport-Security"))
}

package middleware

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestBodyLimit_AllowsSmallBody(t *testing.T) {
	router := gin.New()
	router.Use(BodyLimit(1024)) // 1KB limit
	router.POST("/test", func(c *gin.Context) {
		body, err := io.ReadAll(c.Request.Body)
		if err != nil {
			c.JSON(http.StatusRequestEntityTooLarge, gin.H{"error": "body too large"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"size": len(body)})
	})

	smallBody := bytes.NewBufferString(`{"name": "test"}`)
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/test", smallBody)
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestBodyLimit_RejectsOversizedBody(t *testing.T) {
	router := gin.New()
	router.Use(BodyLimit(100)) // 100 bytes limit
	router.POST("/test", func(c *gin.Context) {
		_, err := io.ReadAll(c.Request.Body)
		if err != nil {
			c.JSON(http.StatusRequestEntityTooLarge, gin.H{"error": "body too large"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"success": true})
	})

	// Create a body larger than the limit
	largeBody := strings.NewReader(strings.Repeat("x", 200))
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/test", largeBody)
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	// MaxBytesReader causes the handler to receive an error when reading
	assert.Equal(t, http.StatusRequestEntityTooLarge, w.Code)
}

func TestBodyLimit_NilBodyIsAllowed(t *testing.T) {
	router := gin.New()
	router.Use(BodyLimit(100))
	router.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"success": true})
	})

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/test", http.NoBody)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestBodyLimit_ExactlyAtLimit(t *testing.T) {
	router := gin.New()
	router.Use(BodyLimit(50))
	router.POST("/test", func(c *gin.Context) {
		body, err := io.ReadAll(c.Request.Body)
		if err != nil {
			c.JSON(http.StatusRequestEntityTooLarge, gin.H{"error": "body too large"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"size": len(body)})
	})

	// Body exactly at limit should be allowed
	exactBody := strings.NewReader(strings.Repeat("x", 50))
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/test", exactBody)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

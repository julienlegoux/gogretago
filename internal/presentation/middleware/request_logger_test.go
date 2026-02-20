package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestRequestLogger_SetsRequestIdHeader(t *testing.T) {
	router := gin.New()
	router.Use(RequestLogger())
	router.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"success": true})
	})

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/test", http.NoBody)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	requestID := w.Header().Get("X-Request-Id")
	assert.NotEmpty(t, requestID, "X-Request-Id header should be set")
	assert.Len(t, requestID, 36, "should be a UUID v4 (36 chars with hyphens)")
}

func TestRequestLogger_SetsRequestIdInContext(t *testing.T) {
	var contextRequestID string

	router := gin.New()
	router.Use(RequestLogger())
	router.GET("/test", func(c *gin.Context) {
		val, _ := c.Get("requestId")
		contextRequestID, _ = val.(string)
		c.JSON(http.StatusOK, gin.H{"success": true})
	})

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/test", http.NoBody)
	router.ServeHTTP(w, req)

	assert.NotEmpty(t, contextRequestID)
	assert.Equal(t, contextRequestID, w.Header().Get("X-Request-Id"))
}

func TestRequestLogger_UniqueRequestIdsPerRequest(t *testing.T) {
	router := gin.New()
	router.Use(RequestLogger())
	router.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"success": true})
	})

	w1 := httptest.NewRecorder()
	req1 := httptest.NewRequest(http.MethodGet, "/test", http.NoBody)
	router.ServeHTTP(w1, req1)

	w2 := httptest.NewRecorder()
	req2 := httptest.NewRequest(http.MethodGet, "/test", http.NoBody)
	router.ServeHTTP(w2, req2)

	id1 := w1.Header().Get("X-Request-Id")
	id2 := w2.Header().Get("X-Request-Id")
	assert.NotEqual(t, id1, id2, "each request should get a unique ID")
}

func TestRequestLogger_WorksWithDifferentHTTPMethods(t *testing.T) {
	router := gin.New()
	router.Use(RequestLogger())
	handler := func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"success": true})
	}
	router.GET("/test", handler)
	router.POST("/test", handler)
	router.DELETE("/test", handler)

	methods := []string{http.MethodGet, http.MethodPost, http.MethodDelete}
	for _, method := range methods {
		t.Run(method, func(t *testing.T) {
			w := httptest.NewRecorder()
			req := httptest.NewRequest(method, "/test", http.NoBody)
			router.ServeHTTP(w, req)

			assert.Equal(t, http.StatusOK, w.Code)
			assert.NotEmpty(t, w.Header().Get("X-Request-Id"))
		})
	}
}

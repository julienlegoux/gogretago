package middleware

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRateLimiter_AllowsRequestsWithinLimit(t *testing.T) {
	router := gin.New()
	router.Use(RateLimiter(60)) // 60 requests per minute
	router.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"success": true})
	})

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/test", http.NoBody)
	req.RemoteAddr = "192.168.1.1:1234"
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestRateLimiter_BlocksExcessiveRequests(t *testing.T) {
	router := gin.New()
	// Set a very low burst: 2 requests per minute (burst = 2)
	router.Use(RateLimiter(2))
	router.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"success": true})
	})

	// Send burst+1 requests quickly - the initial burst should be allowed,
	// then subsequent requests should be blocked
	var lastCode int
	blocked := false
	for i := 0; i < 10; i++ {
		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/test", http.NoBody)
		req.RemoteAddr = "10.0.0.1:1234"
		router.ServeHTTP(w, req)
		lastCode = w.Code
		if w.Code == http.StatusTooManyRequests {
			blocked = true
			break
		}
	}

	assert.True(t, blocked, "rate limiter should block excessive requests")
	assert.Equal(t, http.StatusTooManyRequests, lastCode)
}

func TestRateLimiter_ReturnsCorrectErrorResponse(t *testing.T) {
	router := gin.New()
	router.Use(RateLimiter(1)) // 1 request per minute (burst = 1)
	router.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"success": true})
	})

	// First request should pass
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/test", http.NoBody)
	req.RemoteAddr = "10.0.0.2:1234"
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	// Subsequent requests should be rate limited
	var rateLimitedResponse *httptest.ResponseRecorder
	for i := 0; i < 10; i++ {
		w = httptest.NewRecorder()
		req = httptest.NewRequest(http.MethodGet, "/test", http.NoBody)
		req.RemoteAddr = "10.0.0.2:1234"
		router.ServeHTTP(w, req)
		if w.Code == http.StatusTooManyRequests {
			rateLimitedResponse = w
			break
		}
	}

	require.NotNil(t, rateLimitedResponse, "should have received a 429 response")

	var body map[string]interface{}
	err := json.Unmarshal(rateLimitedResponse.Body.Bytes(), &body)
	require.NoError(t, err)
	assert.Equal(t, false, body["success"])

	errObj := body["error"].(map[string]interface{})
	assert.Equal(t, "RATE_LIMITED", errObj["code"])
	assert.NotEmpty(t, errObj["message"])
}

func TestRateLimiter_DifferentIPsHaveSeparateLimits(t *testing.T) {
	router := gin.New()
	router.Use(RateLimiter(1)) // 1 request per minute
	router.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"success": true})
	})

	// First IP - should be allowed
	w1 := httptest.NewRecorder()
	req1 := httptest.NewRequest(http.MethodGet, "/test", http.NoBody)
	req1.RemoteAddr = "10.0.0.10:1234"
	router.ServeHTTP(w1, req1)
	assert.Equal(t, http.StatusOK, w1.Code)

	// Second IP - should also be allowed (separate rate limiter)
	w2 := httptest.NewRecorder()
	req2 := httptest.NewRequest(http.MethodGet, "/test", http.NoBody)
	req2.RemoteAddr = "10.0.0.11:1234"
	router.ServeHTTP(w2, req2)
	assert.Equal(t, http.StatusOK, w2.Code)
}

func TestRateLimiter_AbortsOnRateLimit(t *testing.T) {
	handlerCalled := false
	router := gin.New()
	router.Use(RateLimiter(1))
	router.GET("/test", func(c *gin.Context) {
		handlerCalled = true
		c.JSON(http.StatusOK, gin.H{"success": true})
	})

	// First request passes
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/test", http.NoBody)
	req.RemoteAddr = "10.0.0.20:1234"
	router.ServeHTTP(w, req)
	assert.True(t, handlerCalled)

	// Wait for rate limit to kick in
	handlerCalled = false
	for i := 0; i < 10; i++ {
		w = httptest.NewRecorder()
		req = httptest.NewRequest(http.MethodGet, "/test", http.NoBody)
		req.RemoteAddr = "10.0.0.20:1234"
		handlerCalled = false
		router.ServeHTTP(w, req)
		if w.Code == http.StatusTooManyRequests {
			assert.False(t, handlerCalled, "handler should not be called when rate limited")
			return
		}
	}
	t.Fatal("expected rate limiting to kick in")
}

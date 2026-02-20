package middleware

import (
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

type ipLimiter struct {
	limiter *rate.Limiter
}

type rateLimiterStore struct {
	mu       sync.Mutex
	limiters map[string]*ipLimiter
	rate     rate.Limit
	burst    int
}

func newRateLimiterStore(r rate.Limit, burst int) *rateLimiterStore {
	return &rateLimiterStore{
		limiters: make(map[string]*ipLimiter),
		rate:     r,
		burst:    burst,
	}
}

func (s *rateLimiterStore) getLimiter(ip string) *rate.Limiter {
	s.mu.Lock()
	defer s.mu.Unlock()

	if l, exists := s.limiters[ip]; exists {
		return l.limiter
	}

	limiter := rate.NewLimiter(s.rate, s.burst)
	s.limiters[ip] = &ipLimiter{limiter: limiter}
	return limiter
}

// RateLimiter creates a rate limiting middleware.
// limit is the number of requests allowed per minute.
func RateLimiter(requestsPerMinute int) gin.HandlerFunc {
	// Convert requests per minute to rate.Limit (requests per second)
	r := rate.Limit(float64(requestsPerMinute) / 60.0)
	store := newRateLimiterStore(r, requestsPerMinute)

	return func(c *gin.Context) {
		ip := c.ClientIP()
		limiter := store.getLimiter(ip)

		if !limiter.Allow() {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"success": false,
				"error": gin.H{
					"code":    "RATE_LIMITED",
					"message": "Too many requests, please try again later",
				},
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

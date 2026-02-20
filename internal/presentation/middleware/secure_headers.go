package middleware

import "github.com/gin-gonic/gin"

// SecureHeaders sets security-related HTTP response headers matching covoitapi's secureHeaders() behavior.
func SecureHeaders() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("X-Content-Type-Options", "nosniff")
		c.Header("X-Frame-Options", "SAMEORIGIN")
		c.Header("Strict-Transport-Security", "max-age=63072000; includeSubDomains")
		c.Header("X-XSS-Protection", "0")
		c.Header("Content-Security-Policy", "default-src 'self'")
		c.Header("Referrer-Policy", "no-referrer")
		c.Header("X-Permitted-Cross-Domain-Policies", "none")
		c.Header("X-Download-Options", "noopen")
		c.Next()
	}
}

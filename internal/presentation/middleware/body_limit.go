package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// BodyLimit restricts the request body size to maxBytes.
// Returns 413 Payload Too Large if the body exceeds the limit.
func BodyLimit(maxBytes int64) gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.Body != nil {
			c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, maxBytes)
		}
		c.Next()
	}
}

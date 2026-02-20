package middleware

import (
	"fmt"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// RequestLogger logs each request with method, path, status, duration, and sets X-Request-Id header.
func RequestLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID := uuid.New().String()
		c.Set("requestId", requestID)
		c.Header("X-Request-Id", requestID)

		start := time.Now()

		c.Next()

		duration := time.Since(start)
		status := c.Writer.Status()
		method := c.Request.Method
		path := c.Request.URL.Path

		if os.Getenv("GIN_MODE") == "release" {
			fmt.Fprintf(os.Stdout,
				`{"level":"info","message":"request","timestamp":"%s","requestId":"%s","method":"%s","path":"%s","status":%d,"duration_ms":%d}`+"\n",
				time.Now().Format(time.RFC3339), requestID, method, path, status, duration.Milliseconds(),
			)
		} else {
			fmt.Fprintf(os.Stdout,
				"\033[2m%s\033[0m \033[32mINFO \033[0m \033[2m[%s]\033[0m %s %s %d %dms\n",
				time.Now().Format(time.RFC3339), requestID[:8], method, path, status, duration.Milliseconds(),
			)
		}
	}
}

package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/lgxju/gogretago/internal/domain/services"
)

// AuthMiddleware validates JWT tokens and sets user context
func AuthMiddleware(jwtService services.JwtService) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Try Authorization header first, then x-auth-token
		authHeader := c.GetHeader("Authorization")
		token := ""
		if strings.HasPrefix(authHeader, "Bearer ") {
			token = authHeader[7:]
		}
		if token == "" {
			token = c.GetHeader("x-auth-token")
		}

		if token == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"error": gin.H{
					"code":    "UNAUTHORIZED",
					"message": "Authentication token is required",
				},
			})
			c.Abort()
			return
		}

		// Remove "Bearer " prefix if still present
		token = strings.TrimPrefix(token, "Bearer ")

		payload, err := jwtService.Verify(token)
		if err != nil || payload == nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"error": gin.H{
					"code":    "INVALID_TOKEN",
					"message": "Invalid or expired token",
				},
			})
			c.Abort()
			return
		}

		c.Set("userId", payload.UserID)
		c.Set("role", payload.Role)
		c.Next()
	}
}

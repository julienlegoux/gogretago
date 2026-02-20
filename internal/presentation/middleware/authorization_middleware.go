package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lgxju/gogretago/internal/domain/authorization"
)

// RequireRole creates a middleware that checks the user's role against the required minimum role.
// Uses hierarchical role system: USER(1) < DRIVER(2) < ADMIN(3).
func RequireRole(roles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userRole, exists := c.Get("role")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"error": gin.H{
					"code":    "UNAUTHORIZED",
					"message": "Authentication required",
				},
			})
			c.Abort()
			return
		}

		roleStr, ok := userRole.(string)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"error": gin.H{
					"code":    "UNAUTHORIZED",
					"message": "Authentication required",
				},
			})
			c.Abort()
			return
		}

		userLevel := authorization.GetRoleLevel(roleStr)

		// Find minimum required level among the specified roles
		minRequired := 999
		for _, r := range roles {
			level := authorization.GetRoleLevel(r)
			if level < minRequired {
				minRequired = level
			}
		}

		if userLevel < minRequired {
			c.JSON(http.StatusForbidden, gin.H{
				"success": false,
				"error": gin.H{
					"code":    "FORBIDDEN",
					"message": "Insufficient permissions",
				},
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

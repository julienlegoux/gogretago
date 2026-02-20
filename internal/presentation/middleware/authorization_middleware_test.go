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

func setupAuthzRouter(role interface{}, setRole bool, requiredRoles ...string) (*gin.Engine, *httptest.ResponseRecorder) {
	router := gin.New()

	// Middleware to simulate auth context
	router.Use(func(c *gin.Context) {
		if setRole {
			c.Set("role", role)
		}
		c.Next()
	})

	router.Use(RequireRole(requiredRoles...))

	router.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"success": true})
	})

	return router, httptest.NewRecorder()
}

func TestRequireRole_AdminAccessesAdminRoute(t *testing.T) {
	router, w := setupAuthzRouter("ADMIN", true, "ADMIN")
	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestRequireRole_AdminAccessesDriverRoute(t *testing.T) {
	router, w := setupAuthzRouter("ADMIN", true, "DRIVER")
	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code, "ADMIN (level 3) should access DRIVER (level 2) route")
}

func TestRequireRole_UserBlockedFromDriverRoute(t *testing.T) {
	router, w := setupAuthzRouter("USER", true, "DRIVER")
	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusForbidden, w.Code)

	var body map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &body)
	require.NoError(t, err)
	assert.Equal(t, false, body["success"])

	errObj := body["error"].(map[string]interface{})
	assert.Equal(t, "FORBIDDEN", errObj["code"])
}

func TestRequireRole_NoRoleInContext(t *testing.T) {
	router, w := setupAuthzRouter(nil, false, "USER")
	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)

	var body map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &body)
	require.NoError(t, err)

	errObj := body["error"].(map[string]interface{})
	assert.Equal(t, "UNAUTHORIZED", errObj["code"])
}

func TestRequireRole_InvalidRoleType(t *testing.T) {
	// Set role as int instead of string
	router, w := setupAuthzRouter(42, true, "USER")
	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)

	var body map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &body)
	require.NoError(t, err)

	errObj := body["error"].(map[string]interface{})
	assert.Equal(t, "UNAUTHORIZED", errObj["code"])
}

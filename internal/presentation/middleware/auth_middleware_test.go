package middleware

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/lgxju/gogretago/internal/domain/services"
	"github.com/lgxju/gogretago/internal/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func init() {
	gin.SetMode(gin.TestMode)
}

func setupAuthTest(t *testing.T, mockJwt *mocks.MockJwtService) (*gin.Engine, *httptest.ResponseRecorder) {
	t.Helper()
	router := gin.New()
	router.Use(AuthMiddleware(mockJwt))
	router.GET("/test", func(c *gin.Context) {
		userId, _ := c.Get("userId")
		role, _ := c.Get("role")
		c.JSON(http.StatusOK, gin.H{"userId": userId, "role": role})
	})
	return router, httptest.NewRecorder()
}

func TestAuthMiddleware_ValidBearerToken(t *testing.T) {
	mockJwt := mocks.NewMockJwtService(t)
	mockJwt.EXPECT().Verify("valid-token").Return(&services.JwtPayload{
		UserID: "user-123",
		Role:   "ADMIN",
	}, nil)

	router, w := setupAuthTest(t, mockJwt)
	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Set("Authorization", "Bearer valid-token")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var body map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &body)
	require.NoError(t, err)
	assert.Equal(t, "user-123", body["userId"])
	assert.Equal(t, "ADMIN", body["role"])
}

func TestAuthMiddleware_ValidXAuthToken(t *testing.T) {
	mockJwt := mocks.NewMockJwtService(t)
	mockJwt.EXPECT().Verify("x-token-value").Return(&services.JwtPayload{
		UserID: "user-456",
		Role:   "DRIVER",
	}, nil)

	router, w := setupAuthTest(t, mockJwt)
	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Set("x-auth-token", "x-token-value")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var body map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &body)
	require.NoError(t, err)
	assert.Equal(t, "user-456", body["userId"])
	assert.Equal(t, "DRIVER", body["role"])
}

func TestAuthMiddleware_MissingToken(t *testing.T) {
	mockJwt := mocks.NewMockJwtService(t)
	// No Verify call expected since no token provided

	router, w := setupAuthTest(t, mockJwt)
	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)

	var body map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &body)
	require.NoError(t, err)
	assert.Equal(t, false, body["success"])

	errObj := body["error"].(map[string]interface{})
	assert.Equal(t, "UNAUTHORIZED", errObj["code"])
}

func TestAuthMiddleware_InvalidToken(t *testing.T) {
	mockJwt := mocks.NewMockJwtService(t)
	mockJwt.EXPECT().Verify("bad-token").Return(nil, fmt.Errorf("token is invalid"))

	router, w := setupAuthTest(t, mockJwt)
	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Set("Authorization", "Bearer bad-token")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)

	var body map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &body)
	require.NoError(t, err)
	assert.Equal(t, false, body["success"])

	errObj := body["error"].(map[string]interface{})
	assert.Equal(t, "INVALID_TOKEN", errObj["code"])
}

func TestAuthMiddleware_NilPayload(t *testing.T) {
	mockJwt := mocks.NewMockJwtService(t)
	mockJwt.EXPECT().Verify("nil-payload-token").Return(nil, nil)

	router, w := setupAuthTest(t, mockJwt)
	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Set("Authorization", "Bearer nil-payload-token")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)

	var body map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &body)
	require.NoError(t, err)

	errObj := body["error"].(map[string]interface{})
	assert.Equal(t, "INVALID_TOKEN", errObj["code"])
}

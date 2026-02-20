package services

import (
	"os"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/lgxju/gogretago/config"
	domainservices "github.com/lgxju/gogretago/internal/domain/services"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const testJWTSecret = "a-very-long-secret-key-for-testing-purposes-at-least-32-chars"

func setupJwtEnv(t *testing.T, expiresIn string) {
	t.Helper()
	os.Setenv("JWT_SECRET", testJWTSecret)
	os.Setenv("JWT_EXPIRES_IN", expiresIn)
	// Reset the config so it reloads from env
	config.Load()
	t.Cleanup(func() {
		os.Unsetenv("JWT_SECRET")
		os.Unsetenv("JWT_EXPIRES_IN")
	})
}

func TestJwtSign_ProducesValidToken(t *testing.T) {
	setupJwtEnv(t, "24h")
	svc := NewJwtService()

	payload := domainservices.JwtPayload{UserID: "user-123", Role: "ADMIN"}
	token, err := svc.Sign(payload)
	require.NoError(t, err)
	assert.NotEmpty(t, token)

	// Token should be parseable
	parsed, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		return []byte(testJWTSecret), nil
	})
	require.NoError(t, err)
	assert.True(t, parsed.Valid)
}

func TestJwtVerify_ValidToken(t *testing.T) {
	setupJwtEnv(t, "24h")
	svc := NewJwtService()

	payload := domainservices.JwtPayload{UserID: "user-456", Role: "DRIVER"}
	token, err := svc.Sign(payload)
	require.NoError(t, err)

	result, err := svc.Verify(token)
	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, "user-456", result.UserID)
	assert.Equal(t, "DRIVER", result.Role)
}

func TestJwtVerify_ExpiredToken(t *testing.T) {
	setupJwtEnv(t, "24h")
	svc := NewJwtService()

	// Manually create an expired token
	claims := jwt.MapClaims{
		"userId": "user-789",
		"role":   "USER",
		"exp":    time.Now().Add(-1 * time.Hour).Unix(),
		"iat":    time.Now().Add(-2 * time.Hour).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(testJWTSecret))
	require.NoError(t, err)

	result, err := svc.Verify(tokenString)
	assert.Error(t, err, "expired token should return error")
	assert.Nil(t, result)
}

func TestJwtVerify_InvalidSignature(t *testing.T) {
	setupJwtEnv(t, "24h")
	svc := NewJwtService()

	// Sign with a different secret
	claims := jwt.MapClaims{
		"userId": "user-789",
		"role":   "USER",
		"exp":    time.Now().Add(1 * time.Hour).Unix(),
		"iat":    time.Now().Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte("a-completely-different-secret-key-than-the-test-one"))
	require.NoError(t, err)

	result, err := svc.Verify(tokenString)
	assert.Error(t, err, "token signed with different secret should fail")
	assert.Nil(t, result)
}

func TestJwtVerify_MalformedToken(t *testing.T) {
	setupJwtEnv(t, "24h")
	svc := NewJwtService()

	result, err := svc.Verify("garbage")
	assert.Error(t, err, "malformed token should return error")
	assert.Nil(t, result)
}

func TestJwtVerify_MissingRoleClaim(t *testing.T) {
	setupJwtEnv(t, "24h")
	svc := NewJwtService()

	// Create token without role claim
	claims := jwt.MapClaims{
		"userId": "user-no-role",
		"exp":    time.Now().Add(1 * time.Hour).Unix(),
		"iat":    time.Now().Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(testJWTSecret))
	require.NoError(t, err)

	result, err := svc.Verify(tokenString)
	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, "user-no-role", result.UserID)
	assert.Equal(t, "USER", result.Role, "missing role should default to USER")
}

func TestJwtSign_ExpirationFormats(t *testing.T) {
	tests := []struct {
		name      string
		expiresIn string
		minDelta  time.Duration
		maxDelta  time.Duration
	}{
		{"24 hours", "24h", 23*time.Hour + 59*time.Minute, 24*time.Hour + 1*time.Minute},
		{"7 days", "7d", 7*24*time.Hour - 1*time.Minute, 7*24*time.Hour + 1*time.Minute},
		{"30 minutes", "30m", 29 * time.Minute, 31 * time.Minute},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			setupJwtEnv(t, tt.expiresIn)
			svc := NewJwtService()

			payload := domainservices.JwtPayload{UserID: "user-exp", Role: "USER"}
			tokenString, err := svc.Sign(payload)
			require.NoError(t, err)

			// Parse the token to inspect exp claim
			parsed, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
				return []byte(testJWTSecret), nil
			})
			require.NoError(t, err)

			claims, ok := parsed.Claims.(jwt.MapClaims)
			require.True(t, ok)

			expFloat, ok := claims["exp"].(float64)
			require.True(t, ok)

			expTime := time.Unix(int64(expFloat), 0)
			delta := time.Until(expTime)

			assert.True(t, delta >= tt.minDelta, "expiration too soon: delta=%v, minDelta=%v", delta, tt.minDelta)
			assert.True(t, delta <= tt.maxDelta, "expiration too far: delta=%v, maxDelta=%v", delta, tt.maxDelta)
		})
	}
}

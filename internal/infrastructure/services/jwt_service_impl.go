package services

import (
	"fmt"
	"regexp"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/lgxju/gogretago/config"
	"github.com/lgxju/gogretago/internal/domain/services"
)

// JwtServiceImpl implements JwtService using golang-jwt
type JwtServiceImpl struct {
	secret    string
	expiresIn string
}

// NewJwtService creates a new JwtServiceImpl
func NewJwtService() services.JwtService {
	cfg := config.Get()
	return &JwtServiceImpl{
		secret:    cfg.JWTSecret,
		expiresIn: cfg.JWTExpiresIn,
	}
}

// Sign creates a new JWT token
func (s *JwtServiceImpl) Sign(payload services.JwtPayload) (string, error) {
	exp := s.calculateExpiration()

	claims := jwt.MapClaims{
		"userId": payload.UserID,
		"exp":    exp,
		"iat":    time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.secret))
}

// Verify validates a JWT token and returns the payload
func (s *JwtServiceImpl) Verify(tokenString string) (*services.JwtPayload, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(s.secret), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		userID, ok := claims["userId"].(string)
		if !ok {
			return nil, fmt.Errorf("invalid token payload")
		}
		return &services.JwtPayload{UserID: userID}, nil
	}

	return nil, fmt.Errorf("invalid token")
}

// calculateExpiration parses the expiresIn string and returns Unix timestamp
func (s *JwtServiceImpl) calculateExpiration() int64 {
	now := time.Now()
	re := regexp.MustCompile(`^(\d+)(h|d|m)$`)
	match := re.FindStringSubmatch(s.expiresIn)

	if len(match) != 3 {
		return now.Add(24 * time.Hour).Unix()
	}

	value, _ := strconv.Atoi(match[1])
	unit := match[2]

	switch unit {
	case "h":
		return now.Add(time.Duration(value) * time.Hour).Unix()
	case "d":
		return now.Add(time.Duration(value) * 24 * time.Hour).Unix()
	case "m":
		return now.Add(time.Duration(value) * time.Minute).Unix()
	default:
		return now.Add(24 * time.Hour).Unix()
	}
}

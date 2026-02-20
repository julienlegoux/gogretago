package services

// JwtPayload represents the JWT token payload
type JwtPayload struct {
	UserID string
	Role   string
}

// JwtService defines the interface for JWT operations
type JwtService interface {
	Sign(payload JwtPayload) (string, error)
	Verify(token string) (*JwtPayload, error)
}

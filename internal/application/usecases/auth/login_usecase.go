package auth

import (
	"context"

	"github.com/lgxju/gogretago/internal/application/dtos"
	domainerrors "github.com/lgxju/gogretago/internal/domain/errors"
	"github.com/lgxju/gogretago/internal/domain/repositories"
	"github.com/lgxju/gogretago/internal/domain/services"
)

// LoginUseCase handles user login logic
type LoginUseCase struct {
	userRepository  repositories.UserRepository
	passwordService services.PasswordService
	jwtService      services.JwtService
}

// NewLoginUseCase creates a new LoginUseCase with its dependencies
func NewLoginUseCase(
	userRepository repositories.UserRepository,
	passwordService services.PasswordService,
	jwtService services.JwtService,
) *LoginUseCase {
	return &LoginUseCase{
		userRepository:  userRepository,
		passwordService: passwordService,
		jwtService:      jwtService,
	}
}

// Execute performs the login use case
func (uc *LoginUseCase) Execute(ctx context.Context, input dtos.LoginInput) (*dtos.AuthResponse, error) {
	// Find user by email
	user, err := uc.userRepository.FindByEmail(ctx, input.Email)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, domainerrors.NewInvalidCredentialsError()
	}

	// Verify password
	valid, err := uc.passwordService.Verify(input.Password, user.Password)
	if err != nil {
		return nil, err
	}
	if !valid {
		return nil, domainerrors.NewInvalidCredentialsError()
	}

	// Generate JWT token
	token, err := uc.jwtService.Sign(services.JwtPayload{UserID: user.ID})
	if err != nil {
		return nil, err
	}

	return &dtos.AuthResponse{
		UserID: user.ID,
		Token:  token,
	}, nil
}

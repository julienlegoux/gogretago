package auth

import (
	"context"

	"github.com/lgxju/gogretago/internal/application/dtos"
	"github.com/lgxju/gogretago/internal/domain/entities"
	domainerrors "github.com/lgxju/gogretago/internal/domain/errors"
	"github.com/lgxju/gogretago/internal/domain/repositories"
	"github.com/lgxju/gogretago/internal/domain/services"
)

// RegisterUseCase handles user registration logic
type RegisterUseCase struct {
	userRepository  repositories.UserRepository
	passwordService services.PasswordService
	emailService    services.EmailService
	jwtService      services.JwtService
}

// NewRegisterUseCase creates a new RegisterUseCase with its dependencies
func NewRegisterUseCase(
	userRepository repositories.UserRepository,
	passwordService services.PasswordService,
	emailService services.EmailService,
	jwtService services.JwtService,
) *RegisterUseCase {
	return &RegisterUseCase{
		userRepository:  userRepository,
		passwordService: passwordService,
		emailService:    emailService,
		jwtService:      jwtService,
	}
}

// Execute performs the registration use case
func (uc *RegisterUseCase) Execute(ctx context.Context, input dtos.RegisterInput) (*dtos.AuthResponse, error) {
	// Check if user already exists
	exists, err := uc.userRepository.ExistsByEmail(ctx, input.Email)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, domainerrors.NewUserAlreadyExistsError(input.Email)
	}

	// Hash password
	hashedPassword, err := uc.passwordService.Hash(input.Password)
	if err != nil {
		return nil, err
	}

	// Create user
	user, err := uc.userRepository.Create(ctx, entities.CreateUserData{
		Email:     input.Email,
		Password:  hashedPassword,
		FirstName: input.FirstName,
		LastName:  input.LastName,
		Phone:     input.Phone,
	})
	if err != nil {
		return nil, err
	}

	// Send welcome email (fire and forget, don't fail registration)
	_ = uc.emailService.SendWelcomeEmail(user.Email, user.FirstName)

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

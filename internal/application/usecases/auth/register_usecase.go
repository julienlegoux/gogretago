package auth

import (
	"context"

	"github.com/lgxju/gogretago/internal/application/dtos"
	"github.com/lgxju/gogretago/internal/domain/entities"
	domainerrors "github.com/lgxju/gogretago/internal/domain/errors"
	"github.com/lgxju/gogretago/internal/domain/repositories"
	"github.com/lgxju/gogretago/internal/domain/services"
)

type RegisterUseCase struct {
	authRepository  repositories.AuthRepository
	passwordService services.PasswordService
	emailService    services.EmailService
	jwtService      services.JwtService
}

func NewRegisterUseCase(
	authRepository repositories.AuthRepository,
	passwordService services.PasswordService,
	emailService services.EmailService,
	jwtService services.JwtService,
) *RegisterUseCase {
	return &RegisterUseCase{
		authRepository:  authRepository,
		passwordService: passwordService,
		emailService:    emailService,
		jwtService:      jwtService,
	}
}

func (uc *RegisterUseCase) Execute(ctx context.Context, input dtos.RegisterInput) (*dtos.AuthResponse, error) {
	exists, err := uc.authRepository.ExistsByEmail(ctx, input.Email)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, domainerrors.NewUserAlreadyExistsError(input.Email)
	}

	hashedPassword, err := uc.passwordService.Hash(input.Password)
	if err != nil {
		return nil, err
	}

	auth, user, err := uc.authRepository.CreateWithUser(ctx,
		entities.CreateAuthData{Email: input.Email, Password: hashedPassword},
		entities.CreateUserData{FirstName: nil, LastName: nil, Phone: nil},
	)
	if err != nil {
		return nil, err
	}

	// Send welcome email (fire and forget)
	name := "there"
	if user.FirstName != nil {
		name = *user.FirstName
	}
	_ = uc.emailService.SendWelcomeEmail(auth.Email, name)

	token, err := uc.jwtService.Sign(services.JwtPayload{UserID: user.ID, Role: auth.Role})
	if err != nil {
		return nil, err
	}

	return &dtos.AuthResponse{
		UserID: user.ID,
		Token:  token,
	}, nil
}

package auth

import (
	"context"

	"github.com/lgxju/gogretago/internal/application/dtos"
	domainerrors "github.com/lgxju/gogretago/internal/domain/errors"
	"github.com/lgxju/gogretago/internal/domain/repositories"
	"github.com/lgxju/gogretago/internal/domain/services"
)

type LoginUseCase struct {
	authRepository  repositories.AuthRepository
	userRepository  repositories.UserRepository
	passwordService services.PasswordService
	jwtService      services.JwtService
}

func NewLoginUseCase(
	authRepository repositories.AuthRepository,
	userRepository repositories.UserRepository,
	passwordService services.PasswordService,
	jwtService services.JwtService,
) *LoginUseCase {
	return &LoginUseCase{
		authRepository:  authRepository,
		userRepository:  userRepository,
		passwordService: passwordService,
		jwtService:      jwtService,
	}
}

func (uc *LoginUseCase) Execute(ctx context.Context, input dtos.LoginInput) (*dtos.AuthResponse, error) {
	auth, err := uc.authRepository.FindByEmail(ctx, input.Email)
	if err != nil {
		return nil, err
	}
	if auth == nil {
		return nil, domainerrors.NewInvalidCredentialsError()
	}

	valid, err := uc.passwordService.Verify(input.Password, auth.Password)
	if err != nil {
		return nil, err
	}
	if !valid {
		return nil, domainerrors.NewInvalidCredentialsError()
	}

	user, err := uc.userRepository.FindByAuthRefID(ctx, auth.RefID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, domainerrors.NewInvalidCredentialsError()
	}

	token, err := uc.jwtService.Sign(services.JwtPayload{UserID: user.ID, Role: auth.Role})
	if err != nil {
		return nil, err
	}

	return &dtos.AuthResponse{
		UserID: user.ID,
		Token:  token,
	}, nil
}

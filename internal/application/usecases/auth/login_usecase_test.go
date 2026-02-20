package auth

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/lgxju/gogretago/internal/application/dtos"
	"github.com/lgxju/gogretago/internal/domain/entities"
	domainerrors "github.com/lgxju/gogretago/internal/domain/errors"
	"github.com/lgxju/gogretago/internal/domain/services"
	"github.com/lgxju/gogretago/internal/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestLogin_Success(t *testing.T) {
	ctx := context.Background()
	authRepo := mocks.NewMockAuthRepository(t)
	userRepo := mocks.NewMockUserRepository(t)
	passwordSvc := mocks.NewMockPasswordService(t)
	jwtSvc := mocks.NewMockJwtService(t)

	auth := &entities.Auth{
		ID:       "auth-1",
		RefID:    100,
		Email:    "user@example.com",
		Password: "hashed-password",
		Role:     "USER",
	}
	user := &entities.PublicUser{
		User: entities.User{
			ID:        "user-1",
			RefID:     200,
			AuthRefID: 100,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		Email: "user@example.com",
	}

	authRepo.EXPECT().FindByEmail(ctx, "user@example.com").Return(auth, nil)
	passwordSvc.EXPECT().Verify("secret123", "hashed-password").Return(true, nil)
	userRepo.EXPECT().FindByAuthRefID(ctx, int64(100)).Return(user, nil)
	jwtSvc.EXPECT().Sign(services.JwtPayload{UserID: "user-1", Role: "USER"}).Return("jwt-token", nil)

	uc := NewLoginUseCase(authRepo, userRepo, passwordSvc, jwtSvc)
	result, err := uc.Execute(ctx, dtos.LoginInput{
		Email:    "user@example.com",
		Password: "secret123",
	})

	require.NoError(t, err)
	assert.Equal(t, "user-1", result.UserID)
	assert.Equal(t, "jwt-token", result.Token)
}

func TestLogin_EmailNotFound(t *testing.T) {
	ctx := context.Background()
	authRepo := mocks.NewMockAuthRepository(t)
	userRepo := mocks.NewMockUserRepository(t)
	passwordSvc := mocks.NewMockPasswordService(t)
	jwtSvc := mocks.NewMockJwtService(t)

	authRepo.EXPECT().FindByEmail(ctx, "unknown@example.com").Return(nil, nil)

	uc := NewLoginUseCase(authRepo, userRepo, passwordSvc, jwtSvc)
	result, err := uc.Execute(ctx, dtos.LoginInput{
		Email:    "unknown@example.com",
		Password: "secret123",
	})

	assert.Nil(t, result)
	require.Error(t, err)
	var credErr *domainerrors.InvalidCredentialsError
	assert.True(t, errors.As(err, &credErr))
}

func TestLogin_WrongPassword(t *testing.T) {
	ctx := context.Background()
	authRepo := mocks.NewMockAuthRepository(t)
	userRepo := mocks.NewMockUserRepository(t)
	passwordSvc := mocks.NewMockPasswordService(t)
	jwtSvc := mocks.NewMockJwtService(t)

	auth := &entities.Auth{
		ID:       "auth-1",
		RefID:    100,
		Email:    "user@example.com",
		Password: "hashed-password",
		Role:     "USER",
	}

	authRepo.EXPECT().FindByEmail(ctx, "user@example.com").Return(auth, nil)
	passwordSvc.EXPECT().Verify("wrong-password", "hashed-password").Return(false, nil)

	uc := NewLoginUseCase(authRepo, userRepo, passwordSvc, jwtSvc)
	result, err := uc.Execute(ctx, dtos.LoginInput{
		Email:    "user@example.com",
		Password: "wrong-password",
	})

	assert.Nil(t, result)
	require.Error(t, err)
	var credErr *domainerrors.InvalidCredentialsError
	assert.True(t, errors.As(err, &credErr))
}

func TestLogin_UserProfileMissing(t *testing.T) {
	ctx := context.Background()
	authRepo := mocks.NewMockAuthRepository(t)
	userRepo := mocks.NewMockUserRepository(t)
	passwordSvc := mocks.NewMockPasswordService(t)
	jwtSvc := mocks.NewMockJwtService(t)

	auth := &entities.Auth{
		ID:       "auth-1",
		RefID:    100,
		Email:    "user@example.com",
		Password: "hashed-password",
		Role:     "USER",
	}

	authRepo.EXPECT().FindByEmail(ctx, "user@example.com").Return(auth, nil)
	passwordSvc.EXPECT().Verify("secret123", "hashed-password").Return(true, nil)
	userRepo.EXPECT().FindByAuthRefID(ctx, int64(100)).Return(nil, nil)

	uc := NewLoginUseCase(authRepo, userRepo, passwordSvc, jwtSvc)
	result, err := uc.Execute(ctx, dtos.LoginInput{
		Email:    "user@example.com",
		Password: "secret123",
	})

	assert.Nil(t, result)
	require.Error(t, err)
	var credErr *domainerrors.InvalidCredentialsError
	assert.True(t, errors.As(err, &credErr))
}

func TestLogin_AuthRepoError(t *testing.T) {
	ctx := context.Background()
	authRepo := mocks.NewMockAuthRepository(t)
	userRepo := mocks.NewMockUserRepository(t)
	passwordSvc := mocks.NewMockPasswordService(t)
	jwtSvc := mocks.NewMockJwtService(t)

	repoErr := errors.New("database connection failed")
	authRepo.EXPECT().FindByEmail(ctx, "user@example.com").Return(nil, repoErr)

	uc := NewLoginUseCase(authRepo, userRepo, passwordSvc, jwtSvc)
	result, err := uc.Execute(ctx, dtos.LoginInput{
		Email:    "user@example.com",
		Password: "secret123",
	})

	assert.Nil(t, result)
	require.Error(t, err)
	assert.Equal(t, repoErr, err)
}

func TestLogin_PasswordServiceError(t *testing.T) {
	ctx := context.Background()
	authRepo := mocks.NewMockAuthRepository(t)
	userRepo := mocks.NewMockUserRepository(t)
	passwordSvc := mocks.NewMockPasswordService(t)
	jwtSvc := mocks.NewMockJwtService(t)

	auth := &entities.Auth{
		ID:       "auth-1",
		RefID:    100,
		Email:    "user@example.com",
		Password: "hashed-password",
		Role:     "USER",
	}

	svcErr := errors.New("bcrypt internal error")
	authRepo.EXPECT().FindByEmail(ctx, "user@example.com").Return(auth, nil)
	passwordSvc.EXPECT().Verify("secret123", "hashed-password").Return(false, svcErr)

	uc := NewLoginUseCase(authRepo, userRepo, passwordSvc, jwtSvc)
	result, err := uc.Execute(ctx, dtos.LoginInput{
		Email:    "user@example.com",
		Password: "secret123",
	})

	assert.Nil(t, result)
	require.Error(t, err)
	assert.Equal(t, svcErr, err)
}

func TestLogin_JwtSignError(t *testing.T) {
	ctx := context.Background()
	authRepo := mocks.NewMockAuthRepository(t)
	userRepo := mocks.NewMockUserRepository(t)
	passwordSvc := mocks.NewMockPasswordService(t)
	jwtSvc := mocks.NewMockJwtService(t)

	auth := &entities.Auth{
		ID:       "auth-1",
		RefID:    100,
		Email:    "user@example.com",
		Password: "hashed-password",
		Role:     "USER",
	}
	user := &entities.PublicUser{
		User: entities.User{
			ID:        "user-1",
			RefID:     200,
			AuthRefID: 100,
		},
		Email: "user@example.com",
	}

	jwtErr := errors.New("signing key not found")
	authRepo.EXPECT().FindByEmail(ctx, "user@example.com").Return(auth, nil)
	passwordSvc.EXPECT().Verify("secret123", "hashed-password").Return(true, nil)
	userRepo.EXPECT().FindByAuthRefID(ctx, int64(100)).Return(user, nil)
	jwtSvc.EXPECT().Sign(services.JwtPayload{UserID: "user-1", Role: "USER"}).Return("", jwtErr)

	uc := NewLoginUseCase(authRepo, userRepo, passwordSvc, jwtSvc)
	result, err := uc.Execute(ctx, dtos.LoginInput{
		Email:    "user@example.com",
		Password: "secret123",
	})

	assert.Nil(t, result)
	require.Error(t, err)
	assert.Equal(t, jwtErr, err)
}

// Ensure mock is used (silence unused import warnings if needed)
var _ mock.TestingT = (*testing.T)(nil)

package auth

import (
	"context"
	"errors"
	"testing"

	"github.com/lgxju/gogretago/internal/application/dtos"
	"github.com/lgxju/gogretago/internal/domain/entities"
	domainerrors "github.com/lgxju/gogretago/internal/domain/errors"
	"github.com/lgxju/gogretago/internal/domain/services"
	"github.com/lgxju/gogretago/internal/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestRegister_Success(t *testing.T) {
	ctx := context.Background()
	authRepo := mocks.NewMockAuthRepository(t)
	passwordSvc := mocks.NewMockPasswordService(t)
	emailSvc := mocks.NewMockEmailService(t)
	jwtSvc := mocks.NewMockJwtService(t)

	firstName := "John"
	auth := &entities.Auth{
		ID:    "auth-1",
		RefID: 100,
		Email: "new@example.com",
		Role:  "USER",
	}
	user := &entities.PublicUser{
		User: entities.User{
			ID:        "user-1",
			RefID:     200,
			FirstName: &firstName,
			AuthRefID: 100,
		},
		Email: "new@example.com",
	}

	authRepo.EXPECT().ExistsByEmail(ctx, "new@example.com").Return(false, nil)
	passwordSvc.EXPECT().Hash("password123").Return("hashed-pw", nil)
	authRepo.EXPECT().CreateWithUser(ctx,
		entities.CreateAuthData{Email: "new@example.com", Password: "hashed-pw"},
		entities.CreateUserData{FirstName: nil, LastName: nil, Phone: nil},
	).Return(auth, user, nil)
	emailSvc.EXPECT().SendWelcomeEmail("new@example.com", "John").Return(nil)
	jwtSvc.EXPECT().Sign(services.JwtPayload{UserID: "user-1", Role: "USER"}).Return("jwt-token", nil)

	uc := NewRegisterUseCase(authRepo, passwordSvc, emailSvc, jwtSvc)
	result, err := uc.Execute(ctx, dtos.RegisterInput{
		Email:           "new@example.com",
		Password:        "password123",
		ConfirmPassword: "password123",
	})

	require.NoError(t, err)
	assert.Equal(t, "user-1", result.UserID)
	assert.Equal(t, "jwt-token", result.Token)
}

func TestRegister_EmailAlreadyExists(t *testing.T) {
	ctx := context.Background()
	authRepo := mocks.NewMockAuthRepository(t)
	passwordSvc := mocks.NewMockPasswordService(t)
	emailSvc := mocks.NewMockEmailService(t)
	jwtSvc := mocks.NewMockJwtService(t)

	authRepo.EXPECT().ExistsByEmail(ctx, "existing@example.com").Return(true, nil)

	uc := NewRegisterUseCase(authRepo, passwordSvc, emailSvc, jwtSvc)
	result, err := uc.Execute(ctx, dtos.RegisterInput{
		Email:           "existing@example.com",
		Password:        "password123",
		ConfirmPassword: "password123",
	})

	assert.Nil(t, result)
	require.Error(t, err)
	var alreadyExistsErr *domainerrors.UserAlreadyExistsError
	assert.True(t, errors.As(err, &alreadyExistsErr))
}

func TestRegister_PasswordHashError(t *testing.T) {
	ctx := context.Background()
	authRepo := mocks.NewMockAuthRepository(t)
	passwordSvc := mocks.NewMockPasswordService(t)
	emailSvc := mocks.NewMockEmailService(t)
	jwtSvc := mocks.NewMockJwtService(t)

	hashErr := errors.New("hashing failed")
	authRepo.EXPECT().ExistsByEmail(ctx, "new@example.com").Return(false, nil)
	passwordSvc.EXPECT().Hash("password123").Return("", hashErr)

	uc := NewRegisterUseCase(authRepo, passwordSvc, emailSvc, jwtSvc)
	result, err := uc.Execute(ctx, dtos.RegisterInput{
		Email:           "new@example.com",
		Password:        "password123",
		ConfirmPassword: "password123",
	})

	assert.Nil(t, result)
	require.Error(t, err)
	assert.Equal(t, hashErr, err)
}

func TestRegister_CreateWithUserError(t *testing.T) {
	ctx := context.Background()
	authRepo := mocks.NewMockAuthRepository(t)
	passwordSvc := mocks.NewMockPasswordService(t)
	emailSvc := mocks.NewMockEmailService(t)
	jwtSvc := mocks.NewMockJwtService(t)

	createErr := errors.New("database constraint violation")
	authRepo.EXPECT().ExistsByEmail(ctx, "new@example.com").Return(false, nil)
	passwordSvc.EXPECT().Hash("password123").Return("hashed-pw", nil)
	authRepo.EXPECT().CreateWithUser(ctx,
		entities.CreateAuthData{Email: "new@example.com", Password: "hashed-pw"},
		entities.CreateUserData{FirstName: nil, LastName: nil, Phone: nil},
	).Return(nil, nil, createErr)

	uc := NewRegisterUseCase(authRepo, passwordSvc, emailSvc, jwtSvc)
	result, err := uc.Execute(ctx, dtos.RegisterInput{
		Email:           "new@example.com",
		Password:        "password123",
		ConfirmPassword: "password123",
	})

	assert.Nil(t, result)
	require.Error(t, err)
	assert.Equal(t, createErr, err)
}

func TestRegister_EmailFailureDoesNotBlock(t *testing.T) {
	ctx := context.Background()
	authRepo := mocks.NewMockAuthRepository(t)
	passwordSvc := mocks.NewMockPasswordService(t)
	emailSvc := mocks.NewMockEmailService(t)
	jwtSvc := mocks.NewMockJwtService(t)

	firstName := "Jane"
	auth := &entities.Auth{
		ID:    "auth-1",
		RefID: 100,
		Email: "new@example.com",
		Role:  "USER",
	}
	user := &entities.PublicUser{
		User: entities.User{
			ID:        "user-1",
			RefID:     200,
			FirstName: &firstName,
			AuthRefID: 100,
		},
		Email: "new@example.com",
	}

	authRepo.EXPECT().ExistsByEmail(ctx, "new@example.com").Return(false, nil)
	passwordSvc.EXPECT().Hash("password123").Return("hashed-pw", nil)
	authRepo.EXPECT().CreateWithUser(ctx,
		entities.CreateAuthData{Email: "new@example.com", Password: "hashed-pw"},
		entities.CreateUserData{FirstName: nil, LastName: nil, Phone: nil},
	).Return(auth, user, nil)
	emailSvc.EXPECT().SendWelcomeEmail("new@example.com", "Jane").Return(errors.New("SMTP failure"))
	jwtSvc.EXPECT().Sign(services.JwtPayload{UserID: "user-1", Role: "USER"}).Return("jwt-token", nil)

	uc := NewRegisterUseCase(authRepo, passwordSvc, emailSvc, jwtSvc)
	result, err := uc.Execute(ctx, dtos.RegisterInput{
		Email:           "new@example.com",
		Password:        "password123",
		ConfirmPassword: "password123",
	})

	require.NoError(t, err)
	assert.Equal(t, "user-1", result.UserID)
	assert.Equal(t, "jwt-token", result.Token)
}

func TestRegister_JwtSignError(t *testing.T) {
	ctx := context.Background()
	authRepo := mocks.NewMockAuthRepository(t)
	passwordSvc := mocks.NewMockPasswordService(t)
	emailSvc := mocks.NewMockEmailService(t)
	jwtSvc := mocks.NewMockJwtService(t)

	auth := &entities.Auth{
		ID:    "auth-1",
		RefID: 100,
		Email: "new@example.com",
		Role:  "USER",
	}
	user := &entities.PublicUser{
		User: entities.User{
			ID:        "user-1",
			RefID:     200,
			FirstName: nil,
			AuthRefID: 100,
		},
		Email: "new@example.com",
	}

	jwtErr := errors.New("signing key error")
	authRepo.EXPECT().ExistsByEmail(ctx, "new@example.com").Return(false, nil)
	passwordSvc.EXPECT().Hash("password123").Return("hashed-pw", nil)
	authRepo.EXPECT().CreateWithUser(ctx,
		entities.CreateAuthData{Email: "new@example.com", Password: "hashed-pw"},
		entities.CreateUserData{FirstName: nil, LastName: nil, Phone: nil},
	).Return(auth, user, nil)
	emailSvc.EXPECT().SendWelcomeEmail("new@example.com", "there").Return(nil)
	jwtSvc.EXPECT().Sign(services.JwtPayload{UserID: "user-1", Role: "USER"}).Return("", jwtErr)

	uc := NewRegisterUseCase(authRepo, passwordSvc, emailSvc, jwtSvc)
	result, err := uc.Execute(ctx, dtos.RegisterInput{
		Email:           "new@example.com",
		Password:        "password123",
		ConfirmPassword: "password123",
	})

	assert.Nil(t, result)
	require.Error(t, err)
	assert.Equal(t, jwtErr, err)
}

func TestRegister_WelcomeEmailUsesDefaultName(t *testing.T) {
	ctx := context.Background()
	authRepo := mocks.NewMockAuthRepository(t)
	passwordSvc := mocks.NewMockPasswordService(t)
	emailSvc := mocks.NewMockEmailService(t)
	jwtSvc := mocks.NewMockJwtService(t)

	auth := &entities.Auth{
		ID:    "auth-1",
		RefID: 100,
		Email: "new@example.com",
		Role:  "USER",
	}
	user := &entities.PublicUser{
		User: entities.User{
			ID:        "user-1",
			RefID:     200,
			FirstName: nil, // no first name set
			AuthRefID: 100,
		},
		Email: "new@example.com",
	}

	authRepo.EXPECT().ExistsByEmail(ctx, "new@example.com").Return(false, nil)
	passwordSvc.EXPECT().Hash("password123").Return("hashed-pw", nil)
	authRepo.EXPECT().CreateWithUser(ctx,
		entities.CreateAuthData{Email: "new@example.com", Password: "hashed-pw"},
		entities.CreateUserData{FirstName: nil, LastName: nil, Phone: nil},
	).Return(auth, user, nil)
	// Expect "there" as the default name when FirstName is nil
	emailSvc.EXPECT().SendWelcomeEmail("new@example.com", "there").Return(nil)
	jwtSvc.EXPECT().Sign(services.JwtPayload{UserID: "user-1", Role: "USER"}).Return("jwt-token", nil)

	uc := NewRegisterUseCase(authRepo, passwordSvc, emailSvc, jwtSvc)
	result, err := uc.Execute(ctx, dtos.RegisterInput{
		Email:           "new@example.com",
		Password:        "password123",
		ConfirmPassword: "password123",
	})

	require.NoError(t, err)
	assert.Equal(t, "user-1", result.UserID)
	assert.Equal(t, "jwt-token", result.Token)
}

// Ensure mock import is used
var _ mock.TestingT = (*testing.T)(nil)

package user

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/lgxju/gogretago/internal/domain/entities"
	domainerrors "github.com/lgxju/gogretago/internal/domain/errors"
	"github.com/lgxju/gogretago/internal/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetUser_Success(t *testing.T) {
	ctx := context.Background()
	userRepo := mocks.NewMockUserRepository(t)

	firstName := "John"
	lastName := "Doe"
	user := &entities.PublicUser{
		User: entities.User{
			ID:           "user-1",
			RefID:        200,
			FirstName:    &firstName,
			LastName:     &lastName,
			AnonymizedAt: nil,
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
		},
		Email: "john@example.com",
	}

	userRepo.EXPECT().FindByID(ctx, "user-1").Return(user, nil)

	uc := NewGetUserUseCase(userRepo)
	result, err := uc.Execute(ctx, "user-1")

	require.NoError(t, err)
	assert.Equal(t, "user-1", result.ID)
	assert.Equal(t, "john@example.com", result.Email)
	assert.Equal(t, "John", *result.FirstName)
}

func TestGetUser_NotFound(t *testing.T) {
	ctx := context.Background()
	userRepo := mocks.NewMockUserRepository(t)

	userRepo.EXPECT().FindByID(ctx, "nonexistent").Return(nil, nil)

	uc := NewGetUserUseCase(userRepo)
	result, err := uc.Execute(ctx, "nonexistent")

	assert.Nil(t, result)
	require.Error(t, err)
	var notFoundErr *domainerrors.UserNotFoundError
	assert.True(t, errors.As(err, &notFoundErr))
}

func TestGetUser_AnonymizedUser(t *testing.T) {
	ctx := context.Background()
	userRepo := mocks.NewMockUserRepository(t)

	anonymizedAt := time.Now()
	user := &entities.PublicUser{
		User: entities.User{
			ID:           "user-1",
			RefID:        200,
			AnonymizedAt: &anonymizedAt,
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
		},
		Email: "anon@example.com",
	}

	userRepo.EXPECT().FindByID(ctx, "user-1").Return(user, nil)

	uc := NewGetUserUseCase(userRepo)
	result, err := uc.Execute(ctx, "user-1")

	assert.Nil(t, result)
	require.Error(t, err)
	var notFoundErr *domainerrors.UserNotFoundError
	assert.True(t, errors.As(err, &notFoundErr))
}

func TestGetUser_RepoError(t *testing.T) {
	ctx := context.Background()
	userRepo := mocks.NewMockUserRepository(t)

	repoErr := errors.New("database timeout")
	userRepo.EXPECT().FindByID(ctx, "user-1").Return(nil, repoErr)

	uc := NewGetUserUseCase(userRepo)
	result, err := uc.Execute(ctx, "user-1")

	assert.Nil(t, result)
	require.Error(t, err)
	assert.Equal(t, repoErr, err)
}

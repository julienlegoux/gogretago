package user

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/lgxju/gogretago/internal/application/dtos"
	"github.com/lgxju/gogretago/internal/domain/entities"
	domainerrors "github.com/lgxju/gogretago/internal/domain/errors"
	"github.com/lgxju/gogretago/internal/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUpdateUser_Success(t *testing.T) {
	ctx := context.Background()
	userRepo := mocks.NewMockUserRepository(t)

	existingFirstName := "Old"
	existingUser := &entities.PublicUser{
		User: entities.User{
			ID:        "user-1",
			RefID:     200,
			FirstName: &existingFirstName,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		Email: "user@example.com",
	}

	input := dtos.UpdateProfileInput{
		FirstName: "John",
		LastName:  "Doe",
		Phone:     "1234567890",
	}

	updatedFirstName := "John"
	updatedLastName := "Doe"
	updatedPhone := "1234567890"
	updatedUser := &entities.PublicUser{
		User: entities.User{
			ID:        "user-1",
			RefID:     200,
			FirstName: &updatedFirstName,
			LastName:  &updatedLastName,
			Phone:     &updatedPhone,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		Email: "user@example.com",
	}

	userRepo.EXPECT().FindByID(ctx, "user-1").Return(existingUser, nil)
	userRepo.EXPECT().Update(ctx, "user-1", entities.UpdateUserData{
		FirstName: &input.FirstName,
		LastName:  &input.LastName,
		Phone:     &input.Phone,
	}).Return(updatedUser, nil)

	uc := NewUpdateUserUseCase(userRepo)
	result, err := uc.Execute(ctx, "user-1", input)

	require.NoError(t, err)
	assert.Equal(t, "user-1", result.ID)
	assert.Equal(t, "John", *result.FirstName)
	assert.Equal(t, "Doe", *result.LastName)
	assert.Equal(t, "1234567890", *result.Phone)
}

func TestUpdateUser_NotFound(t *testing.T) {
	ctx := context.Background()
	userRepo := mocks.NewMockUserRepository(t)

	userRepo.EXPECT().FindByID(ctx, "nonexistent").Return(nil, nil)

	uc := NewUpdateUserUseCase(userRepo)
	result, err := uc.Execute(ctx, "nonexistent", dtos.UpdateProfileInput{
		FirstName: "John",
		LastName:  "Doe",
		Phone:     "1234567890",
	})

	assert.Nil(t, result)
	require.Error(t, err)
	var notFoundErr *domainerrors.UserNotFoundError
	assert.True(t, errors.As(err, &notFoundErr))
}

func TestUpdateUser_RepoError(t *testing.T) {
	ctx := context.Background()
	userRepo := mocks.NewMockUserRepository(t)

	existingUser := &entities.PublicUser{
		User: entities.User{
			ID:        "user-1",
			RefID:     200,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		Email: "user@example.com",
	}

	input := dtos.UpdateProfileInput{
		FirstName: "John",
		LastName:  "Doe",
		Phone:     "1234567890",
	}

	repoErr := errors.New("update failed")
	userRepo.EXPECT().FindByID(ctx, "user-1").Return(existingUser, nil)
	userRepo.EXPECT().Update(ctx, "user-1", entities.UpdateUserData{
		FirstName: &input.FirstName,
		LastName:  &input.LastName,
		Phone:     &input.Phone,
	}).Return(nil, repoErr)

	uc := NewUpdateUserUseCase(userRepo)
	result, err := uc.Execute(ctx, "user-1", input)

	assert.Nil(t, result)
	require.Error(t, err)
	assert.Equal(t, repoErr, err)
}

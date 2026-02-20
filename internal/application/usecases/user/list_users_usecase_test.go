package user

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/lgxju/gogretago/internal/domain/entities"
	"github.com/lgxju/gogretago/internal/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestListUsers_Success(t *testing.T) {
	ctx := context.Background()
	userRepo := mocks.NewMockUserRepository(t)

	firstName := "Alice"
	users := []entities.PublicUser{
		{
			User: entities.User{
				ID:        "user-1",
				RefID:     100,
				FirstName: &firstName,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			Email: "alice@example.com",
		},
		{
			User: entities.User{
				ID:        "user-2",
				RefID:     101,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			Email: "bob@example.com",
		},
	}

	userRepo.EXPECT().FindAll(ctx).Return(users, nil)

	uc := NewListUsersUseCase(userRepo)
	result, err := uc.Execute(ctx)

	require.NoError(t, err)
	assert.Len(t, result, 2)
	assert.Equal(t, "user-1", result[0].ID)
	assert.Equal(t, "user-2", result[1].ID)
}

func TestListUsers_Empty(t *testing.T) {
	ctx := context.Background()
	userRepo := mocks.NewMockUserRepository(t)

	userRepo.EXPECT().FindAll(ctx).Return([]entities.PublicUser{}, nil)

	uc := NewListUsersUseCase(userRepo)
	result, err := uc.Execute(ctx)

	require.NoError(t, err)
	assert.Empty(t, result)
}

func TestListUsers_RepoError(t *testing.T) {
	ctx := context.Background()
	userRepo := mocks.NewMockUserRepository(t)

	repoErr := errors.New("connection refused")
	userRepo.EXPECT().FindAll(ctx).Return(nil, repoErr)

	uc := NewListUsersUseCase(userRepo)
	result, err := uc.Execute(ctx)

	assert.Nil(t, result)
	require.Error(t, err)
	assert.Equal(t, repoErr, err)
}

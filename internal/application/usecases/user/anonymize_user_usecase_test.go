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

func TestAnonymizeUser_Success(t *testing.T) {
	ctx := context.Background()
	userRepo := mocks.NewMockUserRepository(t)

	user := &entities.PublicUser{
		User: entities.User{
			ID:        "user-1",
			RefID:     200,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		Email: "user@example.com",
	}

	userRepo.EXPECT().FindByID(ctx, "user-1").Return(user, nil)
	userRepo.EXPECT().Anonymize(ctx, "user-1").Return(nil)

	uc := NewAnonymizeUserUseCase(userRepo)
	err := uc.Execute(ctx, "user-1")

	require.NoError(t, err)
}

func TestAnonymizeUser_NotFound(t *testing.T) {
	ctx := context.Background()
	userRepo := mocks.NewMockUserRepository(t)

	userRepo.EXPECT().FindByID(ctx, "nonexistent").Return(nil, nil)

	uc := NewAnonymizeUserUseCase(userRepo)
	err := uc.Execute(ctx, "nonexistent")

	require.Error(t, err)
	var notFoundErr *domainerrors.UserNotFoundError
	assert.True(t, errors.As(err, &notFoundErr))
}

func TestAnonymizeUser_RepoError(t *testing.T) {
	ctx := context.Background()
	userRepo := mocks.NewMockUserRepository(t)

	repoErr := errors.New("database error")
	userRepo.EXPECT().FindByID(ctx, "user-1").Return(nil, repoErr)

	uc := NewAnonymizeUserUseCase(userRepo)
	err := uc.Execute(ctx, "user-1")

	require.Error(t, err)
	assert.Equal(t, repoErr, err)
}

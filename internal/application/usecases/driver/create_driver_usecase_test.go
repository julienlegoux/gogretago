package driver

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

func TestCreateDriver_Success(t *testing.T) {
	ctx := context.Background()
	driverRepo := mocks.NewMockDriverRepository(t)
	userRepo := mocks.NewMockUserRepository(t)
	authRepo := mocks.NewMockAuthRepository(t)

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

	createdDriver := &entities.Driver{
		ID:            "driver-1",
		RefID:         300,
		DriverLicense: "DL-12345",
		UserRefID:     200,
	}

	userRepo.EXPECT().FindByID(ctx, "user-1").Return(user, nil)
	driverRepo.EXPECT().FindByUserRefID(ctx, int64(200)).Return(nil, nil)
	driverRepo.EXPECT().Create(ctx, entities.CreateDriverData{
		DriverLicense: "DL-12345",
		UserRefID:     200,
	}).Return(createdDriver, nil)
	authRepo.EXPECT().UpdateRole(ctx, int64(100), "DRIVER").Return(nil)

	uc := NewCreateDriverUseCase(driverRepo, userRepo, authRepo)
	result, err := uc.Execute(ctx, "user-1", dtos.CreateDriverInput{
		DriverLicense: "DL-12345",
	})

	require.NoError(t, err)
	assert.Equal(t, "driver-1", result.ID)
	assert.Equal(t, "DL-12345", result.DriverLicense)
	assert.Equal(t, int64(200), result.UserRefID)
}

func TestCreateDriver_UserNotFound(t *testing.T) {
	ctx := context.Background()
	driverRepo := mocks.NewMockDriverRepository(t)
	userRepo := mocks.NewMockUserRepository(t)
	authRepo := mocks.NewMockAuthRepository(t)

	userRepo.EXPECT().FindByID(ctx, "nonexistent").Return(nil, nil)

	uc := NewCreateDriverUseCase(driverRepo, userRepo, authRepo)
	result, err := uc.Execute(ctx, "nonexistent", dtos.CreateDriverInput{
		DriverLicense: "DL-12345",
	})

	assert.Nil(t, result)
	require.Error(t, err)
	var notFoundErr *domainerrors.UserNotFoundError
	assert.True(t, errors.As(err, &notFoundErr))
}

func TestCreateDriver_AlreadyExists(t *testing.T) {
	ctx := context.Background()
	driverRepo := mocks.NewMockDriverRepository(t)
	userRepo := mocks.NewMockUserRepository(t)
	authRepo := mocks.NewMockAuthRepository(t)

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

	existingDriver := &entities.Driver{
		ID:            "driver-existing",
		RefID:         300,
		DriverLicense: "DL-OLD",
		UserRefID:     200,
	}

	userRepo.EXPECT().FindByID(ctx, "user-1").Return(user, nil)
	driverRepo.EXPECT().FindByUserRefID(ctx, int64(200)).Return(existingDriver, nil)

	uc := NewCreateDriverUseCase(driverRepo, userRepo, authRepo)
	result, err := uc.Execute(ctx, "user-1", dtos.CreateDriverInput{
		DriverLicense: "DL-12345",
	})

	assert.Nil(t, result)
	require.Error(t, err)
	var alreadyExistsErr *domainerrors.DriverAlreadyExistsError
	assert.True(t, errors.As(err, &alreadyExistsErr))
}

func TestCreateDriver_RepoError(t *testing.T) {
	ctx := context.Background()
	driverRepo := mocks.NewMockDriverRepository(t)
	userRepo := mocks.NewMockUserRepository(t)
	authRepo := mocks.NewMockAuthRepository(t)

	repoErr := errors.New("database error")
	userRepo.EXPECT().FindByID(ctx, "user-1").Return(nil, repoErr)

	uc := NewCreateDriverUseCase(driverRepo, userRepo, authRepo)
	result, err := uc.Execute(ctx, "user-1", dtos.CreateDriverInput{
		DriverLicense: "DL-12345",
	})

	assert.Nil(t, result)
	require.Error(t, err)
	assert.Equal(t, repoErr, err)
}

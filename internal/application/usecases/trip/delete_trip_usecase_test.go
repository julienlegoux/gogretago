package trip

import (
	"context"
	"errors"
	"testing"

	"github.com/lgxju/gogretago/internal/domain/entities"
	domainerrors "github.com/lgxju/gogretago/internal/domain/errors"
	"github.com/lgxju/gogretago/internal/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDeleteTrip_Success(t *testing.T) {
	ctx := context.Background()
	tripRepo := mocks.NewMockTripRepository(t)
	driverRepo := mocks.NewMockDriverRepository(t)

	existingTrip := &entities.Trip{
		ID:          "trip-1",
		RefID:       500,
		Kms:         450,
		Seats:       3,
		DriverRefID: 300,
		CarRefID:    400,
	}
	driver := &entities.Driver{
		ID:            "driver-1",
		RefID:         300,
		DriverLicense: "DL-12345",
		UserRefID:     200,
	}

	tripRepo.EXPECT().FindByID(ctx, "trip-1").Return(existingTrip, nil)
	driverRepo.EXPECT().FindByUserID(ctx, "user-1").Return(driver, nil)
	tripRepo.EXPECT().Delete(ctx, "trip-1").Return(nil)

	uc := NewDeleteTripUseCase(tripRepo, driverRepo)
	err := uc.Execute(ctx, "trip-1", "user-1")

	require.NoError(t, err)
}

func TestDeleteTrip_TripNotFound(t *testing.T) {
	ctx := context.Background()
	tripRepo := mocks.NewMockTripRepository(t)
	driverRepo := mocks.NewMockDriverRepository(t)

	tripRepo.EXPECT().FindByID(ctx, "nonexistent").Return(nil, nil)

	uc := NewDeleteTripUseCase(tripRepo, driverRepo)
	err := uc.Execute(ctx, "nonexistent", "user-1")

	require.Error(t, err)
	var notFoundErr *domainerrors.TripNotFoundError
	assert.True(t, errors.As(err, &notFoundErr))
}

func TestDeleteTrip_DriverNotFound(t *testing.T) {
	ctx := context.Background()
	tripRepo := mocks.NewMockTripRepository(t)
	driverRepo := mocks.NewMockDriverRepository(t)

	existingTrip := &entities.Trip{
		ID:          "trip-1",
		RefID:       500,
		Kms:         450,
		Seats:       3,
		DriverRefID: 300,
		CarRefID:    400,
	}

	tripRepo.EXPECT().FindByID(ctx, "trip-1").Return(existingTrip, nil)
	driverRepo.EXPECT().FindByUserID(ctx, "user-1").Return(nil, nil)

	uc := NewDeleteTripUseCase(tripRepo, driverRepo)
	err := uc.Execute(ctx, "trip-1", "user-1")

	require.Error(t, err)
	var notFoundErr *domainerrors.DriverNotFoundError
	assert.True(t, errors.As(err, &notFoundErr))
}

func TestDeleteTrip_NotOwner(t *testing.T) {
	ctx := context.Background()
	tripRepo := mocks.NewMockTripRepository(t)
	driverRepo := mocks.NewMockDriverRepository(t)

	existingTrip := &entities.Trip{
		ID:          "trip-1",
		RefID:       500,
		Kms:         450,
		Seats:       3,
		DriverRefID: 300, // owned by driver with RefID 300
		CarRefID:    400,
	}
	differentDriver := &entities.Driver{
		ID:            "driver-2",
		RefID:         999, // different RefID than the trip's DriverRefID
		DriverLicense: "DL-OTHER",
		UserRefID:     201,
	}

	tripRepo.EXPECT().FindByID(ctx, "trip-1").Return(existingTrip, nil)
	driverRepo.EXPECT().FindByUserID(ctx, "user-2").Return(differentDriver, nil)

	uc := NewDeleteTripUseCase(tripRepo, driverRepo)
	err := uc.Execute(ctx, "trip-1", "user-2")

	require.Error(t, err)
	var forbiddenErr *domainerrors.ForbiddenError
	assert.True(t, errors.As(err, &forbiddenErr))
}

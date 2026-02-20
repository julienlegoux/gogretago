package car

import (
	"context"
	"errors"
	"testing"

	"github.com/lgxju/gogretago/internal/domain/entities"
	domainerrors "github.com/lgxju/gogretago/internal/domain/errors"
	"github.com/lgxju/gogretago/internal/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestDeleteCar_Success(t *testing.T) {
	ctx := context.Background()
	carID := "car-1"
	userID := "user-1"

	existing := &entities.Car{ID: carID, RefID: 1, LicensePlate: "ABC-123", ModelRefID: 30, DriverRefID: 10}
	driver := &entities.Driver{ID: "driver-1", RefID: 10, UserRefID: 1}

	carRepo := mocks.NewMockCarRepository(t)
	driverRepo := mocks.NewMockDriverRepository(t)

	carRepo.EXPECT().FindByID(mock.Anything, carID).Return(existing, nil)
	driverRepo.EXPECT().FindByUserID(mock.Anything, userID).Return(driver, nil)
	carRepo.EXPECT().Delete(mock.Anything, carID).Return(nil)

	uc := NewDeleteCarUseCase(carRepo, driverRepo)
	err := uc.Execute(ctx, carID, userID)

	assert.NoError(t, err)
}

func TestDeleteCar_NotFound(t *testing.T) {
	ctx := context.Background()
	carID := "car-nonexistent"
	userID := "user-1"

	carRepo := mocks.NewMockCarRepository(t)
	driverRepo := mocks.NewMockDriverRepository(t)

	carRepo.EXPECT().FindByID(mock.Anything, carID).Return(nil, nil)

	uc := NewDeleteCarUseCase(carRepo, driverRepo)
	err := uc.Execute(ctx, carID, userID)

	assert.Error(t, err)
	var notFoundErr *domainerrors.CarNotFoundError
	assert.True(t, errors.As(err, &notFoundErr))
}

func TestDeleteCar_DriverNotFound(t *testing.T) {
	ctx := context.Background()
	carID := "car-1"
	userID := "user-nonexistent"

	existing := &entities.Car{ID: carID, RefID: 1, LicensePlate: "ABC-123", ModelRefID: 30, DriverRefID: 10}

	carRepo := mocks.NewMockCarRepository(t)
	driverRepo := mocks.NewMockDriverRepository(t)

	carRepo.EXPECT().FindByID(mock.Anything, carID).Return(existing, nil)
	driverRepo.EXPECT().FindByUserID(mock.Anything, userID).Return(nil, nil)

	uc := NewDeleteCarUseCase(carRepo, driverRepo)
	err := uc.Execute(ctx, carID, userID)

	assert.Error(t, err)
	var driverNotFoundErr *domainerrors.DriverNotFoundError
	assert.True(t, errors.As(err, &driverNotFoundErr))
}

func TestDeleteCar_NotOwner(t *testing.T) {
	ctx := context.Background()
	carID := "car-1"
	userID := "user-2"

	existing := &entities.Car{ID: carID, RefID: 1, LicensePlate: "ABC-123", ModelRefID: 30, DriverRefID: 10}
	driver := &entities.Driver{ID: "driver-2", RefID: 99, UserRefID: 2}

	carRepo := mocks.NewMockCarRepository(t)
	driverRepo := mocks.NewMockDriverRepository(t)

	carRepo.EXPECT().FindByID(mock.Anything, carID).Return(existing, nil)
	driverRepo.EXPECT().FindByUserID(mock.Anything, userID).Return(driver, nil)

	uc := NewDeleteCarUseCase(carRepo, driverRepo)
	err := uc.Execute(ctx, carID, userID)

	assert.Error(t, err)
	var forbiddenErr *domainerrors.ForbiddenError
	assert.True(t, errors.As(err, &forbiddenErr))
}

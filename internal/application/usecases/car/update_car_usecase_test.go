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

func strPtr(s string) *string {
	return &s
}

func TestUpdateCar_Success(t *testing.T) {
	ctx := context.Background()
	carID := "car-1"
	userID := "user-1"
	newPlate := "NEW-999"
	newModel := "Camry"
	newBrandID := "brand-1"

	existing := &entities.Car{ID: carID, RefID: 1, LicensePlate: "OLD-123", ModelRefID: 30, DriverRefID: 10}
	driver := &entities.Driver{ID: "driver-1", RefID: 10, UserRefID: 1}
	brand := &entities.Brand{ID: "brand-1", RefID: 20, Name: "Toyota"}
	model := &entities.VehicleModel{ID: "model-2", RefID: 31, Name: "Camry", BrandRefID: 20}
	updatedCar := &entities.Car{ID: carID, RefID: 1, LicensePlate: "NEW-999", ModelRefID: 31, DriverRefID: 10}

	carRepo := mocks.NewMockCarRepository(t)
	modelRepo := mocks.NewMockModelRepository(t)
	brandRepo := mocks.NewMockBrandRepository(t)
	driverRepo := mocks.NewMockDriverRepository(t)

	carRepo.EXPECT().FindByID(mock.Anything, carID).Return(existing, nil)
	driverRepo.EXPECT().FindByUserID(mock.Anything, userID).Return(driver, nil)
	carRepo.EXPECT().ExistsByLicensePlate(mock.Anything, newPlate).Return(false, nil)
	brandRepo.EXPECT().FindByID(mock.Anything, newBrandID).Return(brand, nil)
	modelRepo.EXPECT().FindByNameAndBrand(mock.Anything, "Camry", int64(20)).Return(model, nil)
	carRepo.EXPECT().Update(mock.Anything, carID, entities.UpdateCarData{
		LicensePlate: &newPlate,
		ModelRefID:   &model.RefID,
	}).Return(updatedCar, nil)

	uc := NewUpdateCarUseCase(carRepo, modelRepo, brandRepo, driverRepo)
	result, err := uc.Execute(ctx, carID, userID, UpdateCarData{
		LicensePlate: &newPlate,
		Model:        &newModel,
		BrandID:      &newBrandID,
	})

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "NEW-999", result.LicensePlate)
	assert.Equal(t, int64(31), result.ModelRefID)
}

func TestUpdateCar_NotFound(t *testing.T) {
	ctx := context.Background()
	carID := "car-nonexistent"
	userID := "user-1"

	carRepo := mocks.NewMockCarRepository(t)
	modelRepo := mocks.NewMockModelRepository(t)
	brandRepo := mocks.NewMockBrandRepository(t)
	driverRepo := mocks.NewMockDriverRepository(t)

	carRepo.EXPECT().FindByID(mock.Anything, carID).Return(nil, nil)

	uc := NewUpdateCarUseCase(carRepo, modelRepo, brandRepo, driverRepo)
	result, err := uc.Execute(ctx, carID, userID, UpdateCarData{})

	assert.Nil(t, result)
	assert.Error(t, err)
	var notFoundErr *domainerrors.CarNotFoundError
	assert.True(t, errors.As(err, &notFoundErr))
}

func TestUpdateCar_DriverNotFound(t *testing.T) {
	ctx := context.Background()
	carID := "car-1"
	userID := "user-nonexistent"

	existing := &entities.Car{ID: carID, RefID: 1, LicensePlate: "ABC-123", ModelRefID: 30, DriverRefID: 10}

	carRepo := mocks.NewMockCarRepository(t)
	modelRepo := mocks.NewMockModelRepository(t)
	brandRepo := mocks.NewMockBrandRepository(t)
	driverRepo := mocks.NewMockDriverRepository(t)

	carRepo.EXPECT().FindByID(mock.Anything, carID).Return(existing, nil)
	driverRepo.EXPECT().FindByUserID(mock.Anything, userID).Return(nil, nil)

	uc := NewUpdateCarUseCase(carRepo, modelRepo, brandRepo, driverRepo)
	result, err := uc.Execute(ctx, carID, userID, UpdateCarData{})

	assert.Nil(t, result)
	assert.Error(t, err)
	var driverNotFoundErr *domainerrors.DriverNotFoundError
	assert.True(t, errors.As(err, &driverNotFoundErr))
}

func TestUpdateCar_NotOwner(t *testing.T) {
	ctx := context.Background()
	carID := "car-1"
	userID := "user-2"

	// Car belongs to driver with RefID=10, but user's driver has RefID=99
	existing := &entities.Car{ID: carID, RefID: 1, LicensePlate: "ABC-123", ModelRefID: 30, DriverRefID: 10}
	driver := &entities.Driver{ID: "driver-2", RefID: 99, UserRefID: 2}

	carRepo := mocks.NewMockCarRepository(t)
	modelRepo := mocks.NewMockModelRepository(t)
	brandRepo := mocks.NewMockBrandRepository(t)
	driverRepo := mocks.NewMockDriverRepository(t)

	carRepo.EXPECT().FindByID(mock.Anything, carID).Return(existing, nil)
	driverRepo.EXPECT().FindByUserID(mock.Anything, userID).Return(driver, nil)

	uc := NewUpdateCarUseCase(carRepo, modelRepo, brandRepo, driverRepo)
	result, err := uc.Execute(ctx, carID, userID, UpdateCarData{})

	assert.Nil(t, result)
	assert.Error(t, err)
	var forbiddenErr *domainerrors.ForbiddenError
	assert.True(t, errors.As(err, &forbiddenErr))
}

func TestUpdateCar_DuplicateLicensePlate(t *testing.T) {
	ctx := context.Background()
	carID := "car-1"
	userID := "user-1"
	newPlate := "TAKEN-999"

	existing := &entities.Car{ID: carID, RefID: 1, LicensePlate: "ABC-123", ModelRefID: 30, DriverRefID: 10}
	driver := &entities.Driver{ID: "driver-1", RefID: 10, UserRefID: 1}

	carRepo := mocks.NewMockCarRepository(t)
	modelRepo := mocks.NewMockModelRepository(t)
	brandRepo := mocks.NewMockBrandRepository(t)
	driverRepo := mocks.NewMockDriverRepository(t)

	carRepo.EXPECT().FindByID(mock.Anything, carID).Return(existing, nil)
	driverRepo.EXPECT().FindByUserID(mock.Anything, userID).Return(driver, nil)
	carRepo.EXPECT().ExistsByLicensePlate(mock.Anything, newPlate).Return(true, nil)

	uc := NewUpdateCarUseCase(carRepo, modelRepo, brandRepo, driverRepo)
	result, err := uc.Execute(ctx, carID, userID, UpdateCarData{
		LicensePlate: &newPlate,
	})

	assert.Nil(t, result)
	assert.Error(t, err)
	var carExistsErr *domainerrors.CarAlreadyExistsError
	assert.True(t, errors.As(err, &carExistsErr))
}

func TestUpdateCar_SameLicensePlateNoCheck(t *testing.T) {
	ctx := context.Background()
	carID := "car-1"
	userID := "user-1"
	samePlate := "ABC-123"

	existing := &entities.Car{ID: carID, RefID: 1, LicensePlate: "ABC-123", ModelRefID: 30, DriverRefID: 10}
	driver := &entities.Driver{ID: "driver-1", RefID: 10, UserRefID: 1}
	updatedCar := &entities.Car{ID: carID, RefID: 1, LicensePlate: "ABC-123", ModelRefID: 30, DriverRefID: 10}

	carRepo := mocks.NewMockCarRepository(t)
	modelRepo := mocks.NewMockModelRepository(t)
	brandRepo := mocks.NewMockBrandRepository(t)
	driverRepo := mocks.NewMockDriverRepository(t)

	carRepo.EXPECT().FindByID(mock.Anything, carID).Return(existing, nil)
	driverRepo.EXPECT().FindByUserID(mock.Anything, userID).Return(driver, nil)
	// No ExistsByLicensePlate call since the plate hasn't changed
	carRepo.EXPECT().Update(mock.Anything, carID, entities.UpdateCarData{
		LicensePlate: &samePlate,
	}).Return(updatedCar, nil)

	uc := NewUpdateCarUseCase(carRepo, modelRepo, brandRepo, driverRepo)
	result, err := uc.Execute(ctx, carID, userID, UpdateCarData{
		LicensePlate: &samePlate,
	})

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "ABC-123", result.LicensePlate)
}

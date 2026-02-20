package car

import (
	"context"
	"errors"
	"testing"

	"github.com/lgxju/gogretago/internal/application/dtos"
	"github.com/lgxju/gogretago/internal/domain/entities"
	domainerrors "github.com/lgxju/gogretago/internal/domain/errors"
	"github.com/lgxju/gogretago/internal/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCreateCar_Success(t *testing.T) {
	ctx := context.Background()
	userID := "user-1"
	input := dtos.CreateCarInput{
		Model:        "Corolla",
		BrandID:      "brand-1",
		LicensePlate: "ABC-123",
	}

	driver := &entities.Driver{ID: "driver-1", RefID: 10, UserRefID: 1}
	brand := &entities.Brand{ID: "brand-1", RefID: 20, Name: "Toyota"}
	model := &entities.VehicleModel{ID: "model-1", RefID: 30, Name: "Corolla", BrandRefID: 20}
	expectedCar := &entities.Car{ID: "car-1", RefID: 1, LicensePlate: "ABC-123", ModelRefID: 30, DriverRefID: 10}

	carRepo := mocks.NewMockCarRepository(t)
	modelRepo := mocks.NewMockModelRepository(t)
	brandRepo := mocks.NewMockBrandRepository(t)
	driverRepo := mocks.NewMockDriverRepository(t)

	driverRepo.EXPECT().FindByUserID(mock.Anything, userID).Return(driver, nil)
	carRepo.EXPECT().ExistsByLicensePlate(mock.Anything, "ABC-123").Return(false, nil)
	brandRepo.EXPECT().FindByID(mock.Anything, "brand-1").Return(brand, nil)
	modelRepo.EXPECT().FindByNameAndBrand(mock.Anything, "Corolla", int64(20)).Return(model, nil)
	carRepo.EXPECT().Create(mock.Anything, entities.CreateCarData{
		LicensePlate: "ABC-123",
		ModelRefID:   30,
		DriverRefID:  10,
	}).Return(expectedCar, nil)

	uc := NewCreateCarUseCase(carRepo, modelRepo, brandRepo, driverRepo)
	result, err := uc.Execute(ctx, userID, input)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "car-1", result.ID)
	assert.Equal(t, "ABC-123", result.LicensePlate)
	assert.Equal(t, int64(30), result.ModelRefID)
	assert.Equal(t, int64(10), result.DriverRefID)
}

func TestCreateCar_DriverNotFound(t *testing.T) {
	ctx := context.Background()
	userID := "user-nonexistent"
	input := dtos.CreateCarInput{
		Model:        "Corolla",
		BrandID:      "brand-1",
		LicensePlate: "ABC-123",
	}

	carRepo := mocks.NewMockCarRepository(t)
	modelRepo := mocks.NewMockModelRepository(t)
	brandRepo := mocks.NewMockBrandRepository(t)
	driverRepo := mocks.NewMockDriverRepository(t)

	driverRepo.EXPECT().FindByUserID(mock.Anything, userID).Return(nil, nil)

	uc := NewCreateCarUseCase(carRepo, modelRepo, brandRepo, driverRepo)
	result, err := uc.Execute(ctx, userID, input)

	assert.Nil(t, result)
	assert.Error(t, err)
	var driverNotFoundErr *domainerrors.DriverNotFoundError
	assert.True(t, errors.As(err, &driverNotFoundErr))
}

func TestCreateCar_AlreadyExists(t *testing.T) {
	ctx := context.Background()
	userID := "user-1"
	input := dtos.CreateCarInput{
		Model:        "Corolla",
		BrandID:      "brand-1",
		LicensePlate: "ABC-123",
	}

	driver := &entities.Driver{ID: "driver-1", RefID: 10, UserRefID: 1}

	carRepo := mocks.NewMockCarRepository(t)
	modelRepo := mocks.NewMockModelRepository(t)
	brandRepo := mocks.NewMockBrandRepository(t)
	driverRepo := mocks.NewMockDriverRepository(t)

	driverRepo.EXPECT().FindByUserID(mock.Anything, userID).Return(driver, nil)
	carRepo.EXPECT().ExistsByLicensePlate(mock.Anything, "ABC-123").Return(true, nil)

	uc := NewCreateCarUseCase(carRepo, modelRepo, brandRepo, driverRepo)
	result, err := uc.Execute(ctx, userID, input)

	assert.Nil(t, result)
	assert.Error(t, err)
	var carExistsErr *domainerrors.CarAlreadyExistsError
	assert.True(t, errors.As(err, &carExistsErr))
}

func TestCreateCar_BrandNotFound(t *testing.T) {
	ctx := context.Background()
	userID := "user-1"
	input := dtos.CreateCarInput{
		Model:        "Corolla",
		BrandID:      "brand-nonexistent",
		LicensePlate: "ABC-123",
	}

	driver := &entities.Driver{ID: "driver-1", RefID: 10, UserRefID: 1}

	carRepo := mocks.NewMockCarRepository(t)
	modelRepo := mocks.NewMockModelRepository(t)
	brandRepo := mocks.NewMockBrandRepository(t)
	driverRepo := mocks.NewMockDriverRepository(t)

	driverRepo.EXPECT().FindByUserID(mock.Anything, userID).Return(driver, nil)
	carRepo.EXPECT().ExistsByLicensePlate(mock.Anything, "ABC-123").Return(false, nil)
	brandRepo.EXPECT().FindByID(mock.Anything, "brand-nonexistent").Return(nil, nil)

	uc := NewCreateCarUseCase(carRepo, modelRepo, brandRepo, driverRepo)
	result, err := uc.Execute(ctx, userID, input)

	assert.Nil(t, result)
	assert.Error(t, err)
	var brandNotFoundErr *domainerrors.BrandNotFoundError
	assert.True(t, errors.As(err, &brandNotFoundErr))
}

func TestCreateCar_CreatesNewModel(t *testing.T) {
	ctx := context.Background()
	userID := "user-1"
	input := dtos.CreateCarInput{
		Model:        "NewModel",
		BrandID:      "brand-1",
		LicensePlate: "ABC-123",
	}

	driver := &entities.Driver{ID: "driver-1", RefID: 10, UserRefID: 1}
	brand := &entities.Brand{ID: "brand-1", RefID: 20, Name: "Toyota"}
	newModel := &entities.VehicleModel{ID: "model-new", RefID: 31, Name: "NewModel", BrandRefID: 20}
	expectedCar := &entities.Car{ID: "car-1", RefID: 1, LicensePlate: "ABC-123", ModelRefID: 31, DriverRefID: 10}

	carRepo := mocks.NewMockCarRepository(t)
	modelRepo := mocks.NewMockModelRepository(t)
	brandRepo := mocks.NewMockBrandRepository(t)
	driverRepo := mocks.NewMockDriverRepository(t)

	driverRepo.EXPECT().FindByUserID(mock.Anything, userID).Return(driver, nil)
	carRepo.EXPECT().ExistsByLicensePlate(mock.Anything, "ABC-123").Return(false, nil)
	brandRepo.EXPECT().FindByID(mock.Anything, "brand-1").Return(brand, nil)
	// Model not found -> creates new model
	modelRepo.EXPECT().FindByNameAndBrand(mock.Anything, "NewModel", int64(20)).Return(nil, nil)
	modelRepo.EXPECT().Create(mock.Anything, entities.CreateModelData{
		Name:       "NewModel",
		BrandRefID: 20,
	}).Return(newModel, nil)
	carRepo.EXPECT().Create(mock.Anything, entities.CreateCarData{
		LicensePlate: "ABC-123",
		ModelRefID:   31,
		DriverRefID:  10,
	}).Return(expectedCar, nil)

	uc := NewCreateCarUseCase(carRepo, modelRepo, brandRepo, driverRepo)
	result, err := uc.Execute(ctx, userID, input)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "car-1", result.ID)
	assert.Equal(t, int64(31), result.ModelRefID)
}

func TestCreateCar_UsesExistingModel(t *testing.T) {
	ctx := context.Background()
	userID := "user-1"
	input := dtos.CreateCarInput{
		Model:        "Corolla",
		BrandID:      "brand-1",
		LicensePlate: "XYZ-789",
	}

	driver := &entities.Driver{ID: "driver-1", RefID: 10, UserRefID: 1}
	brand := &entities.Brand{ID: "brand-1", RefID: 20, Name: "Toyota"}
	existingModel := &entities.VehicleModel{ID: "model-1", RefID: 30, Name: "Corolla", BrandRefID: 20}
	expectedCar := &entities.Car{ID: "car-2", RefID: 2, LicensePlate: "XYZ-789", ModelRefID: 30, DriverRefID: 10}

	carRepo := mocks.NewMockCarRepository(t)
	modelRepo := mocks.NewMockModelRepository(t)
	brandRepo := mocks.NewMockBrandRepository(t)
	driverRepo := mocks.NewMockDriverRepository(t)

	driverRepo.EXPECT().FindByUserID(mock.Anything, userID).Return(driver, nil)
	carRepo.EXPECT().ExistsByLicensePlate(mock.Anything, "XYZ-789").Return(false, nil)
	brandRepo.EXPECT().FindByID(mock.Anything, "brand-1").Return(brand, nil)
	// Model found -> uses existing model, no Create call on modelRepo
	modelRepo.EXPECT().FindByNameAndBrand(mock.Anything, "Corolla", int64(20)).Return(existingModel, nil)
	carRepo.EXPECT().Create(mock.Anything, entities.CreateCarData{
		LicensePlate: "XYZ-789",
		ModelRefID:   30,
		DriverRefID:  10,
	}).Return(expectedCar, nil)

	uc := NewCreateCarUseCase(carRepo, modelRepo, brandRepo, driverRepo)
	result, err := uc.Execute(ctx, userID, input)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "car-2", result.ID)
	assert.Equal(t, int64(30), result.ModelRefID)
}

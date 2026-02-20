package car

import (
	"context"

	"github.com/lgxju/gogretago/internal/application/dtos"
	"github.com/lgxju/gogretago/internal/domain/entities"
	domainerrors "github.com/lgxju/gogretago/internal/domain/errors"
	"github.com/lgxju/gogretago/internal/domain/repositories"
)

type CreateCarUseCase struct {
	carRepository   repositories.CarRepository
	modelRepository repositories.ModelRepository
	brandRepository repositories.BrandRepository
	driverRepository repositories.DriverRepository
}

func NewCreateCarUseCase(
	carRepository repositories.CarRepository,
	modelRepository repositories.ModelRepository,
	brandRepository repositories.BrandRepository,
	driverRepository repositories.DriverRepository,
) *CreateCarUseCase {
	return &CreateCarUseCase{
		carRepository:   carRepository,
		modelRepository: modelRepository,
		brandRepository: brandRepository,
		driverRepository: driverRepository,
	}
}

func (uc *CreateCarUseCase) Execute(ctx context.Context, userID string, input dtos.CreateCarInput) (*entities.Car, error) {
	driver, err := uc.driverRepository.FindByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}
	if driver == nil {
		return nil, domainerrors.NewDriverNotFoundError(userID)
	}

	exists, err := uc.carRepository.ExistsByLicensePlate(ctx, input.LicensePlate)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, domainerrors.NewCarAlreadyExistsError(input.LicensePlate)
	}

	brand, err := uc.brandRepository.FindByID(ctx, input.BrandID)
	if err != nil {
		return nil, err
	}
	if brand == nil {
		return nil, domainerrors.NewBrandNotFoundError(input.BrandID)
	}

	model, err := uc.modelRepository.FindByNameAndBrand(ctx, input.Model, brand.RefID)
	if err != nil {
		return nil, err
	}
	if model == nil {
		model, err = uc.modelRepository.Create(ctx, entities.CreateModelData{
			Name:       input.Model,
			BrandRefID: brand.RefID,
		})
		if err != nil {
			return nil, err
		}
	}

	return uc.carRepository.Create(ctx, entities.CreateCarData{
		LicensePlate: input.LicensePlate,
		ModelRefID:   model.RefID,
		DriverRefID:  driver.RefID,
	})
}

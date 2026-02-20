package car

import (
	"context"

	"github.com/lgxju/gogretago/internal/domain/entities"
	domainerrors "github.com/lgxju/gogretago/internal/domain/errors"
	"github.com/lgxju/gogretago/internal/domain/repositories"
)

// UpdateCarData holds partial update fields for a car use case.
// All fields are optional to support both PUT (all set) and PATCH (some set).
type UpdateCarData struct {
	Model        *string
	BrandID      *string
	LicensePlate *string
}

type UpdateCarUseCase struct {
	carRepository    repositories.CarRepository
	modelRepository  repositories.ModelRepository
	brandRepository  repositories.BrandRepository
	driverRepository repositories.DriverRepository
}

func NewUpdateCarUseCase(
	carRepository repositories.CarRepository,
	modelRepository repositories.ModelRepository,
	brandRepository repositories.BrandRepository,
	driverRepository repositories.DriverRepository,
) *UpdateCarUseCase {
	return &UpdateCarUseCase{
		carRepository:    carRepository,
		modelRepository:  modelRepository,
		brandRepository:  brandRepository,
		driverRepository: driverRepository,
	}
}

func (uc *UpdateCarUseCase) Execute(ctx context.Context, id string, userID string, input UpdateCarData) (*entities.Car, error) {
	existing, err := uc.carRepository.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if existing == nil {
		return nil, domainerrors.NewCarNotFoundError(id)
	}

	driver, err := uc.driverRepository.FindByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}
	if driver == nil {
		return nil, domainerrors.NewDriverNotFoundError(userID)
	}
	if existing.DriverRefID != driver.RefID {
		return nil, domainerrors.NewForbiddenError("car", id)
	}

	updateData := entities.UpdateCarData{}

	if input.LicensePlate != nil {
		if *input.LicensePlate != existing.LicensePlate {
			exists, err := uc.carRepository.ExistsByLicensePlate(ctx, *input.LicensePlate)
			if err != nil {
				return nil, err
			}
			if exists {
				return nil, domainerrors.NewCarAlreadyExistsError(*input.LicensePlate)
			}
		}
		updateData.LicensePlate = input.LicensePlate
	}

	if input.Model != nil && input.BrandID != nil {
		brand, err := uc.brandRepository.FindByID(ctx, *input.BrandID)
		if err != nil {
			return nil, err
		}
		if brand == nil {
			return nil, domainerrors.NewBrandNotFoundError(*input.BrandID)
		}

		model, err := uc.modelRepository.FindByNameAndBrand(ctx, *input.Model, brand.RefID)
		if err != nil {
			return nil, err
		}
		if model == nil {
			model, err = uc.modelRepository.Create(ctx, entities.CreateModelData{
				Name:       *input.Model,
				BrandRefID: brand.RefID,
			})
			if err != nil {
				return nil, err
			}
		}
		updateData.ModelRefID = &model.RefID
	}

	return uc.carRepository.Update(ctx, id, updateData)
}

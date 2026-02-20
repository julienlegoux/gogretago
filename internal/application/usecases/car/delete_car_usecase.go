package car

import (
	"context"

	domainerrors "github.com/lgxju/gogretago/internal/domain/errors"
	"github.com/lgxju/gogretago/internal/domain/repositories"
)

type DeleteCarUseCase struct {
	carRepository    repositories.CarRepository
	driverRepository repositories.DriverRepository
}

func NewDeleteCarUseCase(
	carRepository repositories.CarRepository,
	driverRepository repositories.DriverRepository,
) *DeleteCarUseCase {
	return &DeleteCarUseCase{
		carRepository:    carRepository,
		driverRepository: driverRepository,
	}
}

func (uc *DeleteCarUseCase) Execute(ctx context.Context, id, userID string) error {
	existing, err := uc.carRepository.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if existing == nil {
		return domainerrors.NewCarNotFoundError(id)
	}

	driver, err := uc.driverRepository.FindByUserID(ctx, userID)
	if err != nil {
		return err
	}
	if driver == nil {
		return domainerrors.NewDriverNotFoundError(userID)
	}
	if existing.DriverRefID != driver.RefID {
		return domainerrors.NewForbiddenError("car", id)
	}

	return uc.carRepository.Delete(ctx, id)
}

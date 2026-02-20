package trip

import (
	"context"

	domainerrors "github.com/lgxju/gogretago/internal/domain/errors"
	"github.com/lgxju/gogretago/internal/domain/repositories"
)

type DeleteTripUseCase struct {
	tripRepository   repositories.TripRepository
	driverRepository repositories.DriverRepository
}

func NewDeleteTripUseCase(
	tripRepository repositories.TripRepository,
	driverRepository repositories.DriverRepository,
) *DeleteTripUseCase {
	return &DeleteTripUseCase{
		tripRepository:   tripRepository,
		driverRepository: driverRepository,
	}
}

func (uc *DeleteTripUseCase) Execute(ctx context.Context, id string, userID string) error {
	existing, err := uc.tripRepository.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if existing == nil {
		return domainerrors.NewTripNotFoundError(id)
	}

	driver, err := uc.driverRepository.FindByUserID(ctx, userID)
	if err != nil {
		return err
	}
	if driver == nil {
		return domainerrors.NewDriverNotFoundError(userID)
	}
	if existing.DriverRefID != driver.RefID {
		return domainerrors.NewForbiddenError("trip", id)
	}

	return uc.tripRepository.Delete(ctx, id)
}

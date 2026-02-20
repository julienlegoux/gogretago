package trip

import (
	"context"

	"github.com/lgxju/gogretago/internal/domain/entities"
	domainerrors "github.com/lgxju/gogretago/internal/domain/errors"
	"github.com/lgxju/gogretago/internal/domain/repositories"
)

type GetTripUseCase struct {
	tripRepository repositories.TripRepository
}

func NewGetTripUseCase(tripRepository repositories.TripRepository) *GetTripUseCase {
	return &GetTripUseCase{
		tripRepository: tripRepository,
	}
}

func (uc *GetTripUseCase) Execute(ctx context.Context, id string) (*entities.Trip, error) {
	trip, err := uc.tripRepository.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if trip == nil {
		return nil, domainerrors.NewTripNotFoundError(id)
	}
	return trip, nil
}

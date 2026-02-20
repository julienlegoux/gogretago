package trip

import (
	"context"

	"github.com/lgxju/gogretago/internal/domain/entities"
	"github.com/lgxju/gogretago/internal/domain/repositories"
)

type ListTripsUseCase struct {
	tripRepository repositories.TripRepository
}

func NewListTripsUseCase(tripRepository repositories.TripRepository) *ListTripsUseCase {
	return &ListTripsUseCase{
		tripRepository: tripRepository,
	}
}

func (uc *ListTripsUseCase) Execute(ctx context.Context, params entities.PaginationParams) (*entities.PaginatedResult[entities.Trip], error) {
	trips, total, err := uc.tripRepository.FindAll(ctx, params.Skip(), params.Take())
	if err != nil {
		return nil, err
	}

	return &entities.PaginatedResult[entities.Trip]{
		Data: trips,
		Meta: entities.BuildPaginationMeta(params, total),
	}, nil
}

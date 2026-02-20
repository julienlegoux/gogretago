package trip

import (
	"context"
	"time"

	"github.com/lgxju/gogretago/internal/application/dtos"
	"github.com/lgxju/gogretago/internal/domain/entities"
	"github.com/lgxju/gogretago/internal/domain/repositories"
)

type FindTripsUseCase struct {
	tripRepository repositories.TripRepository
}

func NewFindTripsUseCase(tripRepository repositories.TripRepository) *FindTripsUseCase {
	return &FindTripsUseCase{
		tripRepository: tripRepository,
	}
}

func (uc *FindTripsUseCase) Execute(ctx context.Context, query dtos.FindTripQuery) ([]entities.Trip, error) {
	filters := entities.TripFilters{
		DepartureCity: query.DepartureCity,
		ArrivalCity:   query.ArrivalCity,
	}

	if query.Date != nil {
		parsed, err := time.Parse("2006-01-02", *query.Date)
		if err != nil {
			return nil, err
		}
		filters.Date = &parsed
	}

	return uc.tripRepository.FindByFilters(ctx, filters)
}

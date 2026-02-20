package car

import (
	"context"

	"github.com/lgxju/gogretago/internal/domain/entities"
	"github.com/lgxju/gogretago/internal/domain/repositories"
)

type ListCarsUseCase struct {
	carRepository repositories.CarRepository
}

func NewListCarsUseCase(carRepository repositories.CarRepository) *ListCarsUseCase {
	return &ListCarsUseCase{
		carRepository: carRepository,
	}
}

func (uc *ListCarsUseCase) Execute(ctx context.Context, params entities.PaginationParams) (*entities.PaginatedResult[entities.Car], error) {
	cars, total, err := uc.carRepository.FindAll(ctx, params.Skip(), params.Take())
	if err != nil {
		return nil, err
	}

	return &entities.PaginatedResult[entities.Car]{
		Data: cars,
		Meta: entities.BuildPaginationMeta(params, total),
	}, nil
}

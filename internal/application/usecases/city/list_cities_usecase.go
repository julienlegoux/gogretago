package city

import (
	"context"

	"github.com/lgxju/gogretago/internal/domain/entities"
	"github.com/lgxju/gogretago/internal/domain/repositories"
)

type ListCitiesUseCase struct {
	cityRepository repositories.CityRepository
}

func NewListCitiesUseCase(cityRepository repositories.CityRepository) *ListCitiesUseCase {
	return &ListCitiesUseCase{
		cityRepository: cityRepository,
	}
}

func (uc *ListCitiesUseCase) Execute(ctx context.Context, params entities.PaginationParams) (*entities.PaginatedResult[entities.City], error) {
	cities, total, err := uc.cityRepository.FindAll(ctx, params.Skip(), params.Take())
	if err != nil {
		return nil, err
	}

	return &entities.PaginatedResult[entities.City]{
		Data: cities,
		Meta: entities.BuildPaginationMeta(params, total),
	}, nil
}

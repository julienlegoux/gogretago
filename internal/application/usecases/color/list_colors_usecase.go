package color

import (
	"context"

	"github.com/lgxju/gogretago/internal/domain/entities"
	"github.com/lgxju/gogretago/internal/domain/repositories"
)

type ListColorsUseCase struct {
	colorRepository repositories.ColorRepository
}

func NewListColorsUseCase(colorRepository repositories.ColorRepository) *ListColorsUseCase {
	return &ListColorsUseCase{
		colorRepository: colorRepository,
	}
}

func (uc *ListColorsUseCase) Execute(ctx context.Context, params entities.PaginationParams) (*entities.PaginatedResult[entities.Color], error) {
	colors, total, err := uc.colorRepository.FindAll(ctx, params.Skip(), params.Take())
	if err != nil {
		return nil, err
	}

	return &entities.PaginatedResult[entities.Color]{
		Data: colors,
		Meta: entities.BuildPaginationMeta(params, total),
	}, nil
}

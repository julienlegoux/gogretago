package brand

import (
	"context"

	"github.com/lgxju/gogretago/internal/domain/entities"
	"github.com/lgxju/gogretago/internal/domain/repositories"
)

type ListBrandsUseCase struct {
	brandRepository repositories.BrandRepository
}

func NewListBrandsUseCase(brandRepository repositories.BrandRepository) *ListBrandsUseCase {
	return &ListBrandsUseCase{
		brandRepository: brandRepository,
	}
}

func (uc *ListBrandsUseCase) Execute(ctx context.Context, params entities.PaginationParams) (*entities.PaginatedResult[entities.Brand], error) {
	brands, total, err := uc.brandRepository.FindAll(ctx, params.Skip(), params.Take())
	if err != nil {
		return nil, err
	}

	return &entities.PaginatedResult[entities.Brand]{
		Data: brands,
		Meta: entities.BuildPaginationMeta(params, total),
	}, nil
}

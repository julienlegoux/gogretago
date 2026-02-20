package brand

import (
	"context"

	"github.com/lgxju/gogretago/internal/application/dtos"
	"github.com/lgxju/gogretago/internal/domain/entities"
	"github.com/lgxju/gogretago/internal/domain/repositories"
)

type CreateBrandUseCase struct {
	brandRepository repositories.BrandRepository
}

func NewCreateBrandUseCase(brandRepository repositories.BrandRepository) *CreateBrandUseCase {
	return &CreateBrandUseCase{
		brandRepository: brandRepository,
	}
}

func (uc *CreateBrandUseCase) Execute(ctx context.Context, input dtos.CreateBrandInput) (*entities.Brand, error) {
	return uc.brandRepository.Create(ctx, entities.CreateBrandData{
		Name: input.Name,
	})
}

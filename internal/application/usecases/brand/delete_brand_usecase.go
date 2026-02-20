package brand

import (
	"context"

	domainerrors "github.com/lgxju/gogretago/internal/domain/errors"
	"github.com/lgxju/gogretago/internal/domain/repositories"
)

type DeleteBrandUseCase struct {
	brandRepository repositories.BrandRepository
}

func NewDeleteBrandUseCase(brandRepository repositories.BrandRepository) *DeleteBrandUseCase {
	return &DeleteBrandUseCase{
		brandRepository: brandRepository,
	}
}

func (uc *DeleteBrandUseCase) Execute(ctx context.Context, id string) error {
	existing, err := uc.brandRepository.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if existing == nil {
		return domainerrors.NewBrandNotFoundError(id)
	}

	return uc.brandRepository.Delete(ctx, id)
}

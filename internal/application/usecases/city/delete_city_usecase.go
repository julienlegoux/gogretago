package city

import (
	"context"

	domainerrors "github.com/lgxju/gogretago/internal/domain/errors"
	"github.com/lgxju/gogretago/internal/domain/repositories"
)

type DeleteCityUseCase struct {
	cityRepository repositories.CityRepository
}

func NewDeleteCityUseCase(cityRepository repositories.CityRepository) *DeleteCityUseCase {
	return &DeleteCityUseCase{
		cityRepository: cityRepository,
	}
}

func (uc *DeleteCityUseCase) Execute(ctx context.Context, id string) error {
	existing, err := uc.cityRepository.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if existing == nil {
		return domainerrors.NewCityNotFoundError(id)
	}

	return uc.cityRepository.Delete(ctx, id)
}

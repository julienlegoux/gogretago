package city

import (
	"context"

	"github.com/lgxju/gogretago/internal/application/dtos"
	"github.com/lgxju/gogretago/internal/domain/entities"
	"github.com/lgxju/gogretago/internal/domain/repositories"
)

type CreateCityUseCase struct {
	cityRepository repositories.CityRepository
}

func NewCreateCityUseCase(cityRepository repositories.CityRepository) *CreateCityUseCase {
	return &CreateCityUseCase{
		cityRepository: cityRepository,
	}
}

func (uc *CreateCityUseCase) Execute(ctx context.Context, input dtos.CreateCityInput) (*entities.City, error) {
	return uc.cityRepository.Create(ctx, entities.CreateCityData{
		CityName: input.CityName,
		Zipcode:  input.Zipcode,
	})
}

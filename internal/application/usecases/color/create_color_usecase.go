package color

import (
	"context"

	"github.com/lgxju/gogretago/internal/application/dtos"
	"github.com/lgxju/gogretago/internal/domain/entities"
	domainerrors "github.com/lgxju/gogretago/internal/domain/errors"
	"github.com/lgxju/gogretago/internal/domain/repositories"
)

type CreateColorUseCase struct {
	colorRepository repositories.ColorRepository
}

func NewCreateColorUseCase(colorRepository repositories.ColorRepository) *CreateColorUseCase {
	return &CreateColorUseCase{
		colorRepository: colorRepository,
	}
}

func (uc *CreateColorUseCase) Execute(ctx context.Context, input dtos.CreateColorInput) (*entities.Color, error) {
	existing, err := uc.colorRepository.FindByName(ctx, input.Name)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return nil, domainerrors.NewColorAlreadyExistsError(input.Name)
	}

	return uc.colorRepository.Create(ctx, entities.CreateColorData{
		Name: input.Name,
		Hex:  input.Hex,
	})
}

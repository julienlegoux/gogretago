package color

import (
	"context"

	"github.com/lgxju/gogretago/internal/application/dtos"
	"github.com/lgxju/gogretago/internal/domain/entities"
	domainerrors "github.com/lgxju/gogretago/internal/domain/errors"
	"github.com/lgxju/gogretago/internal/domain/repositories"
)

type UpdateColorUseCase struct {
	colorRepository repositories.ColorRepository
}

func NewUpdateColorUseCase(colorRepository repositories.ColorRepository) *UpdateColorUseCase {
	return &UpdateColorUseCase{
		colorRepository: colorRepository,
	}
}

func (uc *UpdateColorUseCase) Execute(ctx context.Context, id string, input dtos.UpdateColorInput) (*entities.Color, error) {
	existing, err := uc.colorRepository.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if existing == nil {
		return nil, domainerrors.NewColorNotFoundError(id)
	}

	if input.Name != nil {
		duplicate, err := uc.colorRepository.FindByName(ctx, *input.Name)
		if err != nil {
			return nil, err
		}
		if duplicate != nil && duplicate.ID != id {
			return nil, domainerrors.NewColorAlreadyExistsError(*input.Name)
		}
	}

	return uc.colorRepository.Update(ctx, id, entities.UpdateColorData{
		Name: input.Name,
		Hex:  input.Hex,
	})
}

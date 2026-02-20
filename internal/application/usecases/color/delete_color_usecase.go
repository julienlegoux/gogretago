package color

import (
	"context"

	domainerrors "github.com/lgxju/gogretago/internal/domain/errors"
	"github.com/lgxju/gogretago/internal/domain/repositories"
)

type DeleteColorUseCase struct {
	colorRepository repositories.ColorRepository
}

func NewDeleteColorUseCase(colorRepository repositories.ColorRepository) *DeleteColorUseCase {
	return &DeleteColorUseCase{
		colorRepository: colorRepository,
	}
}

func (uc *DeleteColorUseCase) Execute(ctx context.Context, id string) error {
	existing, err := uc.colorRepository.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if existing == nil {
		return domainerrors.NewColorNotFoundError(id)
	}

	return uc.colorRepository.Delete(ctx, id)
}

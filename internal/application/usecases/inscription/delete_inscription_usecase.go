package inscription

import (
	"context"

	domainerrors "github.com/lgxju/gogretago/internal/domain/errors"
	"github.com/lgxju/gogretago/internal/domain/repositories"
)

type DeleteInscriptionUseCase struct {
	inscriptionRepository repositories.InscriptionRepository
}

func NewDeleteInscriptionUseCase(inscriptionRepository repositories.InscriptionRepository) *DeleteInscriptionUseCase {
	return &DeleteInscriptionUseCase{
		inscriptionRepository: inscriptionRepository,
	}
}

func (uc *DeleteInscriptionUseCase) Execute(ctx context.Context, id string, userID string) error {
	existing, err := uc.inscriptionRepository.FindByIDAndUserID(ctx, id, userID)
	if err != nil {
		return err
	}
	if existing == nil {
		return domainerrors.NewInscriptionNotFoundError(id)
	}

	return uc.inscriptionRepository.Delete(ctx, id)
}

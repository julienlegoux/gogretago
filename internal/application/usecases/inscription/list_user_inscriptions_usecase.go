package inscription

import (
	"context"

	"github.com/lgxju/gogretago/internal/domain/entities"
	"github.com/lgxju/gogretago/internal/domain/repositories"
)

type ListUserInscriptionsUseCase struct {
	inscriptionRepository repositories.InscriptionRepository
}

func NewListUserInscriptionsUseCase(inscriptionRepository repositories.InscriptionRepository) *ListUserInscriptionsUseCase {
	return &ListUserInscriptionsUseCase{
		inscriptionRepository: inscriptionRepository,
	}
}

func (uc *ListUserInscriptionsUseCase) Execute(ctx context.Context, userID string) ([]entities.Inscription, error) {
	return uc.inscriptionRepository.FindByUserID(ctx, userID)
}

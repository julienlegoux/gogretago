package inscription

import (
	"context"

	"github.com/lgxju/gogretago/internal/domain/entities"
	"github.com/lgxju/gogretago/internal/domain/repositories"
)

type ListInscriptionsUseCase struct {
	inscriptionRepository repositories.InscriptionRepository
}

func NewListInscriptionsUseCase(inscriptionRepository repositories.InscriptionRepository) *ListInscriptionsUseCase {
	return &ListInscriptionsUseCase{
		inscriptionRepository: inscriptionRepository,
	}
}

func (uc *ListInscriptionsUseCase) Execute(ctx context.Context, params entities.PaginationParams) (*entities.PaginatedResult[entities.Inscription], error) {
	inscriptions, total, err := uc.inscriptionRepository.FindAll(ctx, params.Skip(), params.Take())
	if err != nil {
		return nil, err
	}

	return &entities.PaginatedResult[entities.Inscription]{
		Data: inscriptions,
		Meta: entities.BuildPaginationMeta(params, total),
	}, nil
}

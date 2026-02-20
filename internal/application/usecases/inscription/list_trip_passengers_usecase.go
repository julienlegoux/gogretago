package inscription

import (
	"context"

	"github.com/lgxju/gogretago/internal/domain/entities"
	"github.com/lgxju/gogretago/internal/domain/repositories"
)

type ListTripPassengersUseCase struct {
	inscriptionRepository repositories.InscriptionRepository
}

func NewListTripPassengersUseCase(inscriptionRepository repositories.InscriptionRepository) *ListTripPassengersUseCase {
	return &ListTripPassengersUseCase{
		inscriptionRepository: inscriptionRepository,
	}
}

func (uc *ListTripPassengersUseCase) Execute(ctx context.Context, tripID string) ([]entities.Inscription, error) {
	return uc.inscriptionRepository.FindByTripID(ctx, tripID)
}

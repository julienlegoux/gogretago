package inscription

import (
	"context"

	"github.com/lgxju/gogretago/internal/application/dtos"
	"github.com/lgxju/gogretago/internal/domain/entities"
	domainerrors "github.com/lgxju/gogretago/internal/domain/errors"
	"github.com/lgxju/gogretago/internal/domain/repositories"
)

type CreateInscriptionUseCase struct {
	inscriptionRepository repositories.InscriptionRepository
	userRepository        repositories.UserRepository
	tripRepository        repositories.TripRepository
}

func NewCreateInscriptionUseCase(
	inscriptionRepository repositories.InscriptionRepository,
	userRepository repositories.UserRepository,
	tripRepository repositories.TripRepository,
) *CreateInscriptionUseCase {
	return &CreateInscriptionUseCase{
		inscriptionRepository: inscriptionRepository,
		userRepository:        userRepository,
		tripRepository:        tripRepository,
	}
}

func (uc *CreateInscriptionUseCase) Execute(ctx context.Context, userID string, input dtos.CreateInscriptionInput) (*entities.Inscription, error) {
	user, err := uc.userRepository.FindByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, domainerrors.NewUserNotFoundError(userID)
	}

	trip, err := uc.tripRepository.FindByID(ctx, input.TripID)
	if err != nil {
		return nil, err
	}
	if trip == nil {
		return nil, domainerrors.NewTripNotFoundError(input.TripID)
	}

	alreadyInscribed, err := uc.inscriptionRepository.ExistsByUserAndTrip(ctx, user.RefID, trip.RefID)
	if err != nil {
		return nil, err
	}
	if alreadyInscribed {
		return nil, domainerrors.NewAlreadyInscribedError(userID, input.TripID)
	}

	count, err := uc.inscriptionRepository.CountByTripRefID(ctx, trip.RefID)
	if err != nil {
		return nil, err
	}
	if count >= trip.Seats {
		return nil, domainerrors.NewNoSeatsAvailableError(input.TripID)
	}

	return uc.inscriptionRepository.Create(ctx, entities.CreateInscriptionData{
		UserRefID: user.RefID,
		TripRefID: trip.RefID,
	})
}

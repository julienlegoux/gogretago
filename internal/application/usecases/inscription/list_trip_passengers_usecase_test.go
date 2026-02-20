package inscription

import (
	"context"
	"fmt"
	"testing"

	"github.com/lgxju/gogretago/internal/domain/entities"
	"github.com/lgxju/gogretago/internal/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestListTripPassengers_Success(t *testing.T) {
	ctx := context.Background()
	tripID := "trip-1"

	inscriptions := []entities.Inscription{
		{ID: "insc-1", RefID: 1, UserRefID: 10, TripRefID: 20},
		{ID: "insc-2", RefID: 2, UserRefID: 11, TripRefID: 20},
	}

	inscriptionRepo := mocks.NewMockInscriptionRepository(t)

	inscriptionRepo.EXPECT().FindByTripID(mock.Anything, tripID).Return(inscriptions, nil)

	uc := NewListTripPassengersUseCase(inscriptionRepo)
	result, err := uc.Execute(ctx, tripID)

	assert.NoError(t, err)
	assert.Len(t, result, 2)
	assert.Equal(t, "insc-1", result[0].ID)
	assert.Equal(t, "insc-2", result[1].ID)
}

func TestListTripPassengers_Empty(t *testing.T) {
	ctx := context.Background()
	tripID := "trip-1"

	inscriptionRepo := mocks.NewMockInscriptionRepository(t)

	inscriptionRepo.EXPECT().FindByTripID(mock.Anything, tripID).Return([]entities.Inscription{}, nil)

	uc := NewListTripPassengersUseCase(inscriptionRepo)
	result, err := uc.Execute(ctx, tripID)

	assert.NoError(t, err)
	assert.Empty(t, result)
}

func TestListTripPassengers_RepoError(t *testing.T) {
	ctx := context.Background()
	tripID := "trip-1"

	inscriptionRepo := mocks.NewMockInscriptionRepository(t)
	inscriptionRepo.EXPECT().FindByTripID(mock.Anything, tripID).Return(nil, fmt.Errorf("database connection failed"))

	uc := NewListTripPassengersUseCase(inscriptionRepo)
	result, err := uc.Execute(ctx, tripID)

	assert.Error(t, err)
	assert.Nil(t, result)
}

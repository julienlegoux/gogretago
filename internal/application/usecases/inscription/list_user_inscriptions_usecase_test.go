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

func TestListUserInscriptions_Success(t *testing.T) {
	ctx := context.Background()
	userID := "user-1"

	inscriptions := []entities.Inscription{
		{ID: "insc-1", RefID: 1, UserRefID: 10, TripRefID: 20},
		{ID: "insc-2", RefID: 2, UserRefID: 10, TripRefID: 21},
	}

	inscriptionRepo := mocks.NewMockInscriptionRepository(t)

	inscriptionRepo.EXPECT().FindByUserID(mock.Anything, userID).Return(inscriptions, nil)

	uc := NewListUserInscriptionsUseCase(inscriptionRepo)
	result, err := uc.Execute(ctx, userID)

	assert.NoError(t, err)
	assert.Len(t, result, 2)
	assert.Equal(t, "insc-1", result[0].ID)
	assert.Equal(t, "insc-2", result[1].ID)
}

func TestListUserInscriptions_Empty(t *testing.T) {
	ctx := context.Background()
	userID := "user-1"

	inscriptionRepo := mocks.NewMockInscriptionRepository(t)

	inscriptionRepo.EXPECT().FindByUserID(mock.Anything, userID).Return([]entities.Inscription{}, nil)

	uc := NewListUserInscriptionsUseCase(inscriptionRepo)
	result, err := uc.Execute(ctx, userID)

	assert.NoError(t, err)
	assert.Empty(t, result)
}

func TestListUserInscriptions_RepoError(t *testing.T) {
	ctx := context.Background()
	userID := "user-1"

	inscriptionRepo := mocks.NewMockInscriptionRepository(t)
	inscriptionRepo.EXPECT().FindByUserID(mock.Anything, userID).Return(nil, fmt.Errorf("database connection failed"))

	uc := NewListUserInscriptionsUseCase(inscriptionRepo)
	result, err := uc.Execute(ctx, userID)

	assert.Error(t, err)
	assert.Nil(t, result)
}

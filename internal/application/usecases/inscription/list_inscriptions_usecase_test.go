package inscription

import (
	"context"
	"testing"

	"github.com/lgxju/gogretago/internal/domain/entities"
	"github.com/lgxju/gogretago/internal/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestListInscriptions_Success(t *testing.T) {
	ctx := context.Background()
	params := entities.PaginationParams{Page: 1, Limit: 20}

	inscriptions := []entities.Inscription{
		{ID: "insc-1", RefID: 1, UserRefID: 10, TripRefID: 20},
		{ID: "insc-2", RefID: 2, UserRefID: 11, TripRefID: 21},
	}

	inscriptionRepo := mocks.NewMockInscriptionRepository(t)

	inscriptionRepo.EXPECT().FindAll(mock.Anything, 0, 20).Return(inscriptions, 2, nil)

	uc := NewListInscriptionsUseCase(inscriptionRepo)
	result, err := uc.Execute(ctx, params)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Len(t, result.Data, 2)
	assert.Equal(t, 2, result.Meta.Total)
	assert.Equal(t, 1, result.Meta.Page)
	assert.Equal(t, 20, result.Meta.Limit)
}

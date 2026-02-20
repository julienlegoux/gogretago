package inscription

import (
	"context"
	"errors"
	"testing"

	"github.com/lgxju/gogretago/internal/domain/entities"
	domainerrors "github.com/lgxju/gogretago/internal/domain/errors"
	"github.com/lgxju/gogretago/internal/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestDeleteInscription_Success(t *testing.T) {
	ctx := context.Background()
	inscriptionID := "insc-1"
	userID := "user-1"

	existing := &entities.Inscription{ID: inscriptionID, RefID: 1, UserRefID: 10, TripRefID: 20}

	inscriptionRepo := mocks.NewMockInscriptionRepository(t)

	inscriptionRepo.EXPECT().FindByIDAndUserID(mock.Anything, inscriptionID, userID).Return(existing, nil)
	inscriptionRepo.EXPECT().Delete(mock.Anything, inscriptionID).Return(nil)

	uc := NewDeleteInscriptionUseCase(inscriptionRepo)
	err := uc.Execute(ctx, inscriptionID, userID)

	assert.NoError(t, err)
}

func TestDeleteInscription_NotFoundOrNotOwner(t *testing.T) {
	ctx := context.Background()
	inscriptionID := "insc-nonexistent"
	userID := "user-1"

	inscriptionRepo := mocks.NewMockInscriptionRepository(t)

	inscriptionRepo.EXPECT().FindByIDAndUserID(mock.Anything, inscriptionID, userID).Return(nil, nil)

	uc := NewDeleteInscriptionUseCase(inscriptionRepo)
	err := uc.Execute(ctx, inscriptionID, userID)

	assert.Error(t, err)
	var notFoundErr *domainerrors.InscriptionNotFoundError
	assert.True(t, errors.As(err, &notFoundErr))
}

func TestDeleteInscription_RepoError(t *testing.T) {
	ctx := context.Background()
	inscriptionID := "insc-1"
	userID := "user-1"

	inscriptionRepo := mocks.NewMockInscriptionRepository(t)

	inscriptionRepo.EXPECT().FindByIDAndUserID(mock.Anything, inscriptionID, userID).Return(nil, errors.New("database error"))

	uc := NewDeleteInscriptionUseCase(inscriptionRepo)
	err := uc.Execute(ctx, inscriptionID, userID)

	assert.Error(t, err)
	assert.Equal(t, "database error", err.Error())
}

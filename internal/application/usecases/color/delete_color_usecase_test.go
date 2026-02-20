package color

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

func TestDeleteColor_Success(t *testing.T) {
	ctx := context.Background()
	colorID := "color-1"

	existing := &entities.Color{ID: colorID, RefID: 1, Name: "Red", Hex: "#FF0000"}

	colorRepo := mocks.NewMockColorRepository(t)

	colorRepo.EXPECT().FindByID(mock.Anything, colorID).Return(existing, nil)
	colorRepo.EXPECT().Delete(mock.Anything, colorID).Return(nil)

	uc := NewDeleteColorUseCase(colorRepo)
	err := uc.Execute(ctx, colorID)

	assert.NoError(t, err)
}

func TestDeleteColor_NotFound(t *testing.T) {
	ctx := context.Background()
	colorID := "color-nonexistent"

	colorRepo := mocks.NewMockColorRepository(t)

	colorRepo.EXPECT().FindByID(mock.Anything, colorID).Return(nil, nil)

	uc := NewDeleteColorUseCase(colorRepo)
	err := uc.Execute(ctx, colorID)

	assert.Error(t, err)
	var notFoundErr *domainerrors.ColorNotFoundError
	assert.True(t, errors.As(err, &notFoundErr))
}

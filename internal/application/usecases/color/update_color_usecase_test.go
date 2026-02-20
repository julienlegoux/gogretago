package color

import (
	"context"
	"errors"
	"testing"

	"github.com/lgxju/gogretago/internal/application/dtos"
	"github.com/lgxju/gogretago/internal/domain/entities"
	domainerrors "github.com/lgxju/gogretago/internal/domain/errors"
	"github.com/lgxju/gogretago/internal/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func strPtr(s string) *string {
	return &s
}

func TestUpdateColor_Success(t *testing.T) {
	ctx := context.Background()
	colorID := "color-1"
	newName := "Crimson"
	newHex := "#DC143C"

	existing := &entities.Color{ID: colorID, RefID: 1, Name: "Red", Hex: "#FF0000"}
	updatedColor := &entities.Color{ID: colorID, RefID: 1, Name: "Crimson", Hex: "#DC143C"}

	colorRepo := mocks.NewMockColorRepository(t)

	colorRepo.EXPECT().FindByID(mock.Anything, colorID).Return(existing, nil)
	colorRepo.EXPECT().FindByName(mock.Anything, "Crimson").Return(nil, nil)
	colorRepo.EXPECT().Update(mock.Anything, colorID, entities.UpdateColorData{
		Name: &newName,
		Hex:  &newHex,
	}).Return(updatedColor, nil)

	uc := NewUpdateColorUseCase(colorRepo)
	result, err := uc.Execute(ctx, colorID, dtos.UpdateColorInput{
		Name: &newName,
		Hex:  &newHex,
	})

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "Crimson", result.Name)
	assert.Equal(t, "#DC143C", result.Hex)
}

func TestUpdateColor_NotFound(t *testing.T) {
	ctx := context.Background()
	colorID := "color-nonexistent"

	colorRepo := mocks.NewMockColorRepository(t)

	colorRepo.EXPECT().FindByID(mock.Anything, colorID).Return(nil, nil)

	uc := NewUpdateColorUseCase(colorRepo)
	result, err := uc.Execute(ctx, colorID, dtos.UpdateColorInput{})

	assert.Nil(t, result)
	assert.Error(t, err)
	var notFoundErr *domainerrors.ColorNotFoundError
	assert.True(t, errors.As(err, &notFoundErr))
}

func TestUpdateColor_DuplicateName(t *testing.T) {
	ctx := context.Background()
	colorID := "color-1"
	duplicateName := "Blue"

	existing := &entities.Color{ID: colorID, RefID: 1, Name: "Red", Hex: "#FF0000"}
	// Another color with the same name but a different ID
	duplicate := &entities.Color{ID: "color-2", RefID: 2, Name: "Blue", Hex: "#0000FF"}

	colorRepo := mocks.NewMockColorRepository(t)

	colorRepo.EXPECT().FindByID(mock.Anything, colorID).Return(existing, nil)
	colorRepo.EXPECT().FindByName(mock.Anything, "Blue").Return(duplicate, nil)

	uc := NewUpdateColorUseCase(colorRepo)
	result, err := uc.Execute(ctx, colorID, dtos.UpdateColorInput{
		Name: &duplicateName,
	})

	assert.Nil(t, result)
	assert.Error(t, err)
	var colorExistsErr *domainerrors.ColorAlreadyExistsError
	assert.True(t, errors.As(err, &colorExistsErr))
}

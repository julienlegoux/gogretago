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

func TestCreateColor_Success(t *testing.T) {
	ctx := context.Background()
	input := dtos.CreateColorInput{
		Name: "Red",
		Hex:  "#FF0000",
	}

	expectedColor := &entities.Color{ID: "color-1", RefID: 1, Name: "Red", Hex: "#FF0000"}

	colorRepo := mocks.NewMockColorRepository(t)

	colorRepo.EXPECT().FindByName(mock.Anything, "Red").Return(nil, nil)
	colorRepo.EXPECT().Create(mock.Anything, entities.CreateColorData{
		Name: "Red",
		Hex:  "#FF0000",
	}).Return(expectedColor, nil)

	uc := NewCreateColorUseCase(colorRepo)
	result, err := uc.Execute(ctx, input)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "color-1", result.ID)
	assert.Equal(t, "Red", result.Name)
	assert.Equal(t, "#FF0000", result.Hex)
}

func TestCreateColor_AlreadyExists(t *testing.T) {
	ctx := context.Background()
	input := dtos.CreateColorInput{
		Name: "Red",
		Hex:  "#FF0000",
	}

	existingColor := &entities.Color{ID: "color-existing", RefID: 1, Name: "Red", Hex: "#FF0000"}

	colorRepo := mocks.NewMockColorRepository(t)

	colorRepo.EXPECT().FindByName(mock.Anything, "Red").Return(existingColor, nil)

	uc := NewCreateColorUseCase(colorRepo)
	result, err := uc.Execute(ctx, input)

	assert.Nil(t, result)
	assert.Error(t, err)
	var colorExistsErr *domainerrors.ColorAlreadyExistsError
	assert.True(t, errors.As(err, &colorExistsErr))
}

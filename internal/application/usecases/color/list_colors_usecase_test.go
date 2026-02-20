package color

import (
	"context"
	"fmt"
	"testing"

	"github.com/lgxju/gogretago/internal/domain/entities"
	"github.com/lgxju/gogretago/internal/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestListColors_Success(t *testing.T) {
	ctx := context.Background()
	params := entities.PaginationParams{Page: 1, Limit: 20}

	colors := []entities.Color{
		{ID: "color-1", RefID: 1, Name: "Red", Hex: "#FF0000"},
		{ID: "color-2", RefID: 2, Name: "Blue", Hex: "#0000FF"},
	}

	colorRepo := mocks.NewMockColorRepository(t)

	colorRepo.EXPECT().FindAll(mock.Anything, 0, 20).Return(colors, 2, nil)

	uc := NewListColorsUseCase(colorRepo)
	result, err := uc.Execute(ctx, params)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Len(t, result.Data, 2)
	assert.Equal(t, 2, result.Meta.Total)
	assert.Equal(t, 1, result.Meta.Page)
	assert.Equal(t, 20, result.Meta.Limit)
}

func TestListColors_Empty(t *testing.T) {
	ctx := context.Background()
	params := entities.PaginationParams{Page: 1, Limit: 20}

	colorRepo := mocks.NewMockColorRepository(t)
	colorRepo.EXPECT().FindAll(mock.Anything, 0, 20).Return([]entities.Color{}, 0, nil)

	uc := NewListColorsUseCase(colorRepo)
	result, err := uc.Execute(ctx, params)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Empty(t, result.Data)
	assert.Equal(t, 0, result.Meta.Total)
	assert.Equal(t, 0, result.Meta.TotalPages)
}

func TestListColors_RepoError(t *testing.T) {
	ctx := context.Background()
	params := entities.PaginationParams{Page: 1, Limit: 20}

	colorRepo := mocks.NewMockColorRepository(t)
	colorRepo.EXPECT().FindAll(mock.Anything, 0, 20).Return(nil, 0, fmt.Errorf("database connection failed"))

	uc := NewListColorsUseCase(colorRepo)
	result, err := uc.Execute(ctx, params)

	assert.Error(t, err)
	assert.Nil(t, result)
}

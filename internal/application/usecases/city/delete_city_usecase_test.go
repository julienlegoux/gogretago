package city

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

func TestDeleteCity_Success(t *testing.T) {
	ctx := context.Background()
	cityID := "city-1"

	existing := &entities.City{ID: cityID, RefID: 1, CityName: "Paris", Zipcode: "75000"}

	cityRepo := mocks.NewMockCityRepository(t)

	cityRepo.EXPECT().FindByID(mock.Anything, cityID).Return(existing, nil)
	cityRepo.EXPECT().Delete(mock.Anything, cityID).Return(nil)

	uc := NewDeleteCityUseCase(cityRepo)
	err := uc.Execute(ctx, cityID)

	assert.NoError(t, err)
}

func TestDeleteCity_NotFound(t *testing.T) {
	ctx := context.Background()
	cityID := "city-nonexistent"

	cityRepo := mocks.NewMockCityRepository(t)

	cityRepo.EXPECT().FindByID(mock.Anything, cityID).Return(nil, nil)

	uc := NewDeleteCityUseCase(cityRepo)
	err := uc.Execute(ctx, cityID)

	assert.Error(t, err)
	var notFoundErr *domainerrors.CityNotFoundError
	assert.True(t, errors.As(err, &notFoundErr))
}

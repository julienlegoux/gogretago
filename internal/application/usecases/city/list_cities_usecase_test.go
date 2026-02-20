package city

import (
	"context"
	"testing"

	"github.com/lgxju/gogretago/internal/domain/entities"
	"github.com/lgxju/gogretago/internal/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestListCities_Success(t *testing.T) {
	ctx := context.Background()
	params := entities.PaginationParams{Page: 1, Limit: 20}

	cities := []entities.City{
		{ID: "city-1", RefID: 1, CityName: "Paris", Zipcode: "75000"},
		{ID: "city-2", RefID: 2, CityName: "Lyon", Zipcode: "69000"},
	}

	cityRepo := mocks.NewMockCityRepository(t)

	cityRepo.EXPECT().FindAll(mock.Anything, 0, 20).Return(cities, 2, nil)

	uc := NewListCitiesUseCase(cityRepo)
	result, err := uc.Execute(ctx, params)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Len(t, result.Data, 2)
	assert.Equal(t, 2, result.Meta.Total)
	assert.Equal(t, 1, result.Meta.Page)
	assert.Equal(t, 20, result.Meta.Limit)
}

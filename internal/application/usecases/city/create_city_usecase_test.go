package city

import (
	"context"
	"testing"

	"github.com/lgxju/gogretago/internal/application/dtos"
	"github.com/lgxju/gogretago/internal/domain/entities"
	"github.com/lgxju/gogretago/internal/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCreateCity_Success(t *testing.T) {
	ctx := context.Background()
	input := dtos.CreateCityInput{
		CityName: "Paris",
		Zipcode:  "75000",
	}

	expectedCity := &entities.City{ID: "city-1", RefID: 1, CityName: "Paris", Zipcode: "75000"}

	cityRepo := mocks.NewMockCityRepository(t)

	cityRepo.EXPECT().Create(mock.Anything, entities.CreateCityData{
		CityName: "Paris",
		Zipcode:  "75000",
	}).Return(expectedCity, nil)

	uc := NewCreateCityUseCase(cityRepo)
	result, err := uc.Execute(ctx, input)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "city-1", result.ID)
	assert.Equal(t, "Paris", result.CityName)
	assert.Equal(t, "75000", result.Zipcode)
}

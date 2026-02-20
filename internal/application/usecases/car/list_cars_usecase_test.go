package car

import (
	"context"
	"fmt"
	"testing"

	"github.com/lgxju/gogretago/internal/domain/entities"
	"github.com/lgxju/gogretago/internal/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestListCars_Success(t *testing.T) {
	ctx := context.Background()
	params := entities.PaginationParams{Page: 1, Limit: 20}

	cars := []entities.Car{
		{ID: "car-1", RefID: 1, LicensePlate: "ABC-123", ModelRefID: 30, DriverRefID: 10},
		{ID: "car-2", RefID: 2, LicensePlate: "XYZ-789", ModelRefID: 31, DriverRefID: 11},
	}

	carRepo := mocks.NewMockCarRepository(t)

	carRepo.EXPECT().FindAll(mock.Anything, 0, 20).Return(cars, 2, nil)

	uc := NewListCarsUseCase(carRepo)
	result, err := uc.Execute(ctx, params)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Len(t, result.Data, 2)
	assert.Equal(t, 2, result.Meta.Total)
	assert.Equal(t, 1, result.Meta.Page)
	assert.Equal(t, 20, result.Meta.Limit)
}

func TestListCars_Empty(t *testing.T) {
	ctx := context.Background()
	params := entities.PaginationParams{Page: 1, Limit: 20}

	carRepo := mocks.NewMockCarRepository(t)
	carRepo.EXPECT().FindAll(mock.Anything, 0, 20).Return([]entities.Car{}, 0, nil)

	uc := NewListCarsUseCase(carRepo)
	result, err := uc.Execute(ctx, params)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Empty(t, result.Data)
	assert.Equal(t, 0, result.Meta.Total)
	assert.Equal(t, 0, result.Meta.TotalPages)
}

func TestListCars_RepoError(t *testing.T) {
	ctx := context.Background()
	params := entities.PaginationParams{Page: 1, Limit: 20}

	carRepo := mocks.NewMockCarRepository(t)
	carRepo.EXPECT().FindAll(mock.Anything, 0, 20).Return(nil, 0, fmt.Errorf("database connection failed"))

	uc := NewListCarsUseCase(carRepo)
	result, err := uc.Execute(ctx, params)

	assert.Error(t, err)
	assert.Nil(t, result)
}

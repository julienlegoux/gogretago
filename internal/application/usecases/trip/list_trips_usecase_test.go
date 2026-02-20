package trip

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/lgxju/gogretago/internal/domain/entities"
	"github.com/lgxju/gogretago/internal/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestListTrips_Success(t *testing.T) {
	ctx := context.Background()
	tripRepo := mocks.NewMockTripRepository(t)

	dateTrip, _ := time.Parse("2006-01-02", "2026-06-15")
	trips := []entities.Trip{
		{
			ID:          "trip-1",
			RefID:       500,
			DateTrip:    dateTrip,
			Kms:         450,
			Seats:       3,
			DriverRefID: 300,
			CarRefID:    400,
		},
		{
			ID:          "trip-2",
			RefID:       501,
			DateTrip:    dateTrip,
			Kms:         200,
			Seats:       2,
			DriverRefID: 301,
			CarRefID:    401,
		},
	}

	params := entities.PaginationParams{Page: 1, Limit: 20}
	// Skip() = (1-1)*20 = 0, Take() = 20
	tripRepo.EXPECT().FindAll(ctx, 0, 20).Return(trips, 2, nil)

	uc := NewListTripsUseCase(tripRepo)
	result, err := uc.Execute(ctx, params)

	require.NoError(t, err)
	assert.Len(t, result.Data, 2)
	assert.Equal(t, "trip-1", result.Data[0].ID)
	assert.Equal(t, "trip-2", result.Data[1].ID)
	assert.Equal(t, 1, result.Meta.Page)
	assert.Equal(t, 20, result.Meta.Limit)
	assert.Equal(t, 2, result.Meta.Total)
	assert.Equal(t, 1, result.Meta.TotalPages)
}

func TestListTrips_Empty(t *testing.T) {
	ctx := context.Background()
	tripRepo := mocks.NewMockTripRepository(t)

	params := entities.PaginationParams{Page: 1, Limit: 10}
	tripRepo.EXPECT().FindAll(ctx, 0, 10).Return([]entities.Trip{}, 0, nil)

	uc := NewListTripsUseCase(tripRepo)
	result, err := uc.Execute(ctx, params)

	require.NoError(t, err)
	assert.Empty(t, result.Data)
	assert.Equal(t, 0, result.Meta.Total)
	assert.Equal(t, 0, result.Meta.TotalPages)
}

func TestListTrips_RepoError(t *testing.T) {
	ctx := context.Background()
	tripRepo := mocks.NewMockTripRepository(t)

	params := entities.PaginationParams{Page: 1, Limit: 20}
	repoErr := errors.New("database error")
	tripRepo.EXPECT().FindAll(ctx, 0, 20).Return(nil, 0, repoErr)

	uc := NewListTripsUseCase(tripRepo)
	result, err := uc.Execute(ctx, params)

	assert.Nil(t, result)
	require.Error(t, err)
	assert.Equal(t, repoErr, err)
}

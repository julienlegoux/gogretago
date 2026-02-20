package trip

import (
	"context"
	"testing"
	"time"

	"github.com/lgxju/gogretago/internal/application/dtos"
	"github.com/lgxju/gogretago/internal/domain/entities"
	"github.com/lgxju/gogretago/internal/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFindTrips_AllFilters(t *testing.T) {
	ctx := context.Background()
	tripRepo := mocks.NewMockTripRepository(t)

	departure := "Paris"
	arrival := "Lyon"
	dateStr := "2026-06-15"
	parsedDate, _ := time.Parse("2006-01-02", dateStr)

	expectedTrips := []entities.Trip{
		{
			ID:          "trip-1",
			RefID:       500,
			DateTrip:    parsedDate,
			Kms:         450,
			Seats:       3,
			DriverRefID: 300,
			CarRefID:    400,
		},
	}

	tripRepo.EXPECT().FindByFilters(ctx, entities.TripFilters{
		DepartureCity: &departure,
		ArrivalCity:   &arrival,
		Date:          &parsedDate,
	}).Return(expectedTrips, nil)

	uc := NewFindTripsUseCase(tripRepo)
	result, err := uc.Execute(ctx, dtos.FindTripQuery{
		DepartureCity: &departure,
		ArrivalCity:   &arrival,
		Date:          &dateStr,
	})

	require.NoError(t, err)
	assert.Len(t, result, 1)
	assert.Equal(t, "trip-1", result[0].ID)
}

func TestFindTrips_NoFilters(t *testing.T) {
	ctx := context.Background()
	tripRepo := mocks.NewMockTripRepository(t)

	expectedTrips := []entities.Trip{
		{
			ID:    "trip-1",
			RefID: 500,
		},
		{
			ID:    "trip-2",
			RefID: 501,
		},
	}

	tripRepo.EXPECT().FindByFilters(ctx, entities.TripFilters{
		DepartureCity: nil,
		ArrivalCity:   nil,
		Date:          nil,
	}).Return(expectedTrips, nil)

	uc := NewFindTripsUseCase(tripRepo)
	result, err := uc.Execute(ctx, dtos.FindTripQuery{
		DepartureCity: nil,
		ArrivalCity:   nil,
		Date:          nil,
	})

	require.NoError(t, err)
	assert.Len(t, result, 2)
}

func TestFindTrips_InvalidDateFilter(t *testing.T) {
	ctx := context.Background()
	tripRepo := mocks.NewMockTripRepository(t)

	badDate := "not-a-valid-date"

	uc := NewFindTripsUseCase(tripRepo)
	result, err := uc.Execute(ctx, dtos.FindTripQuery{
		DepartureCity: nil,
		ArrivalCity:   nil,
		Date:          &badDate,
	})

	assert.Nil(t, result)
	require.Error(t, err)
}

func TestFindTrips_DateOnly(t *testing.T) {
	ctx := context.Background()
	tripRepo := mocks.NewMockTripRepository(t)

	dateStr := "2026-12-25"
	parsedDate, _ := time.Parse("2006-01-02", dateStr)

	expectedTrips := []entities.Trip{
		{
			ID:       "trip-1",
			RefID:    500,
			DateTrip: parsedDate,
			Kms:      100,
			Seats:    4,
		},
	}

	tripRepo.EXPECT().FindByFilters(ctx, entities.TripFilters{
		DepartureCity: nil,
		ArrivalCity:   nil,
		Date:          &parsedDate,
	}).Return(expectedTrips, nil)

	uc := NewFindTripsUseCase(tripRepo)
	result, err := uc.Execute(ctx, dtos.FindTripQuery{
		DepartureCity: nil,
		ArrivalCity:   nil,
		Date:          &dateStr,
	})

	require.NoError(t, err)
	assert.Len(t, result, 1)
	assert.Equal(t, "trip-1", result[0].ID)
	assert.Equal(t, parsedDate, result[0].DateTrip)
}

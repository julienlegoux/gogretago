package trip

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/lgxju/gogretago/internal/domain/entities"
	domainerrors "github.com/lgxju/gogretago/internal/domain/errors"
	"github.com/lgxju/gogretago/internal/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetTrip_Success(t *testing.T) {
	ctx := context.Background()
	tripRepo := mocks.NewMockTripRepository(t)

	dateTrip, _ := time.Parse("2006-01-02", "2026-06-15")
	trip := &entities.Trip{
		ID:          "trip-1",
		RefID:       500,
		DateTrip:    dateTrip,
		Kms:         450,
		Seats:       3,
		DriverRefID: 300,
		CarRefID:    400,
	}

	tripRepo.EXPECT().FindByID(ctx, "trip-1").Return(trip, nil)

	uc := NewGetTripUseCase(tripRepo)
	result, err := uc.Execute(ctx, "trip-1")

	require.NoError(t, err)
	assert.Equal(t, "trip-1", result.ID)
	assert.Equal(t, 450, result.Kms)
	assert.Equal(t, 3, result.Seats)
	assert.Equal(t, dateTrip, result.DateTrip)
}

func TestGetTrip_NotFound(t *testing.T) {
	ctx := context.Background()
	tripRepo := mocks.NewMockTripRepository(t)

	tripRepo.EXPECT().FindByID(ctx, "nonexistent").Return(nil, nil)

	uc := NewGetTripUseCase(tripRepo)
	result, err := uc.Execute(ctx, "nonexistent")

	assert.Nil(t, result)
	require.Error(t, err)
	var notFoundErr *domainerrors.TripNotFoundError
	assert.True(t, errors.As(err, &notFoundErr))
}

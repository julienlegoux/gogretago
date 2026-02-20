//go:build integration

package repositories

import (
	"context"
	"testing"
	"time"

	"github.com/lgxju/gogretago/internal/domain/entities"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// createTripPrerequisites creates all required records for a trip and returns driver RefID, car RefID,
// departure city RefID, and arrival city RefID.
func createTripPrerequisites(t *testing.T) (driverRefID, carRefID, departureCityRefID, arrivalCityRefID int64) {
	t.Helper()
	ctx := context.Background()

	driverRef, modelRef := createCarPrerequisites(t, "trip-driver@example.com")

	// Create car
	carRepo := NewGormCarRepository(testDB)
	car, err := carRepo.Create(ctx, entities.CreateCarData{
		LicensePlate: "TRIP-001",
		ModelRefID:   modelRef,
		DriverRefID:  driverRef,
	})
	require.NoError(t, err)

	// Create departure and arrival cities
	cityRepo := NewGormCityRepository(testDB)
	departure, err := cityRepo.Create(ctx, entities.CreateCityData{CityName: "Paris", Zipcode: "75000"})
	require.NoError(t, err)

	arrival, err := cityRepo.Create(ctx, entities.CreateCityData{CityName: "Lyon", Zipcode: "69000"})
	require.NoError(t, err)

	return driverRef, car.RefID, departure.RefID, arrival.RefID
}

func TestTripRepo_CRUD_Integration(t *testing.T) {
	cleanTables(t)
	t.Cleanup(func() { cleanTables(t) })

	repo := NewGormTripRepository(testDB)
	ctx := context.Background()

	driverRefID, carRefID, departureCityRefID, arrivalCityRefID := createTripPrerequisites(t)

	tripDate := time.Date(2026, 3, 15, 10, 0, 0, 0, time.UTC)

	// Create trip with cities
	trip, err := repo.Create(ctx, entities.CreateTripData{
		DateTrip:    tripDate,
		Kms:         450,
		Seats:       3,
		DriverRefID: driverRefID,
		CarRefID:    carRefID,
		CityRefIDs:  []int64{departureCityRefID, arrivalCityRefID},
	})
	require.NoError(t, err)
	require.NotNil(t, trip)
	assert.NotEmpty(t, trip.ID)
	assert.Greater(t, trip.RefID, int64(0))
	assert.Equal(t, 450, trip.Kms)
	assert.Equal(t, 3, trip.Seats)
	assert.Equal(t, driverRefID, trip.DriverRefID)
	assert.Equal(t, carRefID, trip.CarRefID)

	// FindByID
	found, err := repo.FindByID(ctx, trip.ID)
	require.NoError(t, err)
	require.NotNil(t, found)
	assert.Equal(t, trip.ID, found.ID)
	assert.Equal(t, 450, found.Kms)
	assert.Equal(t, 3, found.Seats)

	// FindByID non-existent
	notFound, err := repo.FindByID(ctx, "00000000-0000-0000-0000-000000000000")
	require.NoError(t, err)
	assert.Nil(t, notFound)

	// FindAll
	trips, total, err := repo.FindAll(ctx, 0, 10)
	require.NoError(t, err)
	assert.Equal(t, 1, total)
	assert.Len(t, trips, 1)

	// Delete (also deletes city_trips associations)
	err = repo.Delete(ctx, trip.ID)
	require.NoError(t, err)

	// Verify deleted
	deleted, err := repo.FindByID(ctx, trip.ID)
	require.NoError(t, err)
	assert.Nil(t, deleted)

	_, total, err = repo.FindAll(ctx, 0, 10)
	require.NoError(t, err)
	assert.Equal(t, 0, total)
}

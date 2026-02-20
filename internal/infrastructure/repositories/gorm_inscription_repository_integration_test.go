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

// createInscriptionPrerequisites creates the full chain of dependencies for an inscription test.
// Returns userRefID, userID (UUID), tripRefID, and tripID (UUID).
func createInscriptionPrerequisites(t *testing.T) (userRefID int64, userID string, tripRefID int64, tripID string) {
	t.Helper()
	ctx := context.Background()

	// Create trip prerequisites (driver, car, cities)
	driverRef, carRef, depCityRef, arrCityRef := createTripPrerequisites(t)

	// Create trip
	tripRepo := NewGormTripRepository(testDB)
	trip, err := tripRepo.Create(ctx, entities.CreateTripData{
		DateTrip:    time.Date(2026, 4, 1, 8, 0, 0, 0, time.UTC),
		Kms:         200,
		Seats:       4,
		DriverRefID: driverRef,
		CarRefID:    carRef,
		CityRefIDs:  []int64{depCityRef, arrCityRef},
	})
	require.NoError(t, err)

	// Create a passenger (separate auth+user from the driver)
	_, passenger := createTestAuthAndUser(t, "passenger@example.com", "Pass", "Enger", "+33611111111")

	return passenger.RefID, passenger.ID, trip.RefID, trip.ID
}

func TestInscriptionRepo_CRUD_Integration(t *testing.T) {
	cleanTables(t)
	t.Cleanup(func() { cleanTables(t) })

	repo := NewGormInscriptionRepository(testDB)
	ctx := context.Background()

	userRefID, userID, tripRefID, tripID := createInscriptionPrerequisites(t)

	// Create inscription
	inscription, err := repo.Create(ctx, entities.CreateInscriptionData{
		UserRefID: userRefID,
		TripRefID: tripRefID,
	})
	require.NoError(t, err)
	require.NotNil(t, inscription)
	assert.NotEmpty(t, inscription.ID)
	assert.Greater(t, inscription.RefID, int64(0))
	assert.Equal(t, userRefID, inscription.UserRefID)
	assert.Equal(t, tripRefID, inscription.TripRefID)
	assert.Equal(t, "ACTIVE", inscription.Status)
	assert.False(t, inscription.CreatedAt.IsZero())

	// FindByID
	found, err := repo.FindByID(ctx, inscription.ID)
	require.NoError(t, err)
	require.NotNil(t, found)
	assert.Equal(t, inscription.ID, found.ID)
	assert.Equal(t, "ACTIVE", found.Status)

	// FindByID non-existent
	notFound, err := repo.FindByID(ctx, "00000000-0000-0000-0000-000000000000")
	require.NoError(t, err)
	assert.Nil(t, notFound)

	// FindByIDAndUserID
	foundByIDAndUser, err := repo.FindByIDAndUserID(ctx, inscription.ID, userID)
	require.NoError(t, err)
	require.NotNil(t, foundByIDAndUser)
	assert.Equal(t, inscription.ID, foundByIDAndUser.ID)

	// FindByIDAndUserID with wrong user
	wrongUser, err := repo.FindByIDAndUserID(ctx, inscription.ID, "00000000-0000-0000-0000-000000000000")
	require.NoError(t, err)
	assert.Nil(t, wrongUser)

	// ExistsByUserAndTrip
	exists, err := repo.ExistsByUserAndTrip(ctx, userRefID, tripRefID)
	require.NoError(t, err)
	assert.True(t, exists)

	exists, err = repo.ExistsByUserAndTrip(ctx, 99999, tripRefID)
	require.NoError(t, err)
	assert.False(t, exists)

	// CountByTripRefID
	count, err := repo.CountByTripRefID(ctx, tripRefID)
	require.NoError(t, err)
	assert.Equal(t, 1, count)

	count, err = repo.CountByTripRefID(ctx, 99999)
	require.NoError(t, err)
	assert.Equal(t, 0, count)

	// FindByUserID
	byUser, err := repo.FindByUserID(ctx, userID)
	require.NoError(t, err)
	assert.Len(t, byUser, 1)
	assert.Equal(t, inscription.ID, byUser[0].ID)

	// FindByUserID with non-existent user returns empty slice
	byNonUser, err := repo.FindByUserID(ctx, "00000000-0000-0000-0000-000000000000")
	require.NoError(t, err)
	assert.Empty(t, byNonUser)

	// FindByTripID
	byTrip, err := repo.FindByTripID(ctx, tripID)
	require.NoError(t, err)
	assert.Len(t, byTrip, 1)
	assert.Equal(t, inscription.ID, byTrip[0].ID)

	// FindByTripID with non-existent trip returns empty slice
	byNonTrip, err := repo.FindByTripID(ctx, "00000000-0000-0000-0000-000000000000")
	require.NoError(t, err)
	assert.Empty(t, byNonTrip)

	// Delete
	err = repo.Delete(ctx, inscription.ID)
	require.NoError(t, err)

	// Verify deleted
	deleted, err := repo.FindByID(ctx, inscription.ID)
	require.NoError(t, err)
	assert.Nil(t, deleted)
}

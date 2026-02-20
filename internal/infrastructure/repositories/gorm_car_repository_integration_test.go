//go:build integration

package repositories

import (
	"context"
	"testing"

	"github.com/lgxju/gogretago/internal/domain/entities"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// createCarPrerequisites creates the required auth, user, driver, brand, and model records.
// Returns the driver RefID and model RefID needed to create a car.
func createCarPrerequisites(t *testing.T, email string) (driverRefID, modelRefID int64) {
	t.Helper()
	ctx := context.Background()

	// Create auth + user
	_, user := createTestAuthAndUser(t, email, "Car", "Owner", "+33600000099")

	// Create driver
	driverRepo := NewGormDriverRepository(testDB)
	driver, err := driverRepo.Create(ctx, entities.CreateDriverData{
		DriverLicense: "DL-CAR-" + email,
		UserRefID:     user.RefID,
	})
	require.NoError(t, err)

	// Create brand
	brandRepo := NewGormBrandRepository(testDB)
	brand, err := brandRepo.Create(ctx, entities.CreateBrandData{Name: "TestBrand-" + email})
	require.NoError(t, err)

	// Create vehicle model
	modelRepo := NewGormModelRepository(testDB)
	model, err := modelRepo.Create(ctx, entities.CreateModelData{
		Name:       "TestModel-" + email,
		BrandRefID: brand.RefID,
	})
	require.NoError(t, err)

	return driver.RefID, model.RefID
}

func TestCarRepo_CRUD_Integration(t *testing.T) {
	cleanTables(t)
	t.Cleanup(func() { cleanTables(t) })

	repo := NewGormCarRepository(testDB)
	ctx := context.Background()

	driverRefID, modelRefID := createCarPrerequisites(t, "car-owner@example.com")

	// Create
	car, err := repo.Create(ctx, entities.CreateCarData{
		LicensePlate: "AB-123-CD",
		ModelRefID:   modelRefID,
		DriverRefID:  driverRefID,
	})
	require.NoError(t, err)
	require.NotNil(t, car)
	assert.NotEmpty(t, car.ID)
	assert.Greater(t, car.RefID, int64(0))
	assert.Equal(t, "AB-123-CD", car.LicensePlate)
	assert.Equal(t, modelRefID, car.ModelRefID)
	assert.Equal(t, driverRefID, car.DriverRefID)

	// FindByID
	found, err := repo.FindByID(ctx, car.ID)
	require.NoError(t, err)
	require.NotNil(t, found)
	assert.Equal(t, car.ID, found.ID)
	assert.Equal(t, "AB-123-CD", found.LicensePlate)

	// FindByID non-existent
	notFound, err := repo.FindByID(ctx, "00000000-0000-0000-0000-000000000000")
	require.NoError(t, err)
	assert.Nil(t, notFound)

	// ExistsByLicensePlate
	exists, err := repo.ExistsByLicensePlate(ctx, "AB-123-CD")
	require.NoError(t, err)
	assert.True(t, exists)

	exists, err = repo.ExistsByLicensePlate(ctx, "ZZ-999-ZZ")
	require.NoError(t, err)
	assert.False(t, exists)

	// Update
	newPlate := "EF-456-GH"
	updated, err := repo.Update(ctx, car.ID, entities.UpdateCarData{LicensePlate: &newPlate})
	require.NoError(t, err)
	require.NotNil(t, updated)
	assert.Equal(t, "EF-456-GH", updated.LicensePlate)

	// Verify update persisted
	found, err = repo.FindByID(ctx, car.ID)
	require.NoError(t, err)
	assert.Equal(t, "EF-456-GH", found.LicensePlate)

	// Delete
	err = repo.Delete(ctx, car.ID)
	require.NoError(t, err)

	// Verify deleted
	deleted, err := repo.FindByID(ctx, car.ID)
	require.NoError(t, err)
	assert.Nil(t, deleted)
}

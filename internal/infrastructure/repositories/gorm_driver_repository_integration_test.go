//go:build integration

package repositories

import (
	"context"
	"testing"

	"github.com/lgxju/gogretago/internal/domain/entities"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDriverRepo_CRUD_Integration(t *testing.T) {
	cleanTables(t)
	t.Cleanup(func() { cleanTables(t) })

	driverRepo := NewGormDriverRepository(testDB)
	ctx := context.Background()

	// Create prerequisite: auth + user
	_, user := createTestAuthAndUser(t, "driver@example.com", "Driver", "Test", "+33600000001")

	// Create driver
	driver, err := driverRepo.Create(ctx, entities.CreateDriverData{
		DriverLicense: "DL-123456",
		UserRefID:     user.RefID,
	})
	require.NoError(t, err)
	require.NotNil(t, driver)
	assert.NotEmpty(t, driver.ID)
	assert.Greater(t, driver.RefID, int64(0))
	assert.Equal(t, "DL-123456", driver.DriverLicense)
	assert.Equal(t, user.RefID, driver.UserRefID)

	// FindByUserRefID
	foundByRefID, err := driverRepo.FindByUserRefID(ctx, user.RefID)
	require.NoError(t, err)
	require.NotNil(t, foundByRefID)
	assert.Equal(t, driver.ID, foundByRefID.ID)
	assert.Equal(t, driver.RefID, foundByRefID.RefID)
	assert.Equal(t, "DL-123456", foundByRefID.DriverLicense)

	// FindByUserRefID non-existent
	notFound, err := driverRepo.FindByUserRefID(ctx, 99999)
	require.NoError(t, err)
	assert.Nil(t, notFound)

	// FindByUserID (uses user UUID, joins with users table)
	foundByUserID, err := driverRepo.FindByUserID(ctx, user.ID)
	require.NoError(t, err)
	require.NotNil(t, foundByUserID)
	assert.Equal(t, driver.ID, foundByUserID.ID)
	assert.Equal(t, "DL-123456", foundByUserID.DriverLicense)

	// FindByUserID non-existent user
	notFound, err = driverRepo.FindByUserID(ctx, "00000000-0000-0000-0000-000000000000")
	require.NoError(t, err)
	assert.Nil(t, notFound)
}

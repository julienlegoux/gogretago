//go:build integration

package repositories

import (
	"context"
	"testing"

	"github.com/lgxju/gogretago/internal/domain/entities"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCityRepo_CRUD_Integration(t *testing.T) {
	cleanTables(t)
	t.Cleanup(func() { cleanTables(t) })

	repo := NewGormCityRepository(testDB)
	ctx := context.Background()

	// Create
	city, err := repo.Create(ctx, entities.CreateCityData{CityName: "Paris", Zipcode: "75000"})
	require.NoError(t, err)
	require.NotNil(t, city)
	assert.NotEmpty(t, city.ID)
	assert.Greater(t, city.RefID, int64(0))
	assert.Equal(t, "Paris", city.CityName)
	assert.Equal(t, "75000", city.Zipcode)

	// FindByCityName
	found, err := repo.FindByCityName(ctx, "Paris")
	require.NoError(t, err)
	require.NotNil(t, found)
	assert.Equal(t, city.ID, found.ID)
	assert.Equal(t, "Paris", found.CityName)
	assert.Equal(t, "75000", found.Zipcode)

	// FindByCityName non-existent
	notFound, err := repo.FindByCityName(ctx, "Atlantis")
	require.NoError(t, err)
	assert.Nil(t, notFound)

	// Create more cities
	_, err = repo.Create(ctx, entities.CreateCityData{CityName: "Lyon", Zipcode: "69000"})
	require.NoError(t, err)
	_, err = repo.Create(ctx, entities.CreateCityData{CityName: "Marseille", Zipcode: "13000"})
	require.NoError(t, err)

	// FindAll
	cities, total, err := repo.FindAll(ctx, 0, 10)
	require.NoError(t, err)
	assert.Equal(t, 3, total)
	assert.Len(t, cities, 3)

	// FindAll with pagination
	cities, total, err = repo.FindAll(ctx, 0, 2)
	require.NoError(t, err)
	assert.Equal(t, 3, total)
	assert.Len(t, cities, 2)

	// Delete
	err = repo.Delete(ctx, city.ID)
	require.NoError(t, err)

	// Verify deleted
	deleted, err := repo.FindByCityName(ctx, "Paris")
	require.NoError(t, err)
	assert.Nil(t, deleted)

	_, total, err = repo.FindAll(ctx, 0, 10)
	require.NoError(t, err)
	assert.Equal(t, 2, total)
}

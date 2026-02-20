//go:build integration

package repositories

import (
	"context"
	"testing"

	"github.com/lgxju/gogretago/internal/domain/entities"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBrandRepo_CRUD_Integration(t *testing.T) {
	cleanTables(t)
	t.Cleanup(func() { cleanTables(t) })

	repo := NewGormBrandRepository(testDB)
	ctx := context.Background()

	// Create
	brand, err := repo.Create(ctx, entities.CreateBrandData{Name: "Toyota"})
	require.NoError(t, err)
	require.NotNil(t, brand)
	assert.NotEmpty(t, brand.ID)
	assert.Greater(t, brand.RefID, int64(0))
	assert.Equal(t, "Toyota", brand.Name)

	// FindByID
	found, err := repo.FindByID(ctx, brand.ID)
	require.NoError(t, err)
	require.NotNil(t, found)
	assert.Equal(t, brand.ID, found.ID)
	assert.Equal(t, brand.RefID, found.RefID)
	assert.Equal(t, "Toyota", found.Name)

	// FindByID non-existent
	notFound, err := repo.FindByID(ctx, "00000000-0000-0000-0000-000000000000")
	require.NoError(t, err)
	assert.Nil(t, notFound)

	// Create more brands for pagination
	_, err = repo.Create(ctx, entities.CreateBrandData{Name: "Honda"})
	require.NoError(t, err)
	_, err = repo.Create(ctx, entities.CreateBrandData{Name: "BMW"})
	require.NoError(t, err)

	// FindAll with pagination
	brands, total, err := repo.FindAll(ctx, 0, 10)
	require.NoError(t, err)
	assert.Equal(t, 3, total)
	assert.Len(t, brands, 3)

	// FindAll with pagination (skip 1, take 1)
	brands, total, err = repo.FindAll(ctx, 1, 1)
	require.NoError(t, err)
	assert.Equal(t, 3, total)
	assert.Len(t, brands, 1)

	// Delete
	err = repo.Delete(ctx, brand.ID)
	require.NoError(t, err)

	// Verify deleted
	deleted, err := repo.FindByID(ctx, brand.ID)
	require.NoError(t, err)
	assert.Nil(t, deleted)

	// Verify total decreased
	_, total, err = repo.FindAll(ctx, 0, 10)
	require.NoError(t, err)
	assert.Equal(t, 2, total)
}

//go:build integration

package repositories

import (
	"context"
	"testing"

	"github.com/lgxju/gogretago/internal/domain/entities"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestColorRepo_CRUD_Integration(t *testing.T) {
	cleanTables(t)
	t.Cleanup(func() { cleanTables(t) })

	repo := NewGormColorRepository(testDB)
	ctx := context.Background()

	// Create
	color, err := repo.Create(ctx, entities.CreateColorData{Name: "Red", Hex: "#FF0000"})
	require.NoError(t, err)
	require.NotNil(t, color)
	assert.NotEmpty(t, color.ID)
	assert.Greater(t, color.RefID, int64(0))
	assert.Equal(t, "Red", color.Name)
	assert.Equal(t, "#FF0000", color.Hex)

	// FindByName
	found, err := repo.FindByName(ctx, "Red")
	require.NoError(t, err)
	require.NotNil(t, found)
	assert.Equal(t, color.ID, found.ID)
	assert.Equal(t, "Red", found.Name)
	assert.Equal(t, "#FF0000", found.Hex)

	// FindByName non-existent
	notFound, err := repo.FindByName(ctx, "Invisible")
	require.NoError(t, err)
	assert.Nil(t, notFound)

	// FindByID
	foundByID, err := repo.FindByID(ctx, color.ID)
	require.NoError(t, err)
	require.NotNil(t, foundByID)
	assert.Equal(t, color.ID, foundByID.ID)
	assert.Equal(t, "Red", foundByID.Name)

	// Update
	newName := "Crimson"
	newHex := "#DC143C"
	updated, err := repo.Update(ctx, color.ID, entities.UpdateColorData{Name: &newName, Hex: &newHex})
	require.NoError(t, err)
	require.NotNil(t, updated)
	assert.Equal(t, "Crimson", updated.Name)
	assert.Equal(t, "#DC143C", updated.Hex)

	// Verify update persisted
	found, err = repo.FindByID(ctx, color.ID)
	require.NoError(t, err)
	assert.Equal(t, "Crimson", found.Name)
	assert.Equal(t, "#DC143C", found.Hex)

	// Create more colors for FindAll
	_, err = repo.Create(ctx, entities.CreateColorData{Name: "Blue", Hex: "#0000FF"})
	require.NoError(t, err)
	_, err = repo.Create(ctx, entities.CreateColorData{Name: "Green", Hex: "#00FF00"})
	require.NoError(t, err)

	// FindAll
	colors, total, err := repo.FindAll(ctx, 0, 10)
	require.NoError(t, err)
	assert.Equal(t, 3, total)
	assert.Len(t, colors, 3)

	// Delete
	err = repo.Delete(ctx, color.ID)
	require.NoError(t, err)

	// Verify deleted
	deleted, err := repo.FindByID(ctx, color.ID)
	require.NoError(t, err)
	assert.Nil(t, deleted)

	_, total, err = repo.FindAll(ctx, 0, 10)
	require.NoError(t, err)
	assert.Equal(t, 2, total)
}

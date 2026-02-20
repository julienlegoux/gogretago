//go:build integration

package repositories

import (
	"context"
	"testing"

	"github.com/lgxju/gogretago/internal/domain/entities"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAuthRepo_CreateWithUser_Integration(t *testing.T) {
	cleanTables(t)
	t.Cleanup(func() { cleanTables(t) })

	repo := NewGormAuthRepository(testDB)
	ctx := context.Background()

	firstName := "John"
	lastName := "Doe"
	phone := "+33612345678"

	auth, user, err := repo.CreateWithUser(ctx,
		entities.CreateAuthData{Email: "john@example.com", Password: "hashed_pw"},
		entities.CreateUserData{FirstName: &firstName, LastName: &lastName, Phone: &phone},
	)

	require.NoError(t, err)
	require.NotNil(t, auth)
	require.NotNil(t, user)

	// Auth should have a valid UUID and RefID
	assert.NotEmpty(t, auth.ID)
	assert.Greater(t, auth.RefID, int64(0))
	assert.Equal(t, "john@example.com", auth.Email)
	assert.Equal(t, "hashed_pw", auth.Password)
	assert.Equal(t, "USER", auth.Role)
	assert.False(t, auth.CreatedAt.IsZero())
	assert.False(t, auth.UpdatedAt.IsZero())

	// User should have a valid UUID and RefID
	assert.NotEmpty(t, user.ID)
	assert.Greater(t, user.RefID, int64(0))
	assert.Equal(t, &firstName, user.FirstName)
	assert.Equal(t, &lastName, user.LastName)
	assert.Equal(t, &phone, user.Phone)
	assert.Equal(t, auth.RefID, user.AuthRefID)
	assert.Equal(t, "john@example.com", user.Email)
}

func TestAuthRepo_FindByEmail_Integration(t *testing.T) {
	cleanTables(t)
	t.Cleanup(func() { cleanTables(t) })

	repo := NewGormAuthRepository(testDB)
	ctx := context.Background()

	// Create a record first
	firstName := "Jane"
	lastName := "Smith"
	_, _, err := repo.CreateWithUser(ctx,
		entities.CreateAuthData{Email: "jane@example.com", Password: "secret"},
		entities.CreateUserData{FirstName: &firstName, LastName: &lastName},
	)
	require.NoError(t, err)

	// Find by existing email
	found, err := repo.FindByEmail(ctx, "jane@example.com")
	require.NoError(t, err)
	require.NotNil(t, found)
	assert.Equal(t, "jane@example.com", found.Email)
	assert.Equal(t, "secret", found.Password)
	assert.Equal(t, "USER", found.Role)
	assert.NotEmpty(t, found.ID)
	assert.Greater(t, found.RefID, int64(0))

	// Find non-existent email returns nil
	notFound, err := repo.FindByEmail(ctx, "nonexistent@example.com")
	require.NoError(t, err)
	assert.Nil(t, notFound)
}

func TestAuthRepo_ExistsByEmail_Integration(t *testing.T) {
	cleanTables(t)
	t.Cleanup(func() { cleanTables(t) })

	repo := NewGormAuthRepository(testDB)
	ctx := context.Background()

	// Create a record
	firstName := "Alice"
	_, _, err := repo.CreateWithUser(ctx,
		entities.CreateAuthData{Email: "alice@example.com", Password: "pw"},
		entities.CreateUserData{FirstName: &firstName},
	)
	require.NoError(t, err)

	// Exists returns true for existing email
	exists, err := repo.ExistsByEmail(ctx, "alice@example.com")
	require.NoError(t, err)
	assert.True(t, exists)

	// Exists returns false for non-existent email
	exists, err = repo.ExistsByEmail(ctx, "nobody@example.com")
	require.NoError(t, err)
	assert.False(t, exists)
}

func TestAuthRepo_UpdateRole_Integration(t *testing.T) {
	cleanTables(t)
	t.Cleanup(func() { cleanTables(t) })

	repo := NewGormAuthRepository(testDB)
	ctx := context.Background()

	// Create a record
	firstName := "Bob"
	auth, _, err := repo.CreateWithUser(ctx,
		entities.CreateAuthData{Email: "bob@example.com", Password: "pw"},
		entities.CreateUserData{FirstName: &firstName},
	)
	require.NoError(t, err)
	assert.Equal(t, "USER", auth.Role)

	// Update role to ADMIN
	err = repo.UpdateRole(ctx, auth.RefID, "ADMIN")
	require.NoError(t, err)

	// Verify the role was persisted
	updated, err := repo.FindByEmail(ctx, "bob@example.com")
	require.NoError(t, err)
	require.NotNil(t, updated)
	assert.Equal(t, "ADMIN", updated.Role)
}

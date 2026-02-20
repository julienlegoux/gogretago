//go:build integration

package repositories

import (
	"context"
	"testing"

	"github.com/lgxju/gogretago/internal/domain/entities"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// createTestAuthAndUser is a helper that creates an auth+user pair for user repository tests.
func createTestAuthAndUser(t *testing.T, email, first, last, phone string) (*entities.Auth, *entities.PublicUser) {
	t.Helper()
	authRepo := NewGormAuthRepository(testDB)
	ctx := context.Background()

	firstName := first
	lastName := last
	ph := phone

	auth, user, err := authRepo.CreateWithUser(ctx,
		entities.CreateAuthData{Email: email, Password: "hashed"},
		entities.CreateUserData{FirstName: &firstName, LastName: &lastName, Phone: &ph},
	)
	require.NoError(t, err)
	return auth, user
}

func TestUserRepo_FindByID_Integration(t *testing.T) {
	cleanTables(t)
	t.Cleanup(func() { cleanTables(t) })

	repo := NewGormUserRepository(testDB)
	ctx := context.Background()

	_, user := createTestAuthAndUser(t, "findme@example.com", "Find", "Me", "+33600000001")

	// Find by ID
	found, err := repo.FindByID(ctx, user.ID)
	require.NoError(t, err)
	require.NotNil(t, found)
	assert.Equal(t, user.ID, found.ID)
	assert.Equal(t, user.RefID, found.RefID)
	assert.Equal(t, "Find", *found.FirstName)
	assert.Equal(t, "Me", *found.LastName)
	assert.Equal(t, "+33600000001", *found.Phone)
	assert.Equal(t, "findme@example.com", found.Email)

	// Find non-existent returns nil
	notFound, err := repo.FindByID(ctx, "00000000-0000-0000-0000-000000000000")
	require.NoError(t, err)
	assert.Nil(t, notFound)
}

func TestUserRepo_FindAll_Integration(t *testing.T) {
	cleanTables(t)
	t.Cleanup(func() { cleanTables(t) })

	repo := NewGormUserRepository(testDB)
	ctx := context.Background()

	// Create multiple users
	createTestAuthAndUser(t, "user1@example.com", "User", "One", "+33600000001")
	createTestAuthAndUser(t, "user2@example.com", "User", "Two", "+33600000002")
	createTestAuthAndUser(t, "user3@example.com", "User", "Three", "+33600000003")

	users, err := repo.FindAll(ctx)
	require.NoError(t, err)
	assert.Len(t, users, 3)

	// Verify all users have emails populated
	emails := make(map[string]bool)
	for _, u := range users {
		emails[u.Email] = true
	}
	assert.True(t, emails["user1@example.com"])
	assert.True(t, emails["user2@example.com"])
	assert.True(t, emails["user3@example.com"])
}

func TestUserRepo_Update_Integration(t *testing.T) {
	cleanTables(t)
	t.Cleanup(func() { cleanTables(t) })

	repo := NewGormUserRepository(testDB)
	ctx := context.Background()

	_, user := createTestAuthAndUser(t, "update@example.com", "Old", "Name", "+33600000000")

	newFirst := "New"
	newLast := "Updated"
	newPhone := "+33699999999"

	updated, err := repo.Update(ctx, user.ID, entities.UpdateUserData{
		FirstName: &newFirst,
		LastName:  &newLast,
		Phone:     &newPhone,
	})
	require.NoError(t, err)
	require.NotNil(t, updated)
	assert.Equal(t, "New", *updated.FirstName)
	assert.Equal(t, "Updated", *updated.LastName)
	assert.Equal(t, "+33699999999", *updated.Phone)

	// Verify persisted
	found, err := repo.FindByID(ctx, user.ID)
	require.NoError(t, err)
	assert.Equal(t, "New", *found.FirstName)
	assert.Equal(t, "Updated", *found.LastName)
	assert.Equal(t, "+33699999999", *found.Phone)
}

func TestUserRepo_Anonymize_Integration(t *testing.T) {
	cleanTables(t)
	t.Cleanup(func() { cleanTables(t) })

	repo := NewGormUserRepository(testDB)
	ctx := context.Background()

	_, user := createTestAuthAndUser(t, "anon@example.com", "Anon", "User", "+33600000000")

	err := repo.Anonymize(ctx, user.ID)
	require.NoError(t, err)

	// Verify fields are cleared and AnonymizedAt is set
	found, err := repo.FindByID(ctx, user.ID)
	require.NoError(t, err)
	require.NotNil(t, found)
	assert.Nil(t, found.FirstName)
	assert.Nil(t, found.LastName)
	assert.Nil(t, found.Phone)
	assert.NotNil(t, found.AnonymizedAt)
}

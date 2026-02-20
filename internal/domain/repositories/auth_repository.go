package repositories

import (
	"context"

	"github.com/lgxju/gogretago/internal/domain/entities"
)

// AuthRepository defines the interface for auth persistence operations
type AuthRepository interface {
	FindByEmail(ctx context.Context, email string) (*entities.Auth, error)
	CreateWithUser(ctx context.Context, authData entities.CreateAuthData, userData entities.CreateUserData) (*entities.Auth, *entities.PublicUser, error)
	ExistsByEmail(ctx context.Context, email string) (bool, error)
	UpdateRole(ctx context.Context, refID int64, role string) error
}

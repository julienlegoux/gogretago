package repositories

import (
	"context"

	"github.com/lgxju/gogretago/internal/domain/entities"
)

// UserRepository defines the interface for user persistence operations
type UserRepository interface {
	FindByID(ctx context.Context, id string) (*entities.User, error)
	FindByEmail(ctx context.Context, email string) (*entities.User, error)
	Create(ctx context.Context, data entities.CreateUserData) (*entities.User, error)
	ExistsByEmail(ctx context.Context, email string) (bool, error)
}

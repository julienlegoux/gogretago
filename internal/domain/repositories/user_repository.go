package repositories

import (
	"context"

	"github.com/lgxju/gogretago/internal/domain/entities"
)

// UserRepository defines the interface for user persistence operations
type UserRepository interface {
	FindAll(ctx context.Context) ([]entities.PublicUser, error)
	FindByID(ctx context.Context, id string) (*entities.PublicUser, error)
	FindByAuthRefID(ctx context.Context, authRefID int64) (*entities.PublicUser, error)
	Update(ctx context.Context, id string, data entities.UpdateUserData) (*entities.PublicUser, error)
	Delete(ctx context.Context, id string) error
	Anonymize(ctx context.Context, id string) error
}

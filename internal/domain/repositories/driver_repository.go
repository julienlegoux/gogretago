package repositories

import (
	"context"

	"github.com/lgxju/gogretago/internal/domain/entities"
)

// DriverRepository defines the interface for driver persistence operations
type DriverRepository interface {
	FindByUserRefID(ctx context.Context, userRefID int64) (*entities.Driver, error)
	FindByUserID(ctx context.Context, userID string) (*entities.Driver, error)
	Create(ctx context.Context, data entities.CreateDriverData) (*entities.Driver, error)
}

package repositories

import (
	"context"

	"github.com/lgxju/gogretago/internal/domain/entities"
)

// BrandRepository defines the interface for brand persistence operations
type BrandRepository interface {
	FindAll(ctx context.Context, skip, take int) ([]entities.Brand, int, error)
	FindByID(ctx context.Context, id string) (*entities.Brand, error)
	Create(ctx context.Context, data entities.CreateBrandData) (*entities.Brand, error)
	Delete(ctx context.Context, id string) error
}

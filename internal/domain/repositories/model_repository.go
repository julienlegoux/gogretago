package repositories

import (
	"context"

	"github.com/lgxju/gogretago/internal/domain/entities"
)

// ModelRepository defines the interface for car model persistence operations
type ModelRepository interface {
	FindAll(ctx context.Context) ([]entities.VehicleModel, error)
	FindByID(ctx context.Context, id string) (*entities.VehicleModel, error)
	FindByNameAndBrand(ctx context.Context, name string, brandRefID int64) (*entities.VehicleModel, error)
	Create(ctx context.Context, data entities.CreateModelData) (*entities.VehicleModel, error)
}

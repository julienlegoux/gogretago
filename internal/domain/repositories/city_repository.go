package repositories

import (
	"context"

	"github.com/lgxju/gogretago/internal/domain/entities"
)

// CityRepository defines the interface for city persistence operations
type CityRepository interface {
	FindAll(ctx context.Context, skip, take int) ([]entities.City, int, error)
	FindByID(ctx context.Context, id string) (*entities.City, error)
	FindByCityName(ctx context.Context, name string) (*entities.City, error)
	Create(ctx context.Context, data entities.CreateCityData) (*entities.City, error)
	Delete(ctx context.Context, id string) error
}

package repositories

import (
	"context"

	"github.com/lgxju/gogretago/internal/domain/entities"
)

// TripRepository defines the interface for trip persistence operations
type TripRepository interface {
	FindAll(ctx context.Context, skip, take int) ([]entities.Trip, int, error)
	FindByID(ctx context.Context, id string) (*entities.Trip, error)
	FindByFilters(ctx context.Context, filters entities.TripFilters) ([]entities.Trip, error)
	Create(ctx context.Context, data entities.CreateTripData) (*entities.Trip, error)
	Delete(ctx context.Context, id string) error
}

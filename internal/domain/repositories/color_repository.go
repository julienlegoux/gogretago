package repositories

import (
	"context"

	"github.com/lgxju/gogretago/internal/domain/entities"
)

// ColorRepository defines the interface for color persistence operations
type ColorRepository interface {
	FindAll(ctx context.Context, skip, take int) ([]entities.Color, int, error)
	FindByID(ctx context.Context, id string) (*entities.Color, error)
	FindByName(ctx context.Context, name string) (*entities.Color, error)
	Create(ctx context.Context, data entities.CreateColorData) (*entities.Color, error)
	Update(ctx context.Context, id string, data entities.UpdateColorData) (*entities.Color, error)
	Delete(ctx context.Context, id string) error
}
